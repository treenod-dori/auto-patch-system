package service

import (
	"auto-patch-system/reservations/entity"
	"auto-patch-system/reservations/repository"
	"auto-patch-system/utils"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"
)

type ReservationService struct {
	reservationRepository repository.ReservationRepository
}

func NewReservationService(reservationRepository repository.ReservationRepository) ReservationService {
	return ReservationService{reservationRepository: reservationRepository}
}

type ErrorInfo struct {
	Env   string
	Error error
}

func (s ReservationService) SaveOnlyNotification(patchDate string) error {
	uploadDate := time.Now().Format("2006-01-02")

	reservation := entity.Reservation{
		UploadDate:      uploadDate,
		ReservationDate: patchDate,
		FileName:        "어드민 설정만 전달합니다.",
		PatchData:       []byte("NOT EXISTS QUERY"),
		Success:         0,
	}

	saveErr := s.reservationRepository.SaveReservation(reservation)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

func (s ReservationService) GetReservations() (error, []entity.Reservation) {
	list, selectErr := s.reservationRepository.GetReservations()
	if selectErr != nil {
		return selectErr, nil
	}

	return nil, list
}

func (s ReservationService) SaveReservation(uploadDate string, patchDate string, files []*multipart.FileHeader, titles []string) error {
	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		// 파일 내용을 메모리에 읽기
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		reservation := entity.Reservation{
			FileName:        titles[i],
			PatchData:       fileBytes,
			ReservationDate: patchDate,
			UploadDate:      uploadDate,
			Success:         0,
		}

		createErr := s.reservationRepository.SaveReservation(reservation)
		if createErr != nil {
			return err
		}
	}

	return nil
}

func (s ReservationService) GetReservationByDate(date string) error {
	reservations, selectErr := s.reservationRepository.GetReservationByDate(date)
	if selectErr != nil {
		return selectErr
	}

	if len(reservations) == 0 {
		log.Println("No patches found for today.")
		// 패치할게 없는 경우에 대한 처리
		err := utils.SendSlackMessage("https://hooks.slack.com/services/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", "#only서버파트", "사전QA 패치 알림봇", "오늘 날짜에 패치할 내용이 없습니다.")
		if err != nil {
			log.Printf("슬랙 알림 실패 : %v", err)
			return err
		}
		return nil
	}

	// 패치할 내용은 없지만, 알림은 쏴야 할 경우
	if len(reservations) == 1 && string(reservations[0].FileName) == "어드민 설정만 전달합니다." {
		// webhook api url
		// https://hooks.slack.com/services/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
		err := utils.SendSlackMessage("https://hooks.slack.com/services/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", "#test", "사전QA 패치 알림봇", "어드민 설정만 전달 필요합니다. <https://treenod.atlassian.net/wiki/spaces/pokopokopang/pages/72213266433/2025+QA|여기를 클릭>하여 확인해주세요.")
		if err != nil {
			s.reservationRepository.UpdateReservationStatus(date, 2)
			log.Printf("슬랙 알림 실패 : %v", err)
			return err
		}

		err = s.reservationRepository.UpdateReservationStatus(date, 1)
		if err != nil {
			log.Printf("Error updating reservations: %v", err)
			return err
		}
		return nil
	}

	patchErr := s.patchAllReservationToMySQL(date, reservations)
	if patchErr != nil {
		return patchErr
	}

	return nil
}

// TODO error 핸들링 처리를 해주기
func (s ReservationService) patchAllReservationToMySQL(todayDate string, reservations []entity.Reservation) error {
	configList := []string{"sandbox", "pre_qa", "build_qa", "build_qa2"}
	for _, config := range configList {
		mySQLConfig, _ := utils.NewMySQLConfig(config)
		connectErr := utils.InitMySQL(mySQLConfig)
		if connectErr != nil {
			log.Fatalf("Failed to connect to RDS MySQL: %v", connectErr)
			return nil
		}

		mysql := utils.GetMySqlDB()
		tx, err := mysql.Begin()
		if err != nil {
			continue
		}

		for _, patch := range reservations {
			// SQL 쿼리 생성
			decodeString, _ := base64.StdEncoding.DecodeString(base64.StdEncoding.EncodeToString(patch.PatchData))

			// titles[i] -> 특정 파일명을 가진 쿼리라면 단순 문법 오류만 아니라 추가 검증이 되도록 한다??
			_, err = tx.Exec(string(decodeString))
			if err != nil {
				return err
			}
		}

		// 실행 성공한 경우 커밋하고 다음으로 넘어간다.
		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction for , %v", err)
			s.reservationRepository.UpdateReservationStatus(todayDate, 2)
			return err
		} else {
			continue
		}
	}

	s.reservationRepository.UpdateReservationStatus(todayDate, 1)
	utils.SendSlackMessage("https://hooks.slack.com/services/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", "#test", "사전QA 패치 알림봇", "오늘 항목의 사전QA 패치가 성공적으로 완료되었습니다. <https://treenod.atlassian.net/wiki/spaces/pokopokopang/pages/72213266433/2025+QA|여기를 클릭>하여 확인해주세요.")
	return nil
}

func (s ReservationService) DeleteReservation(fileName string, reservationDate string) error {
	reservedFile, selectErr := s.reservationRepository.GetReservationByDateAndFileName(reservationDate, fileName)
	if selectErr != nil {
		log.Println("Error selecting patch file:", selectErr)
		return selectErr
	}

	// 파일 삭제
	err := s.reservationRepository.DeleteReservation(reservedFile)
	if err != nil {
		return fmt.Errorf("failed to delete patch file: %w", err)
	}

	return nil
}
