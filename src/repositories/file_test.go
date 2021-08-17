package repositories

import (
	"database/sql"
	"file-uploader/src/models"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
)

type FileRepositorySuite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository FileRepository
	person     *models.File
}

func (s *FileRepositorySuite) currentDate() string {
	return time.Now().Format(DateTimeLayout)
}

func (s *FileRepositorySuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)
	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	require.NoError(s.T(), err)
	s.repository = NewFileRepository(s.DB)
}

func (s *FileRepositorySuite) AfterTest(_, _ string) {
	assert.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *FileRepositorySuite) TestRepositoryGet() {
	var (
		id   = uuid.NewString()
		name = "test-name"
	)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `files` WHERE id=? AND (expired_at IS NULL OR expired_at>=?)")).
		WithArgs(id, s.currentDate()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(id, name))

	res, err := s.repository.Get(id)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), &models.File{ID: id, Name: name}, res)
}

func (s *FileRepositorySuite) TestRepositoryGetNotFound() {
	var (
		id   = uuid.NewString()
		name = "test-name"
	)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `files` WHERE id=? AND (expired_at IS NULL OR expired_at>=?)")).
		WithArgs(id, s.currentDate()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow("", name))

	_, err := s.repository.Get(id)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), fmt.Sprintf("File not found with id %s!", id), err.Error())
}

func (s *FileRepositorySuite) TestRepositoryCreate() {
	var (
		id       = uuid.NewString()
		name     = "test-name"
		metaData = models.MetaData{
			MimeType: "",
			Size:     0,
		}
		bucketPath = "test"
		ownerId    = "test"
		provider   = "s3"
	)
	model := models.NewFile(id, name, ownerId, bucketPath, provider, metaData)
	model.CreatedAt = time.Now()
	s.mock.ExpectBegin()
	s.mock.
		ExpectExec(regexp.QuoteMeta(
			"INSERT INTO `files` (`id`,`name`,`meta_data`,`owner_id`,`bucket_path`,`provider`,`created_at`,`expired_at`) VALUES (?,?,?,?,?,?,?,?)")).
		WithArgs(id, name, `{"mime_type":"","size":0}`, ownerId, bucketPath, provider, model.CreatedAt, nil).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectCommit()

	res, err := s.repository.Create(model)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), model, res)
}

func (s *FileRepositorySuite) TestRepositoryGetAll() {
	var (
		id   = uuid.NewString()
		name = "test-name"
	)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `files` WHERE expired_at IS NULL OR expired_at>=?")).
		WithArgs(s.currentDate()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(id, name))

	res, err := s.repository.GetAll()

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), []*models.File{{ID: id, Name: name}}, res)
}

func (s *FileRepositorySuite) TestRepositoryFindByOwnerAndId() {
	var (
		id      = uuid.NewString()
		name    = "test-name"
		ownerId = "test"
	)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `files` WHERE id=? AND owner_id = ? AND (expired_at IS NULL OR expired_at>=?)")).
		WithArgs(id, ownerId, s.currentDate()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "owner_id"}).
			AddRow(id, name, ownerId))

	res, err := s.repository.FindByOwnerAndId(id, ownerId)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), &models.File{ID: id, Name: name, OwnerID: ownerId}, res)
}

func (s *FileRepositorySuite) TestRepositoryNotFoundFindByOwnerAndId() {
	var (
		id      = uuid.NewString()
		name    = "test-name"
		ownerId = "test"
	)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `files` WHERE id=? AND owner_id = ? AND (expired_at IS NULL OR expired_at>=?)")).
		WithArgs(id, ownerId, s.currentDate()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "owner_id"}).
			AddRow("", name, ownerId))

	_, err := s.repository.FindByOwnerAndId(id, ownerId)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), fmt.Sprintf("File not found with id %s and with owner XXXXXXX. Please verify you are the correct owner of that file or file doesn't exist!", id), err.Error())
}

func (s *FileRepositorySuite) TestRepositoryGetAllByOwner() {
	var (
		id      = uuid.NewString()
		name    = "test-name"
		ownerId = "test"
	)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `files` WHERE owner_id = ? AND (expired_at IS NULL OR expired_at >= ?) LIMIT 500")).
		WithArgs(ownerId, s.currentDate()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "owner_id"}).
			AddRow(id, name, ownerId))

	res, err := s.repository.GetAllByOwner(ownerId, SearchParam{})

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), []*models.File{{ID: id, Name: name, OwnerID: ownerId}}, res)
}

func (s *FileRepositorySuite) TestRepositoryGetExpiredFiles() {
	var (
		id      = uuid.NewString()
		name    = "test-name"
		ownerId = "test"
	)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `files` WHERE expired_at < ?")).
		WithArgs(s.currentDate()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "owner_id"}).
			AddRow(id, name, ownerId))

	res, err := s.repository.GetExpiredFiles()

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), []*models.File{{ID: id, Name: name, OwnerID: ownerId}}, res)
}

func (s *FileRepositorySuite) TestRepositoryDelete() {
	var (
		id      = uuid.NewString()
		ownerId = "test"
	)
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		"DELETE FROM `files` WHERE id=? AND owner_id = ?")).
		WithArgs(id, ownerId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.repository.Delete(id, ownerId)

	assert.NoError(s.T(), err)
}
func (s *FileRepositorySuite) TestRepositoryDeleteAll() {
	ids := []string{uuid.NewString(), uuid.NewString()}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		"DELETE FROM `files` WHERE `files`.`id` IN (?,?)")).
		WithArgs(ids[0], ids[1]).
		WillReturnResult(sqlmock.NewResult(1, 2))
	s.mock.ExpectCommit()
	err := s.repository.DeleteAll(ids)

	assert.NoError(s.T(), err)
}
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(FileRepositorySuite))
}
