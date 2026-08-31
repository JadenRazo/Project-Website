# Contributing

Small, focused pull requests are welcome.

1. Open an issue before a broad architecture or dependency migration.
2. Do not include production data, credentials, private endpoints, or generated dependency directories.
3. Run the relevant frontend or Go checks documented in `README.md`.
4. Explain user-visible behavior, security implications, and known limitations in the pull request.
5. Keep claims tied to a reproducible test, measurement, or source path.

The repository has known Go package debt recorded by `scripts/ci/build-ratchet.sh`. A repaired package must be removed from that baseline in the same change; a new failure must not be added to it.
