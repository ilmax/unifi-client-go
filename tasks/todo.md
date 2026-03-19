# Todo

## Acceptance Criteria
- The OpenAPI preprocessor rewrites the firewall discriminator-based schemas that are still flattened by `oapi-codegen`.
- The rewrite stays curated and safe: if a target family is already fixed upstream, the tool skips it.
- Known upstream oddities in the firewall protocol family are handled without breaking generation.
- Tests cover the firewall rewrite behavior, including a mismatched-parent case.
- Verification steps are recorded.

## Plan
- [x] Map firewall discriminator targets and identify any parent-ref mismatches that need special handling.
- [x] Extend the OpenAPI preprocessor and tests to rewrite the firewall union families safely.
- [x] Verify with tests/build and record the results.

## Working Notes
- The firewall policy graph has nested polymorphic families: action, source/destination filters, port/IP filters, protocol scope, protocol families, and schedule.
- At least one upstream schema appears inconsistent: `Firewall policy IPv4 protocol number` inherits from `Firewall policy IPv6 protocol` while also being referenced from the IPv4 discriminator mapping.
- The existing DNS rewrite machinery should be reusable if we add support for synthetic child overrides when a mapped schema does not inherit from the expected base.

## Progress
- [ ] In progress

## Results
- Confirmed the firewall policy graph needs curated preprocessing for the nested discriminator families used by `Create or update firewall policy` and `Firewall policy`.
- Confirmed a special-case upstream mismatch: `Firewall policy IPv4 protocol` maps `PROTOCOL_NUMBER` to `Firewall policy IPv4 protocol number`, but that schema currently inherits from `Firewall policy IPv6 protocol`.
- Extended the preprocessor to rewrite the firewall action, source/destination filter, port/IP filter, protocol scope, protocol family, named protocol, protocol preset, and schedule unions.
- Added a firewall regression fixture that covers the mismatched IPv4 protocol-number case by cloning it to a synthetic IPv4-specific child while leaving the shared IPv6 branch intact.
- Verified with `go test ./internal/openapipatch ./cmd/openapi-preprocess`, `go test ./...`, and `go build ./...`.
