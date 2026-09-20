#!/usr/bin/env python3
# Copyright (c) 2026 Ravi Sharma

"""Summarize pgbench latency logs in milliseconds."""

from __future__ import annotations

import statistics
import sys
from pathlib import Path


def percentile(values: list[float], rank: float) -> float:
    if not values:
        return 0.0
    index = (len(values) - 1) * rank
    lower = int(index)
    upper = min(lower + 1, len(values) - 1)
    weight = index - lower
    return values[lower] + (values[upper] - values[lower]) * weight


def main() -> int:
    if len(sys.argv) != 2:
        print("usage: benchmark_metrics.py PG_BENCH_LOG", file=sys.stderr)
        return 2

    values: list[float] = []
    for line in Path(sys.argv[1]).read_text(encoding="utf-8").splitlines():
        fields = line.split()
        if len(fields) < 3 or fields[0].startswith("#"):
            continue
        try:
            values.append(float(fields[2]) / 1000.0)
        except ValueError:
            continue

    values.sort()
    if not values:
        print("samples=0")
        return 1

    print(f"samples={len(values)}")
    print(f"average_ms={statistics.fmean(values):.3f}")
    print(f"p50_ms={percentile(values, 0.50):.3f}")
    print(f"p95_ms={percentile(values, 0.95):.3f}")
    print(f"p99_ms={percentile(values, 0.99):.3f}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
