package scanner

import (
	"fmt"

	"gitlab.com/polarsquad/eks-auth-sync/internal/aws"
	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"
	"gitlab.com/polarsquad/eks-auth-sync/internal/scanner/file"
	"gitlab.com/polarsquad/eks-auth-sync/internal/scanner/iam"
	"gitlab.com/polarsquad/eks-auth-sync/internal/scanner/ssm"

	"github.com/spf13/afero"
)

type Scanner struct {
	Name   string          `yaml:"name"`
	Type   string          `yaml:"type"`
	Static mapping.All     `yaml:"static"`
	File   file.ScanConfig `yaml:"file"`
	IAM    iam.ScanConfig  `yaml:"iam"`
	SSM    ssm.ScanConfig  `yaml:"ssm"`
	AWS    aws.Config      `yaml:"aws"`
}

type API struct {
	FS  afero.Fs
	AWS *aws.API
}

func (s *Scanner) Validate() error {
	var err error
	switch s.Type {
	case "file":
		err = s.File.Validate()
	case "static":
		if s.Static.IsEmpty() {
			err = fmt.Errorf("no static mappings specified")
		}
	case "iam":
		err = s.IAM.Validate()
	case "ssm":
		err = s.SSM.Validate()
	case "":
		err = fmt.Errorf("no type specified")
	default:
		err = fmt.Errorf("unknown type %s", s.Type)
	}
	if err != nil {
		return fmt.Errorf("validation of scanner %s failed: %w", s.name(), err)
	}
	return nil
}

func (s *Scanner) Scan(api *API) (*mapping.All, error) {
	var ms *mapping.All
	var err error
	switch s.Type {
	case "file":
		ms, err = file.Scan(&s.File, api.FS)
	case "static":
		ms = &s.Static
	case "iam":
		ms, err = iam.Scan(&s.IAM, api.AWS)
	case "ssm":
		ms, err = ssm.Scan(&s.SSM, api.AWS)
	default:
		err = fmt.Errorf("invalid type [%s]", s.Type)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan %s: %w", s.name(), err)
	}
	return ms, nil
}

func (s *Scanner) UsesAWS() bool {
	switch s.Type {
	case "iam", "ssm":
		return true
	}
	return false
}

func (s *Scanner) name() string {
	if s.Name == "" {
		return s.Type
	}
	return s.Name
}
