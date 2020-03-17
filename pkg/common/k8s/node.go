package k8s

import (
	_ "context"
	"errors"

	//kube "k8s.io/client-go/kubernetes"
	_ "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
)

func ListNodes() ([]models.Node, error) {
	var result []models.Node
	opts := metav1.ListOptions{}
	nodes, err := corev1api.Nodes().List(opts)
	if err != nil {
		return result, err
	}
	for _, node := range nodes.Items {
		uid := string(node.ObjectMeta.UID)
		name := node.ObjectMeta.Name
		runFrom := node.ObjectMeta.CreationTimestamp.Time
		var newNode models.Node

		// Allocatable resources
		cpu, ok := node.Status.Allocatable.Cpu().AsInt64()
		if !ok {
			return result, errors.New("K8s parse cpu format error")
		}
		cpuCores := float64(cpu)
		mem, ok := node.Status.Allocatable.Memory().AsInt64()
		if !ok {
			return result, errors.New("K8s parse memory format error")
		}
		memo := float64(mem) / 1024 / 1024
		availRes := models.NewResource(cpuCores, memo)

		if spot := node.ObjectMeta.Labels["node-spot"]; spot != "false" {
			// TODO: get price
			newNode = models.NewDemandNode(uid, name, availRes, 1., runFrom)
		} else {
			newNode = models.NewConsistNode(uid, name, availRes, runFrom)
		}
		result = append(result, newNode)

	}
	return result, nil
}

// Get uid according to name
func getNodeID(name string) (string, error) {
	var result string
	opts := metav1.GetOptions{}
	node, err := corev1api.Nodes().Get(name, opts)
	if err != nil {
		return result, err
	}
	result = string(node.ObjectMeta.UID)
	return result, nil
}
