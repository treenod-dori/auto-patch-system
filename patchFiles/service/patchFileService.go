package service

import (
	"auto-patch-system/patchFiles/entity"
	"auto-patch-system/patchFiles/repository"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

type PatchFileService struct {
	patchFileRepository repository.PatchFileRepository
}

func NewPatchFileService(patchFileRepository repository.PatchFileRepository) PatchFileService {
	return PatchFileService{patchFileRepository: patchFileRepository}
}

func (s PatchFileService) SavePatchFile(file entity.PatchFile) error {
	// 1. file의 reservationDate와 title이 이미 존재하는지 확인한다.
	isExist := s.patchFileRepository.IsExistPatchData(file.Title, file.ReservationDate)

	// 2. 만약 있다면 이미 저장된 거라는 메시지를 반환한다.
	if isExist {
		// nil이 아닌 다른 response를 반환한다.
		return errors.New("already exist")
	}

	// 3. 없다면 저장한다.
	saveErr := s.patchFileRepository.SavePatchFile(file)
	if saveErr != nil {
		return saveErr
	}

	return nil
}

func (s PatchFileService) FindPatchFilesByDate(patchDate string) ([]entity.PatchFile, error) {
	patchFiles := s.patchFileRepository.GetPatchFileListByDate(patchDate)
	if len(patchFiles) == 0 {
		return nil, errors.New("no data found")
	}

	return patchFiles, nil
}

func (s PatchFileService) MakeMergedPatchFile(patchFiles []entity.PatchFile, patchDate string) (string, error) {
	// 병합된 파일 생성
	tempFilePath := "merged_patch_data.sql"
	file, _ := os.Create(tempFilePath)
	defer file.Close()

	var (
		firstTitle string
		isFirst    = true
	)

	for _, patchFile := range patchFiles {
		// 첫 번째 데이터의 제목 저장
		if isFirst {
			firstTitle = patchFile.Title
			isFirst = false
		}

		// BLOB 데이터를 파일에 쓰기
		if _, err := file.Write(patchFile.PatchData); err != nil {
			log.Println("Error writing to file:", err)
			continue
		}
	}

	// 날짜 형식 변환
	parsedDate, _ := time.Parse("2006-01-02", patchDate)
	formattedDate := parsedDate.Format("060102")
	fileName := fmt.Sprintf("%s_POKOPOKO_GAMEDB_%s.sql", formattedDate, firstTitle)
	if len(patchFiles) > 1 {
		fileName = fmt.Sprintf("%s_POKOPOKO_GAMEDB_%s_ETC.sql", formattedDate, firstTitle)
	}

	// 파일 이름 변경
	finalFilePath := fileName
	if err := os.Rename(tempFilePath, finalFilePath); err != nil {
		return "", fmt.Errorf("failed to rename file: %w", err)
	}

	return finalFilePath, nil
}
