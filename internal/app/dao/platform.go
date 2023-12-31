package dao

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"replite_web/internal/app/utils"
	"runtime"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/proto"
)

func init() {
	local := new(LocalPlatForm)
	local.MechineType = OSType(runtime.GOOS)
	//default is the CPU handle
	local.CoreType = []Core{CPU}
	local.IP, _ = utils.GetLocalIP()
	localPlatForm = local
}

type OSType string

type Core string

const (
	GPU  Core = "GPU"
	CPU  Core = "CPU"
	None Core = "none"
)
const (
	Undefiend OSType = "undefined"
	Linux     OSType = "linux"
	Windows   OSType = "windows"
)

var localPlatForm *LocalPlatForm

var remotePlatForm map[OSType]PlatForm

type PlatForm interface {
	GetCoreType() []Core
	GetOSType() OSType
	// GetExecPrefix() string
	PushTask(Operate) error
	GetParamsPrefix() string
	GetParamsSuffix() string
	GetIP() string
}

type BasePlatForm struct {
	IP       string
	CoreType []Core
	// // the exec prefix , the resprent the function to execution prefix
	// Command      string
	MechineType  OSType
	ParamsPrefix string
	ParamsSuffix string
}

// PushTask meaning the execute environment is the local,the must take the params  prefix and suffix for executing
func PushTask(operate Operate) {
	//default to push the linux cpu compute
	platform := GetLocalPlatForm()
	platform.PushTask(operate)
}
func (base *BasePlatForm) GetOSType() OSType {
	return base.MechineType
}
func (base *BasePlatForm) GetIP() string {
	return base.IP
}

func (base *BasePlatForm) GetCoreType() []Core {
	return base.CoreType
}

//	func (base *BasePlatForm) GetExecPrefix() string {
//		return base.Command
//	}
func (base *BasePlatForm) GetParamsPrefix() string {
	return base.ParamsPrefix
}
func (base *BasePlatForm) GetParamsSuffix() string {
	return base.ParamsSuffix
}

func (base *BasePlatForm) PushTask(op Operate) error {
	return errors.New("error system call")
}

type RemotePlatForm struct {
	BasePlatForm
}

const TaskTopic = "task"

func getTaskWriter() *kafka.Writer {
	return GetTopicConn(TaskTopic)
}

const default_kafka_times = 5 * time.Second

func (remote *RemotePlatForm) PushTask(op Operate) error {
	ctx, cancel := context.WithTimeout(context.Background(), default_kafka_times)
	defer cancel()
	documentID := primitive.NewObjectID()
	image := new(TaskImage)
	image.ID = documentID[:]
	var exec = make([]string, 0, 3)
	funcMap := GetFunctionDao().GetFuncMap(op.GetFunction())
	exec = append(exec, funcMap.Command)
	exec = append(exec, op.GetParams()...)
	image.Exec = exec
	imageBytes, err := proto.Marshal(image)
	if err != nil {
		log.Printf("protocol Buffer marshal failed:%s", err.Error())
		return err
	}
	task := new(TaskInfo)
	task.PlatForm = remote
	task.Operate = op
	task.State = Ongoing
	task.ID = documentID
	err = GetTaskDao().CreateTask(*task)
	if err != nil {
		return err
	}
	//the topic represent the task environment of executing****************************
	topic := utils.MergeStr(string(funcMap.Type), "-", string(funcMap.OSType), "-", string(funcMap.Additional), "-", string(funcMap.Function))
	err = getTaskWriter().WriteMessages(ctx, kafka.Message{
		Key:   []byte(topic),
		Value: imageBytes,
	})
	if err != nil {
		//RollBack this task
		GetTaskDao().DeleteTask(*task)

		log.Printf("push kafka failed:%s", err.Error())
		return err
	}
	//create task success
	return nil
}

type LocalPlatForm struct {
	BasePlatForm
}

func (local *LocalPlatForm) PushTask(op Operate) error {
	if op.GetOperateType() != ShortTerm {
		return errors.New("local platform only support the short term task ")
	}
	documentID := primitive.NewObjectID()
	funcMap := GetFunctionDao().GetFuncMap(op.GetFunction())
	// args := make([]string, 0, 3)
	// args = append(args, funcMap.Command)
	// args = append(args, op.GetParams()...)
	task := new(TaskInfo)
	task.Operate = op
	task.PlatForm = local
	task.State = Ongoing
	task.ID = documentID
	err := GetTaskDao().CreateTask(*task)
	if err != nil {
		return err
	}
	// starting the goroutinue to execute the operate
	go func(id primitive.ObjectID) {
		//TODO the params if is the file ,need to convert file location for this system
		cmd := exec.Command(funcMap.Command, op.GetParams()...)
		msg, err := cmd.Output()
		state := Success
		if err != nil {
			log.Println("执行Task出错:", msg)
			state = Failed
		}
		err = GetTaskDao().UpdateTask(id, bson.M{"message": fmt.Sprintf("任务执行失败:%s", msg), "state": state})
		if err != nil {
			log.Printf("update task state(id:%s,message:%s,state:%s) error:%s", string(id[:]), msg, strconv.Itoa(int(state)), err.Error())
		}
		// err := cmd.Run()
	}(documentID)
	return nil
}

func GetLocalPlatForm() *LocalPlatForm {
	return localPlatForm
}

// // func (local *LocalPlatForm) GetOSType() OSType {

// // }
// // func (base *LocalPlatForm) GetCoreType() Core {

// // }

// type WindowsPlatForm struct {
// 	BasePlatForm
// }

func GetWindowsPlatform() *RemotePlatForm {
	return remotePlatForm[Windows].(*RemotePlatForm)
}

// func (windows *WindowsPlatForm) PushTask(op Operate) error {
// 	cmd := []string{windows.Command}
// }

// type LinuxPlatForm struct {
// 	BasePlatForm
// }

func GetLinuxPlatform() *RemotePlatForm {
	return remotePlatForm[Linux].(*RemotePlatForm)
}

// func (linux *LinuxPlatForm) PushTask(op Operate) error {

// }
