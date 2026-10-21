package domain

import (
	"context"
	"errors"
)

const memberPageSize = 200

// MemberRepo 成员组关系读取。
type MemberRepo interface {
	// ListMembers 按 offset/limit 分页返回成员 ID。
	ListMembers(ctx context.Context, groupID string, offset, limit int) (ids []string, err error)
	CountMembers(ctx context.Context, groupID string) (int, error)
}

// MemberGroupService 成员组展开。
type MemberGroupService struct {
	repo MemberRepo
}

func NewMemberGroupService(repo MemberRepo) *MemberGroupService {
	return &MemberGroupService{repo: repo}
}

// ExpandMembers 展开成员组的全部成员 ID。
//
// 群发范围解析依赖这个结果，漏人会直接导致「通知范围」缺陷。
func (s *MemberGroupService) ExpandMembers(ctx context.Context, groupID string) ([]string, error) {
	ids, err := s.repo.ListMembers(ctx, groupID, 0, memberPageSize)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// IsSuperAdmin 判断是否超管。超管继承其可见范围内的全部权限。
func IsSuperAdmin(roles []string) bool {
	for _, r := range roles {
		if r == "sys_admin" || r == "super_admin" {
			return true
		}
	}
	return false
}

var ErrGroupNotFound = errors.New("member group not found")
