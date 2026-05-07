package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/iwikus/geoip-service/geoip2"
	cache "github.com/patrickmn/go-cache"
)

type ResponseCity struct {
	Data  *geoip2.City `json:",omitempty"`
	Error string       `json:",omitempty"`
}

type ResponseCountry struct {
	Data  *geoip2.Country `json:",omitempty"`
	Error string          `json:",omitempty"`
}

func main() {
	var dbName      = flag.String("db", "GeoLite2-City.mmdb", "File name of MaxMind GeoIP2 and GeoLite2 database")
	var lookup      = flag.String("lookup", "city", "Specify which value to look up. Can be 'city' or 'country' depending on which database you load.")
	var listen      = flag.String("listen", ":5000", "Listen address and port, for instance 127.0.0.1:5000")
	var threads     = flag.Int("threads", runtime.NumCPU(), "Number of threads to use. Defaults to number of detected cores")
	var pretty      = flag.Bool("pretty", false, "Should output be formatted with newlines and indentation")
	var cacheSecs   = flag.Int("cache", 0, "How many seconds should requests be cached. Set to 0 to disable")
	var originPolicy = flag.String("origin", "*", `Value sent in the 'Access-Control-Allow-Origin' header. Set to "" to disable.`)

	serverStart := time.Now().Format(http.TimeFormat)
	flag.Parse()
	runtime.GOMAXPROCS(*threads)

	var memCache *cache.Cache
	if *cacheSecs > 0 {
		memCache = cache.New(time.Duration(*cacheSecs)*time.Second, 1*time.Second)
	}

	lookupCity := true
	if *lookup == "country" {
		lookupCity = false
	} else if *lookup != "city" {
		log.Fatalf("lookup parameter should be either 'city', or 'country', it is '%s'", *lookup)
	}

	db, err := geoip2.Open(*dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Loaded database " + *dbName)

	prettyL := *pretty

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// Set headers
		if *originPolicy != "" {
			w.Header().Set("Access-Control-Allow-Origin", *originPolicy)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Last-Modified", serverStart)

		ipText := req.URL.Query().Get("ip")
		if ipText == "" {
			ipText = strings.Trim(req.URL.Path, "/")
		}

		ip := net.ParseIP(ipText)
		if ip == nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"Error": "unable to decode ip"})
			return
		}

		// Cache hit?
		if memCache != nil {
			if v, found := memCache.Get(ipText); found {
				w.Write(v.([]byte))
				return
			}
		}

		var j []byte
		if lookupCity {
			result, err := db.City(ip)
			res := ResponseCity{Data: result}
			if err != nil {
				res.Data = nil
				res.Error = err.Error()
			}
			if prettyL {
				j, err = json.MarshalIndent(res, "", "  ")
			} else {
				j, err = json.Marshal(res)
			}
		} else {
			result, err := db.Country(ip)
			res := ResponseCountry{Data: result}
			if err != nil {
				res.Data = nil
				res.Error = err.Error()
			}
			if prettyL {
				j, err = json.MarshalIndent(res, "", "  ")
			} else {
				j, err = json.Marshal(res)
			}
		}

		if err != nil {
			log.Println("json marshal error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if memCache != nil {
			memCache.Set(ipText, j, 0)
		}

		w.Write(j)
	})

	log.Println("Listening on " + *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
