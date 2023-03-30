
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// Controller is the main controller for the Kubernetes controller
type Controller struct {
	client    kubernetes.Interface
	dynamic   dynamic.Interface
	prometheus *prometheus.GaugeVec
}

// NewController creates a new controller
func NewController(kubeconfig string) (*Controller, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dynamic, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	prometheus := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "microservice_logs",
			Help: "Logs for microservices",
		},
		[]string{"name", "namespace"},
	)

	return &Controller{
		client:    client,
		dynamic:   dynamic,
		prometheus: prometheus,
	}, nil
}

// Run starts the controller
func (c *Controller) Run(ctx context.Context) error {
	// Register the Prometheus metric
	prometheus.MustRegister(c.prometheus)

	// Create a new informer for the microservice resource
	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return c.dynamic.Resource(schema.GroupVersionResource{
					Group:    "example.com",
					Version:  "v1",
					Resource: "microservices",
				}).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return c.dynamic.Resource(schema.GroupVersionResource{
					Group:    "example.com",
					Version:  "v1",
					Resource: "microservices",
				}).Watch(options)
			},
		},
		&unstructured.Unstructured{},
		0,
		cache.Indexers{},
	)

	// Set up the event handlers
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			c.handleAdd(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			c.handleUpdate(oldObj, newObj)
		},
		DeleteFunc: func(obj interface{}) {
			c.handleDelete(obj)
		},
	})

	// Start the informer
	go informer.Run(ctx.Done())

	// Wait for the informer to sync
	if !cache.WaitForCacheSync(ctx.Done(), informer.HasSynced) {
		return fmt.Errorf("timed out waiting for cache to sync")
	}

	return nil
}

// handleAdd is called when a new microservice is added
func (c *Controller) handleAdd(obj interface{}) {
	ms := obj.(*unstructured.Unstructured)
	name := ms.GetName()
	namespace := ms.GetNamespace()

	// Log the microservice
	c.prometheus.WithLabelValues(name, namespace).Set(time.Now().Unix())
}

// handleUpdate is called when a microservice is updated
func (c *Controller) handleUpdate(oldObj, newObj interface{}) {
	ms := newObj.(*unstructured.Unstructured)
	name := ms.GetName()
	namespace := ms.GetNamespace()

	// Log the microservice
	c.prometheus.WithLabelValues(name, namespace).Set(time.Now().Unix())
}

// handleDelete is called when a microservice is deleted
func (c *Controller) handleDelete(obj interface{}) {
	ms, ok := obj.(*unstructured.Unstructured)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			return
		}
		ms, ok = tombstone.Obj.(*unstructured.Unstructured)
		if !ok {
			return
		}
	}

	name := ms.GetName()
	namespace := ms.GetNamespace()

	// Log the microservice
	c.prometheus.WithLabelValues(name, namespace).Set(0)
}

func main() {
	// Create a new controller
	controller, err := NewController("/path/to/kubeconfig")
	if err != nil {
		panic(err)
	}

	// Run the controller
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := controller.Run(ctx); err != nil {
		panic(err)
	}
}
      