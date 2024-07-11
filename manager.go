package waterlevel

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuominmin/gorm-helper"
	"gorm.io/gorm"
	"sync"
)

var registeredTasks = make(map[string]struct{})
var mu sync.Mutex
var ErrTaskAlreadyRegistered = errors.New("task already registered")

type Manager interface {
	Load(ctx context.Context, taskName string) (watermark uint64, err error)
	Save(ctx context.Context, taskName string, watermark uint64) error
}

type manager struct {
	db *gorm.DB
}

func NewManager(db *gorm.DB, taskName string) (Manager, error) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := registeredTasks[taskName]; exists {
		return nil, ErrTaskAlreadyRegistered
	}

	registeredTasks[taskName] = struct{}{}
	return &manager{
		db: db,
	}, nil
}

func (r *manager) Load(ctx context.Context, taskName string) (watermark uint64, err error) {
	var wm *Watermark
	if wm, err = gormhelper.First[Watermark](r.db, ctx,
		gormhelper.WithWhere("task_name = ?", taskName),
		gormhelper.WithIgnore(),
	); err != nil {
		return 0, fmt.Errorf("error loading watermark for task %s: %w", taskName, err)
	}

	if wm == nil {
		return 0, nil
	}

	return wm.Watermark, nil
}

func (r *manager) Save(ctx context.Context, taskName string, watermark uint64) error {
	defaultData := &Watermark{
		TaskName:  taskName,
		Watermark: watermark,
	}

	updateData := map[string]interface{}{
		"watermark": watermark,
	}

	err := gormhelper.UpdateOrCreate[Watermark](r.db, ctx, defaultData, updateData,
		gormhelper.WithWhere("task_name = ?", taskName))
	if err != nil {
		return fmt.Errorf("error saving watermark for task %s: %w", taskName, err)
	}

	return nil
}
