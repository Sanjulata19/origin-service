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
	s3Bucket     string
	s3Prefix     string
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
		&s3Bucket,
		"s3-bucket",
		"",
		"",
		"The s3 bucket to use as a source",
	)
	rootCmd.Flags().StringVarP(
		&s3Prefix,
		"s3-prefix-folder",
		"",
		"site-deployments",
		"The s3 folder to find sites in",
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
