package accounts

type Service interface {
	Create(input CreateAccountInput) Account
	GetAll() []Account
	GetByID(id int64) (Account, bool)
	ExistsByID(id int64) bool
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(input CreateAccountInput) Account {
	return s.repo.Create(input)
}

func (s *service) GetAll() []Account {
	return s.repo.GetAll()
}

func (s *service) GetByID(id int64) (Account, bool) {
	return s.repo.GetByID(id)
}

func (s *service) ExistsByID(id int64) bool {
	return s.repo.ExistsByID(id)
}
