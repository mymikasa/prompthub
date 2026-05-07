package domain

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Prompt struct {
	ID                  int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	WorkspaceID         int64          `gorm:"uniqueIndex:idx_ws_slug;index;not null" json:"workspace_id"`
	CreatedBy           int64          `gorm:"index;not null" json:"created_by"`
	Title               string         `gorm:"type:varchar(200);not null" json:"title"`
	Slug                string         `gorm:"type:varchar(200);uniqueIndex:idx_ws_slug;not null" json:"slug"`
	Description         string         `gorm:"type:text" json:"description"`
	Body                string         `gorm:"type:longtext;not null" json:"body"`
	MessageFormat       string         `gorm:"type:varchar(20);not null" json:"message_format"`
	Visibility          string         `gorm:"type:varchar(20);not null" json:"visibility"`
	Status              string         `gorm:"type:varchar(20);index;not null" json:"status"`
	TargetProvider      string         `gorm:"type:varchar(50)" json:"target_provider"`
	TargetModel         string         `gorm:"type:varchar(100)" json:"target_model"`
	DefaultTemperature  *float64       `json:"default_temperature"`
	DefaultMaxTokens    *int           `json:"default_max_tokens"`
	UsageNotes          string         `gorm:"type:text" json:"usage_notes"`
	CreatedAt           time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"index" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

type PromptVersion struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PromptID  int64     `gorm:"uniqueIndex:idx_prompt_ver;index;not null" json:"prompt_id"`
	Version   int       `gorm:"uniqueIndex:idx_prompt_ver;not null" json:"version"`
	Snapshot  string    `gorm:"type:json;not null" json:"snapshot"`
	CreatedBy int64     `gorm:"index;not null" json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type PromptVariable struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PromptID     int64     `gorm:"uniqueIndex:idx_prompt_var_name;index;not null" json:"prompt_id"`
	Name         string    `gorm:"type:varchar(100);uniqueIndex:idx_prompt_var_name;not null" json:"name"`
	Label        string    `gorm:"type:varchar(200)" json:"label"`
	Description  string    `gorm:"type:text" json:"description"`
	Required     bool      `gorm:"not null;default:false" json:"required"`
	DefaultValue string    `gorm:"type:text" json:"default_value"`
	ExampleValue string    `gorm:"type:text" json:"example_value"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PromptTestCase struct {
	ID               int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PromptID         int64     `gorm:"index;not null" json:"prompt_id"`
	Name             string    `gorm:"type:varchar(200);not null" json:"name"`
	VariableValues   string    `gorm:"type:json" json:"variable_values"`
	ExpectedBehavior string    `gorm:"type:text" json:"expected_behavior"`
	ExpectedOutput   string    `gorm:"type:text" json:"expected_output"`
	CreatedBy        int64     `gorm:"index;not null" json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type PromptRun struct {
	ID                      int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	PromptID                int64        `gorm:"index;not null" json:"prompt_id"`
	PromptVersionID         int64        `gorm:"index;not null" json:"prompt_version_id"`
	TestCaseID              sql.NullInt64 `gorm:"index" json:"test_case_id"`
	Provider                string       `gorm:"type:varchar(50)" json:"provider"`
	Model                   string       `gorm:"type:varchar(100)" json:"model"`
	InputVariables          string       `gorm:"type:json" json:"input_variables"`
	RenderedPromptSnapshot  string       `gorm:"type:json" json:"rendered_prompt_snapshot"`
	OutputText              string       `gorm:"type:longtext" json:"output_text"`
	Latency                 int          `json:"latency"`
	TokenUsage              string       `gorm:"type:json" json:"token_usage"`
	ErrorMessage            string       `gorm:"type:text" json:"error_message"`
	CreatedBy               int64        `gorm:"index;not null" json:"created_by"`
	CreatedAt               time.Time    `gorm:"index" json:"created_at"`
}
