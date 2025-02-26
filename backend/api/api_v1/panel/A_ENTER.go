package panel

import "sun-panel/lib/storage"

type ApiPanel struct {
	ItemIcon      *ItemIcon
	UserConfig    UserConfig
	UsersApi      UsersApi
	ItemIconGroup ItemIconGroup
}

// InitApiSystem initializes the API system with required dependencies
func InitApiPanel(storageInstance *storage.RcloneStorage) *ApiPanel {
	return &ApiPanel{
		ItemIcon:      NewItemIcon(*storageInstance),
		UserConfig:    UserConfig{},
		UsersApi:      UsersApi{},
		ItemIconGroup: ItemIconGroup{},
	}
}
