package k8s

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKubeConfigPathResolution(t *testing.T) {
	homeDir := "/home/example"

	assert.Equal(
		t,
		filepath.Join(homeDir, ".kube/config"),
		resolveKubeConfigPath(homeDir, ""),
	)
	assert.Equal(
		t,
		filepath.Join(homeDir, "kubeconfig.yaml"),
		resolveKubeConfigPath(homeDir, "~/kubeconfig.yaml"),
	)
	assert.Equal(
		t,
		"kubeconfig.yaml",
		resolveKubeConfigPath(homeDir, "kubeconfig.yaml"),
	)
}
