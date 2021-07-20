package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(batCmd)
	rootCmd.AddCommand(cfgCmd)
	rootCmd.AddCommand(dirCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(fileCmd)
	rootCmd.AddCommand(goCmd)
	rootCmd.AddCommand(gqlCmd)
	rootCmd.AddCommand(grpcCmd)
}

func initConfig() {
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
