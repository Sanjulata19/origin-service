package config

type StaticHost struct {
	S3Source   S3Source
	HostSuffix string
}

type S3Source struct {
	BucketId         string
	SiteFolderPrefix string
}
