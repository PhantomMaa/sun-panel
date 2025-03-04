package panelApiStructs

import (
	"sun-panel/api/common/apiData/commonApiStructs"
	"sun-panel/internal/repository"
)

type ItemIconEditRequest struct {
	repository.ItemIcon
	IconJson string
}

type ItemIconSaveSortRequest struct {
	SortItems       []commonApiStructs.SortRequestItem `json:"sortItems"`
	ItemIconGroupId uint                               `json:"itemIconGroupId"`
}

type ItemIconGetSiteFaviconReq struct {
	Url string `form:"url" json:"url"`
}

type ItemIconGetSiteFaviconResp struct {
	IconUrl string `json:"iconUrl"`
}
