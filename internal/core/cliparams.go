package core

import (
	"flag"
	"fmt"
	"io"
	"log"
)

const (
	appName            = "eks-auth-sync"
	exampleInvocations = `
Examples:
  # Read config from a file and print out the scanned results.
  eks-auth-sync -config configfile.yaml

  # Read config from STDIN and print out the scanned results.
  eks-auth-sync -config -

  # Read config from a file and commit the scanned results to Kubernetes
  eks-auth-sync -config configfile.yaml -commit`
)

type cliParams struct {
	fls        *flag.FlagSet
	configFile string
	commit     bool
	version    bool
}

func (c *cliParams) setup(out io.Writer) {
	fls := flag.NewFlagSet(appName, flag.ContinueOnError)
	fls.StringVar(&c.configFile, "config", "", "Path to the YAML config file. Set to '-' to read from STDIN.")
	fls.BoolVar(&c.commit, "commit", false, "If used, the scanned results are committed to Kubernetes.")
	fls.BoolVar(&c.version, "version", false, "Print out the version information")
	fls.SetOutput(out)
	fls.Usage = func() {
		if _, err := fmt.Fprintf(fls.Output(), "Usage of %s:\n", fls.Name()); err != nil {
			log.Printf("failed to print usage: %s", err)
		}
		fls.PrintDefaults()
		if _, err := fmt.Fprintln(fls.Output(), exampleInvocations); err != nil {
			log.Printf("failed to print usage: %s", err)
		}
	}
	c.fls = fls
}

func (c *cliParams) parse(args []string) error {
	return c.fls.Parse(args)
}

func (c *cliParams) printUsage() {
	c.fls.Usage()
}

func (c *cliParams) validate() error {
	if c.configFile == "" {
		return fmt.Errorf("no config file specified")
	}
	return nil
}
