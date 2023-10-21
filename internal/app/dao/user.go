package dao

type User interface {
	CreateUser(user *UserInfo) error
	UpdateUser(user *UserInfo) error
	DeleteUser(user *UserInfo) error
	QueryUser(user *UserInfo) (UserInfo, error)
	QueryUsers(page int, pageNumber int) ([]*UserInfo, error)
}

func GetUserDao() User {
	return getUserDao()
}
