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
	DefaultTemperature *float64  `json:"defaultTemperature"`
	DefaultMaxTokens   *int      `json:"defaultMaxTokens"`
	UsageNotes         string    `json:"usageNotes"`
	Tags               []string  `json:"tags"`
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
