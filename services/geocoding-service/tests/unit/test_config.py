"""Copyright (c) 2026 Ravi Sharma."""

import os
import pytest
from app.config import Settings


@pytest.mark.unit
class TestConfig:
    def test_default_settings(self, monkeypatch):
        monkeypatch.delenv("ELASTICSEARCH_URL", raising=False)
        monkeypatch.delenv("PORT", raising=False)
        monkeypatch.delenv("HOST", raising=False)

        settings = Settings()
        assert settings.app_name == "geocoding-service"
        assert settings.elasticsearch_url == "http://localhost:9200"
        assert settings.port == 8000
        assert settings.host == "0.0.0.0"

    def test_custom_environment_settings(self, monkeypatch):
        monkeypatch.setenv("ELASTICSEARCH_URL", "http://es-cluster:9200")
        monkeypatch.setenv("PORT", "8005")
        monkeypatch.setenv("HOST", "127.0.0.1")

        settings = Settings(
            elasticsearch_url=os.getenv("ELASTICSEARCH_URL"),
            port=int(os.getenv("PORT", "8000")),
            host=os.getenv("HOST", "0.0.0.0"),
        )
        assert settings.elasticsearch_url == "http://es-cluster:9200"
        assert settings.port == 8005
        assert settings.host == "127.0.0.1"
