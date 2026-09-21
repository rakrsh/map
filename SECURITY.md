# Security Scanning and Reporting

## Local checks

Install Semgrep and run the local static scan:

```bash
python -m pip install semgrep
make scan-sast
```

Run the combined local security checks when Docker is available:

```bash
make scan-security
```

Run dependency, license, and filesystem checks with Trivy:

```bash
make scan-dependencies
make scan-licenses
```

Run a baseline DAST scan against a running service:

```bash
make dev-up
make scan-dast TARGET=http://localhost:8000/health
```

## CI coverage

Pull requests and pushes to `main` run:

- Semgrep SAST
- Gitleaks secret detection
- Trivy filesystem vulnerability, license, and misconfiguration scanning
- ZAP baseline DAST against the local geocoding health endpoint
- SonarCloud analysis when the repository variable `SONAR_ENABLED=true` and secrets `SONAR_TOKEN` and `SONAR_HOST_URL` are configured

A scan with a High or Critical vulnerability, a detected secret, a restricted license, or a failed quality gate blocks the security job.

## Findings and remediation

1. Reproduce the finding locally with the command shown in the workflow or Makefile.
2. Fix the source, dependency, configuration, or secret exposure rather than suppressing the rule.
3. Add a focused regression test where behavior changes.
4. Document an unavoidable finding with its rule ID, rationale, owner, expiry date, and compensating control before adding a narrow suppression.
5. Never commit credentials, tokens, private keys, `.env` files, scan reports containing secrets, or generated local data.

Report suspected vulnerabilities privately through the repository's GitHub security advisory process. Do not open a public issue containing an exploitable secret or vulnerability details.
