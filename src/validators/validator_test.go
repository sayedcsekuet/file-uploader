package validators

import (
	"file-uploader/src/models"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestValidateStructHavingError(t *testing.T) {
	err := NewAppValidator().Validate(models.File{
		ID:         "",
		Name:       "",
		MetaData:   nil,
		OwnerID:    "",
		BucketPath: "",
		Provider:   "",
		CreatedAt:  time.Time{},
		ExpiredAt:  models.NullTime{},
	})
	assert.Equal(t, 4, len(err.(validator.ValidationErrors)))
}

func TestValidateStruct(t *testing.T) {
	err := NewAppValidator().Validate(models.File{
		ID:         "7e26b96f-9ff5-408e-90cf-eb0a9a7c05bf",
		Name:       "",
		MetaData:   nil,
		OwnerID:    "adsf",
		BucketPath: "sadf",
		Provider:   "sadf",
		CreatedAt:  time.Time{},
		ExpiredAt:  models.NullTime{},
	})
	assert.NoError(t, err)
}
func TestValidateValueHavingError(t *testing.T) {
	err := NewAppValidator().ValidateValue(2, "gt=5")
	assert.Error(t, err)
}
