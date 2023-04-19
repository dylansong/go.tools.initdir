package main

import (
	"fmt"
	"github.com/dylansong/go.tools.initdir/appconfig"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	var configPath string

	var rootCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize project directory structure",
		Long:  `This command initializes the project directory structure based on a JSON or YAML configuration file.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Initializing project directory structure...")
			config, err := appconfig.ParseConfig(configPath)
			if err != nil {
				fmt.Printf("Error parsing config file: %v\n", err)
				os.Exit(1)
			}

			basePath := "."
			if len(args) > 0 {
				basePath = args[0]
			}
			err = appconfig.CreateDirectories(basePath, config.Directories)
			if err != nil {
				fmt.Printf("Error creating directories: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Project directory structure initialized successfully!")
		},
	}

	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to the JSON or YAML configuration file")
	rootCmd.MarkFlagRequired("config")

	rootCmd.Execute()
}
