package accounts

type Service interface {
	Create(input CreateAccountInput) (Account, error)
	GetAll() ([]Account, error)
	GetByID(id int64) (Account, bool, error)
	ExistsByID(id int64) (bool, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(input CreateAccountInput) (Account, error) {
	return s.repo.Create(input)
}

func (s *service) GetAll() ([]Account, error) {
	return s.repo.GetAll()
}

func (s *service) GetByID(id int64) (Account, bool, error) {
	return s.repo.GetByID(id)
}

func (s *service) ExistsByID(id int64) (bool, error) {
	return s.repo.ExistsByID(id)
}
