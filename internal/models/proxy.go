package models

import (
	"net"
	"time"
)

type ProxyItem struct {
	Ip      string
	Port    string
	Code    string
	Country string
	//Anonymity string
	//CreatedAt time.Time
	//UpdatedAt time.Time
}

// Proxy from https://www.sslproxies.org/
type Proxy struct {
	ProxyId   int
	Ip        net.IP
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
