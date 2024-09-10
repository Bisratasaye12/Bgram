package interfaces

import "mime/multipart"



type CloudinaryServiceInterface interface {
	UploadImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}