package property

import (
	"mime/multipart"
	"os"

	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_property-api/constants"
	"github.com/superbkibbles/realestate_property-api/domain/property"
	"github.com/superbkibbles/realestate_property-api/domain/query"
	"github.com/superbkibbles/realestate_property-api/repository/db"
	"github.com/superbkibbles/realestate_property-api/utils/date_utils"
	"github.com/superbkibbles/realestate_property-api/utils/file_utils"
)

type Service interface {
	Create(property.Property) (*property.Property, rest_errors.RestErr)
	Get(sort string, asc bool, local string) (property.Properties, rest_errors.RestErr)
	GetByID(string, local string) (*property.Property, rest_errors.RestErr)
	Search(query query.EsQuery, sort string, asc bool, local string) (property.Properties, rest_errors.RestErr)
	Update(id string, updateRequest property.EsUpdate) (*property.Property, rest_errors.RestErr)
	UploadMedia(files []*multipart.FileHeader, propertyID string) rest_errors.RestErr
	DeleteMedia(propertyID string, mediaID string) rest_errors.RestErr
	UploadProperyPic(id string, fileHeader *multipart.FileHeader) (*property.Property, rest_errors.RestErr)
	GetActive(sort string, asc bool, local string) (property.Properties, rest_errors.RestErr)
	GetDeactive(sort string, asc bool, local string) (property.Properties, rest_errors.RestErr)
	Translate(id string, translateProperty property.TranslateProperty, local string) (*property.Property, rest_errors.RestErr)
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

func (s *service) Translate(id string, translateProperty property.TranslateProperty, local string) (*property.Property, rest_errors.RestErr) {
	translateProperty.Local = local
	translateProperty.PropertyID = id
	if err := translateProperty.Validate(); err != nil {
		return nil, err
	}
	p, err := s.dbRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	ts, err := s.dbRepo.GetTranslateById(id, local)
	if err != nil {
		return nil, err
	}
	if ts.PropertyID != "" {
		// Update The property if there is record
		tpp, err := s.dbRepo.UpdateTranslate(ts.ID, translateProperty)
		if err != nil {
			return nil, err
		}
		return tpp.Marshal(p), nil
	}

	// Insert If there is no record
	if err := s.dbRepo.Translate(translateProperty); err != nil {
		return nil, err
	}

	return translateProperty.Marshal(p), nil
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

func (s *service) Get(sort string, asc bool, local string) (property.Properties, rest_errors.RestErr) {
	properties, err := s.dbRepo.Get(sort, asc)
	if err != nil {
		return nil, err
	}
	if local == "en" || local == "" {
		return properties, nil
	}

	ts, err := s.dbRepo.GetAllTranslated(local)
	if err != nil {
		return nil, err
	}

	return ts.Marshal(properties), nil
}

func (s *service) GetActive(sort string, asc bool, local string) (property.Properties, rest_errors.RestErr) {
	properties, err := s.dbRepo.GetActive(sort, asc)
	if err != nil {
		return nil, err
	}
	if local == "en" || local == "" {
		return properties, nil
	}
	ts, err := s.dbRepo.GetActiveLocal(local)
	if err != nil {
		return nil, err
	}

	return ts.Marshal(properties), nil
}

func (s *service) GetDeactive(sort string, asc bool, local string) (property.Properties, rest_errors.RestErr) {
	properties, err := s.dbRepo.GetDeactive(sort, asc)
	if err != nil {
		return nil, err
	}
	if local == "en" || local == "" {
		return properties, nil
	}
	ts, err := s.dbRepo.GetDeactiveLocal(local)
	if err != nil {
		return nil, err
	}

	return ts.Marshal(properties), nil
}

func (s *service) GetByID(id string, local string) (*property.Property, rest_errors.RestErr) {
	p, err := s.dbRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if local == "en" || local == "" {
		return p, nil
	}

	tp, err := s.dbRepo.GetTranslateById(id, local)
	if err != nil {
		return nil, err
	}
	return tp.Marshal(p), nil
}

func (s *service) Search(query query.EsQuery, sort string, asc bool, local string) (property.Properties, rest_errors.RestErr) {
	properties, err := s.dbRepo.Search(query, sort, asc)
	if err != nil {
		return nil, err
	}

	if local == "en" || local == "" {
		return properties, nil
	}
	ts, err := s.dbRepo.GetDeactiveLocal(local)
	if err != nil {
		return nil, err
	}

	return ts.Marshal(properties), nil
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
			video.Url = os.Getenv(constants.PUBLIC_API_KEY) + "assets/" + p.ID + "/" + v
			video.FileType = ext
			videos = append(videos, video)
		} else {
			visual.Url = os.Getenv(constants.PUBLIC_API_KEY) + "assets/" + p.ID + "/" + v
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
	p, err := s.dbRepo.GetByID(propertyID)
	if err != nil {
		return err
	}

	file_utils.DeleteFile(mediaID, propertyID)

	var visuals []property.Visual
	var videos []property.Video

	for _, v := range p.Visuals {
		if v.Url == os.Getenv(constants.PUBLIC_API_KEY)+"assets/"+p.ID+"/"+mediaID {
			continue
		}
		visuals = append(visuals, v)
	}

	for _, v := range p.Videos {
		if v.Url == os.Getenv(constants.PUBLIC_API_KEY)+"assets/"+p.ID+"/"+mediaID {
			continue
		}
		videos = append(videos, v)
	}

	return s.dbRepo.UploadMedia(visuals, videos, propertyID)
}

func (srv *service) UploadProperyPic(propertyID string, fileHeader *multipart.FileHeader) (*property.Property, rest_errors.RestErr) {
	p, err := srv.dbRepo.GetByID(propertyID)
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
	p.PropertyPic = os.Getenv(constants.PUBLIC_API_KEY) + "assets/" + p.ID + "/" + filePath

	var es property.EsUpdate
	update := property.UpdatePropertyRequest{
		Field: "property_pic",
		Value: p.PropertyPic,
	}

	es.Fields = append(es.Fields, update)

	srv.dbRepo.Update(propertyID, es)

	return p, nil
}
