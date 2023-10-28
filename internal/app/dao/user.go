package dao

type User interface {
	CreateUser(user *UserInfo) error
	UpdateUser(user *UserInfo) error
	DeleteUser(user *UserInfo) error
	QueryUser(user *UserInfo) (UserInfo, error)
	QueryUsers(page int, pageNumber int) ([]*UserInfo, error)
	FilterUsers(filterTempalte *UserFilterTemplate) ([]*UserInfo, error)
}

type UserFilterTemplate struct {
	Username    string `json:"username" bson:"username"`
	PhoneNumber string `json:"phoneNumber" bson:"phoneNumber"`
	Start       int64  `json:"start" bson:"start"`
	End         int64  `json:"end" bson:"end"`
	Department  string `json:"department" bson:"department"`
	RealName    string `json:"realName" bson:"realName"`
	Authority   string `json:"authority" bson:"authority"`
	Page        int64  `json:"page"`
	PageNumber  int64  `json:"pageNumber"`
}

func GetUserDao() User {
	return getUserDao()
}
