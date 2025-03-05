# Configuration

This document outlines the environment variables and configuration files required to run the DHCP REST API.

## Environment Variables

The project uses a `.env` file to manage configuration settings. Below is an example `.env` file:

```env
# .env

# Path to the ISC DHCP configuration file
DHCP_CONF_PATH=/etc/dhcp/dhcpd.conf

# Path to the ISC DHCP interfaces configuration file
INTERFACES_CONF_PATH=/etc/default/isc-dhcp-server

# Secret token for API authentication
TOKEN_SECRET=your-secret-token

# Optional: DHCP service name (if needed for service reloads)
DHCP_SERVICE_NAME=isc-dhcp-server
```

### Explanation

- **DHCP_CONF_PATH:**  
  The file path to your ISC DHCP configuration file.

- **TOKEN_SECRET:**  
  A secret token used to secure all API endpoints. Ensure you choose a strong token.

- **INTERFACES_CONF_PATH:**  
  The file path for the network interfaces configuration.

## Updating Configuration

If your configuration files are located in different paths, update the values in your `.env` file accordingly. This approach makes it easy to adapt the API to various environments without changing the code.
