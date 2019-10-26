Terraform Provider
==================


- **Note:** This is a [fork](https://github.com/terraform-providers/terraform-provider-aws).
  - This was forked to add one more optional field `user_data_param_rendered` for `aws_instance` resource. 
  - This field is Optional and modifying this field doesn't force a new resource creation.
  - Rather than showing a checksum of the user_data, this fork will render the output as the value of `user_data_param_rendered` key.
  
  - Example: 
  
  ```hcl
  resource "aws_instance" "some_instance" {
  ami                    = "${data.aws_ami.some_ami.id}"
  instance_type          = "t2.small"
  iam_instance_profile   = "${var.instance_profile}"
  user_data_param_rendered = "" // ** Yes, this has to be explicitely defined. **
  user_data = <<EOF
  #cloud-config
  bootcmd:
  - echo "===========starting userdata=========="
  - export FQDN=${format("server02d.%s", count.index + 1, var.r53_zone_name)}
  - hostname $FQDN
  - hostname > /etc/hostname
  EOF
  ---------  snip  ---------
  }
  ```
  - Output: 
  ```bash
   ~ module.some_module.aws_instance.some_instance[1]
      user_data_param_rendered: "" =>  "#cloud-config\nbootcmd:\n  - echo \"===========starting userdata==========\"\n  - export FQDN=server01.some-r53-zone-name.com\n  - hostname $FQDN\n hostname > /etc/hostname\n"
  
  ```
**`user_data_param_rendered` is not atomic and hence only useful in viewing the rendered template in a `terraform plan`. Since the value of the field `user_data` of type `map[string]*schema.Schema` needs to be shown at another field `user_data_param_rendered`, a lame attempt is made using a variable of type string global to `resourceAwsInstance` function inside `terraform-provider-aws/aws/resource_aws_instance.go`, and because of no atomicity, interpolated values like `count.index` in user_data cannot be relied upon. This will affect only if `count > 1`. If there is a way we can get the value of count.index which is stored somewhere in any of the Resource schemas or through any of the helper functions, we can use it to store it in a slice and get it by the index. As of now, I cannot find a documentation that states this with my limited knowledge.**


- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.10+
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

Developing the Provider
---------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (please check the [requirements](https://github.com/terraform-providers/terraform-provider-aws#requirements) before proceeding).

*Note:* This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard GOPATH (i.e `$HOME/development/terraform-providers/`).

Clone repository to: `$HOME/development/terraform-providers/`

```sh
$ mkdir -p $HOME/development/terraform-providers/; cd $HOME/development/terraform-providers/
$ git clone git@github.com:terraform-providers/terraform-provider-aws
...
```

Enter the provider directory and run `make tools`. This will install the needed tools for the provider.

```sh
$ make tools
```

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-aws
...
```

Using the Provider
----------------------

To use a released provider in your Terraform environment, run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the provider. To specify a particular provider version when installing released providers, see the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).

To instead use a custom-built provider in your Terraform environment (e.g. the provider binary from the build instructions above), follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

For either installation method, documentation about the provider specific configuration options can be found on the [provider's website](https://www.terraform.io/docs/providers/aws/index.html).

Testing the Provider
---------------------------

In order to test the provider, you can run `make test`.

*Note:* Make sure no `AWS_ACCESS_KEY_ID` or `AWS_SECRET_ACCESS_KEY` variables are set, and there's no `[default]` section in the AWS credentials file `~/.aws/credentials`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run. Please read [Running an Acceptance Test](https://github.com/terraform-providers/terraform-provider-aws/blob/master/.github/CONTRIBUTING.md#running-an-acceptance-test) in the contribution guidelines for more information on usage.

```sh
$ make testacc
```

Contributing
---------------------------

Terraform is the work of thousands of contributors. We appreciate your help!

To contribute, please read the contribution guidelines: [Contributing to Terraform - AWS Provider](.github/CONTRIBUTING.md)

Issues on GitHub are intended to be related to bugs or feature requests with provider codebase. See https://www.terraform.io/docs/extend/community/index.html for a list of community resources to ask questions about Terraform.

