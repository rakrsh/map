#!/usr/bin/env python3
"""H3 indexing adapter for the shared point workload."""

from __future__ import annotations

import os
import time

try:
    import h3
except ImportError as error:
    raise SystemExit("Install the H3 adapter with: pip install h3") from error

scale = int(os.getenv("BENCHMARK_SCALE", "100000"))
start = time.perf_counter()
indexes = [h3.latlng_to_cell((index % 180) - 90, (index % 360) - 180, 9) for index in range(scale)]
elapsed = time.perf_counter() - start
print("H3 benchmark")
print(f"scale={scale}")
print(f"unique_cells={len(set(indexes))}")
print(f"elapsed_seconds={elapsed:.6f}")
print(f"cells_per_second={scale / elapsed:.2f}")
print("H3 measures hierarchical indexing, not transactional storage or SQL spatial joins.")
