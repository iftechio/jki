package build

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
	"time"

	"github.com/containerd/console"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	controlapi "github.com/moby/buildkit/api/services/control"
	bkclient "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/filesync"
	"github.com/moby/buildkit/util/progress/progressui"
	fsutiltypes "github.com/tonistiigi/fsutil/types"
	"golang.org/x/sync/errgroup"
)

func writeSolveStatusToChannel(displayCh chan *bkclient.SolveStatus) func(jsonmessage.JSONMessage) {
	return func(msg jsonmessage.JSONMessage) {
		var resp controlapi.StatusResponse

		if msg.ID != "moby.buildkit.trace" {
			return
		}

		var dt []byte
		// ignoring all messages that are not understood
		if err := json.Unmarshal(*msg.Aux, &dt); err != nil {
			return
		}
		if err := resp.Unmarshal(dt); err != nil {
			return
		}

		s := bkclient.SolveStatus{}
		for _, v := range resp.Vertexes {
			s.Vertexes = append(s.Vertexes, &bkclient.Vertex{
				Digest:    v.Digest,
				Inputs:    v.Inputs,
				Name:      v.Name,
				Started:   v.Started,
				Completed: v.Completed,
				Error:     v.Error,
				Cached:    v.Cached,
			})
		}
		for _, v := range resp.Statuses {
			s.Statuses = append(s.Statuses, &bkclient.VertexStatus{
				ID:        v.ID,
				Vertex:    v.Vertex,
				Name:      v.Name,
				Total:     v.Total,
				Current:   v.Current,
				Timestamp: v.Timestamp,
				Started:   v.Started,
				Completed: v.Completed,
			})
		}
		for _, v := range resp.Logs {
			s.Logs = append(s.Logs, &bkclient.VertexLog{
				Vertex:    v.Vertex,
				Stream:    int(v.Stream),
				Data:      v.Msg,
				Timestamp: v.Timestamp,
			})
		}

		displayCh <- &s
	}
}

func resetUIDAndGID(_ string, s *fsutiltypes.Stat) bool {
	s.Uid = 0
	s.Gid = 0
	return true
}

func trySession(contextDir string) (*session.Session, error) {
	s, err := session.NewSession(context.TODO(), path.Base(contextDir), contextDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %s", err)
	}
	return s, nil
}

func (o *BuildOptions) runBuildKit(ctx context.Context, buildOpts types.ImageBuildOptions) error {
	s, err := trySession(o.context)
	if err != nil {
		return err
	}
	var dockerfileDir string
	if len(o.dockerFileName) == 0 {
		dockerfileDir = o.context
	} else {
		dockerfileDir = path.Dir(o.dockerFileName)
	}
	s.Allow(filesync.NewFSSyncProvider([]filesync.SyncedDir{
		{
			Name: "context",
			Dir:  o.context,
			Map:  resetUIDAndGID,
		},
		{
			Name: "dockerfile",
			Dir:  dockerfileDir,
		},
	}))
	s.Allow(NewAuthProvider(o.allRegistries))

	eg, ctx := errgroup.WithContext(ctx)
	dialSession := func(ctx context.Context, proto string, meta map[string][]string) (net.Conn, error) {
		return o.dockerClient.DialHijack(ctx, "/session", proto, meta)
	}
	eg.Go(func() error {
		return s.Run(context.TODO(), dialSession)
	})

	eg.Go(func() error {
		defer s.Close()
		buildOpts.Version = types.BuilderBuildKit
		buildOpts.RemoteContext = "client-session"
		buildOpts.SessionID = s.ID()
		buildOpts.BuildID = time.Now().String()
		buildOpts.Dockerfile = path.Base(o.dockerFileName)

		response, err := o.dockerClient.ImageBuild(ctx, nil, buildOpts)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		done := make(chan struct{})
		displayCh := make(chan *bkclient.SolveStatus)
		defer close(done)
		defer close(displayCh)

		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return o.dockerClient.BuildCancel(ctx, buildOpts.BuildID)
			case <-done:
			}
			return nil
		})

		out := os.Stderr

		var c console.Console
		if cons, err := console.ConsoleFromFile(out); err == nil {
			c = cons
		}
		eg.Go(func() error {
			return progressui.DisplaySolveStatus(ctx, "", c, out, displayCh)
		})
		writeAux := writeSolveStatusToChannel(displayCh)
		termFd, isTerm := term.GetFdInfo(os.Stdout)
		return jsonmessage.DisplayJSONMessagesStream(response.Body, os.Stdout, termFd, isTerm, writeAux)
	})

	return eg.Wait()
}
