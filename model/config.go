package model

type Config struct {
	HTTP     Http     `json:"http"`
	Postgres Postgres `json:"postgres"`
}

type Http struct {
	Port string `json:"port"`
}

type Postgres struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Database  string `json:"database"`
	SSL       string `json:"ssl"`
	Migration string `json:"migration"`
}
