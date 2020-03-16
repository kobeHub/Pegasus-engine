package k8s

import (
	_ "context"

	corev1 "k8s.io/api/core/v1"
	//kube "k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//log "github.com/sirupsen/logrus"
)

func ListNodes() (string, error) {
	var nodes *corev1.NodeList
	var err error
	var result string
	api := client.CoreV1()
	opts := metav1.ListOptions{}
	if nodes, err = api.Nodes().List(opts); err != nil {
		return result, err
	}
	result = nodes.String()
	return result, err
}
