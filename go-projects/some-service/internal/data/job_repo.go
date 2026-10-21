package data

import (
	"context"
	"sync"

	bo "mockrepo/go-projects/some-service/internal/biz/bulk_send_v1"
)

// MemoryJobRepo 进程内任务仓储，供 mock 与测试使用。
type MemoryJobRepo struct {
	mu   sync.RWMutex
	jobs map[string]*bo.Job
}

func NewMemoryJobRepo() *MemoryJobRepo {
	return &MemoryJobRepo{jobs: make(map[string]*bo.Job)}
}

func (r *MemoryJobRepo) Create(_ context.Context, job *bo.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *job
	r.jobs[job.ID] = &cp
	return nil
}

func (r *MemoryJobRepo) Get(_ context.Context, jobID string) (*bo.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	j, ok := r.jobs[jobID]
	if !ok {
		return nil, bo.ErrJobNotFound
	}
	cp := *j
	return &cp, nil
}

func (r *MemoryJobRepo) Update(_ context.Context, job *bo.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.jobs[job.ID]; !ok {
		return bo.ErrJobNotFound
	}
	cp := *job
	r.jobs[job.ID] = &cp
	return nil
}

func (r *MemoryJobRepo) CompareAndSwapStatus(_ context.Context, jobID string, from, to bo.JobStatus) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	j, ok := r.jobs[jobID]
	if !ok {
		return false, bo.ErrJobNotFound
	}
	if j.Status != from {
		return false, nil
	}
	j.Status = to
	return true, nil
}

func (r *MemoryJobRepo) ListByStatus(_ context.Context, status bo.JobStatus, limit int) ([]*bo.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*bo.Job, 0, limit)
	for _, j := range r.jobs {
		if j.Status != status {
			continue
		}
		cp := *j
		out = append(out, &cp)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}
