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

	if job.SentCount >= int32(len(job.Targets)) {
		job.Status = JobStatusFinished
	}
	return uc.repo.Update(ctx, job)
}
