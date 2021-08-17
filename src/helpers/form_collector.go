package helpers

import (
	"file-uploader/src/models"
	"github.com/go-playground/form"
	"github.com/labstack/echo"
	"mime/multipart"
	"regexp"
)

var decoder = form.NewDecoder()

type FormCollector interface {
	CollectFromContext(c echo.Context) ([]*multipart.FileHeader, error)
	CollectFileCreateDataFromContext(c echo.Context) (*models.CreateFileData, error)
	Collect(form *multipart.Form) ([]*multipart.FileHeader, error)
	CollectFileCreateData(form *multipart.Form) (*models.CreateFileData, error)
}
type fileCollector struct {
}

func NewFileCollector() FormCollector {
	return &fileCollector{}
}
func (fc *fileCollector) CollectFromContext(c echo.Context) ([]*multipart.FileHeader, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}
	return fc.Collect(form)
}
func (fc *fileCollector) Collect(form *multipart.Form) ([]*multipart.FileHeader, error) {
	var files []*multipart.FileHeader
	for _, file := range form.File {
		files = append(files, file...)
	}
	return files, nil
}

func (fc *fileCollector) CollectFileCreateDataFromContext(c echo.Context) (*models.CreateFileData, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}
	return fc.CollectFileCreateData(form)
}
func (fc *fileCollector) CollectFileCreateData(form *multipart.Form) (*models.CreateFileData, error) {
	fileData := models.CreateFileData{}
	err := decoder.Decode(&fileData, form.Value)
	if err != nil {
		return nil, err
	}

	var files = map[string]*models.FileData{}
	re := regexp.MustCompile(`[0-9]+`)
	for key, fis := range form.File {
		i := re.FindStringSubmatch(key)
		files[key] = &models.FileData{
			ID:         "",
			FileHeader: fis[0],
		}
		if len(i) > 0 {
			if fd, ok := fileData.Files[i[0]]; ok {
				files[key].ID = fd.ID
			}
		}
	}
	fileData.Files = files
	return &fileData, nil
}
