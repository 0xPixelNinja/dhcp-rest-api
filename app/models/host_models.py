from pydantic import BaseModel
from typing import Optional


class Host(BaseModel):
    """Model for a DHCP host"""

    name: str
    hardware_ethernet: str
    option_routers: str
    option_subnet_mask: str
    fixed_address: str
    option_domain_name_servers: str


class HostUpdate(BaseModel):
    """Model for updating a DHCP host Only the fields that need to be updated are required"""

    name: Optional[str] = None
    hardware_ethernet: Optional[str] = None
    option_routers: Optional[str] = None
    option_subnet_mask: Optional[str] = None
    fixed_address: Optional[str] = None
    option_domain_name_servers: Optional[str] = None
