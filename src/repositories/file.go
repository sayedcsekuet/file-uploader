package repositories

import (
	"file-uploader/src/errors"
	"file-uploader/src/models"
	"fmt"
	"github.com/araddon/dateparse"
	"gorm.io/gorm"
	"time"
)

type SearchParam struct {
	Name        string
	CreatedDate string
	Offset      int
	Limit       int
}

type FileRepository interface {
	Get(id string) (*models.File, error)
	Create(data *models.File) (*models.File, error)
	GetAll() ([]*models.File, error)
	GetAllByOwner(ownerId string, params SearchParam) ([]*models.File, error)
	GetExpiredFiles() ([]*models.File, error)
	FindByOwnerAndId(id, ownerId string) (*models.File, error)
	Delete(id, ownerId string) error
	DeleteAll(ids []string) error
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db}
}

const DateTimeLayout = "2006-01-02 15:04:05"
const DateLayout = "2006-01-02"

type fileRepository struct {
	db *gorm.DB
}

func (fr *fileRepository) GetExpiredFiles() ([]*models.File, error) {
	var a []*models.File
	if err := fr.db.
		Where("expired_at < ?", fr.currentDateTime()).
		Find(&a).Error; err != nil {
		return nil, err
	}

	return a, nil
}

func (fr *fileRepository) Delete(id, ownerId string) error {
	var a *models.File
	return fr.db.Where("id=? AND owner_id = ?", id, ownerId).Delete(&a).Error
}
func (fr *fileRepository) DeleteAll(ids []string) error {
	var a []*models.File
	if err := fr.db.Delete(&a, ids).Error; err != nil {
		return err
	}
	return nil
}
func (fr *fileRepository) Create(data *models.File) (*models.File, error) {
	if err := fr.db.Create(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (fr *fileRepository) GetAll() ([]*models.File, error) {
	var a []*models.File
	if err := fr.db.Where("expired_at IS NULL OR expired_at>=?", fr.currentDateTime()).Find(&a).Error; err != nil {
		return nil, err
	}

	return a, nil
}

func (fr *fileRepository) GetAllByOwner(ownerId string, params SearchParam) ([]*models.File, error) {
	var a []*models.File
	condition := "owner_id = @owner_id AND (expired_at IS NULL OR expired_at >= @expired_at)"
	searParams := map[string]interface{}{"owner_id": ownerId, "expired_at": fr.currentDateTime()}
	if params.Limit == 0 {
		params.Limit = 500
	}
	if params.Name != "" {
		condition = condition + " AND name LIKE @name"
		searParams["name"] = "%" + params.Name + "%"
	}
	if params.CreatedDate != "" {
		condition = condition + " AND DATE(created_at)=@created_at"
		cDate, _ := dateparse.ParseAny(params.CreatedDate)
		searParams["created_at"] = cDate.Format(DateLayout)
	}
	if err := fr.db.
		Where(condition, searParams).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&a).Error; err != nil {
		return nil, err
	}

	return a, nil
}

func (fr *fileRepository) Get(id string) (*models.File, error) {
	var a *models.File
	if err := fr.db.Where("id=? AND (expired_at IS NULL OR expired_at>=?)", id, fr.currentDateTime()).Find(&a).Error; err != nil {
		return nil, err
	}

	if a.ID == "" {
		return nil, errors.NewKnown(
			404,
			fmt.Sprintf("File not found with id %s!", id),
		)
	}
	return a, nil
}

func (fr *fileRepository) FindByOwnerAndId(id, ownerId string) (*models.File, error) {
	var a *models.File
	if err := fr.db.Where("id=? AND owner_id = ? AND (expired_at IS NULL OR expired_at>=?)", id, ownerId, fr.currentDateTime()).Find(&a).Error; err != nil {
		return nil, err
	}
	if a.ID == "" {
		return nil, errors.NewKnown(
			404,
			fmt.Sprintf("File not found with id %s and with owner XXXXXXX. Please verify you are the correct owner of that file or file doesn't exist!", id),
		)
	}
	return a, nil
}

func (fr *fileRepository) currentDateTime() string {
	return time.Now().Format(DateTimeLayout)
}
