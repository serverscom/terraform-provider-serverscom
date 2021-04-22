---
page_title: "Servers.com User-data"
---

# User-data

User-data is a helpful tool to get rid of routine operations after server provisioning. You can get a ready-to-use server with additional software installed and configured according to your specification. The feature is built upon the cloud-init package for Linux operating systems. Cloud-init serves for performing operations while server's initialization, its behavior is defined by a special type of content - user-data. To manage post-install operations via Terraform, we have implemented the "user_data" argument in the [serverscom_dedicated_server](https://registry.terraform.io/providers/serverscom/serverscom/latest/docs/resources/dedicated_server) resource.

The tutorial below will show you in practice how to configure your resource with operations performed after provisioning. As an example, we will create a server with a Go web application deployed and SSH key access. Please note: the provisioned server will be billed according to your plan.

## Preparing
To perform the actions described further, you need to have a Terraform client and Servers.com account. The Go application will be taken from the [Hashicorp repository](https://github.com/hashicorp/learn-terraform-provisioning/tree/cloudinit/).

First of all, you need to make a clone of the repository:
```
$ git clone -b cloudinit https://github.com/hashicorp/learn-terraform-provisioning
```
Open the directory containing the clone:
```
$ cd learn-terraform-provisioning
```
## SSH key creation
**Mac or Linux terminal**

To generate a local SSH key associated with the new terraform user, make the following command in the terminal:
```
$ ssh-keygen -t rsa -C "email@example.com" -f ./tf-cloud-init
```
- `tf-cloud-init` - name given to the key;
- `f `- according to this flag, `tf-cloud-init` and `tf-cloud-init.pub` files will be created in the current directory;
- `email@example.com` - your email.

**Windows Putty client**

To generate an SSH key in Windows, use Putty and follow [these instructions](https://www.ssh.com/academy/ssh/putty/windows/puttygen).

## Adding public SSH key to the script
Open the `scripts/add-ssh-web-app.yaml` file in the cloned repository. Copy and paste the public key from the `tf-cloud-init.pub` into the `ssh_authorized_keys` parameter. Further, you will need to refer to this script in the server's `user_data` argument.
```
##...
users:
  - default
  - name: terraform
    gecos: terraform
    primary_group: hashicorp
    sudo: ALL=(ALL) NOPASSWD:ALL
    groups: users, admin
    ssh_import_id:
    lock_passwd: false
    ssh_authorized_keys:
      -  # "Your public SSH key from the tf-cloud-init.pub"
##...
```
For more details, you can read the [cloud-init documentation](https://cloudinit.readthedocs.io/en/latest/topics/format.html).

## Adding the script to the user-data
Go to the directory containing the `main.tf `file. In our case, it's the `serverscom-terraform`:
```
$ cd serverscom-terraform
```
Open the `main.tf` file and the `template_file` block and fill the `user_data` as it's described below:
```
data "template_file" "user_data" {
  template = file("../scripts/add-ssh-web-app.yaml")
 
resource "serverscom_dedicated_server" "node_1" {
  hostname             = "node-1"
  location             = "SJC1"
  server_model         = "Dell R440 / 2xIntel Xeon Silver-4114 / 32 GB RAM / 1x480 GB SSD"
  ram_size             = 32
 
  operating_system     = "Ubuntu 20.04-server x86_64"
 
  private_uplink       = "Private 10 Gbps with redundancy"
  public_uplink        = "Public 10 Gbps with redundancy"
  bandwidth            = "200000 GB"
  # ...
  # Some parameters are not displayed here to shorten the specification.
  # You can see the complete example of the resource in the relevant section of the documentation.
  # ...
  user_data            = data.template_file.user_data.rendered
   
}
```
The `template_file` block refers to the `add-ssh-web-app.yaml` that will be processed as the `user-data` that will be initialized after provisioning.

Save changes and make commands to initialize and apply the configuration:
```
$ terraform init && terraform apply
```
After a successful operation, you will have a dedicated server with a configured SSH key and Go web application.

## Checking the result

To connect to the server via SSH, use the following command in the terraform console:
```
$ ssh terraform@$(terraform output -raw public_ip) -i ../tf-cloud-init
```
- `public_ip` - public IPv4 address of the server.

Thus, your SSH access is established without creating the keys on the provider's side.

Move to the Go directory:
```
$ cd go/src/github.com/hashicorp/learn-go-webapp-demo
```
Start the application:
```
$ go run webapp.go
```
Open you web browser and enter this address: `<server IP>:8080` to see the result.
