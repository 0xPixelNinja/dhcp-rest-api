# DHCP REST API

A lightweight REST API for managing ISC DHCP Server configuration, specifically designed for Proxmox environments.

## Background

This project started as a personal solution to a frustrating problem I encountered while setting up virtual machines in my Proxmox homelab. Managing DHCP reservations for VMs was becoming increasingly tedious - I had to SSH into my DHCP server, manually edit configuration files, and restart services every time I wanted to assign a static IP to a new VM.

After dealing with this workflow for months, I decided to build a simple REST API that could automate DHCP host management. What began as a weekend project to solve my own problem has evolved into a tool that I hope others will find useful for their Proxmox and virtualization setups.

## Features

- **Host Management**: Add, update, delete, and list DHCP host reservations
- **Interface Management**: Configure network interfaces for DHCP service
- **Security**: Token-based authentication with configurable security headers
- **Rate Limiting**: Protection against abuse with configurable rate limits
- **Logging**: Comprehensive logging for debugging and monitoring
- **Lightweight**: Single binary deployment with minimal dependencies

## Why This Project?

Setting up DHCP for Proxmox VMs typically involves:
1. Manually editing `/etc/dhcp/dhcpd.conf`
2. Restarting the DHCP service
3. Hoping you didn't make a syntax error
4. Repeating this process for every VM

With this REST API, you can:
- Programmatically manage DHCP reservations
- Integrate with VM provisioning workflows
- Avoid manual file editing and potential syntax errors
- Automate IP assignment for new VMs

## Installation

### Prerequisites

- Ubuntu/Debian Linux system
- Root access
- ISC DHCP Server (`isc-dhcp-server` package)

### Manual Installation

#### Step 1: Install ISC DHCP Server

```bash
sudo apt update
sudo apt install isc-dhcp-server
```

#### Step 2: Create Basic DHCP Configuration

Create or update `/etc/dhcp/dhcpd.conf`:

```bash
sudo nano /etc/dhcp/dhcpd.conf
```

Add basic configuration:

```
# Example subnet (adjust for your network)
subnet 0.0.0.0 netmask 0.0.0.0 {
        deny-unknown-clients;
        authoritative;
        default-lease-time 21600000;
        max-lease-time 432000000;
}

# Host declarations will be managed by the REST API
```

#### Step 3: Configure Network Interfaces

Edit `/etc/default/isc-dhcp-server`:

```bash
sudo nano /etc/default/isc-dhcp-server
```

Set the interface for DHCP service:

```
# Replace with your network interface
INTERFACESv4="vmbr0"
INTERFACESv6=""
```

#### Step 4: Download and Install the Binary

```bash
# Download the latest binary
sudo wget -O /usr/local/bin/dhcp-rest-api https://raw.githubusercontent.com/0xPixelNinja/dhcp-rest-api/refs/heads/main/bin/dhcp-rest-api-linux

# Make it executable
sudo chmod +x /usr/local/bin/dhcp-rest-api
```

#### Step 5: Create Configuration Directory and Token

```bash
# Create configuration directory
sudo mkdir -p /etc/dhcp-rest-api

# Generate a secure token (or use your own)
echo "your-secure-token-here" | sudo tee /etc/dhcp-rest-api/token

# Secure the token file
sudo chmod 600 /etc/dhcp-rest-api/token
```

#### Step 6: Create Systemd Service

Create `/etc/systemd/system/dhcp-rest-api.service`:

```bash
sudo nano /etc/systemd/system/dhcp-rest-api.service
```

Add the following content:

```ini
[Unit]
Description=DHCP REST API Service
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/dhcp-rest-api
Environment=TOKEN_FILE_PATH=/etc/dhcp-rest-api/token
Environment=PORT=8080
Environment=DHCP_CONF_PATH=/etc/dhcp/dhcpd.conf
Environment=INTERFACES_CONF_PATH=/etc/default/isc-dhcp-server
Environment=ENVIRONMENT=production
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

#### Step 7: Start the Service

```bash
# Reload systemd and start the service
sudo systemctl daemon-reload
sudo systemctl enable dhcp-rest-api
sudo systemctl start dhcp-rest-api

# Check service status
sudo systemctl status dhcp-rest-api
```

#### Step 8: Test the Installation

```bash
# Test health endpoint
curl http://localhost:8080/health

# Test authenticated endpoint (replace YOUR_TOKEN with your token)
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/hosts/
```

## Integration with Proxmox

This API works perfectly with Proxmox VM provisioning workflows. You can:

1. Create a VM in Proxmox
2. Get the VM's MAC address
3. Use this API to create a DHCP reservation
4. Start the VM with a guaranteed IP address

## Troubleshooting

### Service Issues

```bash
# Check service status
sudo systemctl status dhcp-rest-api

# View service logs
sudo journalctl -u dhcp-rest-api -f

# Check DHCP server status
sudo systemctl status isc-dhcp-server
```

## Contributing

This project is open to contributions! Feel free to:
- Report bugs
- Suggest features
- Submit pull requests
- Improve documentation

## Support

If you find this project helpful, please consider starring it on GitHub. For issues or questions, please open a GitHub issue.

## License

This project is licensed under the MIT License - see the LICENSE file for details.