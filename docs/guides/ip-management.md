---
page_title: "IP management"
---
# IP management

Some network operations, such as a DNS records update, require some time to be completed. The Servers.com Terraform provider gives an opportunity to retrieve IP-addresses of a dedicated server before the provisioning process starts. It allows simultaneously configuring DNS records, networks, or software provisioning tools while the server is being prepared and is not yet ready to use.

We will consider a Terraform configuration with simultaneous dedicated servers being provisioned, and will set up the DNS A records for these hosts.
## Preparing

This guide requires the Servers.com Terraform provider set up, and a Servers.com account to be present from your side. Please note, steps performed further will lead to the creation of a dedicated server that will be billed according to your plan.

You will need the following resources and data sources:

- [serverscom_dedicated_server](https://registry.terraform.io/providers/serverscom/serverscom/latest/docs/resources/dedicated_server) - resource to complete the server's specification;

- [serverscom_subnetwork](https://registry.terraform.io/providers/serverscom/serverscom/latest/docs/resources/subnetwork) - resource for network management;

- [serverscom_network_pool](https://registry.terraform.io/providers/serverscom/serverscom/latest/docs/data-sources/network_pool) - data source to obtain the network pool information.
## Writing the specification

You need to write data section in the specification to point a necessary network pool with its ID (obtained through a support request):

```

data "serverscom_network_pool" "my_pool" {
  by {
    id = "my_pool_id"
  }
}

```

Next, it's necessary to define the subnetwork size that will be allocated from the pool. The default subnetwork size is /29 (shown in the example below), you should leave this setting as is in most of the cases:

```

resource "serverscom_subnetwork" "my_network_1" {
    pool_id = serverscom.serverscom_network_pool.my_pool.id
    subnet = 29
}

```

Specify the A record's parameters:

```

resource "some_dns_resource", "my_network_a_record_1" {
    type = "A"
    data = cidrhost(serverscom.serverscom_subnetwork.my_network_1.cidr, 4)
}

```

- `type` - type of the record, in our case it's the A record;

- `4` value in the `data` parameter means that the fourth IP-address of the subnetwork will be taken for the record.

Finally, you should complete the server's configuration:

```

resource "serverscom_dedicated_server" "my_dedicated_server_1" {
    ...
    public_network_id = serverscom_subnetwork.my_network_1.id
	...
}

```

As a result, you will get the following specification for two dedicated servers provisioning with the A record set up for each one:

```

data "serverscom_network_pool" "my_pool" {
  by {
    id = "my_pool_id"
  }
}
resource "serverscom_subnetwork" "my_network_1" {
    pool_id = serverscom.serverscom_network_pool.my_pool.id
    subnet = 29
}
resource "some_dns_resource", "my_network_a_record_1" {
    type = "A"
    data = cidrhost(serverscom.serverscom_subnetwork.my_network_1.cidr, 4)
}
resource "serverscom_subnetwork" "my_network_2" {
    pool_id = serverscom.serverscom_network_pool.my_pool.id
    subnet = 29
}
resource "some_dns_resource", "my_network_a_record_2" {
    type = "A"
    data = cidrhost(serverscom.serverscom_subnetwork.my_network_2.cidr, 4)
}
resource "serverscom_dedicated_server" "my_dedicated_server_1" {
    ...
    public_network_id = serverscom_subnetwork.my_network_1.id
}
resource "serverscom_dedicated_server" "my_dedicated_server_2" {
    ...
    public_network_id = serverscom_subnetwork.my_network_2.id
}

```
