package file

import (
	"fmt"

	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"

	"github.com/spf13/afero"
)

type ScanConfig struct {
	Path string `yaml:"path"`
}

func (s *ScanConfig) Validate() error {
	if s.Path == "" {
		return fmt.Errorf("no path specified")
	}
	return nil
}

func Scan(c *ScanConfig, fs afero.Fs) (ms *mapping.All, err error) {
	var bs []byte
	ms = &mapping.All{}

	bs, err = afero.ReadFile(fs, c.Path)
	if err != nil {
		return
	}

	err = ms.FromYAML(bs)
	return
}
