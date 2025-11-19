# GitHub Release Notes for v0.2.0-beta

**Copy this content when creating the GitHub release**

---

## 🎉 v0.2.0-beta - Python SDK v0.35.0 Feature Parity

This beta release achieves **full feature parity** with the official Groq Python SDK v0.35.0, making the Groq Go SDK the most comprehensive and up-to-date Go client for the Groq API.

### 🚀 Major Features

#### Compound AI Support
Multi-model orchestration for sophisticated AI workflows:
- Custom model selection (answering + reasoning models)
- Tool configuration (web search, code interpreter, Wolfram)
- Executed tool results tracking (browser, code, search)

#### Documents & Citations (RAG)
Provide context documents with automatic citation tracking:
- Text and JSON document sources
- Document and function call citations
- Citation options control

#### Reasoning Models
Advanced reasoning with configurable output:
- Reasoning effort levels (none, low, medium, high)
- Reasoning format options (hidden, raw, parsed)
- Reasoning content in responses

#### Enhanced Streaming Metadata
Rich metadata in streaming responses via `XGroq` field:
- Request ID tracking
- Debug information (token IDs and strings)
- Per-model usage breakdown for compound AI
- Error reporting in streams

### 📦 Additional Features

- **Search Settings**: Fine-grained web search control (country, domains, images)
- **New Request Parameters**: `MaxCompletionTokens`, `ServiceTier`, `DisableToolValidation`
- **9 New Model Constants**: Compound AI, Llama 4, Kimi, GPT-OSS, Qwen models

### 📚 New Examples

Three comprehensive examples added:
- `groq/examples/compound_ai/` - Multi-model workflows with tools
- `groq/examples/documents_rag/` - Document-based context with citations
- `groq/examples/reasoning/` - Advanced reasoning capabilities

### ✅ Quality Assurance

- All tests passing (100% pass rate)
- All examples compile and run
- Documentation complete
- No breaking changes
- Backward compatible with v0.1.0

### 📖 Documentation

- [CHANGELOG.md](CHANGELOG.md) - Complete changelog
- [docs/v0.2.0-beta-summary.md](docs/v0.2.0-beta-summary.md) - Detailed release summary
- [docs/python-sdk-analysis-v0.35.0.md](docs/python-sdk-analysis-v0.35.0.md) - Python SDK analysis
- [README.md](README.md) - Updated with new features

### 🔄 Migration from v0.1.0

**No migration needed!** All existing code works as-is. New features are opt-in via optional fields.

### 🧪 Beta Testing

We encourage users to test these new features and report any issues. This is a beta release to gather feedback before the stable v0.2.0 release.

### 📦 Installation

```bash
go get github.com/ZaguanLabs/groq-go/groq@v0.2.0-beta
```

### 🐛 Report Issues

https://github.com/ZaguanLabs/groq-go/issues

---

**Full Changelog**: https://github.com/ZaguanLabs/groq-go/compare/v0.1.0...v0.2.0-beta
