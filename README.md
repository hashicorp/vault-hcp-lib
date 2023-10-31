# Vault Library: HCP Vault Library

The HCP Vault library is a standalone backend library for use with [Hashicorp
Vault](https://www.github.com/hashicorp/vault).

Please note: We take Vault's security and our users' trust very seriously. If
you believe you have found a security issue in Vault, please responsibly
disclose by contacting us at [security@hashicorp.com](mailto:security@hashicorp.com).

## Quick Links

- [Vault Website](https://developer.hashicorp.com/vault)
- [Vault Project GitHub](https://www.github.com/hashicorp/vault)

## Getting Started

This is a [Vault plugin](https://developer.hashicorp.com/vault/docs/plugins)
and is meant to work with Vault. This guide assumes you have already installed
Vault and have a basic understanding of how Vault works.

Otherwise, first read this guide on how to [get started with
Vault](https://developer.hashicorp.com/vault/tutorials/getting-started/getting-started-install).

## Usage

This library is currently built into Vault and can be used by a Vault CLI client.

```sh
$ vault hcp connect
```

It authenticates users or machines to HCP using either provided arguments or retrieved HCP token through
browser login. A successful authentication results in an HCP token and an HCP Vault address being
locally cached.

The default authentication method is an interactive one, redirecting users to the HCP login browser.
If a non-interactive option is supplied, it can be used if provided with a service principal credential
generated through the HCP portal with the necessary capabilities to access the organization, project, and
HCP Vault cluster chosen.

```sh
$ vault hcp connect -non-interactive=true -client-id=client-id-value -secret-id=secret-id-value
```

Additionally, the organization identification, project identification, and cluster name can be passed in to
directly connect to a specific HCP Vault cluster without interacting with the CLI.

```sh
$ vault hcp connect -non-interactive=true -client-id=client-id-value -secret-id=secret-id-value -organization-id=org-UUID -project-id=proj-UUID -cluster-id=cluster-name
```

In order to clean up the cache with the HCP credentials used to connect to a HCP Vault cluster, you can use the disconnect subcommand:

```sh
$ vault hcp disconnect
```

For more information about supported subcommands and options, refer to the [documentation](https://add-documentation-here).

## How to contribute

Thanks for considering contributing to this project. Unfortunately, HashiCorp does not currently accept new contributions for this project.

## License

This code is released under the Mozilla Public License 2.0. Please see [LICENSE](https://github.com/hashicorp/terraform-aws-hcp-consul/blob/main/LICENSE) for more details.