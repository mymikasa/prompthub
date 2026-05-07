package ioc

import (
	"fmt"

	"github.com/mymikasa/prompthub/internal/repo/dao/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.User{},
		&model.Workspace{},
		&model.WorkspaceMember{},
		&model.Prompt{},
		&model.PromptVersion{},
		&model.PromptVariable{},
		&model.PromptTestCase{},
		&model.PromptRun{},
		&model.Tag{},
		&model.PromptTag{},
		&model.ProviderConfig{},
	)
	if err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}
