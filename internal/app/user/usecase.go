package user

// UseCase ...
type UseCase interface {
	FindUser(email string) (int64, error)
	InsertUser(user user) error
	GetUserByEmail(email string, user *user) error
	GetUserByID(userID string, user *user) error
	UpdateUser(userID string, user *user) error
	UpdateFullName(userID string, user *user) error
	UpdatePassword(userID string, user *user) error
}

type useCase struct {
	repo Repository
}

// NewUseCase ...
func NewUseCase() UseCase {
	return &useCase{
		repo: NewRepository(),
	}
}

func (instance *useCase) FindUser(email string) (int64, error) {
	count, err := instance.repo.FindUser(email)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (instance *useCase) InsertUser(anUser user) error {
	err := instance.repo.InsertUser(anUser)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) GetUserByEmail(email string, anUser *user) error {
	err := instance.repo.GetUserByEmail(email, anUser)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) GetUserByID(userID string, anUser *user) error {
	err := instance.repo.GetUserByID(userID, anUser)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) UpdateFullName(userID string, anUser *user) error {
	err := instance.repo.UpdateFullName(userID, anUser)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) UpdatePassword(userID string, anUser *user) error {
	err := instance.repo.UpdatePassword(userID, anUser)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) UpdateUser(userID string, anUser *user) error {
	err := instance.repo.UpdateUser(userID, anUser)
	if err != nil {
		return err
	}
	return nil
}
