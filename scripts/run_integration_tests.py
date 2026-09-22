# Copyright (c) 2026 Ravi Sharma

"""
Ephemeral Container Integration Test Orchestrator.
Spins up isolated test containers, polls healthchecks, runs integration test suites,
and guarantees clean teardown.
"""

import os
import sys
import time
import subprocess
import argparse


COMPOSE_FILE = "docker-compose.test.yml"


def run_cmd(cmd: list[str], cwd: str | None = None, check: bool = True) -> subprocess.CompletedProcess:
    print(f"[Exec] {' '.join(cmd)}")
    res = subprocess.run(cmd, cwd=cwd, text=True, capture_output=True)
    if res.stdout:
        print(res.stdout)
    if res.stderr:
        print(res.stderr, file=sys.stderr)
    if check and res.returncode != 0:
        raise subprocess.CalledProcessError(res.returncode, cmd)
    return res


def start_containers():
    print("\n[Integration] Starting ephemeral test containers...")
    run_cmd(["docker", "compose", "-f", COMPOSE_FILE, "up", "-d"])


def stop_containers():
    print("\n[Integration] Tearing down ephemeral test containers and volumes...")
    run_cmd(["docker", "compose", "-f", COMPOSE_FILE, "down", "-v", "--remove-orphans"], check=False)


def wait_for_health(timeout_sec: int = 60):
    print("\n[Integration] Waiting for ephemeral containers to become healthy...")
    start_time = time.time()
    
    while time.time() - start_time < timeout_sec:
        res = run_cmd(["docker", "compose", "-f", COMPOSE_FILE, "ps", "--format", "json"], check=False)
        output = res.stdout.strip()
        # In modern compose, check if any container is still starting/unhealthy
        res_unhealthy = run_cmd([
            "docker", "compose", "-f", COMPOSE_FILE, "ps", "--filter", "health=starting"
        ], check=False)
        
        if "starting" not in res_unhealthy.stdout.lower() and output:
            print("[Integration] All test containers are ready and healthy!")
            return True
        
        print("[Integration] Waiting for healthchecks...")
        time.sleep(3)
    
    raise TimeoutError(f"Containers did not become healthy within {timeout_sec} seconds")


def run_go_integration_tests() -> bool:
    print("\n[Integration] Running Go integration tests...")
    env = os.environ.copy()
    env["TEST_REDIS_ADDR"] = "localhost:6380"
    env["TEST_POSTGRES_ADDR"] = "localhost:5433"

    # Routing service
    res_routing = subprocess.run(
        ["go", "test", "-tags=integration", "-v", "./tests/integration/..."],
        cwd="services/routing-service",
        env=env,
        text=True,
    )
    if res_routing.returncode != 0:
        return False

    # Tile service
    res_tile = subprocess.run(
        ["go", "test", "-tags=integration", "-v", "./tests/integration/..."],
        cwd="services/tile-service",
        env=env,
        text=True,
    )
    return res_tile.returncode == 0


def run_python_integration_tests() -> bool:
    print("\n[Integration] Running Python integration tests...")
    env = os.environ.copy()
    env["TEST_ELASTICSEARCH_URL"] = "http://localhost:9201"

    res = subprocess.run(
        [sys.executable, "-m", "pytest", "tests/integration", "-m", "integration", "-v"],
        cwd="services/geocoding-service",
        env=env,
        text=True,
    )
    return res.returncode == 0


def main():
    parser = argparse.ArgumentParser(description="Run integration tests with ephemeral containers")
    parser.add_argument("--skip-teardown", action="store_true", help="Leave containers running after tests")
    args = parser.parse_args()

    overall_success = False
    try:
        start_containers()
        wait_for_health()
        
        go_ok = run_go_integration_tests()
        py_ok = run_python_integration_tests()
        
        overall_success = go_ok and py_ok
    except Exception as e:
        print(f"[Error] Integration test execution failed: {e}", file=sys.stderr)
        overall_success = False
    finally:
        if not args.skip_teardown:
            stop_containers()

    if overall_success:
        print("\n[Integration] SUCCESS: All integration tests passed!")
        sys.exit(0)
    else:
        print("\n[Integration] FAILURE: One or more integration tests failed.", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
