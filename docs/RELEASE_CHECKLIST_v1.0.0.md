# v1.0.0 Release Checklist

## Pre-Release Verification

- [x] Version updated to `1.0.0` in `groq/constants.go`
- [x] Version updated to `1.0.0` in `VERSION` file
- [x] `CHANGELOG.md` updated with v1.0.0 release notes
- [x] `README.md` updated with v1.0.0 features
- [x] All tests passing (`go test ./...`)
- [x] Race detection clean (`go test -race ./...`)
- [x] Release documentation created (`docs/RELEASE_v1.0.0.md`)

## Release Files

- [x] `groq/constants.go` - Version = "1.0.0"
- [x] `VERSION` - 1.0.0
- [x] `CHANGELOG.md` - v1.0.0 release notes with Python SDK v1.0.0 sync
- [x] `README.md` - Updated with v1.0.0 features
- [x] `docs/RELEASE_v1.0.0.md` - Release documentation
- [x] `groq/types/chat.go` - XGroq alignment + request param additions

## v1.0.0 Highlights (Python SDK v1.0.0 Synchronization)

### Chat

- [x] `ChatCompletionMessageParam.FunctionCall` (deprecated)
- [x] `ChatCompletionMessageParam.Reasoning`
- [x] `CreateChatCompletionRequest.Functions` (deprecated)
- [x] `CreateChatCompletionRequest.FunctionCall` (deprecated)

### XGroq / Metadata Alignment

- [x] `XGroq` (non-streaming) aligned with Python SDK v1.0.0
- [x] `XGroqStream` (streaming) aligned with Python SDK v1.0.0
- [x] `XGroqUsage` contains `dram_cached_tokens` and `sram_cached_tokens`

## Git Commands for Release

```bash
# Verify all changes
git status
git diff

# Stage all changes
git add -A

# Commit with release message
git commit -m "Release v1.0.0 - Python SDK v1.0.0 Parity"

# Create and push tag
git tag -a v1.0.0 -m "v1.0.0 - Python SDK v1.0.0 Parity"
git push origin main
git push origin v1.0.0
```

## Post-Release

- [ ] Verify tag on GitHub
- [ ] Check pkg.go.dev for new version
- [ ] Update any dependent projects
