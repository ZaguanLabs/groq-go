package types

// Role represents the role of a message author
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
	RoleFunction  Role = "function"
)

// FinishReason represents the reason a chat completion finished
type FinishReason string

const (
	FinishReasonStop          FinishReason = "stop"
	FinishReasonLength        FinishReason = "length"
	FinishReasonToolCalls     FinishReason = "tool_calls"
	FinishReasonContentFilter FinishReason = "content_filter"
	FinishReasonFunctionCall  FinishReason = "function_call" // Deprecated but still used?
)

// ModelID represents available Groq model identifiers
type ModelID string

const (
	// Compound AI models
	ModelCompoundBeta     ModelID = "compound-beta"
	ModelCompoundBetaMini ModelID = "compound-beta-mini"

	// Llama models
	ModelLlama31_8BInstant     ModelID = "llama-3.1-8b-instant"
	ModelLlama33_70BVersatile  ModelID = "llama-3.3-70b-versatile"
	ModelLlama4Maverick17B128E ModelID = "meta-llama/llama-4-maverick-17b-128e-instruct"
	ModelLlama4Scout17B16E     ModelID = "meta-llama/llama-4-scout-17b-16e-instruct"
	ModelLlamaGuard412B        ModelID = "meta-llama/llama-guard-4-12b"

	// Gemma models
	ModelGemma29BIT ModelID = "gemma2-9b-it"

	// Kimi models
	ModelKimiK2Instruct ModelID = "moonshotai/kimi-k2-instruct"

	// OpenAI OSS models
	ModelGPTOSS120B ModelID = "openai/gpt-oss-120b"
	ModelGPTOSS20B  ModelID = "openai/gpt-oss-20b"

	// Qwen models
	ModelQwen332B ModelID = "qwen/qwen3-32b"
)
