package systemApi

type MonitorGetDiskStateByPathReq struct {
	Path string `form:"path" json:"path"`
}
