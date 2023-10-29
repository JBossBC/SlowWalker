package dao

type User interface {
	CreateUser(user *UserInfo) error
	UpdateUser(user *UserInfo) error
	DeleteUser(user *UserInfo) error
	FilterUsers(userInfo *UserInfo) ([]*UserInfo, error)
	QueryUsers(filterTempalte *UserInfo) ([]*UserInfo, error)
}

func GetUserDao() User {
	return getUserDao()
}
