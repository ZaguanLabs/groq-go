# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.1] - 2025-11-21

### Added
- **Multimodal Content Parts**: Support for document content in chat messages (Groq API v0.36.0)
  - `ContentPart` interface for text, image, and document content
  - `ContentPartText` for text content parts
  - `ContentPartImage` for image URL content parts
  - `ContentPartDocument` for structured JSON document content (NEW)
  - Ability to send multiple content types in a single message
- **New Types**:
  - `ContentPartDocument_Document` with `Data` (map[string]interface{}) and optional `ID`
  - Full support for multimodal messages with mixed content types
- **Examples**: Added `document_content` example demonstrating JSON document usage
- **Tests**: 5 new comprehensive tests for content part serialization and multimodal messages
- **Documentation**: 
  - `patches/PATCH_ANALYSIS.md` - Analysis of Groq Python SDK v0.36.0 changes
  - `groq/examples/document_content/README.md` - Complete feature documentation

### Changed
- Updated `ChatCompletionMessageParam.Content` documentation to clarify multimodal support
- Enhanced README.md with multimodal content feature

### Quality
- **Test Coverage**: Maintained at 73.5%
- **New Tests**: 5 content part tests (all passing)
- **Total Tests**: 140+ comprehensive tests
- **Pass Rate**: 100%
- **Backward Compatibility**: ✅ Fully maintained (string content still works)

### Synchronization
- Synchronized with Groq Python SDK v0.36.0 API changes
- Implemented document content parts matching Python SDK functionality

## [0.2.0] - 2025-11-19

### Added
- **Compound AI Support**: Multi-model orchestration with custom model selection and tool configuration
  - `CompoundCustom` type for configuring answering and reasoning models
  - `CompoundCustomTools` for enabling web search, code interpreter, and Wolfram tools
  - `ExecutedTool` types for browser, code, and search results
- **Enhanced Streaming Metadata**: `XGroq` field in streaming responses
  - Request ID tracking
  - Debug information (token IDs and strings)
  - Per-model usage breakdown for compound AI
  - Error reporting in streams
- **Documents & Citations**: RAG-like workflows with citation support
  - `Document` type for providing text or JSON context
  - `Annotation` types for document and function citations
  - `CitationOptions` parameter for enabling citations
- **Reasoning Model Support**: Advanced reasoning capabilities
  - `ReasoningEffort` parameter (none, default, low, medium, high)
  - `ReasoningFormat` parameter (hidden, raw, parsed)
  - `IncludeReasoning` flag
  - `Reasoning` field in responses
- **Search Settings**: Fine-grained web search control
  - Country-specific search
  - Domain filtering (include/exclude)
  - Image inclusion
- **New Request Parameters**:
  - `MaxCompletionTokens` for token limit control
  - `ServiceTier` for tier selection (auto, on_demand, flex, performance)
  - `DisableToolValidation` flag
  - `Metadata` and `Store` fields (future support)
- **New Model Constants**: Added 9 new model identifiers including Compound AI, Llama 4, Kimi, GPT-OSS, and Qwen models
- **Examples**: Added comprehensive examples for Compound AI, Documents/RAG, and Reasoning models
- **Comprehensive Test Suite**: 135+ tests achieving 73.5% overall coverage
  - 100% coverage for all 6 resource packages (chat, models, embeddings, audio, files, batches)
  - 100% coverage for option and retry packages
  - 79.3% coverage for main client package
  - Integration tests with mock HTTP servers
  - Error path and edge case testing
  - Context cancellation testing
  - Streaming API tests with SSE parsing

### Changed
- Enhanced `ChatCompletionMessage` and `ChatCompletionChunkDelta` with new fields
- Improved `XGroq` type with complete metadata structure
- Updated streaming responses to include usage breakdown

### Quality
- **Test Coverage**: 73.5% overall (up from 38.9%)
- **Test Count**: 135+ comprehensive tests
- **Pass Rate**: 100%
- **Race Conditions**: 0
- **Audit Grade**: A- (91%)

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
