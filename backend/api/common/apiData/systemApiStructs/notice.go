package systemApiStructs

type NoticeGetListByDisplayTypeReq struct {
	DisplayType []int `form:"displayType" json:"displayType"`
}
