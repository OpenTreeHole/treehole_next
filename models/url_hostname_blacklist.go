package models

type UrlHostnameBlacklist struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Hostname string `json:"hostname" gorm:"size:255;not null"`
}
