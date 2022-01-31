package core

import (
	"bufio"
	"io"
	"log"

	intaws "gitlab.com/polarsquad/eks-auth-sync/internal/aws"
	"gitlab.com/polarsquad/eks-auth-sync/internal/k8s"
	"gitlab.com/polarsquad/eks-auth-sync/internal/scanner"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Kubernetes k8s.Config         `yaml:"kubernetes"`
	AWS        intaws.Config      `yaml:"aws"`
	Scanners   []*scanner.Scanner `yaml:"scanners"`
}

func (c *Config) FromYAML(yamlInput io.Reader) error {
	return yaml.NewDecoder(bufio.NewReader(yamlInput)).Decode(c)
}

func (c *Config) FromYAMLFile(fs afero.Fs, filename string) error {
	file, err := fs.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close file %s: %s", filename, err)
		}
	}()

	return c.FromYAML(file)
}

func (c *Config) UsesAWS() bool {
	for _, s := range c.Scanners {
		if s.UsesAWS() {
			return true
		}
	}
	return false
}
