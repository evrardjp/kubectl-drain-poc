package main

import (
	"context"
	"fmt"
	"log"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/drain"
)

func main() {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", "/home/jean-philippe/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	drainer := &drain.Helper{
		Client:              clientset,
		Force:               true,
		DeleteLocalData:     true,
		IgnoreAllDaemonSets: true,
		ErrOut:              os.Stderr,
		Out:                 os.Stdout,
	}

	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), "kind-control-plane", metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	if err := drain.RunCordonOrUncordon(drainer, node, true); err != nil {
		log.Fatal("Error cordonning")

	}

	log.Print("Done")
	if err := drain.RunNodeDrain(drainer, "kind-control-plane"); err != nil {
		log.Fatal("Error draining")
	}

}
