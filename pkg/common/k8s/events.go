package k8s

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	//"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	//"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/api/core/v1"
)

// Watch the pods phase. When a new pod is created and
// stay `Pending`, call the cluster scale up.
/*
func WatchPods(ctx context.Context, srvc chan error) {
	opts := metav1.ListOptions{}
	watcher, err := corev1api.Pods("").Watch(opts)
	if err != nil {
		log.Error("Get pods watcher: ", err)
		srvc<-err
		return
	}
	resultChan := watcher.ResultChan()
	for event := range resultChan {
		pod, ok := event.Object.(*v1.Pod)
		if !ok {
			log.Error("Get event pod object error: unexpected type")
			srvc<-err
			return
		}
		switch event.Type {
		case watch.Added:
			podPhase := pod.Status.Phase
			switch podPhase {
			case v1.PodRunning:
				log.WithFields(log.Fields{
					"Pod name": pod.Name,
					"Namespace": pod.Namespace,
				}).Info("Pod schedule successfully!")
			case v1.PodPending:
				log.WithFields(log.Fields{
					"Pod name": pod.Name,
					"Namespace": pod.Namespace,
				}).Info("Pod pending because resource short")
				// TODO: compute proprely node type, call cloud provider
				// api to scale up cluster
			}
		case watch.Deleted:
			// TODO: Delete a pod after a duration, check if there is a more
			// efficient assignments
			log.WithFields(log.Fields{
				"Pod name": pod.Name,
				"Namespace": pod.Namespace,
			}).Info("Pod deleted waiting for a cluster steady state...")
		}
		// Return early if context canceled
		select {
		case <-ctx.Done():
			return
		}
	}
}*/

func WatchPods(ctx context.Context, srvc chan bool) {
	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod, ok := obj.(*v1.Pod)
			if !ok {
				log.Error("Get event pod object error: unexpected type")
				srvc<-false
				return
			}
			podPhase := pod.Status.Phase
			switch podPhase {
			case v1.PodRunning:
				log.WithFields(log.Fields{
					"Pod name": pod.Name,
					"Namespace": pod.Namespace,
				}).Info("Pod schedule successfully!")
			case v1.PodPending:
				log.WithFields(log.Fields{
					"Pod name": pod.Name,
					"Namespace": pod.Namespace,
				}).Info("Pod pending because resource short")
				// TODO: compute proprely node type, call cloud provider
				// api to scale up cluster
			}
		},
		DeleteFunc: func(obj interface{}) {
			// TODO: Delete a pod after a duration, check if there is a more
			// efficient assignments
			pod, ok := obj.(*v1.Pod)
			if !ok {
				log.Error("Get event pod object error: unexpected type")
				srvc<-false
				return
			}
			log.WithFields(log.Fields{
				"Pod name": pod.Name,
				"Namespace": pod.Namespace,
			}).Info("Pod deleted waiting for a cluster steady state...")
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			/*pod, ok := newObj.(*v1.Pod)
			if !ok {
				log.Error("Get event pod object error: unexpected type")
				srvc<-false
				return
			}
			log.WithFields(log.Fields{
				"Pod name": pod.Name,
				"Namespace": pod.Namespace,
			}).Info("Pod update waiting for a cluster steady state...")*/
		},
	}
	_, controller := cache.NewInformer(
		buildWatchList("pods"),
		&v1.Pod{},
		time.Millisecond*100,
		handlers,
	)
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}

// Utils
func buildWatchList(resource string) *cache.ListWatch {
	return cache.NewListWatchFromClient(
		corev1api.RESTClient(),
		resource,
		v1.NamespaceAll,
		fields.Everything())
}
