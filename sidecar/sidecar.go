package main

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net"
	"time"
)

const CONTROL_ADDR = "localhost:2999"

func getKubeClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
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
		fmt.Println("Error dialing connection to control", CONTROL_ADDR, ":", err.Error())
		return
	}
	conn.Write([]byte(remoteAddress))
	conn.Close()
}

func main() {
	clientset, err := getKubeClient()
	if err != nil {
		panic(err.Error())
	}

	for {
		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		time.Sleep(10 * time.Second)
	}
	go dialControl("localhost:3000")
}
