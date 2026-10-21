package bulk_send_v1

import (
	"context"
	"fmt"
	"time"
)

// CreateJobInput 创建群发1任务。
type CreateJobInput struct {
	OrgID      string
	GroupID    string
	BotID      string
	Name       string
	NotifyMode NotifyMode
	MarkTag    *BindingMarkTagInfo
	Targets    []Target
}

// CreateJob 校验并落库任务。
//
// 本函数只负责「保存配置」：NotifyMode 与 BindingMarkTagInfo 会原样落库，
// 供详情页展示与后续执行阶段读取。它不触发任何下发或打标动作。
func (uc *UseCase) CreateJob(ctx context.Context, in CreateJobInput) (*Job, error) {
	if in.OrgID == "" || in.BotID == "" || in.Name == "" {
		return nil, fmt.Errorf("%w: orgId/botId/name required", ErrInvalidInput)
	}
	if in.NotifyMode != NotifyModeManual && in.NotifyMode != NotifyModeNoConfirm {
		return nil, fmt.Errorf("%w: unknown notifyMode %d", ErrInvalidInput, in.NotifyMode)
	}
	if len(in.Targets) == 0 {
		return nil, fmt.Errorf("%w: targets required", ErrInvalidInput)
	}
	if in.MarkTag != nil {
		if in.MarkTag.TagType != MarkTagTypeWecom && in.MarkTag.TagType != MarkTagTypeSystem {
			return nil, fmt.Errorf("%w: unknown tagType %d", ErrInvalidInput, in.MarkTag.TagType)
		}
		if len(in.MarkTag.TagIDs) == 0 {
			return nil, fmt.Errorf("%w: tagIds required", ErrInvalidInput)
		}
	}

	now := time.UnixMilli(uc.clock.Now())
	job := &Job{
		ID:             newID("job"),
		OrgID:          in.OrgID,
		GroupID:        in.GroupID,
		BotID:          in.BotID,
		Name:           in.Name,
		Status:         JobStatusCreated,
		NotifyMode:     in.NotifyMode,
		BindingMarkTag: in.MarkTag,
		Targets:        in.Targets,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := uc.repo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("create job: %w", err)
	}
	return job, nil
}

var idSeq uint64

// newID 生成带前缀的简易 ID，仅供 mock 使用。
func newID(prefix string) string {
	idSeq++
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixNano(), idSeq)
}
