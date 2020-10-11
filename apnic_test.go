package country_cidr

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

const apnicNum = 77

var lan = countryT{
	Name: "LAN",
	Cidrs: []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	},
}

func Test_Line2CidrV4(t *testing.T) {
	cases := []struct {
		line    string
		country string
		cidr    string
		valid   bool
	}{
		{"apnic|CN|ipv4|43.254.228.0|1024|20140729|allocated", "CN", "43.254.228.0/22", true},
		{"apnic|TW|asn|7532|8|19970322|allocated", "", "", false},
		{"apnic|AU|ipv6|2401:700::|32|20110606|allocated", "", "", false},
		{"# statement of the location in which any specific resource may", "", "", false},
	}

	for _, c := range cases {
		country, cidr, valid := line2CidrV4(c.line)
		assert.Equal(t, c.country, country)
		assert.Equal(t, c.cidr, cidr)
		assert.Equal(t, c.valid, valid)
	}
}

func Test_Bindata(t *testing.T) {
	jsonBytes, err := Asset(dataFile)
	assert.NoError(t, err)

	var cs []countryT
	err = json.Unmarshal(jsonBytes, &cs)
	assert.NoError(t, err)
	assert.Equal(t, apnicNum+1, len(cs)) //lan included
}

func Test_ParseAPNIC(t *testing.T) {
	apnicURL := "https://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest"

	resp, err := http.Get(apnicURL)
	assert.NoError(t, err)
	defer resp.Body.Close()
	apnic := apnicParse(resp.Body)
	assert.Equal(t, apnicNum, len(apnic))

	apnic = append(apnic, lan)
	jsonBytes, err := json.MarshalIndent(apnic, "", "    ")
	assert.NoError(t, err)

	if _, ok := os.LookupEnv("GEN_DATA_FILE"); ok {
		f, err := os.Create("apnic.json")
		assert.NoError(t, err)
		defer f.Close()
		f.Write(jsonBytes)
	}
}
