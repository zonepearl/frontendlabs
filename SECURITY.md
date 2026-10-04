# Security policy

## Reporting a vulnerability

Please report security issues privately through GitHub:
**Security → Report a vulnerability** on this repository
(private vulnerability reporting is enabled). Do not open a public issue.

You can expect an acknowledgement within 7 days. Include the affected URL or
file, steps to reproduce, and the impact you see.

## Scope

- The site generator (`cmd/`, `internal/`), templates and scripts in this repository
- The deployed site at https://frontendlabs.xyz
- The build and deploy pipeline (`.github/workflows/`)

## How this repository is protected

- `main` accepts changes only through pull requests whose `e2e` check passed;
  force pushes, deletion and non-linear history are blocked for everyone,
  including administrators (`.github/rulesets/`).
- Tags cannot be moved or deleted once published.
- Secret scanning with push protection, Dependabot alerts and security updates
  are enabled; workflow actions are pinned to commit SHAs and run with
  read-only tokens by default.
