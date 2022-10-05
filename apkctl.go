/*
 * Copyright (c) 2022, WSO2 LLC. (https://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

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
	RootCmd.AddCommand(cmd.VersionCmd)
	execute()
}
