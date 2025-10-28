from fastapi import FastAPI, Header, HTTPException
from pydantic import BaseModel
from datetime import datetime
import jwt
import ssl
import uvicorn

SECRET_KEY = "super_secret_key"
ALGORITHM = "HS256"

app = FastAPI()
latest_metrics = None


class Metrics(BaseModel):
    cpu_usage: float
    ram_usage: float
    timestamp: str


@app.post("/metrics")
async def receive_metrics(data: Metrics, authorization: str = Header(None)):
    if not authorization or not authorization.startswith("Bearer "):
        raise HTTPException(status_code=403, detail="Missing or invalid token")

    token = authorization.split("Bearer ")[1]

    try:
        jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
    except Exception as e:
        raise HTTPException(status_code=403, detail=f"Invalid token: {e}")

    global latest_metrics
    latest_metrics = data
    return {"message": "Metrics received", "data": data}


@app.get("/metrics")
async def get_metrics():
    if latest_metrics:
        return latest_metrics
    return {"message": "No metrics received yet."}


if __name__ == "__main__":
    context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
    context.load_cert_chain("cert.pem", "key.pem")
    uvicorn.run(app, host="0.0.0.0", port=8443, ssl_keyfile="key.pem", ssl_certfile="cert.pem")
