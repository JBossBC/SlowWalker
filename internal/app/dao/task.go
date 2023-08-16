package dao

import (
	"context"
	"log"
	"replite_web/internal/app/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const default_task_times = 3 * time.Second

const default_task_cache_times = 3 * time.Minute

type TaskState int

const (
	Ongoing TaskState = iota
	Failed
	Success
)

type Task struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	State    TaskState          `json:"state" bson:"state"`
	Message  string             `json:"message" bson:"message"`
	PlatForm PlatForm           `json:"platform" bson:"platform"`
	Operate  Operate            `json:"operate" bson:"operate"`
}

const taskModelName = "task"

func getTaskCollection() *mongo.Collection {
	return getMongoConn().Collection(config.CollectionConfig.Get(taskModelName).(string))
}

func CreateTask(task Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), default_task_times)
	defer cancel()
	_, err := getTaskCollection().InsertOne(ctx, task)
	if err != nil {
		log.Printf("create the task failed:%s", err.Error())
		return err
	}
	if task.State != Ongoing {
		err = Create(task.Operate.GetOperator(), task, default_task_cache_times)
		log.Printf("redis cache task failed:%s", err.Error())
	}
	return nil
}

func DeleteTask(task Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), default_task_times)
	defer cancel()
	_, err := getTaskCollection().DeleteOne(ctx, bson.M{"_id": task.ID})
	return err
}

func UpdateTask(taskid primitive.ObjectID, fields bson.M) error {
	// if task.State != Ongoing {
	// 	return errors.New("the update state must be ongoing")
	// }
	ctx, cancel := context.WithTimeout(context.Background(), default_task_times)
	defer cancel()
	_, err := getTaskCollection().UpdateOne(ctx, bson.M{"_id": taskid}, bson.M{"$set": fields})
	return err
}
