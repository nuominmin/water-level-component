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
	Load(ctx context.Context) (waterLevel uint64, err error)
	Save(ctx context.Context, waterLevel uint64) error
}

type manager struct {
	db       *gorm.DB
	taskName string
}

func NewManager(db *gorm.DB, taskName string) (Manager, error) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := registeredTasks[taskName]; exists {
		return nil, ErrTaskAlreadyRegistered
	}

	registeredTasks[taskName] = struct{}{}
	return &manager{
		db:       db,
		taskName: taskName,
	}, nil
}

func (r *manager) Load(ctx context.Context) (waterLevel uint64, err error) {
	var wm *WaterLevel
	if wm, err = gormhelper.First[WaterLevel](r.db, ctx,
		gormhelper.WithWhere("task_name = ?", r.taskName),
		gormhelper.WithIgnore(),
	); err != nil {
		return 0, fmt.Errorf("error loading water level for task %s: %w", r.taskName, err)
	}

	if wm == nil {
		return 0, nil
	}

	return wm.WaterLevel, nil
}

func (r *manager) Save(ctx context.Context, waterLevel uint64) error {
	defaultData := &WaterLevel{
		TaskName:   r.taskName,
		WaterLevel: waterLevel,
	}

	updateData := map[string]interface{}{
		"water_level": waterLevel,
	}

	err := gormhelper.UpdateOrCreate[WaterLevel](r.db, ctx, defaultData, updateData,
		gormhelper.WithWhere("task_name = ?", r.taskName))
	if err != nil {
		return fmt.Errorf("error saving water level for task %s: %w", r.taskName, err)
	}

	return nil
}
