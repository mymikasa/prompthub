package domain

import "time"

type Prompt struct {
	ID                 int64     `json:"id"`
	WorkspaceID        int64     `json:"workspaceId"`
	CreatedBy          int64     `json:"createdBy"`
	Title              string    `json:"title"`
	Slug               string    `json:"slug"`
	Description        string    `json:"description"`
	Body               string    `json:"body"`
	MessageFormat      string    `json:"messageFormat"`
	Visibility         string    `json:"visibility"`
	Status             string    `json:"status"`
	TargetProvider     string    `json:"targetProvider"`
	TargetModel        string    `json:"targetModel"`
	ProviderConfigID   *int64    `json:"providerConfigId"`
	DefaultTemperature *float64  `json:"defaultTemperature"`
	DefaultMaxTokens   *int      `json:"defaultMaxTokens"`
	UsageNotes         string    `json:"usageNotes"`
	Tags               []string  `json:"tags"`
	CurrentVersionID   int64     `json:"currentVersionId"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

func (p *Prompt) CanView(userID int64) bool {
	if p.Visibility == "workspace" {
		return true
	}
	return p.CreatedBy == userID
}

func (p *Prompt) CanEdit(actor *Actor) bool {
	if actor.IsOwner() {
		return true
	}
	return p.CreatedBy == actor.UserID
}

type PromptList struct {
	Items    []*Prompt
	Total    int64
	Page     int
	PageSize int
}

type VersionSnapshot struct {
	Content        string              `json:"content"`
	Messages       []ChatMessage       `json:"messages"`
	Variables      []PromptVariable    `json:"variables"`
	TargetProvider string              `json:"targetProvider"`
	TargetModel    string              `json:"targetModel"`
	Status         string              `json:"status"`
	Tags           []string            `json:"tags"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PromptVersion struct {
	ID            int64           `json:"id"`
	PromptID      int64           `json:"promptId"`
	VersionNumber int             `json:"versionNumber"`
	ChangeNote    string          `json:"changeNote"`
	Author        string          `json:"author"`
	Snapshot      VersionSnapshot `json:"snapshot"`
	CreatedAt     time.Time       `json:"createdAt"`
}
