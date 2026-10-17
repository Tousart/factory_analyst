package analyst

type Service struct {
	analystRepo AnalystRepository
}

func NewService(
	analystRepo AnalystRepository,
) *Service {
	return &Service{
		analystRepo: analystRepo,
	}
}
