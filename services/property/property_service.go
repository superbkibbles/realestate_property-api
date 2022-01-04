package property

import (
	"mime/multipart"

	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_property-api/domain/property"
	"github.com/superbkibbles/realestate_property-api/domain/query"
	"github.com/superbkibbles/realestate_property-api/repository/db"
	"github.com/superbkibbles/realestate_property-api/utils/date_utils"
	"github.com/superbkibbles/realestate_property-api/utils/file_utils"
)

type Service interface {
	Create(property.Property) (*property.Property, rest_errors.RestErr)
	Get() (property.Properties, rest_errors.RestErr)
	GetByID(string) (*property.Property, rest_errors.RestErr)
	Search(query query.EsQuery) (property.Properties, rest_errors.RestErr)
	Update(id string, updateRequest property.EsUpdate) (*property.Property, rest_errors.RestErr)
	UploadMedia(files []*multipart.FileHeader, propertyID string) rest_errors.RestErr
	DeleteMedia(propertyID string, mediaID string) rest_errors.RestErr
	UploadProperyPic(id string, fileHeader *multipart.FileHeader) (*property.Property, rest_errors.RestErr)
}

type service struct {
	// Add Database Repository
	dbRepo db.DbRepository
}

func NewService(dbRepo db.DbRepository) Service {
	return &service{
		dbRepo: dbRepo,
	}
}

func (s *service) Update(id string, updateRequest property.EsUpdate) (*property.Property, rest_errors.RestErr) {
	if err := updateRequest.Validate(); err != nil {
		return nil, err
	}
	return s.dbRepo.Update(id, updateRequest)
}

func (s *service) Create(p property.Property) (*property.Property, rest_errors.RestErr) {
	if err := p.Validate(); err != nil {
		return nil, err
	}

	p.Status = property.STATUS_ACTIVE
	p.DateCreated = date_utils.GetNowDBFromat()
	newProperty, err := s.dbRepo.Create(p)
	if err != nil {
		return nil, err
	}

	return newProperty, nil
}

func (s *service) Get() (property.Properties, rest_errors.RestErr) {
	return s.dbRepo.Get()
}

func (s *service) GetByID(id string) (*property.Property, rest_errors.RestErr) {
	return s.dbRepo.GetByID(id)
}

func (s *service) Search(query query.EsQuery) (property.Properties, rest_errors.RestErr) {
	return s.dbRepo.Search(query)
}

func (s *service) UploadMedia(files []*multipart.FileHeader, propertyID string) rest_errors.RestErr {
	p, err := s.dbRepo.GetByID(propertyID)
	if err != nil {
		return err
	}
	visuals := p.Visuals
	videos := p.Videos
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			return rest_errors.NewInternalServerErr("Error while trying to open the file", nil)
		}
		defer f.Close()
		var visual property.Visual
		var video property.Video
		v, ext, fileErr := file_utils.SaveFile(file, f, propertyID)
		if fileErr != nil {
			return fileErr
		}
		if ext == "mp4" || ext == "mov" {
			video.Url = "http://localhost:3030/assets/" + p.ID + "/" + v
			video.FileType = ext
			videos = append(videos, video)
		} else {
			visual.Url = "http://localhost:3030/assets/" + p.ID + "/" + v
			visual.FileType = ext
			visuals = append(visuals, visual)
		}
	}

	return s.dbRepo.UploadMedia(visuals, videos, propertyID)
	// LOGIC
	// some saved
	// first one throwed error the others will not be uploaded
}

func (s *service) DeleteMedia(propertyID string, mediaID string) rest_errors.RestErr {
	p, err := s.GetByID(propertyID)
	if err != nil {
		return err
	}

	file_utils.DeleteFile(mediaID, propertyID)

	var visuals []property.Visual
	var videos []property.Video

	for _, v := range p.Visuals {
		if v.Url == "http://localhost:3030/assets/"+p.ID+"/"+mediaID {
			continue
		}
		visuals = append(visuals, v)
	}

	for _, v := range p.Videos {
		if v.Url == "http://localhost:3030/assets/"+p.ID+"/"+mediaID {
			continue
		}
		videos = append(videos, v)
	}

	return s.dbRepo.UploadMedia(visuals, videos, propertyID)
}

func (srv *service) UploadProperyPic(propertyID string, fileHeader *multipart.FileHeader) (*property.Property, rest_errors.RestErr) {
	p, err := srv.GetByID(propertyID)
	if err != nil {
		return nil, err
	}
	if p.PropertyPic != "" {
		file_utils.DeleteFile(p.PropertyPic, propertyID)
		p.PropertyPic = ""
	}
	file, fErr := fileHeader.Open()
	if fErr != nil {
		return nil, rest_errors.NewInternalServerErr("Error while trying to open the file", nil)
	}
	filePath, _, fileErr := file_utils.SaveFile(fileHeader, file, propertyID)
	if fileErr != nil {
		return nil, fileErr
	}
	p.PropertyPic = "http://localhost:3030/assets/" + p.ID + "/" + filePath

	var es property.EsUpdate
	update := property.UpdatePropertyRequest{
		Field: "property_pic",
		Value: p.PropertyPic,
	}

	es.Fields = append(es.Fields, update)

	srv.dbRepo.Update(propertyID, es)

	return p, nil
}
