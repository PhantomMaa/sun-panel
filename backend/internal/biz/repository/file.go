package repository

type File struct {
	BaseModel
	UserId   uint   `gorm:"index" json:"userId"`
	FileName string `gorm:"type:varchar(50)" json:"fileName"`
}

type FileRepo struct{}

type IFileRepo interface {
	Get(userId, id uint) (File, error)
	GetList(userId uint) ([]File, uint, error)
	Delete(userId, id uint) error
}

func NewFileRepo() *FileRepo {
	return &FileRepo{}
}

func (r *FileRepo) AddFile(userId uint, fileName string) (File, error) {
	file := File{
		UserId:   userId,
		FileName: fileName,
	}
	err := Db.Create(&file).Error
	return file, err
}

func (r *FileRepo) Get(userId, id uint) (File, error) {
	var file File
	err := Db.Where("user_id=? AND id=?", userId, id).First(&file).Error
	return file, err
}

func (r *FileRepo) GetList(userId uint) ([]File, uint, error) {
	var list []File
	var count int64
	err := Db.Order("created_at desc").Find(&list, "user_id=?", userId).Count(&count).Error
	return list, uint(count), err
}

func (r *FileRepo) Delete(userId, id uint) error {
	return Db.Delete(&File{}, "id = ? AND user_id = ?", id, userId).Error
}
