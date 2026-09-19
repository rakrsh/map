class GeocoderService:
    def __init__(self):
        pass

    def search(self, query: str):
        return {
            "query": query,
            "results": [],
            "status": "pending_implementation",
        }

    def reverse_geocode(self, latitude: float, longitude: float):
        return {
            "latitude": latitude,
            "longitude": longitude,
            "address": None,
            "status": "pending_implementation",
        }
