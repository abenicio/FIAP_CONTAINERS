package cmd

import (
	"time"
)

// Document ..
type DocumentReceive struct {
	ID             string    `json:"id" structs:"id"`
	FileName       string    `json:"FileName" structs:"filename"`
	Extension      string    `json:"Extension" structs:"extension"`
	Creation       time.Time `json:"creation" structs:"creation"`
	CLASSIFICATION string    `json:"classification" structs:"classification"`
	IMAGE          string    `json:"iamge" structs:"iamge"`
}

type Files struct {
	FileName       string `json:"FileName" structs:"filename"`
	CLASSIFICATION string `json:"classification" structs:"classification"`
}
