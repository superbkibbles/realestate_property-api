package file_utils

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_property-api/src/utils/crypto_utils"
)

// Save
// If user not saved
// Delete Pic

func DeleteFile(fileName string, propertyId string) {
	os.Remove(filepath.Join("clients/visuals/"+propertyId, filepath.Base(fileName)))
}

func SaveFile(header *multipart.FileHeader, file multipart.File, propertyId string) (string, string, rest_errors.RestErr) {

	// Check if file is Pic Or Video
	splitter := strings.Split(header.Filename, ".")
	ext := splitter[len(splitter)-1]
	fileName := crypto_utils.GetMd5(header.Filename+strconv.FormatInt(time.Now().Unix(), 36)) + "." + ext

	_, folderErr := os.Stat(filepath.Join("clients/visuals/" + propertyId))
	if os.IsNotExist(folderErr) {
		if err := os.Mkdir(filepath.Join("clients/visuals/"+propertyId), os.ModeAppend); err != nil {
			return "", "", rest_errors.NewInternalServerErr("Error while creating file", err)
		}
	}

	_, err := os.Stat(filepath.Join("clients/visuals/"+propertyId, filepath.Base(fileName)))
	if os.IsNotExist(err) {
		out, err := os.Create(filepath.Join("clients/visuals/"+propertyId, filepath.Base(fileName)))
		if err != nil {
			return "", "", rest_errors.NewInternalServerErr("Error while creating file", err)
		}

		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			return "", "", rest_errors.NewInternalServerErr("Error while saving Pic", err)
		}
		return fileName, ext, nil
	} else {
		return "", "", rest_errors.NewRestError("File Already exist", http.StatusAlreadyReported, "Already exist", nil)
	}
}

func UpdateFile(header *multipart.FileHeader, file multipart.File, path string, propertyId string) (string, string, rest_errors.RestErr) {
	splittedPath := strings.Split(path, "/")
	fileName := splittedPath[len(splittedPath)-1]
	DeleteFile(fileName, propertyId)
	newFileName, ext, err := SaveFile(header, file, propertyId)
	if err != nil {
		return "", "", err
	}
	return newFileName, ext, nil
}
