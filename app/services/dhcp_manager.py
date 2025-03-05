# app/services/dhcp_manager.py
import os
import re
from typing import List, Dict
from dotenv import load_dotenv

load_dotenv()

DHCP_CONF_PATH = os.getenv("DHCP_CONF_PATH", "/etc/dhcp/dhcpd.conf")

def list_hosts() -> List[Dict]:
    """List all DHCP hosts"""
    hosts = []
    try:
        with open(DHCP_CONF_PATH, "r") as f:
            content = f.read()
    except Exception as e:
        print(f"Error reading file: {e}")
        return hosts

    host_blocks = re.findall(r'host\s+(\S+)\s*{([^}]+)}', content, re.MULTILINE)
    for name, block in host_blocks:
        host_data = {"name": name}
        for line in block.splitlines():
            line = line.strip().strip(';')
            if line.startswith("hardware ethernet"):
                host_data["hardware_ethernet"] = line.split(" ", 2)[2]
            elif line.startswith("option routers"):
                host_data["option_routers"] = line.split(" ", 2)[2]
            elif line.startswith("option subnet-mask"):
                host_data["option_subnet_mask"] = line.split(" ", 2)[2]
            elif line.startswith("fixed-address"):
                host_data["fixed_address"] = line.split(" ", 1)[1]
            elif line.startswith("option domain-name-servers"):
                host_data["option_domain_name_servers"] = line.split(" ", 2)[2]
        hosts.append(host_data)
    return hosts

def add_host(host) -> bool:
    """Add a new DHCP host"""
    host_block = f"""
host {host.name} {{
    hardware ethernet {host.hardware_ethernet};
    option routers {host.option_routers};
    option subnet-mask {host.option_subnet_mask};
    fixed-address {host.fixed_address};
    option domain-name-servers {host.option_domain_name_servers};
}}
"""
    try:
        with open(DHCP_CONF_PATH, "a") as f:
            f.write(host_block)
        return True
    except Exception as e:
        print(f"Error writing to file: {e}")
        return False
