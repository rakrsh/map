"""Copyright (c) 2026 Ravi Sharma."""

import pytest
from fastapi.testclient import TestClient

from app.main import app
from app.services.geocoder import GeocoderService


@pytest.fixture
def test_client():
    return TestClient(app)


@pytest.fixture
def geocoder_service():
    return GeocoderService()
