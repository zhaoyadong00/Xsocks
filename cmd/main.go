package main

import (
	Xsocks "Xsocks-core"
	"Xsocks-core/util/config"
	"Xsocks-core/util/logs"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	configPath = ""

	proxyCommand = &cobra.Command{
		Use:   "start",
		Short: "Run proxy network is a background project",
		Run: func(cmd *cobra.Command, args []string) {
			err := config.LoadConfig(configPath)
			if err != nil {
				logs.Logger.Info("LoadConfig error" , zap.Error(err))
			}
			Xsocks.Start()
		},
	}
)

func main() {
	_ = proxyCommand.Execute()
}


func init() {
	proxyCommand.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config path")
}
