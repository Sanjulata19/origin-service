package cmd

import (
	"github.com/nullserve/origin-service/config"
	"github.com/nullserve/origin-service/service"
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
	domainSuffix      string
	appPrefix         string
	refPrefix         string
	dynamoDBTableName string
	s3Bucket          string
	s3Prefix          string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVarP(
		&domainSuffix,
		"domain-suffix",
		"",
		"nullserve.dev",
		"The domain suffix to be used for routing to sites. For example 123.sites.nullserve.dev has the suffix \"sites.nullserve.dev\"",
	)
	rootCmd.Flags().StringVarP(
		&appPrefix,
		"app-prefix",
		"",
		"app",
		"The subdomain prefix used for named app aliases.",
	)
	rootCmd.Flags().StringVarP(
		&refPrefix,
		"ref-prefix",
		"",
		"ref",
		"The subdomain prefix used for ref id lookups.",
	)
	rootCmd.Flags().StringVarP(
		&s3Bucket,
		"s3-bucket",
		"",
		// FIXME: un-hardcode this default
		"nullserve-api-site-deployments20191125172523931100000001",
		// Required until multiple sources are added
		"The s3 bucket to use as a source (required)",
	)
	rootCmd.Flags().StringVarP(
		&s3Prefix,
		"s3-prefix-folder",
		"",
		"site-deployments",
		"The s3 folder to find sites in",
	)
	rootCmd.Flags().StringVarP(
		&dynamoDBTableName,
		"dynamodb-table-name",
		"",
		// FIXME: un-hardcode this default
		"nullserve-api-cbdec46580e5a391",
		"The dynamodb table to find sites in",
	)
	_ = viper.BindPFlags(rootCmd.Flags())
}

func initConfig() {
	viper.SetEnvPrefix("service")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func root(cmd *cobra.Command, _ []string) {
	cfg := &config.OriginService{}
	cfg.HostSuffix, _ = cmd.Flags().GetString("domain-suffix")
	cfg.AppPrefix, _ = cmd.Flags().GetString("app-prefix")
	cfg.RefPrefix, _ = cmd.Flags().GetString("ref-prefix")
	cfg.S3Source.BucketId, _ = cmd.Flags().GetString("s3-bucket")
	cfg.S3Source.SiteFolderPrefix, _ = cmd.Flags().GetString("s3-prefix-folder")
	cfg.DynamoDBTableName, _ = cmd.Flags().GetString("dynamodb-table-name")
	service.Main(cfg)
}

func Execute() error {
	return rootCmd.Execute()
}
