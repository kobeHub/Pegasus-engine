package k8s

import (
	"errors"

	"github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
	_ "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListRescheduleLabelPods(label string) ([]models.Pod, error) {
	var results []models.Pod
	nss, err := GetWorkNS()
	if err != nil {
		return results, err
	}
	opts := metav1.ListOptions{
		LabelSelector: label,
	}

	for _, ns := range nss {
		pi := corev1api.Pods(ns)
		podList, err := pi.List(opts)
		if err != nil {
			return results, err
		}
		for _, pod := range podList.Items {
			uid := string(pod.ObjectMeta.Name)
			nid := ""
			if nname := pod.Spec.NodeName; nname != "" {
				nid, err = getNodeID(nname)
				if err != nil {
					return results, err
				}
			}

			var (
				cpu = 0.
				mem = 0.
			)
			for _, container := range pod.Spec.Containers {
				ccpu, ok := container.Resources.Limits.Cpu().AsDec().Unscaled()
				if !ok {
					return results, errors.New("Parse container cpu cores error")
				}
				cpu += float64(ccpu) / 1000
				cmom, ok := container.Resources.Limits.Memory().AsDec().Unscaled()
				if !ok {
					return results, errors.New("Parse container memory error")
				}
				mem += float64(cmom) / 1024 / 1024
			}
			status := models.PodPhase(pod.Status.Phase)
			podModel := models.NewPod(uid, cpu, mem, nid, status)
			results = append(results, podModel)
		}
	}
	return results, nil
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
