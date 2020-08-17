package main

import (
	"log"
	"os"

	"gitlab.com/polarsquad/eks-auth-sync/internal/core"
)

func main() {
	if err := mainWithErr(); err != nil {
		log.Fatal(err)
	}
}

func mainWithErr() error {
	return core.NewCore().Run(os.Args[1:])
}
