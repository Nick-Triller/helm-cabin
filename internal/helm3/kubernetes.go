package helm3

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var clientset *kubernetes.Clientset

func connectKubernetes() {
	// Try in-cluster first
	config, err := rest.InClusterConfig()
	if err != nil {
		// Use the current context in kubeconfig
		kubeconfigPath := filepath.Join(homeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			log.Fatalf("Failed to create Kubernetes connection config")
		}
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
