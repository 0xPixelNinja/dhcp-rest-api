# app/routers/interfaces.py
from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel
from app.auth import verify_token
from app.services import interfaces_manager

router = APIRouter()


class InterfaceOperation(BaseModel):
    type: str  # "v4" or "v6"
    interface: str


@router.get("/")
def list_interfaces(token: str = Depends(verify_token)):
    interfaces = interfaces_manager.get_interfaces()
    return {"interfaces": interfaces}


@router.post("/")
def add_interface(op: InterfaceOperation, token: str = Depends(verify_token)):
    if interfaces_manager.add_interface(op.type, op.interface):
        return {
            "message": f"Interface {op.interface} added to INTERFACES{op.type.upper()} successfully."
        }
    else:
        raise HTTPException(status_code=400, detail="Failed to add interface.")


@router.delete("/")
def delete_interface(op: InterfaceOperation, token: str = Depends(verify_token)):
    if interfaces_manager.delete_interface(op.type, op.interface):
        return {
            "message": f"Interface {op.interface} removed from INTERFACES{op.type.upper()} successfully."
        }
    else:
        raise HTTPException(status_code=400, detail="Failed to delete interface.")
