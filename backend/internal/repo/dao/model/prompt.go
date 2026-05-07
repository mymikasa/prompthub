package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Prompt struct {
	ID                  int64          `gorm:"primaryKey;autoIncrement"`
	WorkspaceID         int64          `gorm:"uniqueIndex:idx_ws_slug;index;not null"`
	CreatedBy           int64          `gorm:"index;not null"`
	Title               string         `gorm:"type:varchar(200);not null"`
	Slug                string         `gorm:"type:varchar(200);uniqueIndex:idx_ws_slug;not null"`
	Description         string         `gorm:"type:text"`
	Body                string         `gorm:"type:longtext;not null"`
	MessageFormat       string         `gorm:"type:varchar(20);not null"`
	Visibility          string         `gorm:"type:varchar(20);not null"`
	Status              string         `gorm:"type:varchar(20);index;not null"`
	TargetProvider      string         `gorm:"type:varchar(50)"`
	TargetModel         string         `gorm:"type:varchar(100)"`
	DefaultTemperature  *float64
	DefaultMaxTokens    *int
	UsageNotes          string         `gorm:"type:text"`
	CurrentVersionID    int64          `gorm:"index"`
	CreatedAt           time.Time      `gorm:"index"`
	UpdatedAt           time.Time      `gorm:"index"`
	DeletedAt           gorm.DeletedAt `gorm:"index"`
}

type PromptVersion struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	PromptID  int64     `gorm:"uniqueIndex:idx_prompt_ver;index;not null"`
	Version   int       `gorm:"uniqueIndex:idx_prompt_ver;not null"`
	Snapshot  string    `gorm:"type:json;not null"`
	CreatedBy int64     `gorm:"index;not null"`
	CreatedAt time.Time
}

type PromptVariable struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	PromptID     int64     `gorm:"uniqueIndex:idx_prompt_var_name;index;not null"`
	Name         string    `gorm:"type:varchar(100);uniqueIndex:idx_prompt_var_name;not null"`
	Label        string    `gorm:"type:varchar(200)"`
	Description  string    `gorm:"type:text"`
	Required     bool      `gorm:"not null;default:false"`
	DefaultValue string    `gorm:"type:text"`
	ExampleValue string    `gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PromptTestCase struct {
	ID               int64     `gorm:"primaryKey;autoIncrement"`
	PromptID         int64     `gorm:"index;not null"`
	Name             string    `gorm:"type:varchar(200);not null"`
	VariableValues   string    `gorm:"type:json"`
	ExpectedBehavior string    `gorm:"type:text"`
	ExpectedOutput   string    `gorm:"type:text"`
	CreatedBy        int64     `gorm:"index;not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type PromptRun struct {
	ID                     int64        `gorm:"primaryKey;autoIncrement"`
	PromptID               int64        `gorm:"index;not null"`
	PromptVersionID        int64        `gorm:"index;not null"`
	TestCaseID             sql.NullInt64 `gorm:"index"`
	Provider               string       `gorm:"type:varchar(50)"`
	Model                  string       `gorm:"type:varchar(100)"`
	InputVariables         string       `gorm:"type:json"`
	RenderedPromptSnapshot string       `gorm:"type:json"`
	OutputText             string       `gorm:"type:longtext"`
	Latency                int
	TokenUsage             string       `gorm:"type:json"`
	ErrorMessage           string       `gorm:"type:text"`
	CreatedBy              int64        `gorm:"index;not null"`
	CreatedAt              time.Time    `gorm:"index"`
}
