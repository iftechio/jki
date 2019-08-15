package transferimage

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/iftechio/jki/pkg/cmd/cp"
	cmdutils "github.com/iftechio/jki/pkg/cmd/utils"
	"github.com/iftechio/jki/pkg/image"
	"github.com/iftechio/jki/pkg/registry"
	"github.com/spf13/cobra"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func hasErrPullingContainer(pod apiv1.Pod) bool {
	for _, st := range pod.Status.ContainerStatuses {
		if st.State.Waiting.Reason == "ErrImagePull" || st.State.Waiting.Reason == "ImagePullBackOff" {
			return true
		}
	}
	return false
}

type brokenObject struct {
	Kind  string
	Name  string
	Image string
}
type transferImageOptions struct {
	namespace   string
	kubeClient  *kubernetes.Clientset
	cp          *cp.CopyOptions
	dstRegistry *registry.Registry
}

func (o *transferImageOptions) Complete(f cmdutils.Factory) error {
	o.cp = cp.NewCopyOptions()
	dstReg, registries, err := f.LoadRegistries()
	if err != nil {
		return err
	}
	o.dstRegistry = registries[dstReg]
	o.kubeClient, err = f.KubeClient()
	if err != nil {
		return err
	}
	return o.cp.Complete(f, nil, nil)
}

func newTransferImageOptions() *transferImageOptions {
	return &transferImageOptions{}
}

func (o *transferImageOptions) fixPodSpec(podSpec *apiv1.PodTemplateSpec, it brokenObject, domain string) {
	for i, con := range podSpec.Spec.Containers {
		if con.Image == it.Image {
			// copy to accessable registry
			img := image.FromString(it.Image)
			o.cp.Run([]string{it.Image})

			// replace with new domain
			img.Domain = domain
			podSpec.Spec.Containers[i].Image = img.String()
			fmt.Printf("Transfered %s to %s\n", it.Image, img.String())
		}
	}
}

func (o *transferImageOptions) Run() (err error) {
	fmt.Printf("Searching for deploy/ds to fix in namespace: %s\n", o.namespace)
	var itemsToFix []brokenObject

	pendingPodsList, err := o.kubeClient.CoreV1().Pods(o.namespace).List(metav1.ListOptions{
		FieldSelector: "status.phase=Pending",
	})

	errPullingPods := make([]apiv1.Pod, 0)

	// filter pods
	for _, pod := range pendingPodsList.Items {
		if hasErrPullingContainer(pod) {
			errPullingPods = append(errPullingPods, pod)
		}
	}

	// find k8s objects need to fix image for
	itemsMap := make(map[string]bool)
	for _, pod := range errPullingPods {
		owner := pod.OwnerReferences[0]
		if owner.Kind == "ReplicaSet" {
			rs, err := o.kubeClient.AppsV1().ReplicaSets(o.namespace).Get(owner.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			owner = rs.OwnerReferences[0]
		}
		for _, con := range pod.Status.ContainerStatuses {
			if con.State.Waiting.Reason == "ImagePullBackOff" || con.State.Waiting.Reason == "ErrImagePull" {
				// put unique items
				if !itemsMap[owner.Kind+"/"+owner.Name] {
					itemsMap[owner.Kind+"/"+owner.Name] = true
					itemsToFix = append(itemsToFix, brokenObject{
						Kind:  owner.Kind,
						Name:  owner.Name,
						Image: con.Image,
					})
				}
			}
		}
	}

	deploymentClient := o.kubeClient.AppsV1().Deployments(o.namespace)
	dsClient := o.kubeClient.AppsV1().DaemonSets(o.namespace)
	stsClient := o.kubeClient.AppsV1().StatefulSets(o.namespace)

	if len(itemsToFix) == 0 {
		fmt.Println("Found no image to fix")
		return nil
	}
	for _, it := range itemsToFix {
		fmt.Printf("Transfer %s/%s %s(y/n)?\n", it.Kind, it.Name, it.Image)
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		sentence, err := buf.ReadBytes('\n')
		if err != nil {
			return err
		}
		if strings.ToLower(strings.TrimSpace(string(sentence))) == "y" {
			switch it.Kind {
			case "Deployment":
				deploy, err := deploymentClient.Get(it.Name, metav1.GetOptions{})
				if err != nil {
					return err
				}
				o.fixPodSpec(&deploy.Spec.Template, it, o.dstRegistry.Domain())
				_, err = deploymentClient.Update(deploy)
				if err != nil {
					return err
				}
			case "DaemonSet":
				ds, err := dsClient.Get(it.Name, metav1.GetOptions{})
				if err != nil {
					return err
				}
				o.fixPodSpec(&ds.Spec.Template, it, o.dstRegistry.Domain())
				_, err = dsClient.Update(ds)
				if err != nil {
					return err
				}
			case "StatefulSet":
				sts, err := stsClient.Get(it.Name, metav1.GetOptions{})
				if err != nil {
					return err
				}
				o.fixPodSpec(&sts.Spec.Template, it, o.dstRegistry.Domain())
				_, err = stsClient.Update(sts)
				if err != nil {
					return err
				}
			}
		}
	}
	return err
}

// NewCmdTransferImage create fix image command
func NewCmdTransferImage(f cmdutils.Factory) *cobra.Command {
	o := newTransferImageOptions()
	cmd := &cobra.Command{
		Use:   "transferimage",
		Short: "Auto cp images to an accessable registry and modify deployment image",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutils.CheckError(o.Complete(f))
			cmdutils.CheckError(o.Run())
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&o.namespace, "namespace", "n", "default", "NS")
	return cmd
}
