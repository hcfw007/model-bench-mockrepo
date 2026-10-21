package service

import (
	"context"

	bo "mockrepo/go-projects/some-service/internal/biz/bulk_send_v1"
)

// BroadcastOfficialService 对外暴露的群发1接口。
// 真实实现由 proto/bulk_send_v1.proto 生成，这里保留等价的方法签名。
type BroadcastOfficialService struct {
	uc *bo.UseCase
}

func NewBroadcastOfficialService(uc *bo.UseCase) *BroadcastOfficialService {
	return &BroadcastOfficialService{uc: uc}
}

func (s *BroadcastOfficialService) CreateJob(ctx context.Context, in bo.CreateJobInput) (*bo.Job, error) {
	return s.uc.CreateJob(ctx, in)
}

func (s *BroadcastOfficialService) Confirm(ctx context.Context, jobID string) error {
	return s.uc.Confirm(ctx, jobID)
}

func (s *BroadcastOfficialService) Cancel(ctx context.Context, jobID string) error {
	return s.uc.Cancel(ctx, jobID)
}

func (s *BroadcastOfficialService) SyncJob(ctx context.Context, jobID string) error {
	return s.uc.SyncJob(ctx, jobID)
}

func (s *BroadcastOfficialService) GetJobTagInfos(ctx context.Context, jobID string) ([]bo.TagInfo, error) {
	return s.uc.GetJobTagInfos(ctx, jobID)
}
