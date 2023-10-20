package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task interface {
	CreateTask(task TaskInfo) error
	UpdateTask(taskid primitive.ObjectID, fields bson.M) error
	DeleteTask(task TaskInfo) error
}

func GetTaskDao() Task {
	return getTaskDao()
}
