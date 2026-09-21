"""Copyright (c) 2026 Ravi Sharma."""

import pytest
from unittest.mock import MagicMock
from app.services.geocoder import GeocoderService


@pytest.mark.unit
class TestGeocoderService:
    def test_search_returns_expected_structure(self, geocoder_service: GeocoderService):
        query = "100 Market St, San Francisco"
        result = geocoder_service.search(query)

        assert result["query"] == query
        assert isinstance(result["results"], list)
        assert result["status"] == "pending_implementation"

    def test_reverse_geocode_returns_expected_structure(self, geocoder_service: GeocoderService):
        lat, lon = 37.7749, -122.4194
        result = geocoder_service.reverse_geocode(lat, lon)

        assert result["latitude"] == lat
        assert result["longitude"] == lon
        assert result["address"] is None
        assert result["status"] == "pending_implementation"

    def test_search_with_mocked_elasticsearch_client(self, monkeypatch):
        mock_es = MagicMock()
        mock_es.search.return_value = {
            "hits": {
                "hits": [
                    {"_source": {"name": "Test Place", "lat": 37.77, "lon": -122.41}}
                ]
            }
        }
        service = GeocoderService()
        # Verify isolation with mock double attached
        service.es_client = mock_es
        result = service.search("Test Place")
        assert result["query"] == "Test Place"
