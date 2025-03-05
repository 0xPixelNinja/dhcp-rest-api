# DHCP REST API

This project provides a RESTful API interface for managing ISC DHCP server configurations, allowing you to easily create, update, delete, and list DHCP hosts and network interfaces.

## Quick Links

- [Configuration](CONFIGURATION.md) – Setup and environment details.
- [Usage](USAGE.md) – How to run and interact with the API.
- [Contributing](CONTRIBUTING.md) – Guidelines for contributing to the project.
- [Changelog](CHANGELOG.md) – Record of changes and updates.

## Overview

This project is built with Python and FastAPI. It secures endpoints with token-based authentication and wraps two primary configuration files:
- **DHCP Configuration:** Typically located at `/etc/dhcp/dhcpd.conf`.
- **Interface Configuration:** By default, this is located at `/etc/default/isc-dhcp-server`.

I made this project, for automating the DHCP configurations remotely via API, without accessing the PVE shell