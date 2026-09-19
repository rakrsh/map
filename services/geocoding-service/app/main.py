from fastapi import FastAPI

from app.api.v1.endpoints import router as v1_router

app = FastAPI(title="Geocoding Service")
app.include_router(v1_router, prefix="/api/v1")


@app.get("/health")
def health_check():
    return {"status": "ok"}
