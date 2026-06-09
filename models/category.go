package models

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/lib/pq"
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/publish"
	"github.com/r0vx/x/oss"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model

	Name     string
	Products pq.StringArray `gorm:"type:text[]"`
	Position int
	publish.Status
	publish.Schedule
	publish.Version
}

func (c *Category) PrimarySlug() string {
	return fmt.Sprintf("%v_%v", c.ID, c.Version.Version)
}

func (c *Category) PrimaryColumnValuesBySlug(slug string) map[string]string {
	segs := strings.Split(slug, "_")
	if len(segs) != 2 {
		panic(presets.ErrNotFound("wrong slug"))
	}

	_, err := cast.ToInt64E(segs[0])
	if err != nil {
		panic(presets.ErrNotFound(fmt.Sprintf("wrong slug %q: %v", slug, err)))
	}

	return map[string]string{
		"id":      segs[0],
		"version": segs[1],
	}
}

func (c *Category) GetPublishActions(ctx context.Context, db *gorm.DB, storage oss.StorageInterface) (actions []*publish.PublishAction, err error) {
	return
}

func (c *Category) GetUnPublishActions(ctx context.Context, db *gorm.DB, storage oss.StorageInterface) (actions []*publish.PublishAction, err error) {
	return
}

func (c *Category) PermissionRN() []string {
	return []string{"categories", strconv.Itoa(int(c.ID)), c.Version.Version}
}

// SubCategory 二级分类（用于演示 ButtonForParent 子项排序）
type SubCategory struct {
	gorm.Model
	CategoryID uint
	Name       string
	Position   int
}
