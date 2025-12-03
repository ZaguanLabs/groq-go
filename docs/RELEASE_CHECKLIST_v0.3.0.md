# v0.3.0 Release Checklist ✅

## Pre-Release Verification

- [x] ✅ Version updated to `0.3.0` in `groq/constants.go`
- [x] ✅ Version updated to `0.3.0` in `VERSION` file
- [x] ✅ CHANGELOG.md updated with v0.3.0 release notes
- [x] ✅ README.md updated with v0.3.0 features
- [x] ✅ All 238 tests passing
- [x] ✅ Test coverage at 73.5%+ (maintained)
- [x] ✅ Zero race conditions
- [x] ✅ Release documentation created
- [x] ✅ .gitignore updated (tmp directory excluded)

## Release Files

- [x] `groq/constants.go` - Version = "0.3.0"
- [x] `VERSION` - 0.3.0
- [x] `CHANGELOG.md` - v0.3.0 release notes with Python SDK 0.37.0 sync
- [x] `README.md` - Updated with v0.3.0 features
- [x] `docs/RELEASE_v0.3.0.md` - Comprehensive release documentation
- [x] `groq/types/shared.go` - New token details types
- [x] `groq/types/chat.go` - MCP tools, JSON Schema, cache stats
- [x] `groq/types/audio.go` - Sample rate, URL fields

## New Features

### Python SDK 0.37.0 Synchronization
- [x] ✅ `CompletionUsage.QueueTime` field
- [x] ✅ `CompletionTokensDetails` type with `ReasoningTokens`
- [x] ✅ `PromptTokensDetails` type with `CachedTokens`
- [x] ✅ `ChatCompletion.McpListTools` field
- [x] ✅ `ChatCompletion.ServiceTier` field
- [x] ✅ `ChatCompletion.UsageBreakdown` field
- [x] ✅ `ChatCompletion.XGroq` field
- [x] ✅ `McpListTool` and `McpListToolDef` types
- [x] ✅ `XGroqCacheStats` type with DRAM/SRAM metrics
- [x] ✅ `ResponseFormat.JSONSchema` field
- [x] ✅ `ResponseFormatJSONSchema` type
- [x] ✅ `CreateSpeechRequest.SampleRate` field
- [x] ✅ `CreateTranscriptionRequest.URL` field

### Documentation
- [x] ✅ Release notes document
- [x] ✅ Updated main README
- [x] ✅ Updated CHANGELOG

### Testing
- [x] ✅ All existing tests passing
- [x] ✅ Build successful
- [x] ✅ Race detection clean

## Quality Metrics

| Metric | v0.2.1 | v0.3.0 |
|--------|--------|--------|
| Tests | 140+ | 238 |
| Coverage | 73.5% | 73.5%+ |
| Race Conditions | 0 | 0 |
| Pass Rate | 100% | 100% |

## Git Commands for Release

```bash
# Verify all changes
git status
git diff

# Stage all changes
git add -A

# Commit with release message
git commit -m "Release v0.3.0 - Python SDK 0.37.0 Parity

- Synchronized with Groq Python SDK v0.37.0
- Added CompletionTokensDetails and PromptTokensDetails
- Added MCP tool discovery support (McpListTool, McpListToolDef)
- Added JSON Schema structured output (ResponseFormatJSONSchema)
- Added XGroqCacheStats for hardware cache metrics
- Added SampleRate to CreateSpeechRequest
- Added URL to CreateTranscriptionRequest
- Enhanced ChatCompletion with McpListTools, ServiceTier, UsageBreakdown, XGroq
- 238 tests passing, 73.5%+ coverage, zero race conditions"

# Create and push tag
git tag -a v0.3.0 -m "v0.3.0 - Python SDK 0.37.0 Parity"
git push origin main
git push origin v0.3.0
```

## Post-Release

- [ ] Verify tag on GitHub
- [ ] Check pkg.go.dev for new version
- [ ] Update any dependent projects
