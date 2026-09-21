# Copyright (c) 2026 Ravi Sharma

"""
E2E Test Runner with Playwright & Artifact Logging.
Executes end-to-end multi-service test suites with traces, screenshots, and video recordings.
"""

import os
import sys
import subprocess
import argparse
import time
import urllib.request


def check_endpoint(url: str, timeout: int = 5) -> bool:
    try:
        with urllib.request.urlopen(url, timeout=timeout) as response:
            return response.status == 200
    except Exception:
        return False


def wait_for_services(timeout_sec: int = 30):
    services = [
        ("Geocoding", os.getenv("GEOCODING_URL", "http://localhost:8000") + "/health"),
        ("Routing", os.getenv("ROUTING_URL", "http://localhost:8081") + "/health"),
        ("Tile", os.getenv("TILE_URL", "http://localhost:8082") + "/health"),
    ]
    print("[E2E] Checking backend services readiness...")
    start_time = time.time()
    while time.time() - start_time < timeout_sec:
        all_up = True
        for name, url in services:
            if not check_endpoint(url):
                all_up = False
                break
        if all_up:
            print("[E2E] All backend services are healthy and reachable!")
            return True
        time.sleep(2)
    
    print("[E2E] Warning: Backend services did not respond within timeout. Running browser-isolated E2E tests...")
    return False


def main():
    parser = argparse.ArgumentParser(description="Run Playwright E2E tests with artifact capture")
    parser.add_argument("--headed", action="store_true", help="Run tests in headed browser mode")
    args = parser.parse_args()

    e2e_dir = os.path.join(os.getcwd(), "tests", "e2e")
    if not os.path.exists(os.path.join(e2e_dir, "node_modules")):
        print("[E2E] Installing Playwright dependencies...")
        subprocess.run(["npm", "ci"], cwd=e2e_dir, check=False)
        subprocess.run(["npx", "playwright", "install", "--with-deps", "chromium"], cwd=e2e_dir, check=False)

    wait_for_services(timeout_sec=5)

    cmd = ["npx", "playwright", "test"]
    if args.headed:
        cmd.append("--headed")

    print(f"[E2E] Executing: {' '.join(cmd)}")
    res = subprocess.run(cmd, cwd=e2e_dir, text=True)
    sys.exit(res.returncode)


if __name__ == "__main__":
    main()
