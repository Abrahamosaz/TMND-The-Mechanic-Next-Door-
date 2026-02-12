package utils

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/thexovc/TMND/internal/models"
)

type UploadResult struct {
	URL      string
	FileName string
	PublicID string
}

type CloudinaryUploadParams struct {
	Ctx        context.Context
	File       multipart.File
	FileHeader *multipart.FileHeader
	Folder     string
	User       *models.User
}

type Cloudinary struct {
	URL string
}

var CLOUDINARY_PROFILE_IMAGE_FOLDER = "TMND/profilePictures"
var CLOUDINARY_VEHICLE_IMAGE_FOLDER = "TMND/vehicleImages"
var CLOUDINARY_ALLOWED_FORMATS = []string{"jpg", "jpeg", "png", "webp"}
var CLOUDINARY_PROFILE_PICTURE_TRANSFORMATION = "w_500,h_500,c_limit,q_auto"

func (c *Cloudinary) UploadFileToCloudinary(params *CloudinaryUploadParams) (*UploadResult, error) {
	cloudinaryUrl := c.URL
	cld, err := cloudinary.NewFromURL(cloudinaryUrl)
	if err != nil {
		return nil, fmt.Errorf("cloudinary setup error: %w", err)
	}

	// Get current timestamp
	timestamp := time.Now().Unix()
	// Get file name without extension
	fileName := params.FileHeader.Filename
	if dot := strings.LastIndex(fileName, "."); dot != -1 {
		fileName = fileName[:dot]
	}

	// Create PublicID as "<timestamp>_<filename-without-ext>"
	publicID := fmt.Sprintf("%d_%s", timestamp, fileName)
	uploadResult, err := cld.Upload.Upload(params.Ctx, params.File, uploader.UploadParams{
		PublicID:       publicID,
		Folder:         fmt.Sprintf("%s/%s", params.Folder, params.User.ID),
		AllowedFormats: CLOUDINARY_ALLOWED_FORMATS,
		Transformation: CLOUDINARY_PROFILE_PICTURE_TRANSFORMATION,
	})
	if err != nil {
		return nil, fmt.Errorf("cloudinary upload error: %w", err)
	}

	return &UploadResult{
		URL:      uploadResult.SecureURL,
		FileName: params.FileHeader.Filename,
		PublicID: publicID,
	}, nil
}

func (c *Cloudinary) DeleteFileFromCloudinary(publicID string) (*uploader.DestroyResult, error) {
	cloudinaryUrl := c.URL
	cld, err := cloudinary.NewFromURL(cloudinaryUrl)
	if err != nil {
		return nil, fmt.Errorf("cloudinary setup error: %w", err)
	}

	result, err := cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return nil, fmt.Errorf("cloudinary delete error: %w", err)
	}

	log.Printf("cloudinary delete result: %v", result)
	return result, nil
}
