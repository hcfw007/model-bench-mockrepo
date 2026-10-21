package data

import (
	"context"
	"sync"

	bo "mockrepo/go-projects/some-service/internal/biz/bulk_send_v1"
)

// MemoryExecutor 记录下发指令，不真正发消息。
type MemoryExecutor struct {
	mu        sync.Mutex
	published []string
}

func NewMemoryExecutor() *MemoryExecutor { return &MemoryExecutor{} }

func (e *MemoryExecutor) PublishExecute(_ context.Context, job *bo.Job) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.published = append(e.published, job.ID)
	return nil
}

// Published 返回已下发过的任务 ID。
func (e *MemoryExecutor) Published() []string {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]string, len(e.published))
	copy(out, e.published)
	return out
}
