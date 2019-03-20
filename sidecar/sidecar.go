package main

import (
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"net"
	"os"
)

const CONTROL_ADDR = "localhost:2999"

var MY_IP string = os.Getenv("MY_POD_IP")

type foobar struct {
	name string
	age  int
}

func main() {
	clientset, err := getKubeClientset()
	if err != nil {
		panic(err.Error())
	}

	client := clientset.CoreV1().RESTClient()
	optionsModifier := func(options *metav1.ListOptions) {
		options.FieldSelector = "status.phase=Running"
		options.LabelSelector = "app=app"
	}
	listWatcher := cache.NewFilteredListWatchFromClient(client, "pods", apiv1.NamespaceDefault, optionsModifier)
	store := cache.NewFIFO(cache.MetaNamespaceKeyFunc)
	r := cache.NewReflector(listWatcher, &apiv1.Pod{}, store, 0)
	go r.ListAndWatch(wait.NeverStop)

	for {
		store.Pop(func(obj interface{}) error {
			pod := obj.(*apiv1.Pod)

			// do not attempt to connect our own Pod.
			// ignore Pods unless the condition Ready = true
			// unready might mean the pod is being shut down
			if ip := pod.Status.PodIP; ip != MY_IP && isPodReady(pod) {
				fmt.Printf("Pop. %s. IP: %s\n", pod.Name, ip)
				fmt.Printf("Pop. %v", pod.Status)
				go dialControl(ip + ":3000")
			}

			return nil
		})
	}
}

func getKubeClientset() (*kubernetes.Clientset, error) {
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

func isPodReady(pod *apiv1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == apiv1.PodReady && condition.Status == apiv1.ConditionTrue {
			return true
		}
	}
	return false
}
