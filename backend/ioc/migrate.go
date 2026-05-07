package ioc

import (
	"fmt"

	"github.com/mymikasa/prompthub/internal/domain"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&domain.User{},
		&domain.Workspace{},
		&domain.WorkspaceMember{},
		&domain.Prompt{},
		&domain.PromptVersion{},
		&domain.PromptVariable{},
		&domain.PromptTestCase{},
		&domain.PromptRun{},
		&domain.Tag{},
		&domain.PromptTag{},
		&domain.ProviderConfig{},
	)
	if err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}
