"""Copyright (c) 2026 Ravi Sharma."""

import pytest
from fastapi.testclient import TestClient


@pytest.mark.unit
class TestEndpoints:
    def test_health_check(self, test_client: TestClient):
        response = test_client.get("/health")
        assert response.status_code == 200
        assert response.json() == {"status": "ok"}

    def test_search_endpoint(self, test_client: TestClient):
        response = test_client.get("/api/v1/search?query=Downtown")
        assert response.status_code == 200
        payload = response.json()
        assert payload["query"] == "Downtown"
        assert payload["results"] == []
        assert payload["status"] == "pending_implementation"

    def test_search_endpoint_missing_query(self, test_client: TestClient):
        response = test_client.get("/api/v1/search")
        assert response.status_code == 422

    def test_reverse_endpoint(self, test_client: TestClient):
        response = test_client.get("/api/v1/reverse?lat=37.7749&lon=-122.4194")
        assert response.status_code == 200
        payload = response.json()
        assert payload["latitude"] == 37.7749
        assert payload["longitude"] == -122.4194
        assert payload["status"] == "pending_implementation"

    def test_reverse_endpoint_invalid_parameters(self, test_client: TestClient):
        response = test_client.get("/api/v1/reverse?lat=invalid&lon=-122.4194")
        assert response.status_code == 422
