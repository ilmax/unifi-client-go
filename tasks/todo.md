# Todo

## Acceptance Criteria
- OpenAPI generation (download, preprocess, codegen) is moved into a Makefile target.
- The release workflow delegates spec generation to the Makefile.
- Generation downloads specs from the UniFi developer endpoint.
- The compatibility rewrite is limited to the schema still rejected by the generator.
- Tests and builds pass after the generator and dependency upgrades.

## Plan
- [x] Add a Makefile target that downloads the spec, runs preprocessing, and generates the client.
- [x] Update the release workflow to call the Makefile target.
- [x] Switch downloads to `https://developer.ui.com/network/v<version>/openapi.json`.
- [x] Upgrade to `oapi-codegen` v2.8.0, runtime v1.7.0, and libopenapi v0.40.0.
- [x] Reduce preprocessing to the entity-metadata compatibility rewrite.
- [x] Add an upstream tracking link for removing the remaining workaround.
- [x] Run tests and build all packages.

## Working Notes
- Prefer normalized versions without a leading v when naming files.
- `oapi-codegen` v2.8.0 supports OpenAPI 3.1, but still rejects the UniFi entity-metadata `allOf` aggregate when it merges discriminator-bearing schemas.
- The closest active upstream tracking issue is https://github.com/oapi-codegen/oapi-codegen/issues/1593. The broader OpenAPI 3.1 issue, https://github.com/oapi-codegen/oapi-codegen/issues/373, is closed by v2.8.0.

## Progress
- [x] Complete

## Results
- Added Makefile target for OpenAPI generation and wired release workflow to call it.
- Switched OpenAPI downloads from the GitHub mirror to the UniFi developer endpoint.
- Upgraded the generator and dependencies to their latest compatible releases; the project now targets Go 1.27.x.
- Removed DNS and firewall discriminator union rewrites; only `User defined entity metadata` is normalized for code generation.
- Verified direct generation with `oapi-codegen` v2.8.0 after the minimal rewrite.
- `go test ./...` and `go build ./...` pass.
