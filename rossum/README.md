<a href="https://terraform.io">
    <img src=".github/terraform_logo.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>

# Terraform Provider for AWS fork

This fork integrates these two PRs:

- [Add resource aws_db_snapshot_copy](https://github.com/hashicorp/terraform-provider-aws/pull/9886) in the [add-resource-db-snapshot-copy](https://github.com/rossumai/terraform-provider-aws/tree/add-resource-db-snapshot-copy) branch
- [Data source aws_db_snapshot: Shared and public snapshot fix](https://github.com/hashicorp/terraform-provider-aws/pull/5767) in the [fix-shared-aws_db_snapshot](https://github.com/rossumai/terraform-provider-aws/tree/fix-shared-aws_db_snapshot) branch


I have just rebased them onto the current code and fixed a few problems that cropped up. Thanks to the original authors.

The main branch combines the two PRs and contains a few more commits to make the automatic builds work. The packages can be pulled from the [terraform registry](https://registry.terraform.io/providers/rossumai/aws/)

The [rebasescript.sh](rossum/rebasescript.sh) script can be used to bring this repo up to date with upstream.

\
\
\
\
[![Forums][discuss-badge]][discuss]

[discuss-badge]: https://img.shields.io/badge/discuss-terraform--aws-623CE4.svg?style=flat
[discuss]: https://discuss.hashicorp.com/c/terraform-providers/tf-aws/

- Website: [terraform.io](https://terraform.io)
- Tutorials: [learn.hashicorp.com](https://learn.hashicorp.com/terraform?track=getting-started#getting-started)
- Forum: [discuss.hashicorp.com](https://discuss.hashicorp.com/c/terraform-providers/tf-aws/)
- Chat: [gitter](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing List: [Google Groups](http://groups.google.com/group/terraform-tool)

The Terraform AWS provider is a plugin for Terraform that allows for the full lifecycle management of AWS resources.
This provider is maintained internally by the HashiCorp AWS Provider team.

Please note: We take Terraform's security and our users' trust very seriously. If you believe you have found a security issue in the Terraform AWS Provider, please responsibly disclose by contacting us at security@hashicorp.com.

## Quick Starts

- [Using the provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [Provider development](docs/contributing)

## Documentation

Full, comprehensive documentation is available on the Terraform website:

https://terraform.io/docs/providers/aws/index.html

## Roadmap

Our roadmap for expanding support in Terraform for AWS resources can be found in our [Roadmap](ROADMAP.md) which is published quarterly.

## Frequently Asked Questions

Responses to our most frequently asked questions can be found in our [FAQ](docs/contributing/faq.md )

## Contributing

The Terraform AWS Provider is the work of thousands of contributors. We appreciate your help!

To contribute, please read the contribution guidelines: [Contributing to Terraform - AWS Provider](docs/contributing)
