# Branch Mode MVP

## Usage

Omit `PR` field or set to `0` to run plan/apply directly on branch.

**Note**: `Ref` must be branch name. Commit SHAs not supported.

```bash
curl --request POST 'https://<ATLANTIS_HOST_NAME>/api/plan' \
--header 'X-Atlantis-Token: <ATLANTIS_API_SECRET>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "Repository": "owner/repo-name",
    "Ref": "main",
    "Type": "Github",
    "Paths": [{
      "Directory": ".",
      "Workspace": "default"
    }]
}'
```

## Implementation

- `PR=0` triggers branch mode
- `apiParseAndValidate()`: uses `Ref` for head/base when `PR=0`
- `working_dir.go`: forces branch checkout strategy when `PR=0`
- VCS status updates skipped (both plan/apply)
- Modified files: `api_controller.go`, `working_dir.go`, `api-endpoints.md`

## Limitations

1. **Locking**: All branch plans share `PR=0` lock key → concurrent plans conflict
2. **VCS Status**: No commit status updates
3. **Lock Cleanup**: `UnlockByPull(repo, 0)` unlocks all branch-mode locks
4. **SHA Support**: `Ref` accepts branch names only, not commit SHAs

## Future Improvements (Option 2)

When this feature proves valuable, consider refactoring to explicit branch mode:

1. Add `BranchMode bool` field to `APIRequest`
2. Implement branch-based locking (key: `repo/branch` instead of `repo/PR`)
3. Add branch-specific commit status support
4. Create separate storage for branch plans vs PR plans
5. Add branch-based lock management UI

## Code Locations

- `server/controllers/api_controller.go`
- `server/events/working_dir.go`
- `runatlantis.io/docs/api-endpoints.md`
