package k8s

import (
	_ "errors"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
)

func ListReschedulablePods() (models.Pod, error) {
	var results []models.Pod
	nss, err := GetWorkNS()
	if err != nil {
		return results, err
	}
	opts := metav1.ListOption{
		LabelSelector: "reschedulable=true",
	}
	for _, ns := range nss {
		pi := corev1api.Pods(ns)
		podList, err := pi.List(opts)
		if err != nil {
			return results, err
		}
		for _, pod := range podList.Items {

		}
	}
}

func GetWorkNS() ([]string, error) {
	var results []string
	opts := metav1.ListOptions{LabelSelector: "dispense=pegasus"}
	nsList, err := corev1api.Namespaces().List(opts)
	if err != nil {
		return results, err
	}
	for _, ns := range nsList.Items {
		if ns.Status.Phase == v1.NamespaceActive {
			results = append(results, ns.Name)
		}
	}
	return results, nil
}
