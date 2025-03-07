package panelApi

import (
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/web/model/param/commonApi"
)

type ItemIconEditRequest struct {
	repository.ItemIcon
	IconJson string
}

type ItemIconSaveSortRequest struct {
	SortItems       []commonApi.SortRequestItem `json:"sortItems"`
	ItemIconGroupId uint                        `json:"itemIconGroupId"`
}

type ItemIconGetSiteFaviconReq struct {
	Url string `form:"url" json:"url"`
}

type ItemIconGetSiteFaviconResp struct {
	IconUrl string `json:"iconUrl"`
}
