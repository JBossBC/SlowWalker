package dao

type Department interface {
	QueryDepartments() ([]*DepartmentInfo, error)
	QueryDepartment(department *DepartmentInfo) (DepartmentInfo, error)
	DeleteDepartment(departmentInfo *DepartmentInfo) error
	UpdateDepartment(departmentInfo DepartmentInfo) error
	CreateDepartment(departmentInfo DepartmentInfo) error
}

func GetDepartmentDao() Department {
	return getDepartmentDao()
}
