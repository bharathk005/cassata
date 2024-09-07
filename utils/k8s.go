package utils

import (
	"os"
	"path/filepath"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetDynamicClient() (dynamic.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		// If not running in cluster, use kubeconfig from C:\Users\<user>\.kube
		homeDir, _ := os.UserHomeDir()
		kubeconfig := filepath.Join(homeDir, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	return dynamic.NewForConfig(config)
}

func GetKubernetesClientset() (*kubernetes.Clientset, error) {
	// if running in cluster, use in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		// if not running in cluster, use kubeconfig
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = os.Getenv("HOME") + "/.kube/config"
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func GetApiextensionsClientset() (*apiextensionsclientset.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		// If not running in cluster, use kubeconfig from C:\Users\<user>\.kube
		homeDir, _ := os.UserHomeDir()
		kubeconfig := filepath.Join(homeDir, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	return apiextensionsclientset.NewForConfig(config)
}
