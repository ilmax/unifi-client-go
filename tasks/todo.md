# Todo

## Acceptance Criteria
- Release workflow generates release notes that include `oasdiff` output for API diffs.
- Breaking changes are called out using `oasdiff breaking`.
- Diff summary is reproducible from the OpenAPI specs (old vs new) without manual steps.
- README documents the diff behavior and how previous API versions are detected.

## Plan
- [x] Inspect current release workflow and release notes generation to identify where to inject `oasdiff` outputs.
- [x] Define how to discover the previous API version from the latest release body.
- [x] Update the workflow to download both specs, run `oasdiff diff` and `oasdiff breaking`, and append results to the release notes.
- [x] Add tooling setup for `oasdiff` in the workflow.
- [x] Update README to document the diff output and workflow usage.
- [ ] Verify by running tests and a local dry-run example.

## Working Notes
- Prefer diffing the OpenAPI JSON directly; avoid parsing generated Go to reduce noise.
- Keep output stable for release notes (markdown section with counts + lists).

## Progress
- [ ] In progress

## Results
- Release workflow now installs `oasdiff`, compares specs to the previous release, and appends diff + breaking changes to release notes.
- README documents the `oasdiff`-based release note sections and version parsing.
