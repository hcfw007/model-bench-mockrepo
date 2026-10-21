package bulk_send_v1

// UseCase 群发1业务编排。
type UseCase struct {
	repo     JobRepo
	crm      CRMRepo
	identity IdentityRepo
	executor Executor
	clock    Clock
	tagNames TagNameResolver
}

// Deps 构造依赖，显式列出，便于 wire 生成。
type Deps struct {
	Repo     JobRepo
	CRM      CRMRepo
	Identity IdentityRepo
	Executor Executor
	Clock    Clock
	TagNames TagNameResolver
}

// NewUseCase 构造业务对象。
func NewUseCase(d Deps) *UseCase {
	if d.Clock == nil {
		d.Clock = systemClock{}
	}
	return &UseCase{
		repo:     d.Repo,
		crm:      d.CRM,
		identity: d.Identity,
		executor: d.Executor,
		clock:    d.Clock,
		tagNames: d.TagNames,
	}
}

type systemClock struct{}

func (systemClock) Now() int64 { return nowMillis() }
