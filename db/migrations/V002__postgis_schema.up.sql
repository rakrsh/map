-- PostGIS schema: roads, boundaries, water_bodies, buildings, pois
-- Add GiST spatial indexes and optional H3 index columns for aggregation
-- Copyright (c) 2026 Ravi Sharma

-- Roads (lines)
CREATE TABLE IF NOT EXISTS roads (
    id BIGSERIAL PRIMARY KEY,
    osm_id BIGINT,
    name TEXT,
    classification TEXT,
    geom geometry(LineString,4326),
    h3_index BIGINT,
    created_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_roads_geom ON roads USING GIST (geom);
CREATE INDEX IF NOT EXISTS idx_roads_h3 ON roads USING BTREE (h3_index);

-- Boundaries (polygons)
CREATE TABLE IF NOT EXISTS boundaries (
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    kind TEXT,
    geom geometry(Polygon,4326),
    h3_index BIGINT,
    created_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_boundaries_geom ON boundaries USING GIST (geom);
CREATE INDEX IF NOT EXISTS idx_boundaries_h3 ON boundaries USING BTREE (h3_index);

-- Water bodies (polygons)
CREATE TABLE IF NOT EXISTS water_bodies (
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    kind TEXT,
    geom geometry(Polygon,4326),
    h3_index BIGINT,
    created_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_water_bodies_geom ON water_bodies USING GIST (geom);
CREATE INDEX IF NOT EXISTS idx_water_bodies_h3 ON water_bodies USING BTREE (h3_index);

-- Buildings (polygons)
CREATE TABLE IF NOT EXISTS buildings (
    id BIGSERIAL PRIMARY KEY,
    osm_id BIGINT,
    name TEXT,
    height DOUBLE PRECISION,
    geom geometry(Polygon,4326),
    h3_index BIGINT,
    created_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_buildings_geom ON buildings USING GIST (geom);
CREATE INDEX IF NOT EXISTS idx_buildings_h3 ON buildings USING BTREE (h3_index);

-- POIs (points)
CREATE TABLE IF NOT EXISTS pois (
    id BIGSERIAL PRIMARY KEY,
    osm_id BIGINT,
    name TEXT,
    category TEXT,
    geom geometry(Point,4326),
    h3_index BIGINT,
    created_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_pois_geom ON pois USING GIST (geom);
CREATE INDEX IF NOT EXISTS idx_pois_h3 ON pois USING BTREE (h3_index);

-- NOTE: h3_index is a helpful denormalized bigint (H3 base cell address) used for fast aggregation.
-- Populate h3_index using an external tool or extension (postgis-h3 / h3-pg) when available.
