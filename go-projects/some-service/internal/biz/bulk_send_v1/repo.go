package bulk_send_v1

import (
	"context"
	"errors"
)

var (
	ErrJobNotFound   = errors.New("job not found")
	ErrInvalidStatus = errors.New("invalid job status")
	ErrCASFailed     = errors.New("compare-and-swap failed")
	ErrInvalidInput  = errors.New("invalid input")
)

// JobRepo 群发1任务持久化。
type JobRepo interface {
	Create(ctx context.Context, job *Job) error
	Get(ctx context.Context, jobID string) (*Job, error)
	Update(ctx context.Context, job *Job) error
	// CompareAndSwapStatus 原子改状态，用于避免并发确认/取消互相覆盖。
	CompareAndSwapStatus(ctx context.Context, jobID string, from, to JobStatus) (bool, error)
	ListByStatus(ctx context.Context, status JobStatus, limit int) ([]*Job, error)
}

// CRMRepo 客户标签能力。打标接口本身幂等，重复调用不产生副作用。
type CRMRepo interface {
	// MarkTag 给 orgID 下 externalUserID 对应的客户打一个企微标签。
	MarkTag(ctx context.Context, orgID, externalUserID, tagID string) error
}

// IdentityRepo 托管账号与平台成员的身份关系。
type IdentityRepo interface {
	// GetBotUserIDs 返回 botID -> bot_user.userId。
	// 注意：bot_user.userId 是托管账号维度身份，不是平台登录成员 ID。
	GetBotUserIDs(ctx context.Context, botIDs []string) (map[string]string, error)
}

// Executor 向托管账号下发执行指令。
type Executor interface {
	PublishExecute(ctx context.Context, job *Job) error
}

// Clock 便于测试注入。
type Clock interface {
	Now() int64 // unix millis
}
