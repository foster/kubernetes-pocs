package main

import (
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/cache"
	"k8s.io/apimachinery/pkg/util/wait"
	"net"
	"os"
	// "time"
)

const CONTROL_ADDR = "localhost:2999"

type foobar struct {
	name string
	age  int
}

func main() {
	clientset, err := getKubeClient()
	if err != nil {
		panic(err.Error())
	}

	selector := fields.Everything()
	listWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", apiv1.NamespaceDefault, selector)
	store := cache.NewFIFO( cache.MetaNamespaceKeyFunc )
	r := cache.NewReflector( listWatcher, &apiv1.Pod{}, store, 0 )
	go r.ListAndWatch(wait.NeverStop)

	for {
		store.Pop(func (obj interface{}) error {
			pod := obj.(*apiv1.Pod)
			fmt.Printf("Pop. %s: %s\n", pod.Name, pod.Status.PodIP)
			go dialControl("localhost:3000")

			return nil
		})
	}
}

func getKubeClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if kubeconfig := os.Getenv("KUBECONFIG"); len(kubeconfig) > 0 {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// establish connection to the app's CONTROL socket
// and tell the app to establish a connection to remoteAddr
func dialControl(remoteAddress string) {
	conn, err := net.Dial("tcp", CONTROL_ADDR)
	if err != nil {
		panic(fmt.Sprintf("Error dialing connection to control %s: %v", CONTROL_ADDR, err.Error()))
	}
	conn.Write([]byte(remoteAddress))
	conn.Close()
}
