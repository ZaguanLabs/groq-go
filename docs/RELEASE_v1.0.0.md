# Release v1.0.0

This document summarizes the changes included in **groq-go v1.0.0**.

## Highlights

- Synchronized the Go SDK with the **official Groq Python SDK v1.0.0**.
- Aligned chat types and metadata structures to match Python’s v1.0.0 OpenAPI/Stainless output.

## Changes

### Chat / Types

- Updated `ChatCompletionMessageParam` to support:
  - `function_call` (deprecated, for compatibility)
  - `reasoning` (assistant reasoning when `reasoning_format=parsed`)

- Added deprecated request fields to `CreateChatCompletionRequest`:
  - `functions`
  - `function_call`

### XGroq Metadata Alignment

- Refactored Groq-specific metadata structs to match Python SDK v1.0.0:
  - `XGroq` is used for **non-streaming** responses.
  - `XGroqStream` is used for **streaming** responses.
  - `XGroqUsage` contains hardware cache statistics:
    - `dram_cached_tokens`
    - `sram_cached_tokens`

#### Breaking Change

- Code referencing `XGroq.CacheStats` must be updated to use:
  - `XGroq.Usage.DramCachedTokens`
  - `XGroq.Usage.SramCachedTokens`

## Verification

- `go test ./...` passes.
- `go test -race ./...` passes.

## Versioned Files

- `VERSION`: `1.0.0`
- `groq/constants.go`: `Version = "1.0.0"`
- `CHANGELOG.md`: `1.0.0` entry

## Release Steps (Git)

```bash
# Verify all changes
git status
git diff

# Stage all changes
git add -A

# Commit with release message
git commit -m "Release v1.0.0 - Python SDK v1.0.0 Parity"

# Tag and push
git tag -a v1.0.0 -m "v1.0.0 - Python SDK v1.0.0 Parity"
git push origin main
git push origin v1.0.0
```
