package repository

import (
	"auto-patch-system/patchFiles/entity"
	"auto-patch-system/utils"
	"gorm.io/gorm"
	"log"
)

type PatchFileRepository interface {
	SavePatchFile(patchData entity.PatchFile) error
	IsExistPatchData(fileName, reservationDate string) bool
	GetPatchFileListByDate(patchDate string) []entity.PatchFile
	DeletePatchFile(deleteData entity.PatchFile) error
	GetPatchFile(fileName string, patchDate string) (entity.PatchFile, error)
}

// 실제 구현체
type patchFileRepository struct {
	db *gorm.DB
}

func NewPatchFileRepository() PatchFileRepository {
	return patchFileRepository{
		db: utils.GetMySqlGormDB(),
	}
}

func (p patchFileRepository) SavePatchFile(queryData entity.PatchFile) error {
	err := p.db.Table("patchFiles").Create(&queryData).Error
	if err != nil {
		return err
	}
	return nil
}

func (p patchFileRepository) IsExistPatchData(fileName, reservationDate string) bool {
	result := p.db.Table("patchFiles").
		Select("patchData").
		Where("title = ? AND reservationDate = ?", fileName, reservationDate).
		Find(&entity.PatchFile{})

	return result.RowsAffected > 0
}

func (p patchFileRepository) GetPatchFileListByDate(patchDate string) []entity.PatchFile {
	var list []entity.PatchFile
	err := p.db.Table("patchFiles").
		Where("reservationDate = ?", patchDate).
		Find(&list).Error

	if err != nil {
		return nil
	}
	return list
}

func (p patchFileRepository) DeletePatchFile(data entity.PatchFile) error {
	err := p.db.Table("patchFiles").Delete(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (p patchFileRepository) GetPatchFile(fileName string, patchDate string) (entity.PatchFile, error) {
	var result entity.PatchFile
	err := p.db.Table("patchFiles").
		Where("title = ? AND reservationDate = ?", fileName, patchDate).
		First(&result).Error

	if err != nil {
		log.Print("Error retrieving patch file:", err)
		return result, err
	}
	return result, nil
}
