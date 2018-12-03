package cmd

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Document ..
type Document struct {
	ID             string    `json:"id" structs:"id"`
	FileName       string    `json:"FileName" structs:"filename"`
	Extension      string    `json:"Extension" structs:"extension"`
	Creation       time.Time `json:"creation" structs:"creation"`
	TEXT           string    `json:"text" structs:"text"`
	CLASSIFICATION string    `json:"classification" structs:"classification"`
	IMAGE          string    `json:"iamge" structs:"iamge"`
}

// GenerateNewUUID ..
func (*Document) GenerateNewUUID() string {
	u1 := uuid.Must(uuid.NewV4())
	return u1.String()
}
