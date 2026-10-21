package bulk_send_v1

import (
	"context"
	"fmt"
)

// TagInfo 详情页展示用的标签信息。
type TagInfo struct {
	TagID   string
	TagName string
	TagType MarkTagType
}

// TagNameResolver 按标签库反查标签名。由 data 层注入。
type TagNameResolver func(orgID, tagID string) (string, error)

// GetJobTagInfos 仅供列表/详情页反查标签名展示。
//
// 它只读取 job.BindingMarkTag 并翻译成可读名称，
// 不触发任何打标动作，也不校验标签是否真的存在。
func (uc *UseCase) GetJobTagInfos(ctx context.Context, jobID string) ([]TagInfo, error) {
	job, err := uc.repo.Get(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("get job: %w", err)
	}
	if !job.HasMarkTag() {
		return nil, nil
	}
	out := make([]TagInfo, 0, len(job.BindingMarkTag.TagIDs))
	for _, id := range job.BindingMarkTag.TagIDs {
		name := id
		if uc.tagNames != nil {
			if n, err := uc.tagNames(job.OrgID, id); err == nil && n != "" {
				name = n
			}
		}
		out = append(out, TagInfo{TagID: id, TagName: name, TagType: job.BindingMarkTag.TagType})
	}
	return out, nil
}
