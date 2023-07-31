package dao

import (
	"errors"
)

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

var singlePlatform map[OSType]PlatForm

type PlatForm interface {
	GetCoreType() []Core
	GetOSType() OSType
	GetExecPrefix() string
	PushTask(Operate) error
}

type BasePlatForm struct {
	CoreType []Core
	// the exec prefix , the resprent the function to execution prefix
	Command     string
	MechineType OSType
}

func PushTask(operate Operate) {
	//default to push the linux cpu compute
	platform := GetLinuxPlatform(CPU)
	platform.PushTask(operate)
}
func (base *BasePlatForm) GetOSType() OSType {
	return base.MechineType
}

func (base *BasePlatForm) GetCoreType() []Core {
	return base.CoreType
}

func (base *BasePlatForm) GetExecPrefix() string {
	return base.Command
}

func (base *BasePlatForm) PushTask(op Operate) error {
	return errors.New("error system call")
}

type RemotePlatForm struct {
	BasePlatForm
}

type LocalPlatForm struct {
	BasePlatForm
}

// // func (local *LocalPlatForm) GetOSType() OSType {

// // }
// // func (base *LocalPlatForm) GetCoreType() Core {

// // }

// type WindowsPlatForm struct {
// 	BasePlatForm
// }

func GetWindowsPlatform(core Core) *RemotePlatForm {
	return singlePlatform[Windows].(*RemotePlatForm)
}

// func (windows *WindowsPlatForm) PushTask(op Operate) error {
// 	cmd := []string{windows.Command}
// }

// type LinuxPlatForm struct {
// 	BasePlatForm
// }

func GetLinuxPlatform(core Core) *RemotePlatForm {
	return singlePlatform[Linux].(*RemotePlatForm)
}

// func (linux *LinuxPlatForm) PushTask(op Operate) error {

// }
