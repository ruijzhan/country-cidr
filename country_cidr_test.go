package country_cidr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCountryContains(t *testing.T) {
	cases := []struct {
		Country  string
		IP       string
		Contains bool
	}{
		{
			"CN",
			"114.114.114.114",
			true,
		},
		{
			"CN",
			"8.8.8.8",
			false,
		},
		{
			"NO_SUCH_COUNTRY",
			"8.8.8.8",
			false,
		},
		{
			"CN",
			"invalid_ip",
			false,
		},
		{
			"LAN",
			"192.168.1.183",
			true,
		},
		{
			"LAN",
			"10.100.1.1",
			true,
		},
		{
			"US",
			"52.4.20.24",
			false,
		},
		{
			"CN",
			"118.184.26.113",
			false,
		},
	}
	for _, c := range cases {
		assert.Equal(t, Country(c.Country).ContainsIPstr(c.IP), c.Contains, fmt.Sprintf("Case: %s %s %t", c.Country, c.IP, c.Contains))
	}
}

func TestFrom(t *testing.T) {
	cases := []struct {
		IP      string
		Country string
	}{
		{
			"114.114.114.114",
			"CN",
		},
		{
			"192.168.1.183",
			"LAN",
		},
		{
			"1.1.1.1",
			"AU",
		},
	}

	for _, c := range cases {
		r, err := From(c.IP)
		assert.NoError(t, err)
		assert.Equal(t, r, c.Country)
	}
}
