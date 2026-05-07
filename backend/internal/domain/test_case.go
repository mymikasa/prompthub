package domain

import "time"

type TestCase struct {
	ID               int64             `json:"id"`
	PromptID         int64             `json:"promptId"`
	Name             string            `json:"name"`
	VariableValues   map[string]string `json:"variableValues"`
	ExpectedBehavior string            `json:"expectedBehavior"`
	ExpectedOutput   string            `json:"expectedOutput"`
	CreatedBy        int64             `json:"createdBy"`
	CreatedAt        time.Time         `json:"createdAt"`
	UpdatedAt        time.Time         `json:"updatedAt"`
}
