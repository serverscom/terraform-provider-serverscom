---
page_title: "Servers.com User-data"
---

# User-data

User-data is a helpful tool to get rid of routine operations after server provisioning. You can get a ready-to-use server with additional software installed and configured according to your specification. The feature is built upon the cloud-init package for Linux operating systems. Cloud-init serves for performing operations while server's initialization, its behavior is defined by a special type of content - user-data. To manage post-install operations via Terraform, we have implemented the "user_data" argument in the [serverscom_dedicated_server](https://registry.terraform.io/providers/serverscom/serverscom/latest/docs/resources/dedicated_server) resource.

The tutorial below will show you in practice how to configure your resource with operations performed after provisioning. As an example, we will create a server with a pre-installed Apache engine. Please note: the provisioned server will be billed according to your plan.
## Preparing

You need to have a Terraform client installed and Servers.com account for performing the actions described further. To run the Servers.com provider click on the `USE PROVIDER` button in the upper right corner of the documentation page and follow the instruction. We will use the following script as the user-data argument:

```

#! /bin/bash
sudo apt-get update
sudo apt-get install -y apache2
sudo systemctl start apache2
sudo systemctl enable apache2
echo "The page was created by the user-data" | sudo tee /var/www/html/index.html

```
## Adding the script to the user data

**User-data inserted in the tf file**

Open the file that contains your terraform resource parameters, in our case it's a `main.tf` file. Paste the script to the resource specification and use the format shown in the example. `<< EOF` and `EOF` frame the script within the `user_data` argument.

```

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
 user_data = << EOF
#! /bin/bash
sudo apt-get update
sudo apt-get install -y apache2
sudo systemctl start apache2
sudo systemctl enable apache2
echo "The page was created by the user-data" | sudo tee /var/www/html/index.html
EOF
  
}

```

**User-data located in another file**

If you want to place your script in another file, use the `file()` function. For this guide, we keep our script in the `user-data-apache.sh` file. Its content looks exactly the same as it's shown in the Preparing section. The file is located in one directory with the `main.tf` file. This is how the configuration will look like:
```
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
 user_data = "${file("user-data-apache.sh")}"
  
}
```
## Applying changes

When you have completed a configuration of a server, save changes and make commands to initialize and apply the configuration:

```

$ terraform init && terraform apply

```
## Checking the result

When the server is provisioned, enter your server's IP-address in the browser. As a result, you will see an HTML page with this text:
`The page was created by the user-data`.
