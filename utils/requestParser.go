package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io"
)

// FileUploadRequest: 파일 업로드 요청 데이터 구조체
type FileUploadRequest struct {
	FileBytes []byte
	Title     string
	PatchDate string
}

type QueryParameterRequest struct{}

// Parse: 요청 데이터를 파싱하는 메서드
func (r *FileUploadRequest) ParsePatchData(ctx *gin.Context, maxMemory int64) error {
	// 요청 파싱
	if err := ctx.Request.ParseMultipartForm(maxMemory); err != nil {
		return errors.New("unable to parse form")
	}

	// 파일 읽기
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		return errors.New("file not found in request")
	}
	defer file.Close()

	// 파일 데이터 읽기
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return errors.New("failed to read file")
	}

	// 제목 가져오기
	title := ctx.DefaultPostForm("title", "")
	if title == "" {
		return errors.New("title not provided")
	}

	// patchDate 가져오기
	patchDate := ctx.DefaultPostForm("patchDate", "")
	if patchDate == "" {
		return errors.New("patchDate not provided")
	}

	// 구조체 필드에 저장
	r.FileBytes = fileBytes
	r.Title = title
	r.PatchDate = patchDate

	return nil
}

// ParseQueryParams: 여러 개의 키에 해당하는 쿼리 파라미터를 추출하는 메서드
func (qp *QueryParameterRequest) ParseQueryParams(ctx *gin.Context, keys []string) (map[string]string, error) {
	params := make(map[string]string)
	for _, key := range keys {
		value := ctx.DefaultQuery(key, "")
		params[key] = value
	}
	return params, nil
}
