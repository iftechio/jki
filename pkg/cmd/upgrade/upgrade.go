package upgrade

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/iftechio/jki/pkg/factory"
	"github.com/iftechio/jki/pkg/info"
	"github.com/iftechio/jki/pkg/utils"
)

type releaseAsset struct {
	Name        string    `json:"name"`  // "jki_0.2.4_darwin_amd64.tar.gz"
	State       string    `json:"state"` // "uploaded"
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DownloadURL string    `json:"browser_download_url"`
}
type releaseResponse struct {
	TagName string `json:"tag_name"`
	Assets  []releaseAsset
}

type upgradeOptions struct {
	timeout   time.Duration
	noConfirm bool

	client   *http.Client
	selfPath string
}

func (o *upgradeOptions) Complete(f factory.Factory, cmd *cobra.Command, args []string) error {
	o.client = &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 30,
			TLSHandshakeTimeout:   time.Second * 20,
			ExpectContinueTimeout: time.Second * 10,
		},
	}
	if o.timeout > 0 {
		o.client.Timeout = o.timeout
	}
	myPath, err := os.Executable()
	if err != nil {
		return err
	}
	myPath, err = filepath.EvalSymlinks(myPath)
	if err != nil {
		return err
	}
	o.selfPath = myPath
	return nil
}

func (o *upgradeOptions) Validate(args []string) error {
	return validatePath(o.selfPath)
}

func getLatestRelease(client *http.Client, url string) (*releaseResponse, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	body := &releaseResponse{}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func writeWithProgress(total int64, dst io.Writer, src io.Reader) error {
	var (
		err     error
		written int64
	)
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
				fmt.Printf("\033[2K\rProgress: %.2f%%", float64(written)*100/float64(total))
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	if err == nil {
		return nil
	}
	return err
}

func downloadAsset(client *http.Client, asset releaseAsset) (string, error) {
	fp := filepath.Join(os.TempDir(), asset.Name)
	fi, err := os.Stat(fp)
	if err == nil && fi.Size() == asset.Size {
		// already downloaded
		return fp, nil
	}

	f, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return "", err
	}
	defer f.Close()
	resp, err := client.Get(asset.DownloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download asset: unexpected status: %d", resp.StatusCode)
	}
	err = writeWithProgress(asset.Size, f, resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Printf("Asset has been saved at %s", fp)
	return fp, err
}

func extractAsset(fp string) error {
	f, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(gzf)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch header.Typeflag {
		case tar.TypeDir: // dir
			continue
		case tar.TypeReg: // regular file
			binFp := filepath.Join(os.TempDir(), "jki")
			binFile, err := os.OpenFile(binFp, os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				return err
			}
			defer binFile.Close()
			_, err = io.Copy(binFile, tarReader)
			if err != nil {
				return err
			}
			err = os.Chmod(binFp, 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *upgradeOptions) Run() error {
	release, err := getLatestRelease(o.client, "https://api.github.com/repos/iftechio/jki/releases/latest")
	if err != nil {
		return err
	}
	if release.TagName == "v"+info.Version {
		fmt.Printf("v%s is already the latest version.", info.Version)
		return nil
	}
	intact := true
	for _, asset := range release.Assets {
		if !strings.Contains(asset.Name, runtime.GOOS) {
			continue
		}
		intact = false
		if !o.noConfirm {
			fmt.Printf("Latest version is: %s\n", release.TagName)
			fmt.Printf("Updated at: %s\n", asset.UpdatedAt)
			ans := utils.Prompt("Want to upgrade? (Y/n) ")
			if strings.ToLower(ans) == "n" {
				return nil
			}
		}
		fp, err := downloadAsset(o.client, asset)
		if err != nil {
			return err
		}
		err = extractAsset(fp)
		if err != nil {
			return err
		}
		err = os.Rename(filepath.Join(os.TempDir(), "jki"), o.selfPath)
		if err != nil {
			return err
		}
		fmt.Printf("Successfully upgrade to %s!", release.TagName)
	}
	if intact {
		fmt.Printf("OS not supported: %s\n", runtime.GOOS)
	}
	return nil
}

func NewCmdUpgrade(f factory.Factory) *cobra.Command {
	o := &upgradeOptions{}
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade to the latest version",
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckError(o.Complete(f, cmd, args))
			utils.CheckError(o.Validate(args))
			utils.CheckError(o.Run())
		},
	}

	flags := cmd.Flags()
	flags.DurationVar(&o.timeout, "timeout", 0, "Specify timeout. (Defaults to no limit)")
	flags.BoolVarP(&o.noConfirm, "no-confirm", "y", false, "Answer yes for all questions")
	return cmd
}
