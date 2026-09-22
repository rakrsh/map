"""Copyright (c) 2026 Ravi Sharma."""

import os
import pytest
from elasticsearch import Elasticsearch


@pytest.mark.integration
class TestElasticsearchIntegration:
    @pytest.fixture
    def es_client(self):
        es_url = os.getenv("TEST_ELASTICSEARCH_URL", "http://localhost:9201")
        client = Elasticsearch(es_url, request_timeout=5)
        return client

    def test_cluster_health_and_info(self, es_client: Elasticsearch):
        info = es_client.info()
        assert "version" in info
        assert "name" in info

        health = es_client.cluster.health()
        assert health["status"] in ["green", "yellow"]
