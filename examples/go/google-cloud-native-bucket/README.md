# Pulumi

Example to create a bucket with the native provider

* Make sure you have permissions to create resources in the projects (storage.buckets.create).
## Config

- In the Pulumi.dev.yaml are the configuration variables of the example. native-bucket-go is the namespace

    native-bucket-go:project: is the name of the project where the bucket will be created
    native-bucket-go:bucket: is the name of new bucket


## Running the App

1. Inside the project folder run (to download and install packages):

    ```bash
    go mod download
    go install
    ```
2. Run  `pulumi stack` and create a stack, in the example the stack name is "dev"

3.  Run `pulumi up -y` to preview and deploy changes:

4. Cleanup. after testing.

    ```bash
    $ pulumi destroy
    $ pulumi stack rm
    ```


*  Make sure your account has the correct permissions (storage.buckets.create)
