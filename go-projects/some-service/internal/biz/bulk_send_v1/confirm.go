package bulk_send_v1

import (
	"context"
	"fmt"
	"time"
)

// Confirm 成员确认后下发执行。
func (uc *UseCase) Confirm(ctx context.Context, jobID string) error {
	job, err := uc.repo.Get(ctx, jobID)
	if err != nil {
		return fmt.Errorf("get job: %w", err)
	}
	if job.Status != JobStatusWaitConfirm && job.Status != JobStatusCreated {
		return fmt.Errorf("%w: confirm in status %d", ErrInvalidStatus, job.Status)
	}
	// 取消与确认可能并发，必须用 CAS 占位，否则已取消任务会被重新推进。
	ok, err := uc.repo.CompareAndSwapStatus(ctx, jobID, job.Status, JobStatusSending)
	if err != nil {
		return fmt.Errorf("cas job: %w", err)
	}
	if !ok {
		return ErrCASFailed
	}

	job.Status = JobStatusSending
	job.UpdatedAt = time.UnixMilli(uc.clock.Now())
	if err := uc.executor.PublishExecute(ctx, job); err != nil {
		// 下发失败回滚状态，允许下次重试。
		_, _ = uc.repo.CompareAndSwapStatus(ctx, jobID, JobStatusSending, job.Status)
		return fmt.Errorf("publish execute: %w", err)
	}
	return nil
}

// HandleCardEvent 处理成员在企微手机端点击「确认发送」卡片。
func (uc *UseCase) HandleCardEvent(ctx context.Context, jobID, memberID string, confirmed bool) error {
	if !confirmed {
		// 成员点了「暂不发送」，保持待确认，等下一轮提醒或超时。
		return nil
	}
	return uc.Confirm(ctx, jobID)
}

// Cancel 取消任务。终态任务不可取消。
func (uc *UseCase) Cancel(ctx context.Context, jobID string) error {
	job, err := uc.repo.Get(ctx, jobID)
	if err != nil {
		return fmt.Errorf("get job: %w", err)
	}
	if job.Status.Terminal() {
		return fmt.Errorf("%w: cancel in status %d", ErrInvalidStatus, job.Status)
	}
	ok, err := uc.repo.CompareAndSwapStatus(ctx, jobID, job.Status, JobStatusCanceled)
	if err != nil {
		return fmt.Errorf("cas job: %w", err)
	}
	if !ok {
		return ErrCASFailed
	}
	return nil
}
