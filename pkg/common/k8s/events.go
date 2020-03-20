package k8s

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	//"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/api/core/v1"
)

// Watch the pods phase. When a new pod is created and
// stay `Pending`, call the cluster scale up.
func WatchPods(ctx context.Context, srvc chan bool) {
	handler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod, ok := obj.(*v1.Pod)
			if !ok {
				log.Error("Get event pod object error: unexpected type")
				srvc <- false
				return
			}
			length := len(pod.Status.Conditions)
			if length < 1 {
				log.WithFields(log.Fields{
					"pod_name":  pod.Name,
					"namespace": pod.Namespace,
				}).Warn("pod status not update")
				return
			}
			var (
				podPhase      v1.PodPhase     = pod.Status.Phase
				lastCondition v1.PodCondition = pod.Status.Conditions[length-1]
			)
			switch podPhase {
			case v1.PodRunning:
				log.WithFields(log.Fields{
					"Pod name":  pod.Name,
					"Namespace": pod.Namespace,
				}).Info("Pod schedule successfully!")
			case v1.PodPending:
				if lastCondition.Type == v1.PodScheduled && lastCondition.Status == v1.ConditionFalse {
					log.WithFields(log.Fields{
						"Pod name":  pod.Name,
						"Namespace": pod.Namespace,
						"Message":   lastCondition.Message,
					}).Info("Pod pending because resource short")
					// TODO: compute proprely node type, call cloud provider
					// api to scale up cluster
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			// TODO: Delete a pod after a duration, check if there is a more
			// efficient assignments
			pod, ok := obj.(*v1.Pod)
			if !ok {
				log.Error("Get event pod object error: unexpected type")
				srvc <- false
				return
			}
			log.WithFields(log.Fields{
				"Pod name":  pod.Name,
				"Namespace": pod.Namespace,
			}).Info("Pod deleted waiting for a cluster steady state...")
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			pod, ok := newObj.(*v1.Pod)
			if !ok {
				log.Error("Get event pod object error: unexpected type")
				srvc <- false
				return
			}
			length := len(pod.Status.Conditions)
			if length < 1 {
				return
			}
			podCondition := pod.Status.Conditions[length-1]
			if podCondition.Type == v1.PodScheduled && podCondition.Status == v1.ConditionFalse {
				log.WithFields(log.Fields{
					"pod_name":  pod.Name,
					"namespace": pod.Namespace,
					"pod_phase": pod.Status.Phase,
					"message":   podCondition.Message,
				}).Info("Try schedule pod failed...")
				// TODO: compute proprely node type, call cloud provider
				// api to scale up cluster
			}
		},
	}

	podsInformer := buildInformerFactory("pods", handler)
	stop := make(chan struct{})
	defer close(stop)
	podsInformer.Start(stop)
	for {
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			return
		}
	}
}

// Utils
func buildInformerFactory(resource string,
	eventsHandler cache.ResourceEventHandlerFuncs) informers.SharedInformerFactory {
	kubeInformerFac := informers.NewSharedInformerFactory(client, time.Second*30)
	switch resource {
	case "pods":
		log.Info("Build pod informer")
		podsInformer := kubeInformerFac.Core().V1().Pods().Informer()
		podsInformer.AddEventHandler(eventsHandler)
	}
	return kubeInformerFac
}
