#!/usr/bin/env python3
"""
Simple bounding-box benchmark for PostGIS schema created in db/migrations/V002__postgis_schema.up.sql

Usage:
  python scripts/benchmark_postgis_schema.py --bbox MINX MINY MAXX MAXY --runs 10

Requires: psycopg2 installed in the environment and a reachable DATABASE_URL env var
"""
import os
import time
import argparse
import psycopg2


def run_bbox(conn, bbox, limit=1000):
    minx, miny, maxx, maxy = bbox
    sql = (
        "SELECT id FROM pois "
        "WHERE geom && ST_MakeEnvelope(%s,%s,%s,%s,4326) "
        "LIMIT %s"
    )
    with conn.cursor() as cur:
        cur.execute(sql, (minx, miny, maxx, maxy, limit))
        return cur.fetchall()


def main():
    p = argparse.ArgumentParser()
    p.add_argument('--bbox', nargs=4, type=float, required=True, help='MINX MINY MAXX MAXY')
    p.add_argument('--runs', type=int, default=5)
    p.add_argument('--limit', type=int, default=1000)
    args = p.parse_args()

    dsn = os.environ.get('DATABASE_URL', 'postgres://postgres:postgres@localhost:5432/postgres')
    conn = psycopg2.connect(dsn)

    times = []
    for i in range(args.runs):
        t0 = time.time()
        rows = run_bbox(conn, args.bbox, limit=args.limit)
        t1 = time.time()
        elapsed = (t1 - t0) * 1000.0
        times.append(elapsed)
        print(f'run={i+1} rows={len(rows)} {elapsed:.2f}ms')

    print(f'avg {sum(times)/len(times):.2f}ms min {min(times):.2f}ms max {max(times):.2f}ms')
    conn.close()


if __name__ == '__main__':
    main()
