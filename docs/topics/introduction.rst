.. _introduction:

Setup guide
***********

This guide will help you to set up Servers.com as a Terraform provider. Follow the steps described below to create the correct provider's configuration and run it in a proper way.

1) Download Terraform from the `official page <https://www.terraform.io/downloads.html>`_ and install according to the `instructions <https://learn.hashicorp.com/terraform/getting-started/install.html>`_.

2) After a successful installation, create a directory with any name to store Terraform configuration files. We take "serverscom-terraform" as an example.

3) Create a configuration file "main.tf" (it can have another name). The file must be located in the "serverscom-terraform" directory and have the ".tf" extension.

4) Write the following configuration in the "main.tf":

 .. sourcecode:: terraform

   provider "serverscom" {
     token = "<your API token>"
     endpoint = "https://api.servers.com/v1"
   }

 To see a description of the attributes, go to Servers.com Provider section.

5) Download Servers.com provider binary from `this repository <https://github.com/serverscom/terraform-provider-serverscom/releases>`_. Choose an appropriate archive by its name according to your operating system. All the files are named this way: ``terraform-provider-serverscom-<version X.Y.Z>-<Operating system>-<Architecture>``

 For example: ``terraform-provider-serverscom-v0.1.1-windows-amd64.zip``

6) Extract the archive in the following directory (it depends on an operating system):

 - for Windows: ``%APPDATA%\terraform.d\plugins\windows_amd64``
 - for Linux: ``~/.terraform.d/plugins/linux_amd64``
 - for MacOS: ``~/.terraform.d/plugins/darwin_amd64``

7) Run ``terraform init`` command in the directory "serverscom-terraform" to download a Servers.com plugin, it will be saved in the current directory. Thus, Terraform will be able to interact with Servers.com.

8) Terraform is ready to work with the Servers.com infrastructure, just complete the configuration file with a description of a desired infrastructure.

Here are some basic commands that may be useful:
 - ``terraform plan`` - see what actions Terraform will make to get a described configuration;
 - ``terraform apply`` - start a creating of an infrastructure according to a described configuration.
