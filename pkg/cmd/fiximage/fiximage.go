package fiximage

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iftechio/jki/pkg/cmd/cp"
	cmdutils "github.com/iftechio/jki/pkg/cmd/utils"
	"github.com/iftechio/jki/pkg/image"
	"github.com/iftechio/jki/pkg/registry"
	"github.com/spf13/cobra"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
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
type fixImageOptions struct {
	namespace   string
	cp          *cp.CopyOptions
	dstRegistry *registry.Registry
}

func (o *fixImageOptions) Complete(f cmdutils.Factory) error {
	o.cp = cp.NewCopyOptions()
	dstReg, registries, err := f.LoadRegistries()
	if err != nil {
		return err
	}
	o.dstRegistry = registries[dstReg]
	return o.cp.Complete(f, nil, nil)
}

func newFixImageOptions() *fixImageOptions {
	return &fixImageOptions{}
}

func (o *fixImageOptions) fixPodSpec(podSpec *apiv1.PodTemplateSpec, it brokenObject, domain string) {
	for i, con := range podSpec.Spec.Containers {
		if con.Image == it.Image {
			// copy to accessable registry
			img := image.FromString(it.Image)
			o.cp.Run([]string{it.Image})

			// replace with new domain
			img.Domain = domain
			podSpec.Spec.Containers[i].Image = img.String()
			fmt.Printf("Replace %s to %s\n", it.Image, img.String())
		}
	}
}

func (o *fixImageOptions) Run(namespace *string) (err error) {
	fmt.Printf("Searching for deploy/ds to fix in namespace: %s\n", *namespace)
	var itemsToFix []brokenObject
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	podsList, err := clientset.CoreV1().Pods(*namespace).List(metav1.ListOptions{
		FieldSelector: "status.phase=Pending",
	})

	pods := podsList.Items

	// filter pods
	n := 0
	for _, pod := range pods {
		if hasErrPullingContainer(pod) {
			pods[n] = pod
			n++
		}
	}
	pods = pods[:n]

	// find k8s objects need to fix image for
	for _, pod := range pods {
		owner := pod.OwnerReferences[0]
		if owner.Kind == "ReplicaSet" {
			rs, err := clientset.AppsV1().ReplicaSets(*namespace).Get(owner.Name, metav1.GetOptions{})
			if err != nil {
				panic(err)
			}
			owner = rs.OwnerReferences[0]
		}
		for _, con := range pod.Status.ContainerStatuses {
			if con.State.Waiting.Reason == "ImagePullBackOff" || con.State.Waiting.Reason == "ErrImagePull" {
				itemsToFix = append(itemsToFix, brokenObject{
					Kind:  owner.Kind,
					Name:  owner.Name,
					Image: con.Image,
				})
			}
		}
	}

	// remove duplicated
	itemsMap := make(map[string]bool)
	n = 0
	for _, it := range itemsToFix {
		if _, ok := itemsMap[it.Kind+"/"+it.Name]; !ok {
			itemsToFix[n] = it
			n++
		}
	}
	itemsToFix = itemsToFix[:n]

	deploymentClient := clientset.AppsV1().Deployments(*namespace)
	dsClient := clientset.AppsV1().DaemonSets(*namespace)

	if len(itemsToFix) == 0 {
		fmt.Println("Found no image to fix")
		return nil
	}
	for _, it := range itemsToFix {
		fmt.Printf("Fix %s/%s %s(y/n)?\n", it.Kind, it.Name, it.Image)
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		sentence, err := buf.ReadBytes('\n')
		if err != nil {
			return err
		}
		if strings.ToLower(strings.TrimSpace(string(sentence))) == "y" {
			fmt.Printf("Fixing %s/%s \n", it.Kind, it.Name)
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
			}
		}
	}
	return err
}

// NewCmdFixImage create fix image command
func NewCmdFixImage(f cmdutils.Factory) *cobra.Command {
	var namespace string
	o := newFixImageOptions()
	cmd := &cobra.Command{
		Use:   "fiximage",
		Short: "Auto cp images to an accessable registry and modify deployment image",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutils.CheckError(o.Complete(f))
			cmdutils.CheckError(o.Run(&namespace))
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&namespace, "namespace", "n", "default", "NS")
	return cmd
}
