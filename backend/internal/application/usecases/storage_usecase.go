package usecases

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"babywolf-blog-backend/internal/infrastructure/database"
)

type StorageUseCase interface {
	UploadFile(folder string, fileName string, file multipart.File) (string, error)
}

type storageUseCase struct{}

func NewStorageUseCase() StorageUseCase {
	return &storageUseCase{}
}

func (uc *storageUseCase) UploadFile(folder string, fileName string, file multipart.File) (publicURL string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic uploading file: %v", r)
		}
	}()

	bucketName := "blog_assets"

	buf := bytes.NewBuffer(nil)
	if _, errCopy := io.Copy(buf, file); errCopy != nil {
		return "", fmt.Errorf("failed to read file: %v", errCopy)
	}

	remotePath := folder + "/" + fileName

	// Utilizando nedpals/supabase-go que hace panic en errores
	result := database.SupabaseClient.Storage.From(bucketName).Upload(remotePath, buf, nil)
	
	if result.Message != "" && result.Key == "" {
		return "", fmt.Errorf("failed to upload to supabase: %s", result.Message)
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	publicURL = supabaseURL + "/storage/v1/object/public/" + bucketName + "/" + remotePath

	return publicURL, nil
}
