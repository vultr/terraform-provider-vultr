Terraform Provider for Vultr
==================

- Quickstart Guide: [How to Provision a Vultr Cloud Server with Terraform and Cloud-Init](https://www.vultr.com/docs/provision-a-vultr-cloud-server-with-terraform-and-cloud-init/)
- Vultr Website: https://www.vultr.com
- Terraform Website: https://www.terraform.io

<img src="https://raw.githubusercontent.com/hashicorp/terraform-website/master/public/img/logo-hashicorp.svg" width="600px">

Requirements
------------

-   [Terraform](https://www.terraform.io/downloads.html) 0.12.x+
-   [Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/vultr/terraform-provider-vultr`

```sh
$ mkdir -p $GOPATH/src/github.com/vultr; cd $GOPATH/src/github.com/vultr
$ git clone git@github.com:vultr/terraform-provider-vultr.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/vultr/terraform-provider-vultr
$ make build
```

Using the provider
----------------------

See the [Vultr Provider documentation](website/docs/index.html.markdown) to get started using the Vultr provider.

Please read about [V2 changes from V1](example/V2Changes.md) for a list of new changes made to the Vultr Terraform Provider

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-vultr
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

In order to run a specific acceptance test, use the `TESTARGS` environment variable. For example, the following command will run `TestAccVultrUser_base` acceptance test only:

```sh
$ make testacc TESTARGS='-run=TestAccVultrUser_base'
```
