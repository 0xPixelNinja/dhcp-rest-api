# app/routers/hosts.py
from fastapi import APIRouter, Depends, HTTPException
from app.auth import verify_token
from app.services import dhcp_manager
from pydantic import BaseModel

router = APIRouter()

class Host(BaseModel):
    """Model for a DHCP host"""
    name: str
    hardware_ethernet: str
    option_routers: str
    option_subnet_mask: str
    fixed_address: str
    option_domain_name_servers: str

@router.get("/")
def list_hosts(token: str = Depends(verify_token)):
    """List all DHCP hosts"""
    hosts = dhcp_manager.list_hosts()
    return {"hosts": hosts}

@router.post("/")
def add_host(host: Host, token: str = Depends(verify_token)):
    """Add a new DHCP host"""
    if dhcp_manager.add_host(host):
        return {"message": "Host added successfully"}
    else:
        raise HTTPException(status_code=400, detail="Failed to add host")
