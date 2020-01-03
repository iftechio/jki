package deploy

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Credit to https://github.com/kubernetes/kubectl/blob/e2c59440f3e2c1a58e7012cedcb9b5a7460e279f/pkg/polymorphichelpers/updatepodspec.go#L33
func updatePodSpecForObject(obj runtime.Object, fn func(*corev1.PodSpec) error) (bool, error) {
	switch t := obj.(type) {
	case *corev1.Pod:
		return true, fn(&t.Spec)

	// Deployment
	case *extensionsv1beta1.Deployment:
		return true, fn(&t.Spec.Template.Spec)
	case *appsv1beta1.Deployment:
		return true, fn(&t.Spec.Template.Spec)
	case *appsv1beta2.Deployment:
		return true, fn(&t.Spec.Template.Spec)
	case *appsv1.Deployment:
		return true, fn(&t.Spec.Template.Spec)

	// DaemonSet
	case *extensionsv1beta1.DaemonSet:
		return true, fn(&t.Spec.Template.Spec)
	case *appsv1beta2.DaemonSet:
		return true, fn(&t.Spec.Template.Spec)
	case *appsv1.DaemonSet:
		return true, fn(&t.Spec.Template.Spec)

	// StatefulSet
	case *appsv1beta1.StatefulSet:
		return true, fn(&t.Spec.Template.Spec)
	case *appsv1beta2.StatefulSet:
		return true, fn(&t.Spec.Template.Spec)
	case *appsv1.StatefulSet:
		return true, fn(&t.Spec.Template.Spec)

	// CronJob
	case *batchv1beta1.CronJob:
		return true, fn(&t.Spec.JobTemplate.Spec.Template.Spec)
	case *batchv2alpha1.CronJob:
		return true, fn(&t.Spec.JobTemplate.Spec.Template.Spec)

	default:
		return false, fmt.Errorf("the object is not a pod or does not have a pod template: %T", t)
	}
}
