# Todo

## Acceptance Criteria
- The release process rewrites the DNS policy schemas into a `oneOf`-based union before running `oapi-codegen`.
- The raw upstream spec remains committed under `openapi/unifi-network/<version>.json`.
- The patched codegen input is also committed, and the workflow generates `pkg/network/openapi.gen.go` from it.
- The generated DNS request/response/page item types can model derived DNS policy variants.
- Tests cover the preprocessor behavior.
- Verification steps are recorded.

## Plan
- [x] Add a tested OpenAPI preprocessor that rewrites DNS discriminator schemas into `oneOf` unions.
- [x] Integrate the preprocessor into the release workflow and persist the codegen-ready spec artifact.
- [ ] Update docs and verify with tests/build.

## Working Notes
- `oapi-codegen` only builds discriminator unions when the schema uses `oneOf` or `anyOf`; discriminator + `allOf` inheritance stays flattened.
- A safe rewrite needs a synthetic `... base` schema to avoid introducing a self-referential cycle.
- Keep the change targeted to DNS for now; the upstream spec contains many other discriminator-only schemas.

## Progress
- [x] In progress

## Results
- Added `cmd/openapi-preprocess` plus `internal/openapipatch` to rewrite the two DNS policy schemas into proper `oneOf` unions with synthetic base schemas, avoiding recursive `allOf` cycles.
- Added unit coverage for the rewrite behavior and verified it with `go test ./internal/openapipatch ./cmd/openapi-preprocess`.
- Updated the release workflow to commit a patched `openapi/unifi-network/<version>.codegen.json` artifact and generate `pkg/network/openapi.gen.go` from that file instead of the raw upstream spec.
