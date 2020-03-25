package models

import (
	"time"
)

type ProxyItem struct {
	ProxyId   int
	ProxyIp   string
	ProxyPort int
	Country   Country
	Anonymity string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Proxy from https://www.sslproxies.org/
type Proxy struct {
	ProxyId   int
	Ip        string
	Port      int
	CountryId int
	Anonymity string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Country struct {
	CountryId   int
	CountryName string
	CountryCode string
	CreatedAt   time.Time
}
