# Todo

## Acceptance Criteria
- Handwritten Site Manager code and top-level Site Manager wiring are removed.
- Release branches are named `release/sdk-vX.Y.Z-network-vA.B.C`.
- Release workflow commits the downloaded Network OpenAPI spec into `openapi/unifi-network/<version>.json`.
- Release workflow generates code from the committed spec file and commits both spec and generated code.
- README/examples describe a Network-only generated SDK.
- Verification steps are recorded.

## Plan
- [x] Remove Site Manager code, top-level wiring, and examples/docs references.
- [ ] Update release workflow to use spec-aware branch names and commit the spec file.
- [ ] Rewrite docs to describe the Network-only generated SDK.
- [ ] Verify with generation/build/tests and record results.

## Working Notes
- No Site Manager replacement will be added in this change because there is no confirmed OpenAPI JSON source for it.

## Progress
- [ ] In progress

## Results
- Removed handwritten Site Manager code, top-level cloud-client wiring, and the Site Manager example.
