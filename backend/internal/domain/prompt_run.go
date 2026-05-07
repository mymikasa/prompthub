package domain

import "time"

type PromptRun struct {
	ID                     int64             `json:"id"`
	PromptID               int64             `json:"promptId"`
	PromptVersionID        int64             `json:"promptVersionId"`
	TestCaseID             *int64            `json:"testCaseId"`
	Provider               string            `json:"provider"`
	Model                  string            `json:"model"`
	InputVariables         map[string]string `json:"inputVariables"`
	RenderedPromptSnapshot any               `json:"renderedPromptSnapshot"`
	OutputText             string            `json:"outputText"`
	Latency                int               `json:"latency"`
	TokenUsage             *TokenUsage       `json:"tokenUsage"`
	ErrorMessage           string            `json:"errorMessage"`
	CreatedBy              int64             `json:"createdBy"`
	CreatedAt              time.Time         `json:"createdAt"`
}

type TokenUsage struct {
	PromptTokens     int `json:"promptTokens"`
	CompletionTokens int `json:"completionTokens"`
	TotalTokens      int `json:"totalTokens"`
}
