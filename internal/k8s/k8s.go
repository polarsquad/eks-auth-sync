package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
		kubeConfigPath := resolveKubeConfigPath(os.Getenv("HOME"), config.KubeConfigPath)
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(kubeConfig)
}

func resolveKubeConfigPath(homeDir string, kubeConfigPath string) string {
	if strings.TrimSpace(kubeConfigPath) == "" {
		return filepath.Join(homeDir, ".kube/config")
	}
	if strings.HasPrefix(kubeConfigPath, "~/") {
		return filepath.Join(homeDir, kubeConfigPath[2:])
	}
	return kubeConfigPath
}

func UpdateAWSAuthConfigMap(ctx context.Context, clientset kubernetes.Interface, configMap *k8sv1.ConfigMap) error {
	// Make sure we don't accidentally attempt to modify something that's not meant to be updated.
	if configMap.Name != configMapName {
		return fmt.Errorf("attempted to update a ConfigMap that's not %s: %s", configMapName, configMap.Name)
	}
	if configMap.Namespace != "" && configMap.Namespace != namespace {
		return fmt.Errorf("attempted to update a ConfigMap that's not in %s namespace: %s", namespace, configMap.Name)
	}

	cmAPI := clientset.CoreV1().ConfigMaps(namespace)
	_, err := cmAPI.Get(ctx, configMapName, metav1.GetOptions{})
	if err == nil {
		_, err := cmAPI.Update(ctx, configMap, metav1.UpdateOptions{})
		return err
	} else if errors.IsNotFound(err) {
		_, err := cmAPI.Create(ctx, configMap, metav1.CreateOptions{})
		return err
	}
	return err
}
