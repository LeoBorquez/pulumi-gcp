package main

import (
	"encoding/json"

	container "github.com/pulumi/pulumi-google-native/sdk/go/google/container/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type Config struct {
	Cluster                    string // cluster name
	Description                string // the description of the cluster
	ConfProject                string // the project from which the secret manager will obtain the settings
	Machine                    string // machine type
	NodeName                   string // node name
	ServiceAccount             string
	Network                    string // network name
	SubNetwork                 string // subnetwork name
	PodsSecondaryRangeName     string
	ServicesSecondaryRangeName string
	Location                   string // location where the cluster will be created
	MasterIpv4CidrBlock        string //nolint:words
	Tags                       []string
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// binding configuration object
		var c Config
		conf := config.New(ctx, "native-gke-cluster-go")
		conf.RequireObject("var", &c)
		project := conf.Require("project")
		data := conf.Require("labels")

		var v interface{}
		if err := json.Unmarshal([]byte(data), &v); err != nil {
			return err
		}

		var labels = map[string]string{}
		for k, v := range v.(map[string]interface{}) {
			labels[k] = v.(string)
		}

		// Necessary configurations to create a cluster
		nodeConfig := &container.NodeConfigArgs{
			MachineType: pulumi.StringPtr(c.Machine),
			Labels:      pulumi.ToStringMap(labels),
			Tags:        pulumi.ToStringArray([]string{c.Tags[0]}),
			OauthScopes: pulumi.ToStringArray([]string{"https://www.googleapis.com/auth/cloud-platform"}),
			ShieldedInstanceConfig: &container.ShieldedInstanceConfigArgs{
				EnableSecureBoot:          pulumi.Bool(true),
				EnableIntegrityMonitoring: pulumi.Bool(true),
			},
			ServiceAccount: pulumi.StringPtr(c.ServiceAccount),
		}

		nodePools := &container.NodePoolTypeArray{&container.NodePoolTypeArgs{
			Config:           nodeConfig,
			InitialNodeCount: pulumi.IntPtr(2),
			Name:             pulumi.StringPtr(c.NodeName),
		}}

		// Cidr Blocks
		cidrAODC := &container.CidrBlockArgs{
			CidrBlock:   pulumi.StringPtr("17.0.0.0/8"),
			DisplayName: pulumi.StringPtr("AODC"), // optional
		}

		var ipv4 string
		if c.MasterIpv4CidrBlock != "" { //nolint:words
			ipv4 = c.MasterIpv4CidrBlock //nolint:words
		} else {
			ipv4 = "192.168.23.16/28"
		}

		newCluster, err := container.NewCluster(ctx, c.Cluster, &container.ClusterArgs{
			Name:        pulumi.StringPtr(c.Cluster),
			Project:     pulumi.StringPtr(project),
			Description: pulumi.StringPtr(c.Description),
			Location:    pulumi.StringPtr(c.Location),
			NodePools:   nodePools,
			Network:     pulumi.StringPtr(c.Network),
			ReleaseChannel: &container.ReleaseChannelArgs{
				Channel: container.ReleaseChannelChannelStable,
			},
			Autoscaling: &container.ClusterAutoscalingArgs{
				AutoprovisioningNodePoolDefaults: &container.AutoprovisioningNodePoolDefaultsArgs{
					// Shield vm policy
					ShieldedInstanceConfig: &container.ShieldedInstanceConfigArgs{
						EnableSecureBoot:          pulumi.Bool(true),
						EnableIntegrityMonitoring: pulumi.Bool(true),
					},
				},
			},
			Subnetwork: pulumi.StringPtr(c.SubNetwork),
			MasterAuthorizedNetworksConfig: &container.MasterAuthorizedNetworksConfigArgs{ //nolint:words
				CidrBlocks: &container.CidrBlockArray{cidrAODC},
				Enabled:    pulumi.Bool(true),
			},
			// Private Endpoints
			PrivateClusterConfig: &container.PrivateClusterConfigArgs{
				EnablePrivateEndpoint: pulumi.Bool(false),
				EnablePrivateNodes:    pulumi.Bool(true),
				MasterIpv4CidrBlock:   pulumi.StringPtr(ipv4), //nolint:words
			},
			IpAllocationPolicy: &container.IPAllocationPolicyArgs{
				UseIpAliases:               pulumi.Bool(true),
				CreateSubnetwork:           pulumi.Bool(false),
				ClusterSecondaryRangeName:  pulumi.StringPtr(c.PodsSecondaryRangeName),
				ServicesSecondaryRangeName: pulumi.StringPtr(c.ServicesSecondaryRangeName),
			},
		})
		if err != nil {
			return err
		}

		ctx.Log.Info("New cluster created", &pulumi.LogArgs{
			Resource: newCluster,
		})

		return nil
	})
}
