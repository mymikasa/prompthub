package domain

import (
	"database/sql"
	"time"
)

type Prompt struct {
	ID                  int64     `json:"id"`
	WorkspaceID         int64     `json:"workspace_id"`
	CreatedBy           int64     `json:"created_by"`
	Title               string    `json:"title"`
	Slug                string    `json:"slug"`
	Description         string    `json:"description"`
	Body                string    `json:"body"`
	MessageFormat       string    `json:"message_format"`
	Visibility          string    `json:"visibility"`
	Status              string    `json:"status"`
	TargetProvider      string    `json:"target_provider"`
	TargetModel         string    `json:"target_model"`
	DefaultTemperature  *float64  `json:"default_temperature"`
	DefaultMaxTokens    *int      `json:"default_max_tokens"`
	UsageNotes          string    `json:"usage_notes"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type PromptVersion struct {
	ID        int64     `json:"id"`
	PromptID  int64     `json:"prompt_id"`
	Version   int       `json:"version"`
	Snapshot  string    `json:"snapshot"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type PromptVariable struct {
	ID           int64     `json:"id"`
	PromptID     int64     `json:"prompt_id"`
	Name         string    `json:"name"`
	Label        string    `json:"label"`
	Description  string    `json:"description"`
	Required     bool      `json:"required"`
	DefaultValue string    `json:"default_value"`
	ExampleValue string    `json:"example_value"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PromptTestCase struct {
	ID               int64     `json:"id"`
	PromptID         int64     `json:"prompt_id"`
	Name             string    `json:"name"`
	VariableValues   string    `json:"variable_values"`
	ExpectedBehavior string    `json:"expected_behavior"`
	ExpectedOutput   string    `json:"expected_output"`
	CreatedBy        int64     `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type PromptRun struct {
	ID                      int64        `json:"id"`
	PromptID                int64        `json:"prompt_id"`
	PromptVersionID         int64        `json:"prompt_version_id"`
	TestCaseID              sql.NullInt64 `json:"test_case_id"`
	Provider                string       `json:"provider"`
	Model                   string       `json:"model"`
	InputVariables          string       `json:"input_variables"`
	RenderedPromptSnapshot  string       `json:"rendered_prompt_snapshot"`
	OutputText              string       `json:"output_text"`
	Latency                 int          `json:"latency"`
	TokenUsage              string       `json:"token_usage"`
	ErrorMessage            string       `json:"error_message"`
	CreatedBy               int64        `json:"created_by"`
	CreatedAt               time.Time    `json:"created_at"`
}
