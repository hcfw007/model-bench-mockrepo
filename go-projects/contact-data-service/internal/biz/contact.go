package biz

import (
	"context"
	"sync"
)

// Contact 客户档案。tags 是 crm 侧的标签 ID 列表。
type Contact struct {
	ID             string
	OrgID          string
	ExternalUserID string
	Tags           []string
}

// ContactRepo 客户档案读写。
type ContactRepo interface {
	GetByExternalID(ctx context.Context, orgID, externalUserID string) (*Contact, error)
	AppendTag(ctx context.Context, contactID, tagID string) error
}

// TagService 打标业务。打标必须幂等：同一 (contact, tag) 重复写不产生重复项。
type TagService struct {
	repo ContactRepo
}

func NewTagService(repo ContactRepo) *TagService {
	return &TagService{repo: repo}
}

// Mark 给客户打标签。
func (s *TagService) Mark(ctx context.Context, orgID, externalUserID, tagID string) error {
	c, err := s.repo.GetByExternalID(ctx, orgID, externalUserID)
	if err != nil {
		return err
	}
	return s.repo.AppendTag(ctx, c.ID, tagID)
}

// MemoryContactRepo 进程内客户档案。
type MemoryContactRepo struct {
	mu       sync.Mutex
	contacts map[string]*Contact // key: orgID + "/" + externalUserID
}

func NewMemoryContactRepo() *MemoryContactRepo {
	return &MemoryContactRepo{contacts: make(map[string]*Contact)}
}

func (r *MemoryContactRepo) Put(c *Contact) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *c
	cp.Tags = append([]string(nil), c.Tags...)
	r.contacts[c.OrgID+"/"+c.ExternalUserID] = &cp
}

func (r *MemoryContactRepo) GetByExternalID(_ context.Context, orgID, externalUserID string) (*Contact, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.contacts[orgID+"/"+externalUserID]
	if !ok {
		return nil, ErrContactNotFound
	}
	cp := *c
	cp.Tags = append([]string(nil), c.Tags...)
	return &cp, nil
}

func (r *MemoryContactRepo) AppendTag(_ context.Context, contactID, tagID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, c := range r.contacts {
		if c.ID != contactID {
			continue
		}
		for _, t := range c.Tags {
			if t == tagID {
				return nil // 幂等
			}
		}
		c.Tags = append(c.Tags, tagID)
		return nil
	}
	return ErrContactNotFound
}
