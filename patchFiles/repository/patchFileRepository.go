package repository

import (
	"auto-patch-system/patchFiles/entity"
	"auto-patch-system/utils"
)

type PatchFileRepository interface {
	SavePatchFile(patchData entity.PatchFile) error
	IsExistPatchData(fileName, reservationDate string) bool
	GetPatchFileListByDate(patchDate string) []entity.PatchFile
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
