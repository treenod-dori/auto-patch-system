package service

import (
	"auto-patch-system/reservations/entity"
	"auto-patch-system/reservations/repository"
	"auto-patch-system/utils"
	"encoding/base64"
	"io"
	"log"
	"mime/multipart"
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

func (s ReservationService) SaveOnlyNotification(patchDate string, reservedTime string) error {
	reservation := entity.Reservation{
		ReservationTime: reservedTime,
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
	var reservationList []entity.Reservation
	selectErr := s.reservationRepository.GetReservations(&reservationList)
	if selectErr != nil {
		return selectErr, nil
	}

	return nil, reservationList
}

func (s ReservationService) SaveReservation(patchDate string, files []*multipart.FileHeader, titles []string) error {
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
		reservedTime := "10:00:00"

		reservation := entity.Reservation{
			FileName:        titles[i],
			PatchData:       fileBytes,
			ReservationDate: patchDate,
			ReservationTime: reservedTime,
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
	selectErr, reservations := s.reservationRepository.GetReservationByDate(date)
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
	}

	// 패치할 내용은 없지만, 알림은 쏴야 할 경우
	if len(reservations) == 1 && string(reservations[0].FileName) == "어드민 설정만 전달합니다." {
		// webhook api url
		// https://hooks.slack.com/services/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
		err := utils.SendSlackMessage("https://hooks.slack.com/services/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", "#test", "사전QA 패치 알림봇", "어드민 설정만 전달 필요합니다. <https://treenod.atlassian.net/wiki/spaces/pokopokopang/pages/72213266433/2025+QA|여기를 클릭>하여 확인해주세요.")
		if err != nil {
			log.Printf("슬랙 알림 실패 : %v", err)
			return err
		}

		err = s.reservationRepository.UpdateReservationStatus(date, 2)
		if err != nil {
			log.Printf("Error updating reservations: %v", err)
			return err
		}
	}

	patchErr := s.patchAllReservationToMySQL(date, reservations)
	if patchErr != nil {
		return patchErr
	}

	updateErr := s.reservationRepository.UpdateReservationStatus(date, 1)
	if updateErr != nil {
		log.Printf("Error updating reservations: %v", updateErr)
		return updateErr
	}

	// slack에 성공 메시지 전송

	return nil
}

// TODO error 핸들링 처리를 해주기
func (s ReservationService) patchAllReservationToMySQL(todayDate string, reservations []entity.Reservation) error {
	configList := []string{"sandbox", "pre_qa", "build_qa", "build_qa2"}
	var errInfo []ErrorInfo
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
			// 파일이름이 ANI_LIST를 포함하는 경우에는 특정 validation 체크 로직을 진행한다.
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
		} else {
			continue
		}
	}

	if len(errInfo) > 0 {
		err := s.reservationRepository.UpdateReservationStatus(todayDate, 2)
		if err != nil {
			return err
		}

		// slack에 errInfo를 포함한 메시지 전송
	}

	return nil
}
