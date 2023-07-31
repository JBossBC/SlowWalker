package dao

import (
	"errors"
	"log"
	"os/exec"
	"replite_web/internal/app/utils"
	"runtime"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (remote *RemotePlatForm) PushTask(op Operate) error {

	return nil
}

type LocalPlatForm struct {
	BasePlatForm
}

func (local *LocalPlatForm) PushTask(op Operate) error {
	if op.GetOperateType() != ShortTerm {
		return errors.New("local platform only support the short term task ")
	}
	task := new(Task)
	task.Operate = op
	task.PlatForm = local
	task.State = Ongoing
	task.ID = primitive.NewObjectID()
	err := CreateTask(*task)
	if err != nil {
		return err
	}
	funcMap := GetFuncMap(op.GetFunction())
	args := make([]string, 0, 3)
	args = append(args, funcMap.ExecFile)
	args = append(args, op.GetParams()...)
	// starting the goroutinue to execute the operate
	go func(id primitive.ObjectID) {
		cmd := exec.Command(op.GetCommand(), args...)
		msg, err := cmd.Output()
		state := Success
		if err != nil {
			state = Failed
		}
		err = UpdateTask(id, bson.M{"message": msg, "state": state})
		if err != nil {
			log.Printf("update task state(id:%s,message:%s,state:%s) error:%s", id, msg, state, err.Error())
		}
		// err := cmd.Run()
	}(task.ID)
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
