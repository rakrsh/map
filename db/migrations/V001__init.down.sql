-- Rollback for V001__init.sql
DROP INDEX IF EXISTS idx_places_name_trgm;
DROP INDEX IF EXISTS idx_places_geom;
DROP TABLE IF EXISTS places;

-- Note: we intentionally do not drop PostGIS extensions here, they may be shared.
