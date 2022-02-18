[![Deploy](https://get.pulumi.com/new/button.svg)](https://app.pulumi.com/new)

# VMS

Starting point for building the Pulumi web server sample in Google Cloud.

* Make sure you have permissions to create resources in the projects.


## Config

- In the Pulumi.dev.yaml are the configuration variables of the example

        native-vms-go:project: ID of the project where the resources will be created
        native-vms-go:var: Is an array of configurations (pulumi way)
            disk: The name of the disk to be created for use by the vm
            image: Image to create the disk (url)
            instance: The name of the instance to be created
            machine: Machine type
            pnetwork: The ID of the project from which the network will be obtained (VPC shared host)
            network: The name of the network to be used
            subnetwork: The name of the subnetwork to be used
            script: Startup script
            tags: Different tags that can be added
            zone: The location where the vm will be created

## Running the Example

1. Inside the project folder run (to download and install packages):

    ```
    go mod download
    go install
    ```
2. Run `pulumi stack` and create a stack, in the example the stack name is "dev"

3.  Run `pulumi up -y` to preview and deploy changes:

4. Cleanup. after testing.

    ```
    $ pulumi destroy
    $ pulumi stack rm
    ```
