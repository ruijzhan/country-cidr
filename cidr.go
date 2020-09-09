package country_cidr

import (
	"encoding/json"
	"fmt"
	"github.com/yl2chen/cidranger"
	"log"
	"net"
)

const dataFile = "apnic.json"

type myRanger struct {
	cidranger.Ranger
}

func init() {
	jsonBytes, err := Asset(dataFile)
	if err != nil {
		log.Fatal(err)
	}

	var countries []countryT
	err = json.Unmarshal(jsonBytes, &countries)
	if err != nil {
		log.Fatal(err)
	}
	apnicMap = make(map[string]myRanger)

	for _, country := range countries {
		ranger := cidranger.NewPCTrieRanger()
		for _, cidr := range country.Cidrs {
			_, nw, err := net.ParseCIDR(cidr)
			if err != nil {
				log.Fatal(err)
			}
			ranger.Insert(cidranger.NewBasicRangerEntry(*nw))
		}
		apnicMap[country.Name] = myRanger{ranger}
	}
}

var (
	apnicMap map[string]myRanger
)

func Country(c string) myRanger {
	if ranger, ok := apnicMap[c]; ok {
		return ranger
	}
	return myRanger{cidranger.NewPCTrieRanger()}
}

func (r myRanger) ContainsIPstr(ip string) bool {
	return r.Contains(net.ParseIP(ip))
}

func (r myRanger) Contains(ip net.IP) bool {
	contains, err := r.Ranger.Contains(ip)
	if err != nil {
		return false
	}
	return contains
}

func From(ip string) (string, error) {
	for country, _ := range apnicMap {
		if Country(country).ContainsIPstr(ip) {
			return country, nil
		}
	}
	return "", fmt.Errorf("no country has this IP %s", ip)
}
