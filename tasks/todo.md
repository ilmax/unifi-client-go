# Todo

## Acceptance Criteria
- Release workflow uses `unifi-network-version` as the baseline API version when present and valid.
- Workflow falls back to the latest release parsing only when the file is missing or invalid.
- README documents the baseline file behavior.
- Verification is run or a reason is documented.

## Plan
- [x] Update workflow to prefer `unifi-network-version` for the previous API version and keep the release fallback.
- [x] Update README to document the baseline file.
- [x] Verify changes.

## Working Notes
- Treat missing/invalid `unifi-network-version` as "no baseline".

## Progress
- [x] In progress

## Results
- Workflow uses `unifi-network-version` when valid, and falls back to parsing the latest release body otherwise.
- README documents the baseline file behavior.
- Verified with `go test ./...`.
