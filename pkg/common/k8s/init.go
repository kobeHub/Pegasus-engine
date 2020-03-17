package k8s

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	kube "k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var client *kube.Clientset
var corev1api corev1.CoreV1Interface

func Init(env string) {
	var config *rest.Config
	var err error
	if env == "local" {
		kubeConfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Fatal("Read kubeconfig from default file")
		}
	} else if env == "prod" || env == "test" {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Fatal("Get in cluster config error")
		}
	}

	client, err = kube.NewForConfig(config)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Fatal("Get clientset error")
	}
	corev1api = client.CoreV1()
	log.Info("k8s client initialized successfully")
}
