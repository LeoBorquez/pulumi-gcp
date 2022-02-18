package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	cloudkms "github.com/pulumi/pulumi-google-native/sdk/go/google/cloudkms/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type Config struct {
	PreFixKeyRing string
	PreFixKey     string
	Location      string
	Rotation      int
	NextRotation  int
	Members       []string
	Role          string
}

type Timestamp time.Time

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		var c Config
		conf := config.New(ctx, "native-kms-go")
		conf.RequireObject("var", &c)
		project := conf.Require("project")

		// create random string
		random := randomString(4)
		//  "keyring-prefix-random_characters"
		keyRingID := fmt.Sprintf("keyring-%s-%s", c.PreFixKeyRing, random)
		keyID := fmt.Sprintf("key-%s-%s", c.PreFixKey, random)

		rotation := int64(60 * 60 * c.Rotation * 30)
		nextRotation := time.Now().Add(time.Duration(c.NextRotation) * time.Hour)

		_, err := cloudkms.NewKeyRing(ctx, keyRingID, &cloudkms.KeyRingArgs{
			KeyRingId: pulumi.String(keyRingID),
			Location:  pulumi.StringPtr(c.Location),
			Project:   pulumi.StringPtr(project),
		})
		if err != nil {
			return err
		}

		cryptoKey, err := cloudkms.NewCryptoKey(ctx, keyID, &cloudkms.CryptoKeyArgs{
			CryptoKeyId:      pulumi.String(keyID),
			KeyRingId:        pulumi.String(keyRingID),
			RotationPeriod:   pulumi.StringPtr(strconv.Itoa(int(rotation)) + "s"),
			NextRotationTime: pulumi.StringPtr(nextRotation.Local().Format(time.RFC3339)), // timestamp 2006-01-02T15:04:05Z07:00
			Location:         pulumi.StringPtr(c.Location),
			Project:          pulumi.StringPtr(project),
			Purpose:          cloudkms.CryptoKeyPurposeEncryptDecrypt,

			/*
				Depending on the type of key the encryption algorithm can be configured.

				VersionTemplate: cloudkms.CryptoKeyVersionTemplateArgs{
					Algorithm: cloudkms.CryptoKeyVersionTemplateAlgorithmRsaSignPss3072Sha256,
				},
			*/
		})
		if err != nil {
			return err
		}

		_, err = cloudkms.NewKeyRingCryptoKeyIamPolicy(ctx, "iamPolicy", &cloudkms.KeyRingCryptoKeyIamPolicyArgs{
			Bindings: cloudkms.BindingArray{&cloudkms.BindingArgs{
				Members: pulumi.ToStringArray([]string{c.Members[0]}),
				Role:    pulumi.StringPtr(c.Role),
			}},
			CryptoKeyId: pulumi.String(keyID),
			KeyRingId:   pulumi.String(keyRingID),
			Location:    pulumi.StringPtr(c.Location),
			Project:     pulumi.StringPtr(project),
		})
		if err != nil {
			return err
		}

		ctx.Log.Info("Cripto Key created", &pulumi.LogArgs{
			Resource: cryptoKey,
		})

		return nil
	})
}

func randomString(len int) string {
	random := make([]byte, len)
	for i := 0; i < len; i++ {
		random[i] = byte(randomInt(97, 122))
	}
	return string(random)
}

func randomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}
