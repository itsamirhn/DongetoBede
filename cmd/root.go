package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/itsamirhn/dongetobede/internal/config"
)

var rootCmd = &cobra.Command{
	Use: "dong",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadConfigOrPanic(cmd *cobra.Command) {
	err := config.LoadConfig(cmd)
	if err != nil {
		logrus.WithError(err).Panic("failed to load configurations")
	}
}
