# Copyright (c) 2026 Ravi Sharma

"""
E2E Test Runner with Playwright & Artifact Logging.
Spins up lightweight service runners if not already up, executes Playwright tests,
captures traces, screenshots, video recordings, and tears down cleanly.
"""

import os
import sys
import shutil
import subprocess
import argparse
import time
import urllib.request


def check_endpoint(url: str, timeout: int = 2) -> bool:
    try:
        with urllib.request.urlopen(url, timeout=timeout) as response:
            return response.status == 200
    except Exception:
        return False


def wait_for_services(timeout_sec: int = 15) -> bool:
    services = [
        ("Geocoding", os.getenv("GEOCODING_URL", "http://127.0.0.1:8000") + "/health"),
        ("Routing", os.getenv("ROUTING_URL", "http://127.0.0.1:8081") + "/health"),
        ("Tile", os.getenv("TILE_URL", "http://127.0.0.1:8082") + "/health"),
    ]
    start_time = time.time()
    while time.time() - start_time < timeout_sec:
        if all(check_endpoint(url) for _, url in services):
            return True
        time.sleep(1)
    return False


def main():
    parser = argparse.ArgumentParser(description="Run Playwright E2E tests with artifact capture")
    parser.add_argument("--headed", action="store_true", help="Run tests in headed browser mode")
    parser.add_argument("--external-stack", action="store_true", help="Assume services are already running externally")
    args = parser.parse_args()

    e2e_dir = os.path.join(os.getcwd(), "tests", "e2e")
    spawned_procs = []
    npx_bin = shutil.which("npx") or "npx"

    returncode = 1
    try:
        if not args.external_stack and not wait_for_services(timeout_sec=2):
            print("[E2E] Starting local microservices for E2E testing...")

            # 1. Start Geocoding Service
            geo_env = os.environ.copy()
            geo_proc = subprocess.Popen(
                [sys.executable, "-m", "uvicorn", "app.main:app", "--host", "127.0.0.1", "--port", "8000"],
                cwd="services/geocoding-service",
                env=geo_env,
                stdout=subprocess.DEVNULL,
                stderr=subprocess.DEVNULL,
            )
            spawned_procs.append(("Geocoding", geo_proc))

            # 2. Start Routing Service
            route_env = os.environ.copy()
            route_env["PORT"] = "8081"
            route_proc = subprocess.Popen(
                ["go", "run", "./cmd/server"],
                cwd="services/routing-service",
                env=route_env,
                stdout=subprocess.DEVNULL,
                stderr=subprocess.DEVNULL,
            )
            spawned_procs.append(("Routing", route_proc))

            # 3. Start Tile Service
            tile_env = os.environ.copy()
            tile_env["PORT"] = "8082"
            tile_proc = subprocess.Popen(
                ["go", "run", "./cmd/server"],
                cwd="services/tile-service",
                env=tile_env,
                stdout=subprocess.DEVNULL,
                stderr=subprocess.DEVNULL,
            )
            spawned_procs.append(("Tile", tile_proc))

            print("[E2E] Waiting for microservices to become ready...")
            if not wait_for_services(timeout_sec=20):
                print("[E2E] Error: Microservices failed to start within timeout.", file=sys.stderr)
                sys.exit(1)
            print("[E2E] Microservices are ready.")

        # Run Playwright
        cmd = [npx_bin, "playwright", "test"]
        if args.headed:
            cmd.append("--headed")

        print(f"[E2E] Executing: {' '.join(cmd)}")
        env = os.environ.copy()
        env["GEOCODING_URL"] = "http://127.0.0.1:8000"
        env["ROUTING_URL"] = "http://127.0.0.1:8081"
        env["TILE_URL"] = "http://127.0.0.1:8082"

        res = subprocess.run(cmd, cwd=e2e_dir, env=env, text=True, shell=(os.name == "nt"))
        returncode = res.returncode
    finally:
        for name, proc in spawned_procs:
            print(f"[E2E] Stopping {name} service process (pid: {proc.pid})...")
            proc.terminate()
            try:
                proc.wait(timeout=3)
            except subprocess.TimeoutExpired:
                proc.kill()

    sys.exit(returncode)


if __name__ == "__main__":
    main()
