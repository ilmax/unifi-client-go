# Todo

## Acceptance Criteria
- Release workflow chooses a free SDK version even if `VERSION` lags behind existing tags/branches.
- Release workflow pushes release commit content to `main` as part of the run.
- Existing release/tag flow remains intact.
- Verification steps are recorded.

## Plan
- [x] Update version selection to loop until both tag and release branch are available.
- [x] Push release commit to `main` in addition to release branch and tag.
- [x] Verify and document results.

## Working Notes
- Main issue: `VERSION` on `main` can lag, and previous logic bumped only once.

## Progress
- [x] In progress

## Results
- Workflow now increments SDK patch version in a loop until both tag and `release/` branch names are free.
- Workflow now pushes release commit to both `release/<tag>` and `main`.
- Verified repository tests still pass with `go test ./...`.
