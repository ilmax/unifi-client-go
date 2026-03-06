# Todo

## Acceptance Criteria
- Generating the client from UniFi OpenAPI `10.0.162` succeeds.
- `go build ./...` succeeds after generation.
- Missing runtime dependencies from generated code are added to module metadata.
- Verification steps are recorded.

## Plan
- [x] Reproduce generation/build for API version `10.0.162` locally.
- [x] Add/fix required dependencies for generated `oapi-codegen` output.
- [x] Verify with build/tests.

## Working Notes
- Follow workflow behavior as closely as possible.

## Progress
- [ ] In progress

## Results
- Download of `10.0.162` spec succeeded via `https://raw.githubusercontent.com/beezly/unifi-apis/main/unifi-network/10.0.162.json`.
- Generation with `oapi-codegen` reproduced missing dependency errors for `github.com/oapi-codegen/runtime` and `github.com/oapi-codegen/runtime/types`.
- Added `github.com/oapi-codegen/runtime v1.2.0` to module metadata (`go mod tidy` also added required indirect deps).
- Verified `go build ./...` and `go test ./...` succeed after adding dependency.
