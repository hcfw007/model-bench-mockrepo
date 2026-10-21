package bulk_send_v1

import (
	"context"
	"fmt"
)

// Patrol 扫描长时间停留在 SENDING 的任务，触发一次同步。
func (uc *UseCase) Patrol(ctx context.Context) error {
	jobs, err := uc.repo.ListByStatus(ctx, JobStatusSending, 100)
	if err != nil {
		return fmt.Errorf("list sending: %w", err)
	}
	for _, j := range jobs {
		if err := uc.SyncJob(ctx, j.ID); err != nil {
			// 单个任务失败不影响其他任务。
			j.ErrorMsg = err.Error()
			_ = uc.repo.Update(ctx, j)
		}
	}
	return nil
}
