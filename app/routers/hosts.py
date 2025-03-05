# app/routers/hosts.py
from fastapi import APIRouter, Depends, HTTPException
from app.auth import verify_token
from app.services import dhcp_manager
from app.models.host_models import Host, HostUpdate

router = APIRouter()


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


@router.put("/{name}")
def update_host(name: str, host_update: HostUpdate, token: str = Depends(verify_token)):
    """Update an existing DHCP host"""
    if dhcp_manager.update_host(name, host_update.model_dump(exclude_unset=True)):
        return {"message": "Host updated successfully"}
    else:
        raise HTTPException(status_code=400, detail="Failed to update host")


@router.delete("/{name}")
def delete_host(name: str, token: str = Depends(verify_token)):
    """Delete an existing DHCP host"""
    if dhcp_manager.delete_host(name):
        return {"message": "Host deleted successfully"}
    else:
        raise HTTPException(status_code=400, detail="Failed to delete host")
