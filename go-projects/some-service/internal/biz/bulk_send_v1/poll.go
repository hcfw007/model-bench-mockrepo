package bulk_send_v1

import (
	"context"
	"fmt"
	"time"
)

// SyncJob 从企微拉取执行明细并推进任务状态。
//
// 群发1任务状态只由轮询推进：WAIT_CONFIRM/SENDING -> FINISHED 都在这里判定。
func (uc *UseCase) SyncJob(ctx context.Context, jobID string) error {
	job, err := uc.repo.Get(ctx, jobID)
	if err != nil {
		return fmt.Errorf("get job: %w", err)
	}
	if job.Status.Terminal() {
		return nil
	}
	if job.Status != JobStatusSending {
		return nil
	}

	// 真实实现会调用企微接口拉取每个对象的发送结果；mock 直接按 targets 统计。
	job.SentCount = int32(len(job.SentExternalIDs()))
	job.ResultSynced = true
	job.UpdatedAt = time.UnixMilli(uc.clock.Now())

	if job.SentCount > 0 {
		job.Status = JobStatusFinished
	}
	return uc.repo.Update(ctx, job)
}

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
