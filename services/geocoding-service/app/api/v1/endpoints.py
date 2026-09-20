"""Copyright (c) 2026 Ravi Sharma."""

from fastapi import APIRouter

router = APIRouter()


@router.get("/search")
def search_address(query: str):
    return {"query": query, "results": [], "status": "pending_implementation"}


@router.get("/reverse")
def reverse_geocode(lat: float, lon: float):
    return {"latitude": lat, "longitude": lon, "address": None, "status": "pending_implementation"}
