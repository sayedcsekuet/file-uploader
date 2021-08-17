package models

import (
	"database/sql"
	"encoding/json"
	"gorm.io/datatypes"
	"time"
)

type MetaData struct {
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
}

//NullTime is a wrapper around sql.NullTime
type NullTime struct {
	sql.NullTime
}

//MarshalJSON method is called by json.Marshal,
//whenever it is of type NullString
func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

type File struct {
	ID         string         `json:"id" validate:"uuid4,required"`
	Name       string         `json:"name"`
	MetaData   datatypes.JSON `sql:"type:json" json:"meta_data,omitempty"`
	OwnerID    string         `json:"owner_id" validate:"required"`
	BucketPath string         `json:"bucket_path" validate:"required"`
	Provider   string         `json:"provider" validate:"required"`
	CreatedAt  time.Time      `json:"created_at"`
	ExpiredAt  NullTime       `json:"expired_at,omitempty" validate:"datetime=2006-01-02T15:04:05Z07:00"`
}

func NewFile(id, name, owner, bucketPath, provider string, metaData MetaData) *File {
	metaDatJson, _ := json.Marshal(metaData)
	return &File{
		ID:         id,
		Name:       name,
		MetaData:   metaDatJson,
		OwnerID:    owner,
		BucketPath: bucketPath,
		Provider:   provider,
	}
}
