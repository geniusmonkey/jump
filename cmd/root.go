package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "jump",
	Short: "jump is a tool to store and connect to ssh connections",
	Long: `save your ssh connections and port forwarding configuration
with names and connect using names`,
}

func Execute() {
	initConfig()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
}
