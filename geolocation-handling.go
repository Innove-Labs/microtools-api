package main

import (
	"log"
	"net"
	"github.com/oschwald/geoip2-golang"
	"errors"
)

type GeoIPResponse struct {
	IP       string  `json:"ip"`
	Country  string  `json:"country"`
	Region   string  `json:"region"`
	City     string  `json:"city"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone string  `json:"timezone"`
}

func ValidateIP(ipStr string) (GeoIPResponse, error)  {
	ip := net.ParseIP(ipStr)
		if ip == nil {
			return GeoIPResponse{}, errors.New("Invalid IP address")
		}

	db, err := geoip2.Open("./geolite-2-city.mmdb")
		if err != nil {
			log.Fatal(err)
		}
	defer db.Close()

		record, err := db.City(ip)
		if err != nil {
			return GeoIPResponse{}, errors.New("IP address not found")
		}


		resp := GeoIPResponse{
			IP:       ipStr,
			Country:  record.Country.Names["en"],
			Region:   "",
			City:     record.City.Names["en"],
			Latitude: record.Location.Latitude,
			Longitude: record.Location.Longitude,
			Timezone: record.Location.TimeZone,
		}
		if len(record.Subdivisions) > 0 {
			resp.Region = record.Subdivisions[0].Names["en"]
		}

		return resp, nil
	
}