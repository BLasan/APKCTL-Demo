package main

import (
	"os"

	cmd "github.com/BLasan/APKCTL-Demo/cmd"
	"github.com/spf13/cobra"
)

var cfgFile string

// RootCmd ...
var RootCmd = &cobra.Command{
	Use:   "apkctl",
	Short: "apkctl",
	Long:  "apkctl",
}

func execute() {
	if err := RootCmd.Execute(); err != nil {
		// klog.Errorf("help: %v", err)
		os.Exit(1)
	}
}

func main() {
	RootCmd.AddCommand(cmd.InstallPlatformCmd)
	RootCmd.AddCommand(cmd.CreateCmd)
	RootCmd.AddCommand(cmd.DeleteCmd)
	RootCmd.AddCommand(cmd.GetCmd)
	RootCmd.AddCommand(cmd.UninstallPlatformCmd)
	execute()
}
