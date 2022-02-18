# KMS

Example to create a KMS key and key ring

* Make sure you have permissions to create resources in the projects.

## Config

- In the Pulumi.dev.yaml are the configuration variables of the example, in this example could be use the secret manager

        native-kms-go:project: ID of the project where the resources will be created
        native-kms-go:var: Is an array of configurations (pulumi way)
            preFixKeyRing: prefix of the key ring
            preFixKey: prefix of the key
            location: location where to create the resource
            role: role assigned to the members
            rotation: initial rotation time (in hours) must be at least 24 hours
            nextRotation: in hours

## Running the Example

1. Inside the project folder run (to download and install packages):

    ```
    go mod download
    go install
    ```
2. Run  `pulumi stack` and create a stack, in the example the stack name is "dev"

2.  Run `pulumi up -y` to preview and deploy changes:

3. Cleanup. after testing.

    ```
    $ pulumi destroy
    $ pulumi stack rm
    ```
