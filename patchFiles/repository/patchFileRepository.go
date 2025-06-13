package repository

import (
	"auto-patch-system/patchFiles/entity"
	"auto-patch-system/utils"
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
type patchFileRepository struct{}

func NewPatchFileRepository() PatchFileRepository {
	return patchFileRepository{}
}

func (p patchFileRepository) SavePatchFile(queryData entity.PatchFile) error {
	err := utils.GetSqliteDB().Table("patchFiles").Create(&queryData).Error
	if err != nil {
		return err
	}
	return nil
}

func (p patchFileRepository) IsExistPatchData(fileName, reservationDate string) bool {
	result := utils.GetSqliteDB().Table("patchFiles").
		Select("patchData").
		Where("title = ? AND reservationDate = ?", fileName, reservationDate).
		Find(&entity.PatchFile{})

	return result.RowsAffected > 0
}

func (p patchFileRepository) GetPatchFileListByDate(patchDate string) []entity.PatchFile {
	var list []entity.PatchFile
	err := utils.GetSqliteDB().Table("patchFiles").
		Where("reservationDate = ?", patchDate).
		Find(&list).Error

	if err != nil {
		return nil
	}
	return list
}

func (p patchFileRepository) DeletePatchFile(data entity.PatchFile) error {
	err := utils.GetSqliteDB().Table("patchFiles").Delete(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (p patchFileRepository) GetPatchFile(fileName string, patchDate string) (entity.PatchFile, error) {
	var result entity.PatchFile
	err := utils.GetSqliteDB().Table("patchFiles").
		Where("title = ? AND reservationDate = ?", fileName, patchDate).
		First(&result).Error

	if err != nil {
		log.Print("Error retrieving patch file:", err)
		return result, err
	}
	return result, nil
}
