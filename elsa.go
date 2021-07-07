package elsa

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	advisoryURL = "https://www.nhc.noaa.gov/text/refresh/MIATCPAT5+shtml/041158.shtml"
)

func Summary() (string, error) {
	resp, err := http.Get(advisoryURL)
	if err != nil {
		return "", fmt.Errorf("cannot get report from %s: %w", advisoryURL, err)
	}

	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot get body from %s: %w", advisoryURL, err)
	}

	s := string(text)

	idx1 := strings.Index(s, "SUMMARY OF")
	idx2 := strings.Index(s, "WATCHES AND WARNINGS")
	if idx1 > 0 && idx2 > idx1 {
		s = s[idx1:idx2]
	}

	return s, nil
}

// Distance in miles from hurricane Elsa to the point.
func Distance(ourLat, ourLon float64) (float64, error) {
	s, err := Summary()
	if err != nil {
		return 0, fmt.Errorf("cannot get summary: %w", err)
	}

	re := regexp.MustCompile("LOCATION...([0-9]+\\.[0-9])N ([0-9]+\\.[0-9])W")
	lines := strings.Split(s, "\n")

	for _, line := range lines {
		if !strings.HasPrefix(line, "LOCATION") {
			continue
		}

		m := re.FindStringSubmatch(line)
		if len(m) < 3 {
			return 0, fmt.Errorf("invalid LOCATION")
		}

		lat, err := strconv.ParseFloat(m[1], 64)
		if err != nil {
			return 0, fmt.Errorf("cannot parse latitude: %w", err)
		}

		lon, err := strconv.ParseFloat(m[2], 64)
		if err != nil {
			return 0, fmt.Errorf("cannot parse longitude: %w", err)
		}
		lon = -lon

		return haversineDistance(ourLat, ourLon, lat, lon), nil
	}

	return 0, fmt.Errorf("no LOCATION in summary at %s", advisoryURL)
}

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	rad := math.Pi / 180

	la1 := lat1 * rad
	lo1 := lon1 * rad
	la2 := lat2 * rad
	lo2 := lon2 * rad

	const earthRadius = 3958.8 // miles

	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)
	return 2 * earthRadius * math.Asin(math.Sqrt(h))
}
