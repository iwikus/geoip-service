// IMPORTED FROM github.com/oschwald/geoip2-golang
//
// See GEOIP-LIBRARY-LICENSE for license on this file
//
// Modified slightly for use with standard encoding/json
//
// Package geoip2 provides a wrapper around the maxminddb package for
// easy use with the MaxMind GeoIP2 and GeoLite2 databases. The records for
// the IP address is returned from this package as well-formed structures
// that match the internal layout of data from MaxMind.

package geoip2

import (
	"net"

	"github.com/oschwald/maxminddb-golang"
)

type (
	TheCity struct {
		GeoNameID uint              `maxminddb:"geoname_id" json:"GeoNameID"`
		Names     map[string]string `maxminddb:"names"      json:"Names"`
	}
	Continent struct {
		Code      string            `maxminddb:"code"       json:"Code"`
		GeoNameID uint              `maxminddb:"geoname_id" json:"GeoNameID"`
		Names     map[string]string `maxminddb:"names"      json:"Names"`
	}
	TheCountry struct {
		GeoNameID uint              `maxminddb:"geoname_id" json:"GeoNameID"`
		IsoCode   string            `maxminddb:"iso_code"   json:"IsoCode"`
		Names     map[string]string `maxminddb:"names"      json:"Names"`
	}
	Location struct {
		Latitude  float64 `maxminddb:"latitude"   json:"Latitude"`
		Longitude float64 `maxminddb:"longitude"  json:"Longitude"`
		MetroCode uint    `maxminddb:"metro_code" json:"MetroCode"`
		TimeZone  string  `maxminddb:"time_zone"  json:"TimeZone"`
	}
	Postal struct {
		Code string `maxminddb:"code" json:"Code"`
	}
	RegisteredCountry struct {
		GeoNameID uint              `maxminddb:"geoname_id" json:"GeoNameID"`
		IsoCode   string            `maxminddb:"iso_code"   json:"IsoCode"`
		Names     map[string]string `maxminddb:"names"      json:"Names"`
	}
	RepresentedCountry struct {
		GeoNameID uint              `maxminddb:"geoname_id" json:"GeoNameID"`
		IsoCode   string            `maxminddb:"iso_code"   json:"IsoCode"`
		Names     map[string]string `maxminddb:"names"      json:"Names"`
		Type      string            `maxminddb:"type"       json:"Type"`
	}
	Subdivision struct {
		GeoNameID uint              `maxminddb:"geoname_id" json:"GeoNameID"`
		IsoCode   string            `maxminddb:"iso_code"   json:"IsoCode"`
		Names     map[string]string `maxminddb:"names"      json:"Names"`
	}
	Traits struct {
		IsAnonymousProxy    bool `maxminddb:"is_anonymous_proxy"    json:"IsAnonymousProxy"`
		IsSatelliteProvider bool `maxminddb:"is_satellite_provider" json:"IsSatelliteProvider"`
	}
)

// City corresponds to the data in the GeoIP2/GeoLite2 City databases.
type City struct {
	City               TheCity            `maxminddb:"city"                json:"City"`
	Continent          Continent          `maxminddb:"continent"           json:"Continent"`
	Country            TheCountry         `maxminddb:"country"             json:"Country"`
	Location           Location           `maxminddb:"location"            json:"Location"`
	Postal             Postal             `maxminddb:"postal"              json:"Postal"`
	RegisteredCountry  RegisteredCountry  `maxminddb:"registered_country"  json:"RegisteredCountry"`
	RepresentedCountry RepresentedCountry `maxminddb:"represented_country" json:"RepresentedCountry"`
	Subdivisions       []Subdivision      `maxminddb:"subdivisions"        json:"Subdivisions"`
	Traits             Traits             `maxminddb:"traits"              json:"Traits"`
}

// Country corresponds to the data in the GeoIP2/GeoLite2 Country databases.
type Country struct {
	Continent          Continent          `maxminddb:"continent"           json:"Continent"`
	Country            TheCountry         `maxminddb:"country"             json:"Country"`
	RegisteredCountry  RegisteredCountry  `maxminddb:"registered_country"  json:"RegisteredCountry"`
	RepresentedCountry RepresentedCountry `maxminddb:"represented_country" json:"RepresentedCountry"`
	Subdivisions       []Subdivision      `maxminddb:"subdivisions"        json:"Subdivisions"`
	Traits             Traits             `maxminddb:"traits"              json:"Traits"`
}

// ConnectionType corresponds to the data in the GeoIP2 Connection-Type database.
type ConnectionType struct {
	ConnectionType string `maxminddb:"connection_type" json:"ConnectionType"`
}

// Domain corresponds to the data in the GeoIP2 Domain database.
type Domain struct {
	Domain string `maxminddb:"domain" json:"Domain"`
}

// ISP corresponds to the data in the GeoIP2 ISP database.
type ISP struct {
	AutonomousSystemNumber       uint   `maxminddb:"autonomous_system_number"       json:"AutonomousSystemNumber"`
	AutonomousSystemOrganization string `maxminddb:"autonomous_system_organization" json:"AutonomousSystemOrganization"`
	ISP                          string `maxminddb:"isp"                            json:"ISP"`
	Organization                 string `maxminddb:"organization"                   json:"Organization"`
}

// Reader holds the maxminddb.Reader structure.
type Reader struct {
	mmdbReader *maxminddb.Reader
}

// Open takes a path to a file and returns a Reader or an error.
func Open(file string) (*Reader, error) {
	reader, err := maxminddb.Open(file)
	return &Reader{mmdbReader: reader}, err
}

// FromBytes takes a byte slice and returns a Reader or an error.
func FromBytes(bytes []byte) (*Reader, error) {
	reader, err := maxminddb.FromBytes(bytes)
	return &Reader{mmdbReader: reader}, err
}

// City looks up the city data for the given IP address.
func (r *Reader) City(ipAddress net.IP) (*City, error) {
	var city City
	err := r.mmdbReader.Lookup(ipAddress, &city)
	return &city, err
}

// Country looks up the country data for the given IP address.
func (r *Reader) Country(ipAddress net.IP) (*Country, error) {
	var country Country
	err := r.mmdbReader.Lookup(ipAddress, &country)
	return &country, err
}

// ConnectionType looks up the connection type for the given IP address.
func (r *Reader) ConnectionType(ipAddress net.IP) (*ConnectionType, error) {
	var val ConnectionType
	err := r.mmdbReader.Lookup(ipAddress, &val)
	return &val, err
}

// Domain looks up the domain for the given IP address.
func (r *Reader) Domain(ipAddress net.IP) (*Domain, error) {
	var val Domain
	err := r.mmdbReader.Lookup(ipAddress, &val)
	return &val, err
}

// ISP looks up the ISP data for the given IP address.
func (r *Reader) ISP(ipAddress net.IP) (*ISP, error) {
	var val ISP
	err := r.mmdbReader.Lookup(ipAddress, &val)
	return &val, err
}

// Close releases the resources held by the Reader.
func (r *Reader) Close() {
	r.mmdbReader.Close()
}

