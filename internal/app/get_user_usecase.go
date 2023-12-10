package app

type GetUserUseCase struct {
	UserRepository UserRepository
}

type GetUserInput struct {
	Id UserId
}

func (u *GetUserUseCase) Execute(i GetUserInput) *User {
	return u.UserRepository.FindById(i.Id)
}
