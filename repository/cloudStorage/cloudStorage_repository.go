package cloudstorage

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
)

const (
	indexProperties        = "property"
	indexTranslateProperty = "translate_property"
	typeProperty           = "_doc"
)

type CloudStorage interface {
	Save(file multipart.File, publicID string, folderName string) (*cloudRes, rest_errors.RestErr)
	Delete(publicID string) rest_errors.RestErr
}

type cloudRes struct {
	Url      string
	Ext      string
	PublicID string
}

type cloudStorage struct {
	cloud *cloudinary.Cloudinary
}

func NewRepository(cloud *cloudinary.Cloudinary) CloudStorage {
	return &cloudStorage{
		cloud: cloud,
	}
}

func (repo *cloudStorage) Delete(publicID string) rest_errors.RestErr {
	ctx := context.Background()
	_, err := repo.cloud.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	if err != nil {
		return rest_errors.NewInternalServerErr("Error while trying to Delete Image/Video", err)
	}
	return nil
}

func (repo *cloudStorage) Save(file multipart.File, publicID string, folderName string) (*cloudRes, rest_errors.RestErr) {
	ctx := context.Background()
	var res cloudRes
	resp, err := repo.cloud.Upload.Upload(ctx, file, uploader.UploadParams{PublicID: publicID, Folder: folderName, Tags: []string{"property"}})
	if err != nil {
		return nil, rest_errors.NewInternalServerErr("Cloudinary Error", err)
	}
	res.Url = resp.URL
	res.Ext = resp.Format
	res.PublicID = publicID
	return &res, nil
}
