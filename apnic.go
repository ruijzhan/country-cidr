package country_cidr

import (
	"bufio"
	"io"
	"math"
	"strconv"
	"strings"
)

type countryT struct {
	Name  string   `json:"country"`
	Cidrs []string `json:"cidrs"`
}

func line2CidrV4(line string) (country, cidr string, ok bool) {
	if !strings.HasPrefix(line, "apnic|") {
		return "", "", false
	}
	tokens := strings.Split(line, "|")
	if tokens[2] != "ipv4" {
		return "", "", false
	}
	mask, err := strconv.Atoi(tokens[4])
	if err != nil {
		return "", "", false
	}
	mask = 32 - int(math.Log(float64(mask))/math.Log(2))
	return tokens[1], tokens[3] + "/" + strconv.Itoa(mask), true
}

func apnicParse(r io.Reader) []countryT {
	m := make(map[string]*countryT)
	bReader := bufio.NewReader(r)

	for {
		b, _, err := bReader.ReadLine()
		if err == io.EOF {
			break
		}
		line := string(b)
		country, cidr, ok := line2CidrV4(line)
		if !ok || country == "*" {
			continue
		}

		if _, ok = m[country]; ok {
			m[country].Cidrs = append(m[country].Cidrs, cidr)
		} else {
			m[country] = &countryT{
				Name:  country,
				Cidrs: []string{cidr},
			}
		}
	}

	res := make([]countryT, 0, len(m))
	for _, v := range m {
		res = append(res, *v)
	}
	return res
}
