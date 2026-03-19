# Todo

## Acceptance Criteria
- The generated Network client supports `api_key` authentication via request headers.
- The public constructor in `pkg/network` and `unifi.NewNetwork(...)` expose API-key auth configuration.
- Username/password login is no longer advertised in code examples/docs.
- Tests cover the API-key header injection path and constructor validation where possible.
- Verification steps are recorded.

## Plan
- [x] Add API-key auth support to the generated Network client wrapper.
- [x] Update docs/examples/workflow snippets to show API-key-based usage.
- [x] Verify with tests.

## Working Notes
- Use the generated client's request editor hook instead of modifying generated files directly.

## Progress
- [x] In progress

## Results
- Added `APIKey` support to `network.Config` and wired it into the generated client through `oapi-codegen` request editors using the `X-API-Key` header.
- Updated README, example code, and release-note usage snippets to show API-key authentication and removed the old username/password/site-based examples.
- Added tests for constructor validation, API-key header injection, and timeout defaults.
- Verified with `go test ./...` and a repo-wide search confirming the old login/config references are gone.
