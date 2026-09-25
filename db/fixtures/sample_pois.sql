-- Sample POIs fixture for local benchmarks and CI
-- Insert a small set of points in/near Manhattan for reproducible bbox tests
INSERT INTO pois (osm_id, name, category, geom) VALUES
(1, 'Sample Cafe', 'cafe', ST_SetSRID(ST_MakePoint(-73.9857,40.7484),4326)),
(2, 'Sample Museum', 'museum', ST_SetSRID(ST_MakePoint(-73.9772,40.7527),4326)),
(3, 'Sample Park', 'park', ST_SetSRID(ST_MakePoint(-73.9712,40.7648),4326)),
(4, 'Sample Library', 'library', ST_SetSRID(ST_MakePoint(-73.9818,40.7532),4326));
