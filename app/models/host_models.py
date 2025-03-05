from pydantic import BaseModel

class Host(BaseModel):
    """Model for a DHCP host"""
    name: str
    hardware_ethernet: str
    option_routers: str
    option_subnet_mask: str
    fixed_address: str
    option_domain_name_servers: str