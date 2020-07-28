package file

import (
	"io/ioutil"

	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"
)

type ScanConfig struct {
	FilePath string
}

func Scan(c *ScanConfig) (ms *mapping.All, err error) {
	var bs []byte
	ms = &mapping.All{}

	bs, err = ioutil.ReadFile(c.FilePath)
	if err != nil {
		return
	}

	err = ms.FromYAML(bs)
	return
}
