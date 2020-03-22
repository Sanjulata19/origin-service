package cmd

import (
	"github.com/nullserve/static-host/static_host"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Use:   "static-host",
		Short: "Static Host is a cloud storage web server and static file router",
		Run:   root,
	}
	domainSuffix string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVarP(
		&domainSuffix,
		"domain-suffix",
		"",
		"sites.nullserve.dev",
		"The domain suffix to be used for routing to sites. For example 123.sites.nullserve.dev has the suffix \"sites.nullserve.dev\"",
	)
	rootCmd.Flags().StringVarP(
		&domainSuffix,
		"s3-bucket",
		"",
		"sites.nullserve.dev",
		"The domain suffix to be used for routing to sites. For example 123.sites.nullserve.dev has the suffix \"sites.nullserve.dev\"",
	)
	_ = viper.BindPFlags(rootCmd.Flags())
}

func initConfig() {
	viper.SetEnvPrefix("static_host")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func root(cmd *cobra.Command, args []string) {
	static_host.Main()
}

func Execute() error {
	return rootCmd.Execute()
}
