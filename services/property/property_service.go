package property

import (
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_property-api/domain/property"
	"github.com/superbkibbles/realestate_property-api/domain/query"
	cloudstorage "github.com/superbkibbles/realestate_property-api/repository/cloudStorage"
	"github.com/superbkibbles/realestate_property-api/repository/db"
	"github.com/superbkibbles/realestate_property-api/utils/crypto_utils"
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
	dbRepo    db.DbRepository
	cloudRepo cloudstorage.CloudStorage
}

func NewService(dbRepo db.DbRepository, cloudRepo cloudstorage.CloudStorage) Service {
	return &service{
		dbRepo:    dbRepo,
		cloudRepo: cloudRepo,
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
		res, cloudErr := s.cloudRepo.Save(f, propertyID+crypto_utils.GetMd5(uuid.New().String()), p.ID)
		if err != nil {
			return cloudErr
		}
		if res.Url != "" {
			v := res.Url
			ext := res.Ext
			publicID := res.PublicID

			if ext == "mp4" || ext == "mov" {
				video.Url = v
				video.FileType = ext
				video.PublicID = publicID
				videos = append(videos, video)
			} else {
				visual.Url = v
				visual.FileType = ext
				visual.PublicID = publicID
				visuals = append(visuals, visual)
			}
		}
	}

	return s.dbRepo.UploadMedia(visuals, videos, propertyID)
}

func (s *service) DeleteMedia(propertyID string, mediaID string) rest_errors.RestErr {
	p, err := s.dbRepo.GetByID(propertyID)
	if err != nil {
		return err
	}

	var visuals []property.Visual
	var videos []property.Video

	for _, v := range p.Visuals {
		if v.PublicID == mediaID {
			if v.FileType != "mp4" && v.FileType != "mov" {
				if err := s.cloudRepo.Delete(mediaID); err != nil {
					return err
				}
				continue
			}
		}
		visuals = append(visuals, v)
	}

	for _, v := range p.Videos {
		if v.PublicID == mediaID {
			if v.FileType == "mp4" || v.FileType == "mov" {
				if err := s.cloudRepo.Delete(mediaID); err != nil {
					return err
				}
				continue
			}
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
	res, cloudErr := srv.cloudRepo.Save(file, propertyID+crypto_utils.GetMd5(uuid.New().String()), p.ID)
	filePath := res.Url
	if cloudErr != nil {
		return nil, cloudErr
	}
	p.PropertyPic = filePath

	var es property.EsUpdate
	update := property.UpdatePropertyRequest{
		Field: "property_pic",
		Value: p.PropertyPic,
	}

	es.Fields = append(es.Fields, update)

	srv.dbRepo.Update(propertyID, es)

	return p, nil
}
