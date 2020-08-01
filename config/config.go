package config

type OriginService struct {
	S3Source          S3Source
	HostSuffix        string
	AppPrefix         string
	RefPrefix         string
	DynamoDBTableName string
}

type S3Source struct {
	BucketId         string
	SiteFolderPrefix string
}
