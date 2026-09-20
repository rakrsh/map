"""Copyright (c) 2026 Ravi Sharma."""

import os
from pydantic import BaseModel


class Settings(BaseModel):
    app_name: str = "geocoding-service"
    elasticsearch_url: str = os.getenv("ELASTICSEARCH_URL", "http://localhost:9200")
    host: str = os.getenv("HOST", "0.0.0.0")
    port: int = int(os.getenv("PORT", "8000"))


settings = Settings()
