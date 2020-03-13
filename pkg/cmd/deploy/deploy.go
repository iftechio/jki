package deploy

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/iftechio/jki/pkg/factory"
	"github.com/iftechio/jki/pkg/image"
	"github.com/iftechio/jki/pkg/utils"
)

func formatGroupKind(gk schema.GroupKind) string {
	kind := strings.ToLower(gk.Kind)
	if len(gk.Group) != 0 {
		return kind + "." + gk.Group
	}
	return kind
}

type Options struct {
	container string
	image     string
	dryRun    bool

	namespace string
	spec      string
	builder   *resource.Builder
}

func (o *Options) Complete(f factory.Factory, cmd *cobra.Command, args []string) error {
	// deploy <image>
	// deploy name <image>
	// deploy resource <image>
	// deploy resource/name <image>
	// deploy resource name <image>
	var (
		img  image.Image
		spec string
		err  error
	)

	switch len(args) {
	case 1:
		img = image.FromString(args[0])
		spec = "deployment.apps/" + strings.ToLower(img.Repo)
	case 2:
		img = image.FromString(args[1])
		if strings.ContainsRune(args[0], '/') {
			// resource with name
			spec = args[0]
		} else {
			resourceOrName := args[0]
			// args[0] may be resource or name
			mapper, err := f.ToRESTMapper()
			if err != nil {
				return err
			}
			_, gr := schema.ParseResourceArg(resourceOrName)
			_, err = mapper.ResourceFor(gr.WithVersion(""))
			if err != nil {
				// may be a name
				spec = "deployment.apps/" + resourceOrName
			} else {
				// is a valid resource
				spec = fmt.Sprintf("%s/%s", resourceOrName, strings.ToLower(img.Repo))
			}
		}
	case 3:
		img = image.FromString(args[2])
		spec = args[0] + "/" + args[1]
	default:
		return fmt.Errorf("unknown args: %v", args)
	}
	o.spec = spec
	o.image = img.String()

	o.namespace, _, err = f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	o.builder = f.NewBuilder()
	return nil
}

func (o *Options) Validate() error {
	return nil
}

func (o *Options) Run() error {
	result := o.builder.
		WithScheme(scheme.Scheme, scheme.Scheme.PrioritizedVersionsAllGroups()...).
		ContinueOnError().
		NamespaceParam(o.namespace).DefaultNamespace().
		Flatten().
		ResourceTypeOrNameArgs(false, o.spec).
		Latest().
		Do()

	if err := result.Err(); err != nil {
		return err
	}

	encoder := unstructured.NewJSONFallbackEncoder(scheme.Codecs.LegacyCodec(scheme.Scheme.PrioritizedVersionsAllGroups()...))
	err := result.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}
		before, err := runtime.Encode(encoder, info.Object)
		if err != nil {
			return err
		}
		_, err = updatePodSpecForObject(info.Object, func(spec *v1.PodSpec) error {
			totalContainers := len(spec.InitContainers) + len(spec.Containers)
			if len(o.container) == 0 {
				if totalContainers == 1 && len(spec.Containers) > 0 {
					spec.Containers[0].Image = o.image
					return nil
				}
				return fmt.Errorf("ambiguous container: please specify container name")
			}
			for i, ct := range spec.InitContainers {
				if ct.Name == o.container {
					spec.InitContainers[i].Image = o.image
					return nil
				}
			}
			for i, ct := range spec.Containers {
				if ct.Name == o.container {
					spec.Containers[i].Image = o.image
					return nil
				}
			}
			return fmt.Errorf("container not found: %s", o.container)
		})
		if err != nil {
			return err
		}
		after, err := runtime.Encode(encoder, info.Object)
		if err != nil {
			return err
		}

		patch, err := strategicpatch.CreateTwoWayMergePatch(before, after, info.Object)
		if err != nil {
			return err
		}

		if string(patch) == "{}" || len(patch) == 0 {
			// no changes
			return nil
		}

		gk := info.Mapping.GroupVersionKind.GroupKind()
		if o.dryRun {
			fmt.Printf("%s/%s image updated (dry run)\n", formatGroupKind(gk), info.Name)
			return nil
		}

		_, err = resource.NewHelper(info.Client, info.Mapping).Patch(info.Namespace, info.Name, types.StrategicMergePatchType, patch, nil)
		if err != nil {
			return err
		}
		fmt.Printf("%s/%s image updated\n", formatGroupKind(gk), info.Name)
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to patch image update to pod template: %v", err)
	}
	return nil
}

func NewCmdDeploy(f factory.Factory) *cobra.Command {
	o := Options{}
	cmd := &cobra.Command{
		Use:     "deploy [resource/name] <image>",
		Short:   "Update container image of resources",
		Aliases: []string{"d"},
		Long: `Update existing container image of resources.

  Possible resources include (case insensitive):

  pod (po), deployment (deploy), daemonset (ds), statefulset (sts), cronjob (cj)

If there are multiple containers in the pod, you MUST specify the target container name.`,
		Example: ` # Update image of deployment nginx to 'nginx:alpine'
  jki deploy nginx:alpine

  # Update image of statefulset nginx to 'nginx:1-alpine'
  jki deploy sts nginx:1-alpine

  # Update image of daemonset foo to 'nginx:1-alpine'
  jki deploy ds/foo nginx:alpine

  # Update image of cronjob bar to 'alpine:3.10'
  jki deploy cronjob bar alpine:3.10

  # Update image of app container of deployment nginx to 'nginx:alpine'
  jki deploy -c app nginx:alpine`,
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckError(o.Complete(f, cmd, args))
			utils.CheckError(o.Validate())
			utils.CheckError(o.Run())
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&o.container, "container", "c", "", "The target container name")
	flags.BoolVar(&o.dryRun, "dry-run", false, "If true, only print the object that would be sent, without sending it.")
	return cmd
}
