# Copyright (c) 2026 Ravi Sharma

"""
Code Coverage Verification Engine.
Parses Go and Python test coverage outputs and enforces configured minimum line coverage thresholds.
"""

import os
import re
import sys
import subprocess
import argparse


def check_go_coverage(service_path: str, min_threshold: float = 80.0) -> bool:
    """Run go test with coverage profile and check statement/line coverage percentage."""
    print(f"\n[Go Coverage] Checking {service_path}...")
    
    # Run go test with coverage across internal packages
    cmd = ["go", "test", "-coverprofile=coverage.out", "-covermode=atomic", "./internal/..."]
    result = subprocess.run(cmd, cwd=service_path, capture_output=True, text=True)
    if result.returncode != 0:
        print(f"Go tests failed in {service_path}:\n{result.stderr}\n{result.stdout}")
        return False

    cov_file = "coverage.out"
    full_cov_path = os.path.join(service_path, cov_file)
    if not os.path.exists(full_cov_path):
        print(f"Coverage profile not found at {full_cov_path}")
        return False

    # Run go tool cover -func=coverage.out to get coverage summary
    func_cmd = ["go", "tool", "cover", f"-func={cov_file}"]
    func_result = subprocess.run(func_cmd, cwd=service_path, capture_output=True, text=True)
    print(func_result.stdout)

    # Parse total line: "total:	(statements)	88.5%"
    match = re.search(r"total:\s+\(statements\)\s+([0-9.]+)%", func_result.stdout)
    if not match:
        print(f"Could not parse total coverage from output:\n{func_result.stdout}")
        return False

    total_cov = float(match.group(1))
    print(f"Total coverage for {service_path}: {total_cov:.1f}% (required: {min_threshold:.1f}%)")
    if total_cov < min_threshold:
        print(f"FAIL: Coverage {total_cov:.1f}% is below required {min_threshold:.1f}%")
        return False

    print(f"PASS: {service_path} meets coverage threshold.")
    return True


def check_python_coverage(service_path: str, min_threshold: float = 80.0) -> bool:
    """Run pytest with coverage and verify >= min_threshold."""
    print(f"\n[Python Coverage] Checking {service_path}...")
    cmd = [
        sys.executable,
        "-m",
        "pytest",
        "tests/unit",
        "--cov=app",
        f"--cov-fail-under={int(min_threshold)}",
        "--cov-report=term-missing",
    ]
    result = subprocess.run(cmd, cwd=service_path, capture_output=True, text=True)
    print(result.stdout)
    if result.stderr:
        print(result.stderr)

    if result.returncode != 0:
        print(f"FAIL: Python coverage or tests failed in {service_path}")
        return False

    print(f"PASS: {service_path} meets coverage threshold.")
    return True


def main():
    parser = argparse.ArgumentParser(description="Enforce minimum code coverage thresholds")
    parser.add_argument("--threshold", type=float, default=80.0, help="Minimum coverage percentage (default: 80.0)")
    args = parser.parse_args()

    threshold = args.threshold
    print(f"=== Code Coverage Enforcement Engine (Threshold: {threshold:.1f}%) ===")

    success = True

    # 1. Routing Service (Go)
    if not check_go_coverage("services/routing-service", threshold):
        success = False

    # 2. Tile Service (Go)
    if not check_go_coverage("services/tile-service", threshold):
        success = False

    # 3. Geocoding Service (Python)
    if not check_python_coverage("services/geocoding-service", threshold):
        success = False

    print("\n========================================================")
    if success:
        print(f"SUCCESS: All services meet the {threshold:.1f}% code coverage requirement!")
        sys.exit(0)
    else:
        print(f"FAILURE: One or more services failed code coverage threshold ({threshold:.1f}%).")
        sys.exit(1)


if __name__ == "__main__":
    main()
