package data

import (
	"context"
	"sync"
)

// MarkTagCall 一次打标调用记录，供测试断言是否真的执行了打标。
type MarkTagCall struct {
	OrgID          string
	ExternalUserID string
	TagID          string
}

// MemoryCRMRepo 进程内 crm 打标实现。crm 打标接口本身幂等。
type MemoryCRMRepo struct {
	mu    sync.Mutex
	calls []MarkTagCall
	fail  error
}

func NewMemoryCRMRepo() *MemoryCRMRepo { return &MemoryCRMRepo{} }

func (r *MemoryCRMRepo) MarkTag(_ context.Context, orgID, externalUserID, tagID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.fail != nil {
		return r.fail
	}
	r.calls = append(r.calls, MarkTagCall{OrgID: orgID, ExternalUserID: externalUserID, TagID: tagID})
	return nil
}

// Calls 返回已发生的打标调用快照。
func (r *MemoryCRMRepo) Calls() []MarkTagCall {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]MarkTagCall, len(r.calls))
	copy(out, r.calls)
	return out
}

// SetFailure 注入打标失败，用于验证失败不阻塞任务收尾。
func (r *MemoryCRMRepo) SetFailure(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fail = err
}
