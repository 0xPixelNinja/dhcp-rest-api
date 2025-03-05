# app/auth.py
import os
from fastapi import HTTPException, Security, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from dotenv import load_dotenv

load_dotenv()

security = HTTPBearer()
TOKEN_SECRET = os.getenv("TOKEN_SECRET", "your-secret-token")

def verify_token(credentials: HTTPAuthorizationCredentials = Security(security)):
    """Verify the token"""
    token = credentials.credentials
    if token != TOKEN_SECRET:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Invalid or missing token."
        )
    return token
