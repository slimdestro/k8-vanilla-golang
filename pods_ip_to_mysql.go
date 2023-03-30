/**
	@ gets pod IP when they restart and stores in mysql table\
	@ slimdestro
*/
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// Create the config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a watcher
	watcher, err := clientset.CoreV1().Pods("default").Watch(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Create a channel to receive events
	ch := watcher.ResultChan()

	// Create a ticker to check for new events
	ticker := time.NewTicker(time.Second * 5)

	// Loop forever
	for {
		select {
		case event := <-ch:
			// Check if the event is a restart
			if event.Type == "MODIFIED" && event.Object.(*metav1.Pod).Status.Phase == "Running" {
				// Get the pod
				pod, err := clientset.CoreV1().Pods("default").Get(event.Object.(*metav1.Pod).Name, metav1.GetOptions{})
				if err != nil {
					if errors.IsNotFound(err) {
						log.Printf("Pod %s not found\n", event.Object.(*metav1.Pod).Name)
					} else {
						log.Printf("Error getting pod %s: %s\n", event.Object.(*metav1.Pod).Name, err)
					}
					continue
				}

				// Get the IP address
				ip := pod.Status.PodIP

				// Store the IP address in the database
				db, err := mysql.NewMySQLDriver("username:password@tcp(host:port)/dbname")
				if err != nil {
					log.Fatal(err)
				}
				_, err = db.Exec("INSERT INTO ip_addresses (pod_name, ip_address) VALUES (?, ?)", pod.Name, ip)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Stored IP address %s for pod %s\n", ip, pod.Name)
			}
		case <-ticker.C:
			// Check for new events
		}
	}
}