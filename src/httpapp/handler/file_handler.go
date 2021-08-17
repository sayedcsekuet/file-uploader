package handler

import (
	"file-uploader/src/helpers"
	httpapp "file-uploader/src/httpapp/response"
	"file-uploader/src/repositories"
	"file-uploader/src/services/fileservice"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type FileHandler interface {
	List(c echo.Context) error
	Get(c echo.Context) error
	Create(c echo.Context) error
	Delete(c echo.Context) error
	Download(c echo.Context) error
	Stream(c echo.Context) error
}
type fileHandler struct {
	FileCollector helpers.FormCollector
	FileService   fileservice.FileService
}

func NewFileHandler(fc helpers.FormCollector, fs fileservice.FileService) FileHandler {
	return &fileHandler{fc, fs}
}

type UploadResponse struct {
	Id       string `json:"id,omitempty"`
	FileName string `json:"file_name"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
}

func (fh *fileHandler) List(c echo.Context) error {
	name := c.QueryParam("name")
	createdDate := c.QueryParam("created_date")
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	files, err := fh.FileService.List(
		c.Request().Header.Get("x-api-key"),
		repositories.SearchParam{Name: name, CreatedDate: createdDate, Offset: offset, Limit: limit},
	)
	if err != nil {
		return httpapp.ErrorResponse(err)
	}
	return c.JSON(http.StatusOK, files)
}

func (fh *fileHandler) Create(c echo.Context) error {
	cFileData, err := fh.FileCollector.CollectFileCreateDataFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if len(cFileData.Files) == 0 {
		return c.JSON(http.StatusBadRequest, "File not found in the request")
	}
	response := fh.FileService.Create(c.Request().Header.Get("x-api-key"), cFileData)
	status := http.StatusOK
	for _, r := range response {
		if r.Virus {
			status = http.StatusMultiStatus
		}
	}
	return c.JSON(status, response)
}
func (fh *fileHandler) Get(c echo.Context) error {
	file, err := fh.FileService.GetByOwner(c.Param("id"), c.Request().Header.Get("x-api-key"))
	if err != nil {
		return httpapp.ErrorResponse(err)
	}
	return c.JSON(http.StatusOK, file)
}
func (fh *fileHandler) Delete(c echo.Context) error {
	err := fh.FileService.Delete(c.Param("id"), c.Request().Header.Get("x-api-key"))
	if err != nil {
		return httpapp.ErrorResponse(err)
	}
	return c.JSON(http.StatusOK, "OK")
}

func (fh *fileHandler) Download(c echo.Context) error {
	fileReader, err := fh.FileService.FileReader(
		c.Param("id"),
		c.Request().Header.Get("x-api-key"),
		c.QueryParam("token"),
	)
	if err != nil {
		return httpapp.ErrorResponse(err)
	}
	disposition := c.QueryParam("disposition")
	if disposition == "" {
		disposition = "attachment"
	}
	defer fileReader.Close()
	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("%s; filename=%q", disposition, fileReader.Name))
	return c.Stream(http.StatusOK, fileReader.ContentType(), fileReader)
}

func (fh *fileHandler) Stream(c echo.Context) error {
	fileReader, err := fh.FileService.FileReader(
		c.Param("id"),
		c.Request().Header.Get("x-api-key"),
		c.QueryParam("token"),
	)
	if err != nil {
		return httpapp.ErrorResponse(err)
	}
	defer fileReader.Close()
	return c.Stream(http.StatusOK, fileReader.ContentType(), fileReader)
}
