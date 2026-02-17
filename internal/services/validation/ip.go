package validation

import (
	"errors"
	"log"
	"net"

	"github.com/innovelabs/microtools-go/internal/models"
	"github.com/oschwald/geoip2-golang"
)

// ValidateIP validates an IP address and returns geolocation information
func ValidateIP(ipStr string) (models.GeoIPResponse, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return models.GeoIPResponse{}, errors.New("Invalid IP address")
	}

	db, err := geoip2.Open("./assets/geolite-2-city.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	record, err := db.City(ip)
	if err != nil {
		return models.GeoIPResponse{}, errors.New("IP address not found")
	}

	resp := models.GeoIPResponse{
		IP:        ipStr,
		Country:   record.Country.Names["en"],
		Region:    "",
		City:      record.City.Names["en"],
		Latitude:  record.Location.Latitude,
		Longitude: record.Location.Longitude,
		Timezone:  record.Location.TimeZone,
	}
	if len(record.Subdivisions) > 0 {
		resp.Region = record.Subdivisions[0].Names["en"]
	}

	return resp, nil
}
