package waterlevel

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuominmin/gorm-helper"
	"gorm.io/gorm"
	"sync"
)

var registeredTasks = make(map[string]uint64)
var mu sync.Mutex
var ErrTaskAlreadyRegistered = errors.New("task already registered")

type Manager interface {
	Load() (waterLevel uint64)
	Save(ctx context.Context, waterLevel uint64) error
}

type manager struct {
	db       *gorm.DB
	taskName string
}

func NewManager(ctx context.Context, db *gorm.DB, taskName string) (Manager, error) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := registeredTasks[taskName]; exists {
		return nil, ErrTaskAlreadyRegistered
	}

	wl, err := gormhelper.First[WaterLevel](db, ctx,
		gormhelper.WithWhere("task_name = ?", taskName),
		gormhelper.WithIgnore(),
	)
	if err != nil {
		return nil, fmt.Errorf("error loading water level for task %s: %w", taskName, err)
	}

	var waterLevel uint64
	if wl != nil {
		waterLevel = wl.WaterLevel
	}

	registeredTasks[taskName] = waterLevel
	return &manager{
		db:       db,
		taskName: taskName,
	}, nil
}

func (r *manager) Load() (waterLevel uint64) {
	mu.Lock()
	defer mu.Unlock()
	return registeredTasks[r.taskName]
}

func (r *manager) Save(ctx context.Context, waterLevel uint64) error {
	mu.Lock()
	defer mu.Unlock()

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

	registeredTasks[r.taskName] = waterLevel
	return nil
}
