package infrastructure

import (
	interfaces "BChat/Domain/Interfaces"
	models "BChat/Domain/Models"
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cloud *cloudinary.Cloudinary
}

func NewCloudinaryService(env *models.Env) (interfaces.CloudinaryServiceInterface, error) {
	cld, err := cloudinary.NewFromParams(
		env.CLOUDINARY_CLOUD_NAME,
		env.CLOUDINARY_API_KEY,
		env.CLOUDINARY_API_SECRET,
	)
	if err != nil {
		return nil, err
	}

	return &CloudinaryService{cloud: cld}, nil
}

func (cs *CloudinaryService) UploadImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	ctx := context.Background()
	resp, err := cs.cloud.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: fileHeader.Filename,
		Folder:   "profile_photos",
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}
