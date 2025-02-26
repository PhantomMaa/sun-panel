package system

import (
	"sun-panel/lib/storage"
)

type ApiSystem struct {
	About           About
	LoginApi        LoginApi
	UserApi         UserApi
	FileApi         *FileApi
	NoticeApi       NoticeApi
	ModuleConfigApi ModuleConfigApi
	MonitorApi      MonitorApi
}

// InitApiSystem initializes the API system with required dependencies
func InitApiSystem(storageInstance *storage.RcloneStorage) *ApiSystem {
	return &ApiSystem{
		About:           About{},
		LoginApi:        LoginApi{},
		UserApi:         UserApi{},
		FileApi:         NewFileApi(*storageInstance),
		NoticeApi:       NoticeApi{},
		ModuleConfigApi: ModuleConfigApi{},
		MonitorApi:      MonitorApi{},
	}
}
