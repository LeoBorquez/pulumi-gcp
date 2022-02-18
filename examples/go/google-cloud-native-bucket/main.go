package main

import (
	"fmt"

	storage "github.com/pulumi/pulumi-google-native/sdk/go/google/storage/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// "google-native" is the namespace, is in the config file Pulumi.dev.yaml
		conf := config.New(ctx, "native-bucket-go")
		project := conf.Require("project")
		bucketName := conf.Require("bucket")

		bucket, err := storage.NewBucket(ctx, bucketName, &storage.BucketArgs{
			Name:    pulumi.StringPtr(bucketName),
			Project: pulumi.StringPtr(project),
		})
		if err != nil {
			return err
		}

		// Export the bucket self-link
		ctx.Export("bucketSelfLink", bucket.SelfLink)
		fmt.Println("bucket SelfLink :: ", bucket.SelfLink)

		return nil
	})
}
