package core

import (
	"context"
	"fmt"
	"io"
	"os"

	"gitlab.com/polarsquad/eks-auth-sync/internal/buildinfo"
	"gopkg.in/yaml.v2"

	intaws "gitlab.com/polarsquad/eks-auth-sync/internal/aws"
	"gitlab.com/polarsquad/eks-auth-sync/internal/k8s"
	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"
	"gitlab.com/polarsquad/eks-auth-sync/internal/scanner"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/afero"
	"k8s.io/client-go/kubernetes"
)

type Core struct {
	AppFS      afero.Fs
	Input      io.Reader
	Output     io.Writer
	AWSSession func(*intaws.Config) (*session.Session, error)
	AWSAPI     func(*session.Session, *intaws.Config) *intaws.API
	KubeClient func(*k8s.Config) (kubernetes.Interface, error)
	Context    context.Context
}

func NewCore() *Core {
	return &Core{
		AppFS:      afero.NewReadOnlyFs(afero.NewOsFs()),
		Input:      os.Stdin,
		Output:     os.Stdout,
		AWSSession: intaws.NewSession,
		AWSAPI:     intaws.NewAPI,
		KubeClient: k8s.NewClientset,
		Context:    context.Background(),
	}
}

func (c *Core) Run(args []string) error {
	var params cliParams
	params.setup(c.Output)
	if err := params.parse(args); err != nil {
		return err
	}
	if params.version {
		buildinfo.PrintVersion(c.Output)
		return nil
	}
	if err := params.validate(); err != nil {
		params.printUsage()
		return err
	}

	var config Config
	if err := c.readConfig(&config, params.configFile); err != nil {
		return err
	}
	return c.runFromConfig(&config, params.commit)
}

func (c *Core) readConfig(config *Config, filename string) error {
	if filename == "-" {
		return config.FromYAML(c.Input)
	}
	return config.FromYAMLFile(c.AppFS, filename)
}

func (c *Core) runFromConfig(config *Config, commit bool) error {
	var err error

	if len(config.Scanners) == 0 {
		return fmt.Errorf("no scanners defined")
	}

	// Create an AWS session if needed
	var sess *session.Session
	if config.UsesAWS() {
		sess, err = c.AWSSession(&config.AWS)
		if err != nil {
			return err
		}
	}

	// Scan and collect mappings. Create an AWS API if needed.
	var mappings mapping.All
	for _, sc := range config.Scanners {
		var api scanner.API
		api.FS = c.AppFS
		if sc.UsesAWS() {
			awsConfig := intaws.MergeConfigs(&config.AWS, &sc.AWS)
			api.AWS = c.AWSAPI(sess, awsConfig)
		}
		ms, err := sc.Scan(&api)
		if err != nil {
			return err
		}
		mappings.Append(ms)
	}

	// Save the changes Kubernetes if commit is enabled.
	// Otherwise, just print out the results.
	if commit {
		mappingsConfigMap, err := mappings.ToConfigMap()
		if err != nil {
			return err
		}
		clientset, err := c.KubeClient(&config.Kubernetes)
		if err != nil {
			return err
		}
		return k8s.UpdateAWSAuthConfigMap(c.Context, clientset, mappingsConfigMap)
	} else {
		bs, err := yaml.Marshal(mappings)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(c.Output, string(bs)); err != nil {
			return err
		}
	}

	return nil
}
