-- Rollback for V002__postgis_schema.up.sql
DROP INDEX IF EXISTS idx_pois_h3;
DROP INDEX IF EXISTS idx_pois_geom;
DROP TABLE IF EXISTS pois;

DROP INDEX IF EXISTS idx_buildings_h3;
DROP INDEX IF EXISTS idx_buildings_geom;
DROP TABLE IF EXISTS buildings;

DROP INDEX IF EXISTS idx_water_bodies_h3;
DROP INDEX IF EXISTS idx_water_bodies_geom;
DROP TABLE IF EXISTS water_bodies;

DROP INDEX IF EXISTS idx_boundaries_h3;
DROP INDEX IF EXISTS idx_boundaries_geom;
DROP TABLE IF EXISTS boundaries;

DROP INDEX IF EXISTS idx_roads_h3;
DROP INDEX IF EXISTS idx_roads_geom;
DROP TABLE IF EXISTS roads;
