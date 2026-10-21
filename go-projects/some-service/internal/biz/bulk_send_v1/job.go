package bulk_send_v1

import "time"

// JobStatus 群发1任务状态。与 api/bulk_send_v1.pb.go 中的枚举保持一致。
type JobStatus int32

const (
	JobStatusCreated     JobStatus = 1
	JobStatusWaitConfirm JobStatus = 2
	JobStatusSending     JobStatus = 3
	JobStatusFinished    JobStatus = 4
	JobStatusCanceled    JobStatus = 5
	JobStatusError       JobStatus = 6
	JobStatusSubmitting  JobStatus = 7
)

func (s JobStatus) Terminal() bool {
	return s == JobStatusFinished || s == JobStatusCanceled || s == JobStatusError
}

// NotifyMode 群发完成后的确认方式。
type NotifyMode int32

const (
	// NotifyModeManual 需要托管账号成员在企微手机端手动确认。
	NotifyModeManual NotifyMode = 0
	// NotifyModeNoConfirm 免确认：平台侧自动代成员确认后下发。
	NotifyModeNoConfirm NotifyMode = 1
)

// MarkTagType 自动标签的类型。
type MarkTagType int32

const (
	MarkTagTypeWecom  MarkTagType = 3
	MarkTagTypeSystem MarkTagType = 9
)

// BindingMarkTagInfo 创建任务时绑定的自动标签配置。
type BindingMarkTagInfo struct {
	TagType MarkTagType `json:"tagType"`
	TagIDs  []string    `json:"tagIds"`
}

// Target 单个触达对象。
type Target struct {
	ExternalUserID string
	// Status 0=待发送 1=已发送 2=失败
	Status int32
}

// Job 群发1任务。
type Job struct {
	ID         string
	OrgID      string
	GroupID    string
	BotID      string
	Name       string
	Status     JobStatus
	NotifyMode NotifyMode

	BindingMarkTag *BindingMarkTagInfo
	Targets        []Target
	SentCount      int32
	ResultSynced   bool
	ErrorMsg       string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// SentExternalIDs 返回实际发送成功的对象。打标只能作用于这批对象。
func (j *Job) SentExternalIDs() []string {
	out := make([]string, 0, len(j.Targets))
	for _, t := range j.Targets {
		if t.Status == 1 {
			out = append(out, t.ExternalUserID)
		}
	}
	return out
}

// HasMarkTag 任务是否绑定了可执行的自动标签。
func (j *Job) HasMarkTag() bool {
	return j.BindingMarkTag != nil && len(j.BindingMarkTag.TagIDs) > 0
}
