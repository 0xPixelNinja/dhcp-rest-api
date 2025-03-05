# Configuration

This document outlines the environment variables and configuration files required to run the DHCP REST API.

## Environment Variables

The project uses a `.env` file to manage configuration settings. Below is an example `.env` file:

```env
# .env

# Path to the ISC DHCP configuration file
DHCP_CONF_PATH=/etc/dhcp/dhcpd.conf

# Secret token for API authentication
TOKEN_SECRET=your-secret-token

# Path to the interfaces configuration file
# By default, the interface configuration is located at /etc/default/isc-dhcp-server
INTERFACES_CONF_PATH=/etc/default/isc-dhcp-server
```

### Explanation

- **DHCP_CONF_PATH:**  
  The file path to your ISC DHCP configuration file. The default is `/etc/dhcp/dhcpd.conf`, which is used by the ISC DHCP server to specify on which hosts it should listen.

- **TOKEN_SECRET:**  
  A secret token used to secure all API endpoints. Ensure you choose a strong token.

- **INTERFACES_CONF_PATH:**  
  The file path for the network interfaces configuration. The default is `/etc/default/isc-dhcp-server`, which is used by the ISC DHCP server to specify on which interfaces it should listen.

## Updating Configuration

If your configuration files are located in different paths, update the values in your `.env` file accordingly. This approach makes it easy to adapt the API to various environments without changing the code.
