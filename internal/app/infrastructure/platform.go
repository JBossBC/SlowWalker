package infrastructure

import (
	"errors"
)

type OSType string

type Core string

const (
	GPU Core = "GPU"
	CPU Core = "CPU"
)
const (
	Undefiend OSType = "undefined"
	Linux     OSType = "linux"
	Windows   OSType = "windows"
)

var singlePlatform map[OSType]map[Core]PlatForm

type PlatForm interface {
	GetCoreType() Core
	GetOSType() OSType
	GetExecPrefix(string) string
	PushTask(Operate) error
}

type basePlatForm struct {
	CoreType Core
	// the exec prefix , the resprent the function to execution prefix
	Command     map[string]string
	MechineType OSType
}

func PushTask(operate Operate) {
	//default to push the linux cpu compute
	platform := GetLinuxPlatform(CPU)
	platform.PushTask(operate)
}
func (base *basePlatForm) GetOSType() OSType {
	return base.MechineType
}

func (base *basePlatForm) GetCoreType() Core {
	return base.CoreType
}

func (base *basePlatForm) GetExecPrefix(function string) string {
	return base.Command[function]
}

func (base *basePlatForm) PushTask(op Operate) error {
	return errors.New("error system call")
}

type LocalPlatForm struct {
	basePlatForm
}

func (local *LocalPlatForm) GetOSType() OSType {

}
func (base *LocalPlatForm) GetCoreType() Core {

}

type WindowsPlatForm struct {
	basePlatForm
}

func GetWindowsPlatform(core Core) *WindowsPlatForm {

	return singlePlatform[Windows][core].((*WindowsPlatForm))
}

func (windows *WindowsPlatForm) PushTask(op Operate) error {

}

type LinuxPlatForm struct {
	basePlatForm
}

func GetLinuxPlatform(core Core) *LinuxPlatForm {
	return singlePlatform[Linux][core].((*LinuxPlatForm))
}

func (linux *LinuxPlatForm) PushTask(op Operate) error {

}
