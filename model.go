package waterlevel

import (
	"context"
	"time"
)

// Repository is a Watermarks repo.
type Repository interface {
	Load(ctx context.Context, taskName string) (watermark uint64, err error)
	Save(ctx context.Context, taskName string, watermark uint64) error
}

type Watermark struct {
	ID        int64     `gorm:"<-:create;column:id;primaryKey;autoIncrement"`
	TaskName  string    `gorm:"<-:create;column:task_name;type:varchar(255);not null;comment:'任务的名称，用于区分不同的任务'"`
	Watermark uint64    `gorm:"column:watermark;type:bigint;not null;comment:'记录任务进度的水位线，其具体意义根据实际业务决定'"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime:milli;comment:'记录创建时间'"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime:milli;comment:'记录最后一次更新时间'"`
}

func (c Watermark) TableName() string {
	return "watermark"
}
