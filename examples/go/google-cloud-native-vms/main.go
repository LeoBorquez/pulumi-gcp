package main

import (
	"fmt"

	native "github.com/pulumi/pulumi-google-native/sdk/go/google/compute/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type Config struct {
	Project    string   // project where the vm will be created
	Script     string   // startup script
	Image      string   // image name
	Zone       string   // zone where the vm will be created
	PNetwork   string   // project where the network is located
	Network    string   // network name
	Subnetwork string   // subtnetwork name
	Region     string   // region of the subnetwork
	Instance   string   // instance name
	Disk       string   // disk name
	Machine    string   // machine type
	Tags       []string // tags
}

func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {

		// binding configuration object
		var c Config
		conf := config.New(ctx, "native-vms-go")
		conf.RequireObject("var", &c)
		project := conf.Require("project")

		startupScript := fmt.Sprintf("`#!/bin/bash %s `", c.Script)
		/*
			You can also obtain the networks from a shared vpc project, but due to the organization's policies it is not possible to use this feature.
			The account used must have permissions on the shared vpc project.

			shared, err := compute.GetSharedVPCHostProject(ctx, "shared", pulumi.ID("gns-network-prod-0d38"), nil)
			if err != nil {
				return err
			}
			fmt.Println(shared.URN())
		*/

		// disk creation
		disk, err := native.NewDisk(ctx, c.Disk, &native.DiskArgs{
			Project: pulumi.StringPtr(project),
			/*
				gcloud compute images list --project=gcp-at-apple-image-factory
				Image example projects/ubuntu-os-cloud/global/images/ubuntu-apple-pulumi-01
			*/
			SourceImage: pulumi.StringPtr(c.Image),
			Zone:        pulumi.StringPtr(c.Zone),
		})
		if err != nil {
			ctx.Log.Error("Error creating Disk", &pulumi.LogArgs{
				Resource: disk,
			})
			return err
		}

		// To attach disk
		diskToAttach := &native.AttachedDiskArgs{
			AutoDelete: pulumi.BoolPtr(true),
			Boot:       pulumi.BoolPtr(true),
			/*
				Implicit URL
				Source:     pulumi.StringPtr("https://www.googleapis.com/compute/v1/projects/"project_name"/zones/us-central1-a/disks/"disk_name""),
			*/
			Source: disk.SelfLink,
		}

		/*
			Network discovery, gets a list of available network

			network, err := native.LookupNetwork(ctx, &native.LookupNetworkArgs{
				Network: "network name",
				Project: "project",
			})
			if err != nil {
				return err
			}
		*/

		network := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", c.PNetwork, c.Network)
		subnetwork := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/regions/%s/subnetworks/%s", c.PNetwork, c.Region, c.Subnetwork)

		networkInterface := &native.NetworkInterfaceArgs{
			/*
				Network: pulumi.StringPtr(network.SelfLink), using network discovery, the network's self link can be obtained
				Implicit URL
				Network: pulumi.StringPtr("https://www.googleapis.com/compute/v1/projects/"network project"/global/networks/"network"),
			*/
			Network:    pulumi.StringPtr(network),
			Subnetwork: pulumi.StringPtr(subnetwork),
			/*
				Subnetwork: pulumi.StringPtr(network.Subnetworks[2]), // network returns a list with the subnetworks
				Implicit URL
				Subnetwork: pulumi.StringPtr("https://www.googleapis.com/compute/v1/projects/"network project"/regions/us-central1/subnetworks/"subnetwork"),
			*/
		}

		// metadata to initialize the instance
		metadata := &native.MetadataItemsItemArgs{
			Key:   pulumi.StringPtr("startup-script"),
			Value: pulumi.StringPtr(startupScript),
		}

		// vm creation
		vm, err := native.NewInstance(ctx, c.Instance, &native.InstanceArgs{
			Project:     pulumi.StringPtr(project),
			Name:        pulumi.StringPtr(c.Instance),
			Description: pulumi.StringPtr("A template to create app instances."),
			Tags: &native.TagsArgs{
				Items: pulumi.StringArray{pulumi.String(c.Tags[0]), pulumi.String(c.Tags[1])},
			},
			MachineType:  pulumi.StringPtr(c.Machine),
			Zone:         pulumi.StringPtr(c.Zone),
			CanIpForward: pulumi.BoolPtr(false),
			ShieldedInstanceConfig: &native.ShieldedInstanceConfigArgs{
				EnableSecureBoot:          pulumi.BoolPtr(true),
				EnableVtpm:                pulumi.BoolPtr(true),
				EnableIntegrityMonitoring: pulumi.BoolPtr(true),
			},
			Disks:             &native.AttachedDiskArray{diskToAttach},
			NetworkInterfaces: &native.NetworkInterfaceArray{networkInterface},
			Scheduling: &native.SchedulingArgs{
				AutomaticRestart:  pulumi.BoolPtr(true),
				OnHostMaintenance: native.SchedulingOnHostMaintenancePtr("MIGRATE"),
			},
			Metadata: &native.MetadataArgs{
				Items: &native.MetadataItemsItemArray{metadata},
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("diskSelfLink", disk.SelfLink) // these values can be imported to use in another module
		ctx.Export("vm Description", vm.Description)

		return nil
	})
}
