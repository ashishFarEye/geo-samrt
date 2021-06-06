package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Configurations To Be Exported
type Config struct {
	Server       Server
	LibpostalURL string
	Geocoder     Geocoder
	PG           PG
	S3           S3
	Loop		 Loop
}



type Geocoder struct {
	Google   Google
	Arcgis   Arcgis
	Pelias   Pelias
	Mapquest Mapquest
	ES       ES
}

type Google struct {
	APIkey string
	Qps    int
}

type Arcgis struct {
	AccessId  string
	AccessKey string
}

type Pelias struct {
	Host string
}

type PG struct {
	URL      string
	Username string
	Password string
	DB       string
}

//ServerConfigurations initialization
type Server struct {
	Port int
	Host string
}

type S3 struct {
	Bucket string
}

type Mapquest struct {
	ApiKey string
}

type ES struct {
	Host string
	User string
	Pwd  string
}

type Loop struct {
	Host string
}

func Init() Config {
	conf := Config{}
	err := envconfig.Process("maps", &conf)
	if err != nil {
		panic(err)
	}
	return conf
}
