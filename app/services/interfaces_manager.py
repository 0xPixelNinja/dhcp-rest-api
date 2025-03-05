# app/services/interfaces_manager.py
import os
import re
from dotenv import load_dotenv

load_dotenv()

# Path to the interfaces configuration file.
# You can override this by setting INTERFACES_CONF_PATH in your .env file.
INTERFACES_CONF_PATH = os.getenv("INTERFACES_CONF_PATH", "/etc/default/isc-dhcp-server")


def get_interfaces() -> dict:
    """
    Read interfaces.conf and return a dict with the current interface settings.
    """
    try:
        with open(INTERFACES_CONF_PATH, "r") as f:
            content = f.read()
    except Exception as e:
        print(f"Error reading interfaces config: {e}")
        return {}

    # Extract the values using regex.
    result = {}
    m_v4 = re.search(r'INTERFACESv4="([^"]*)"', content)
    m_v6 = re.search(r'INTERFACESv6="([^"]*)"', content)
    result["v4"] = m_v4.group(1) if m_v4 else ""
    result["v6"] = m_v6.group(1) if m_v6 else ""
    return result


def save_interfaces(interfaces: dict) -> bool:
    """
    Rebuild the interfaces.conf file using the updated interfaces.
    Only the INTERFACESv4 and INTERFACESv6 lines are updated.
    """
    try:
        with open(INTERFACES_CONF_PATH, "r") as f:
            lines = f.readlines()
    except Exception as e:
        print(f"Error reading interfaces config: {e}")
        return False

    new_lines = []
    for line in lines:
        if line.startswith("INTERFACESv4="):
            new_lines.append(f'INTERFACESv4="{interfaces.get("v4", "")}"\n')
        elif line.startswith("INTERFACESv6="):
            new_lines.append(f'INTERFACESv6="{interfaces.get("v6", "")}"\n')
        else:
            new_lines.append(line)

    try:
        with open(INTERFACES_CONF_PATH, "w") as f:
            f.writelines(new_lines)
        return True
    except Exception as e:
        print(f"Error writing interfaces config: {e}")
        return False


def add_interface(interface_type: str, interface: str) -> bool:
    """
    Add an interface to the list (for v4 or v6) if not already present.
    """
    interfaces = get_interfaces()
    key = interface_type.lower()  # expects "v4" or "v6"
    current = interfaces.get(key, "")
    current_list = current.split() if current.strip() != "" else []
    if interface in current_list:
        return True  # already present
    current_list.append(interface)
    interfaces[key] = " ".join(current_list)
    return save_interfaces(interfaces)


def delete_interface(interface_type: str, interface: str) -> bool:
    """
    Remove an interface from the list (for v4 or v6) if present.
    """
    interfaces = get_interfaces()
    key = interface_type.lower()  # expects "v4" or "v6"
    current = interfaces.get(key, "")
    current_list = current.split() if current.strip() != "" else []
    if interface not in current_list:
        return True  # nothing to remove
    current_list.remove(interface)
    interfaces[key] = " ".join(current_list)
    return save_interfaces(interfaces)
