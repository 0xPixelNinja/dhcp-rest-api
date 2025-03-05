# app/main.py
from fastapi import FastAPI
from dotenv import load_dotenv

load_dotenv()

from app.routers import hosts, interfaces

app = FastAPI(title="DHCP REST API")
app.include_router(hosts.router, prefix="/hosts")
app.include_router(interfaces.router, prefix="/interfaces")
