package data

import (
	"context"
	"sync"
)

// MemoryIdentityRepo 进程内身份实现。
//
// 刻意复刻 test 环境的已知数据形态：bot_user.userId 存的是创建时兜底生成的
// ObjectId，并不等于平台登录成员 ID。真实成员 ID 在 im_user 体系里。
// 因此「按 bot_user.userId 查通知」会查不到——这是数据形态问题，
// 不能直接断定成「通知范围漏人」的代码缺陷。
type MemoryIdentityRepo struct {
	mu       sync.RWMutex
	fallback map[string]string // botID -> 兜底 ObjectId
	realUser map[string]string // botID -> 真实平台成员 ID
}

func NewMemoryIdentityRepo() *MemoryIdentityRepo {
	return &MemoryIdentityRepo{
		fallback: make(map[string]string),
		realUser: make(map[string]string),
	}
}

// Seed 登记一个托管账号的身份关系。
func (r *MemoryIdentityRepo) Seed(botID, fallbackID, realMemberID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fallback[botID] = fallbackID
	r.realUser[botID] = realMemberID
}

// GetBotUserIDs 返回 botID -> bot_user.userId（兜底 ID），与真实实现一致。
func (r *MemoryIdentityRepo) GetBotUserIDs(_ context.Context, botIDs []string) (map[string]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]string, len(botIDs))
	for _, id := range botIDs {
		if v, ok := r.fallback[id]; ok {
			out[id] = v
		}
	}
	return out, nil
}

// RealMemberID 仅供测试断言使用：拿到该 bot 真实对应的平台成员 ID。
func (r *MemoryIdentityRepo) RealMemberID(botID string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.realUser[botID]
	return v, ok
}
