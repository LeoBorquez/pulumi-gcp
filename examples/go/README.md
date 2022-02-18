# GCP@Apple

![InnerSource Supported](https://badges.pie.apple.com/badges/custom?t=InnerSource&v=supported&c=green)

The purpose of this repo is to share documentation and example Pulumi scripts for users at Apple.

- [GCP@Apple Site](https://cloudtech.apple.com/documentation/gcp)


## Getting Started


### Setup
- You will need access to an [GCP@Apple][gcp_at_apple] account.

- It is necessary to have the `gcloud` installed. To install, please follow [the install instructions][gcp_cli_install]

All example programs are written in Go 1.17 (for the moment)
- To install Go, please follow [these instructions][go_install]

- To install Pulumi, please follow the [Pulumi install instructions][pulumi_install] (Homebrew)<br>

        ```bash
        $ brew install pulumi
        ```

Homebrew can be used as well [brew.sh][brew_sh]

        ```bash
        $ brew install go
        $ brew install pulumi
        $ brew install cask google-cloud-sdk
        ```

### New Pulumi project
* Make sure you have cli gcp installed in your environment. <br>

1.  Pulumi projects exist in a top-level directory and can be created within an empty directory.

        $ mkdir pulumi-sample-app
        $ cd pulumi-sample-app

2.  Configure/Create pulumi project inside the directory

        $ pulumi login
        $ pulumi new

3. Select `gcp-go` or `google-native-go` (currently in public preview), name your project as you want, stack `<org-name>/sample`, and provide the default gcp region e.g. `us-central-1`

### Common Commands

- `$ pulumi new` create a new project
- `$ pulumi stack` manage your stacks
- `$ pulumi up` preview and create/deploy your infrastructure
- `$ pulumi destroy` delete your infrastructure



[gcp_at_apple]: https://cloudtech.apple.com/documentation/gcp
[go_install]: https://golang.org/doc/install
[brew_sh]: https://brew.sh/
[pulumi_install]: https://www.pulumi.com/docs/get-started/install/
[gcp_cli_install]: https://cloud.google.com/sdk/docs/install
