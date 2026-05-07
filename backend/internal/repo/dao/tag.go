package dao

import (
	"context"

	"github.com/mymikasa/prompthub/internal/repo/dao/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TagDAO interface {
	FindByWorkspace(ctx context.Context, workspaceID int64) ([]*model.Tag, error)
	FindByNames(ctx context.Context, workspaceID int64, names []string) ([]*model.Tag, error)
	FindOrCreateByNames(ctx context.Context, workspaceID int64, names []string) ([]*model.Tag, error)
	FindByPromptID(ctx context.Context, promptID int64) ([]*model.Tag, error)
	ReplacePromptTags(ctx context.Context, promptID int64, tagIDs []int64) error
	DeletePromptTags(ctx context.Context, promptID int64) error
}

type tagDAO struct {
	db *gorm.DB
}

func NewTagDAO(db *gorm.DB) TagDAO {
	return &tagDAO{db: db}
}

func (d *tagDAO) FindByWorkspace(ctx context.Context, workspaceID int64) ([]*model.Tag, error) {
	var items []*model.Tag
	if err := d.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order("name ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (d *tagDAO) FindByNames(ctx context.Context, workspaceID int64, names []string) ([]*model.Tag, error) {
	var items []*model.Tag
	if len(names) == 0 {
		return items, nil
	}
	if err := d.db.WithContext(ctx).Where("workspace_id = ? AND name IN ?", workspaceID, names).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (d *tagDAO) FindOrCreateByNames(ctx context.Context, workspaceID int64, names []string) ([]*model.Tag, error) {
	if len(names) == 0 {
		return nil, nil
	}

	existing, err := d.FindByNames(ctx, workspaceID, names)
	if err != nil {
		return nil, err
	}

	existingSet := make(map[string]struct{}, len(existing))
	for _, t := range existing {
		existingSet[t.Name] = struct{}{}
	}

	var newTags []*model.Tag
	for _, name := range names {
		if _, ok := existingSet[name]; !ok {
			m := &model.Tag{WorkspaceID: workspaceID, Name: name}
			newTags = append(newTags, m)
		}
	}

	if len(newTags) > 0 {
		if err := d.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&newTags).Error; err != nil {
			return nil, err
		}
		return d.FindByNames(ctx, workspaceID, names)
	}

	return existing, nil
}

func (d *tagDAO) FindByPromptID(ctx context.Context, promptID int64) ([]*model.Tag, error) {
	var items []*model.Tag
	if err := d.db.WithContext(ctx).
		Joins("JOIN prompt_tags ON prompt_tags.tag_id = tags.id").
		Where("prompt_tags.prompt_id = ?", promptID).
		Order("tags.name ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (d *tagDAO) ReplacePromptTags(ctx context.Context, promptID int64, tagIDs []int64) error {
	if err := d.DeletePromptTags(ctx, promptID); err != nil {
		return err
	}
	if len(tagIDs) == 0 {
		return nil
	}
	pts := make([]model.PromptTag, len(tagIDs))
	for i, tid := range tagIDs {
		pts[i] = model.PromptTag{PromptID: promptID, TagID: tid}
	}
	return d.db.WithContext(ctx).Create(&pts).Error
}

func (d *tagDAO) DeletePromptTags(ctx context.Context, promptID int64) error {
	return d.db.WithContext(ctx).Where("prompt_id = ?", promptID).Delete(&model.PromptTag{}).Error
}
