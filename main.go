package main

import (
	"fmt"
	"os"

	"github.com/raymondsugiarto/reputation-be/cmd/db"
	cmdserver "github.com/raymondsugiarto/reputation-be/cmd/server"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use: "app",
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringP("schema", "s", "default", "choose schema to run on")
	RootCmd.PersistentFlags().StringP("name", "n", "default", "choose seed to run")

	RootCmd.AddCommand(cmdserver.RestCmd)
	RootCmd.AddCommand(cmdserver.StartRestCmd)
	RootCmd.AddCommand(db.DBCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	fmt.Println("init config")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println("init config")
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".config-api" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".config-api")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
