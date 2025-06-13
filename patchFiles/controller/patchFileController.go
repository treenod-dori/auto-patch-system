package controller

import (
	"auto-patch-system/patchFiles/entity"
	"auto-patch-system/patchFiles/service"
	"auto-patch-system/utils"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type PatchFileController struct {
	patchService service.PatchFileService
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type DownloadRequest struct {
	Files []entity.PatchFile `json:"files"`
}

type ErrorResponse struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Code    string      `json:"code"`
}

func NewPatchFileController(patchService service.PatchFileService) PatchFileController {
	return PatchFileController{patchService: patchService}
}

func (controller *PatchFileController) SaveCrawlingQuery(ctx *gin.Context) {
	var request utils.FileUploadRequest
	err := request.ParsePatchData(ctx, 1<<20)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Unable to parse form",
			Details: err,
			Code:    utils.REQUEST_PARSE_ERROR,
		})
		return
	}

	// 패치 파일 구조체 생성
	patchFile := entity.PatchFile{
		Title:           request.Title,
		ReservationDate: request.PatchDate,
		PatchData:       request.FileBytes,
	}
	saveErr := controller.patchService.SavePatchFile(patchFile)
	if saveErr != nil {
		errResponse := ErrorResponse{
			Message: "Failed to save file",
			Details: saveErr,
			Code:    utils.DB_SAVE_ERROR,
		}
		ctx.JSON(http.StatusBadRequest, &errResponse)
		return
	}

	response := Response{
		Status:  "success",
		Message: "File saved successfully",
	}
	ctx.JSON(http.StatusOK, &response)
}

func (controller *PatchFileController) DownloadFile(ctx *gin.Context) {
	var queryParams utils.QueryParameterRequest
	var fileRequest DownloadRequest

	params, _ := queryParams.ParseQueryParams(ctx, []string{"date"})
	patchDate := params["date"]

	if err := ctx.BindJSON(&fileRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	patchList := fileRequest.Files

	finalFile, makeErr := controller.patchService.MakeMergedPatchFile(patchList, patchDate)
	if makeErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": makeErr.Error(),
		})
		return
	}

	// 파일 다운로드
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", finalFile))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")
	ctx.File(finalFile)

	// 다운로드 후 임시 파일 삭제
	if err := os.Remove(finalFile); err != nil {
		log.Println("Failed to delete temporary file:", err)
	}
}

// TestPatchList는 여러 개의 파일을 테스트 실행하는 메서드
func (controller *PatchFileController) TestAllPatchList(ctx *gin.Context) {
	// MySQL 설정 로드
	dbConfig, err := utils.NewMySQLConfig("sandbox")
	if err != nil {
		log.Println("Failed to connect to RDS MySQL: %v", err)
		ctx.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to load database config",
		})
		return
	}

	err = utils.InitMySQL(dbConfig)
	if err != nil {
		log.Println("Failed to connect to RDS MySQL: %v", err)
		ctx.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to load database config",
		})
		return
	}
	db := utils.GetMySqlDB()

	// 요청 파싱
	err = ctx.Request.ParseMultipartForm(1 << 20) // 최대 1MB
	if err != nil {
		ctx.JSON(400, gin.H{
			"status":  "error",
			"message": "Unable to parse form",
		})
		return
	}

	// 파일과 제목 추출
	files := ctx.Request.MultipartForm.File["blobs"] // 'blobs'라는 키로 파일 읽기
	titles := ctx.Request.MultipartForm.Value["titles"]

	//titles를 보고 만약 특별한 처리라라면 그거에 대한 검증 로직을 할 수 있도록 한다.

	// 트랜잭션 시작
	tx, _ := db.Begin()

	// 파일들 처리
	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			ctx.JSON(500, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Failed to open file: %v", err),
			})
			continue
		}
		defer file.Close()

		// 파일 내용을 메모리에 읽기
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			ctx.JSON(500, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Failed to read file: %v", err),
			})
			continue
		}

		// SQL 쿼리 생성
		decodeString, _ := base64.StdEncoding.DecodeString(base64.StdEncoding.EncodeToString(fileBytes))

		// multiStatements를 사용하여 SQL 쿼리 실행
		start := time.Now()
		_, execErr := tx.Exec(string(decodeString))
		elapsed := time.Since(start)
		log.Printf("쿼리 실행 시간: %s", elapsed)

		if execErr != nil {
			_ = tx.Rollback()
			ctx.JSON(400, gin.H{
				"status":  "error",
				"message": execErr.Error(),
				"details": titles[i],
				"code":    400,
			})
			return
		}
	}

	_ = tx.Rollback()
	// 성공 응답
	response := Response{
		Status:  "success",
		Message: "Files test successfully",
	}
	ctx.JSON(http.StatusOK, &response)
}

func (controller *PatchFileController) GetAllPatchFiles(ctx *gin.Context) {
	var queryParams utils.QueryParameterRequest
	params, _ := queryParams.ParseQueryParams(ctx, []string{"date"})
	// 쿼리 파라미터에서 날짜, 제목 가져오기
	patchDate := params["date"]

	patchList, getErr := controller.patchService.FindPatchFilesByDate(patchDate)
	if getErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": getErr.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, patchList)
}

func (controller *PatchFileController) DeletePatchFile(ctx *gin.Context) {
	title := ctx.PostForm("title")
	patchDate := ctx.PostForm("reservationDate")

	if title == "" || patchDate == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing form fields"})
		return
	}

	// 패치 파일 삭제
	deleteErr := controller.patchService.DeletePatchFile(title, patchDate)
	if deleteErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": deleteErr.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Patch file deleted successfully",
	})
}
