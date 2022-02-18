# GKE

Example to create a GKE cluster with the native provider

* Make sure you have permissions to create resources in the projects.

## Config

- In the Pulumi.dev.yaml are the configuration variables of the example, in this example could be use the secret manager

        native-gke-cluster-go:project: ID of the project where the resources will be created
        native-gke-cluster-go:var: Is an array of configurations (pulumi way)
                cluster: The name of cluster
                description: Description of the cluster
                confProject: The project from which the secret keys will be obtained
                machine: Machine type
                nodeName: Name of the node pool
                serviceAccount: The google cloud platform service account to be used by the node VMs.
                network: In case you do not use the secret manager, you can use the name of the network
                subNetwork: In case you do not use the secret manager, you can use the name of the subnetwork
                location: Location of the cluster (zone)
                podsSecondaryRangeName: The name of the secondary range to be used for the cluster CIDR block
                servicesSecondaryRangeName: The name of the secondary range to be used as for the services CIDR block
                masterIpv4CidrBlock: The IP range in CIDR notation to use for the hosted master network. Optional (default value 172.16.0.0/28)
                labels: Metadata information
                tags: Network tags


## Running the Example

1. Inside the project folder run (to download and install packages):

    ```bash
    go mod download
    go install
    ```
2. Run  `pulumi stack` and create a stack, in the example the stack name is "dev"

2.  Run `pulumi up -y` to preview and deploy changes:

3. Cleanup. after testing.

    ```bash
    pulumi destroy
    pulumi stack rm
    ```
