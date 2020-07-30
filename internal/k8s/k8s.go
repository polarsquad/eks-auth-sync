package k8s

import (
	"fmt"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	configMapName = "aws-auth"
	namespace     = "kube-system"
)

type Config struct {
	InKubeCluster  bool   `yaml:"inKubeCluster"`
	KubeConfigPath string `yaml:"kubeConfigPath"`
}

func NewClientset(config *Config) (kubernetes.Interface, error) {
	var err error
	var kubeConfig *rest.Config
	if config.InKubeCluster {
		kubeConfig, err = rest.InClusterConfig()
	} else {
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", config.KubeConfigPath)
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(kubeConfig)
}

func UpdateAWSAuthConfigMap(clientset kubernetes.Interface, configMap *k8sv1.ConfigMap) error {
	// Make sure we don't accidentally attempt to modify something that's not meant to be updated.
	if configMap.Name != configMapName {
		return fmt.Errorf("attempted to update a ConfigMap that's not %s: %s", configMapName, configMap.Name)
	}
	if configMap.Namespace != "" && configMap.Namespace != namespace {
		return fmt.Errorf("attempted to update a ConfigMap that's not in %s namespace: %s", namespace, configMap.Name)
	}

	cmAPI := clientset.CoreV1().ConfigMaps(namespace)
	_, err := cmAPI.Get(configMapName, metav1.GetOptions{})
	if err == nil {
		_, err := cmAPI.Update(configMap)
		return err
	} else if errors.IsNotFound(err) {
		_, err := cmAPI.Create(configMap)
		return err
	}
	return err
}
