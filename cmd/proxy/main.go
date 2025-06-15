package main

import (
	"github.com/flew1x/grpc-chaos-proxy/internal/adapter/cli"
	_ "github.com/flew1x/grpc-chaos-proxy/internal/core/injector/initiator"
	"os"
)

func main() {
	rootCmd := cli.NewCLI()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
