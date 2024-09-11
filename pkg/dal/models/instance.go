package models

import "gorm.io/gorm"

var (
	Migrate = []interface{}{
		&InstanceObject{},
		&AboutInstanceObject{},
		&UnprocessedPeer{},
		&UnprocessedModeration{},
	}
)

func NewInstanceObject(domain string) *InstanceObject {
	return &InstanceObject{
		Domain:           domain,
		About:            &AboutInstanceObject{},
		UnprocessedPeers: make([]*UnprocessedPeer, 0),
	}
}

type InstanceObject struct {
	gorm.Model

	Domain string               `gorm:"type:varchar(253);uniqueIndex;not null" json:"domain"`
	About  *AboutInstanceObject `gorm:"foreignKey:InstanceObject" json:"about"`

	UnprocessedPeers       []*UnprocessedPeer       `gorm:"foreignKey:InstanceObject"`
	UnprocessedModerations []*UnprocessedModeration `gorm:"foreignKey:InstanceObject"`
}

type AboutInstanceObject struct {
	gorm.Model

	InstanceObject uint

	Title       string
	Version     string
	SourceURL   string
	Description string
	Thumbnail   string
}

type UnprocessedPeer struct {
	gorm.Model

	InstanceObject uint
	Processed      bool
	Peer           string `gorm:"varchar(253)"`
}

type UnprocessedModeration struct {
	gorm.Model

	InstanceObject uint
	Processed      bool
	Moderated      string `gorm:"varchar(253)"`
	Digest         string
	Severity       string
	Comment        string
}
