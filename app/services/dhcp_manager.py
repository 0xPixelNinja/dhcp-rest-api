# app/services/dhcp_manager.py
import os
import re
from typing import List, Dict
from dotenv import load_dotenv
from app.models.host_models import Host, HostUpdate

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

    host_blocks = re.findall(r"host\s+(\S+)\s*{([^}]+)}", content, re.MULTILINE)
    for name, block in host_blocks:
        host_data = {"name": name}
        for line in block.splitlines():
            line = line.strip().strip(";")
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


def update_host(name: str, updates: dict) -> bool:
    try:
        with open(DHCP_CONF_PATH, "r") as f:
            content = f.read()

        pattern = rf"(host\s+{name}\s*{{)([^}}]+)(}})"
        match = re.search(pattern, content, re.MULTILINE | re.DOTALL)
        if not match:
            return False

        _ = match.group(1)
        block_body = match.group(2)
        __ = match.group(3)

        current_fields = {}
        for line in block_body.splitlines():
            line = line.strip().strip(";")
            if not line:
                continue
            if line.startswith("hardware ethernet"):
                current_fields["hardware_ethernet"] = line.split(" ", 2)[2]
            elif line.startswith("option routers"):
                current_fields["option_routers"] = line.split(" ", 2)[2]
            elif line.startswith("option subnet-mask"):
                current_fields["option_subnet_mask"] = line.split(" ", 2)[2]
            elif line.startswith("fixed-address"):
                current_fields["fixed_address"] = line.split(" ", 1)[1]
            elif line.startswith("option domain-name-servers"):
                current_fields["option_domain_name_servers"] = line.split(" ", 2)[2]

        current_fields.update({k: v for k, v in updates.items() if v is not None})

        new_name = updates.get("name", name)

        new_host_block = f"""
host {new_name} {{
    hardware ethernet {current_fields.get('hardware_ethernet', '')};
    option routers {current_fields.get('option_routers', '')};
    option subnet-mask {current_fields.get('option_subnet_mask', '')};
    fixed-address {current_fields.get('fixed_address', '')};
    option domain-name-servers {current_fields.get('option_domain_name_servers', '')};
}}
"""

        new_content = re.sub(
            pattern, new_host_block, content, flags=re.MULTILINE | re.DOTALL
        )
        with open(DHCP_CONF_PATH, "w") as f:
            f.write(new_content)
        return True

    except Exception as e:
        print(f"Error updating host: {e}")
        return False


def delete_host(name: str) -> bool:
    try:
        with open(DHCP_CONF_PATH, "r") as f:
            content = f.read()
        pattern = rf"\s*host\s+{name}\s*{{[^}}]+}}\s*"
        if not re.search(pattern, content, re.MULTILINE):
            return False
        new_content = re.sub(pattern, "", content, flags=re.MULTILINE)
        with open(DHCP_CONF_PATH, "w") as f:
            f.write(new_content)
        return True
    except Exception as e:
        print(f"Error deleting host: {e}")
        return False
