# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-11-18

### Added
- Initial release of the Groq Go SDK.
- Support for **Chat Completions** API with streaming support (SSE).
- Support for **Audio** API (Speech, Transcription, Translation).
- Support for **Embeddings**, **Models**, **Files**, and **Batches** APIs.
- **Multipart/form-data** encoding for file uploads.
- **Optional[T]** type for handling optional and nullable fields in JSON.
- Robust **retry logic** with exponential backoff and `Retry-After` header support.
- **Context** support for request cancellation and timeouts.
- **Idempotency Key** support via request options.
- Helper functions for configuration (`WithAPIKey`, `WithBaseURL`, etc.).
