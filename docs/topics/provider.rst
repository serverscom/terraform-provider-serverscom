.. _provider:

Servers.com Provider
********************

The Servers.com Provider allows to interact with Servers.com services. The provider has to be set up properly before using, we recommend you to get acquainted with the Getting started instruction. To see available resources and its description, use the navigation.

Example Usage
=============

.. sourcecode:: terraform

  # Configure the Servers.com Provider
  provider "serverscom" {
    token = "<your API token>"
    endpoint = "https://api.servers.com/v1"
  }

  # Create a dedicated server
  resource "serverscom_dedicated_server" "node_1" {
    ...
  }

Argument Reference
==================

- ``token`` (Required, string) - The token used to perform API-requests for Servers.com services, it can be issued in the `Customer Portal <https://portal.servers.com/#/profile/api-tokens>`_.
- ``endpoint`` (Optional, string) - The Servers.com API endpoint. In most of cases the default one is used: https://api.servers.com/v1.
