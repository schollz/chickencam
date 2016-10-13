package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// Show sunrise and sunset for first 5 days of June in LA
func main() {

	_, timeToSunset := GetSunset(0)
	_, timeToSunrise := GetSunrise(0)
	if len(os.Args[1:]) == 0 {
		fmt.Println(timeToSunrise.Hours(), timeToSunset.Hours())
	} else {
		numberOfDays := 200
		secondsToEvent := []string{}
		i := 0
		for days := 0; days < numberOfDays; days++ {
			_, timeToEvent := GetSunset(days)
			theTime := int(math.Floor(timeToEvent.Seconds() + 30*60))
			if theTime > 0 {
				secondsToEvent = append(secondsToEvent, strconv.Itoa(theTime))
				i++
			}
		}
		sunset := fmt.Sprintf("unsigned long sunset[%d] = {%s};", numberOfDays, strings.Join(secondsToEvent, ","))

		secondsToEvent = []string{}
		i = 0
		for days := 0; days < numberOfDays; days++ {
			_, timeToEvent := GetSunrise(days)
			theTime := int(math.Floor(timeToEvent.Seconds() - 30*60))
			if theTime > 0 {
				secondsToEvent = append(secondsToEvent, strconv.Itoa(theTime))
				i++
			}
		}
		sunrise := fmt.Sprintf("unsigned long sunrise[%d] = {%s};", len(secondsToEvent), strings.Join(secondsToEvent, ","))

		ioutil.WriteFile("addToIno.txt", []byte(sunset+"\n"+sunrise), 0644)
	}
}

// Taken from https://github.com/keep94/sunrise
func GetSunset(offset int) (time.Time, time.Duration) {
	var s Sunrise

	// Start time is June 1, 2013 PST
	location, _ := time.LoadLocation("America/New_York")
	startTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, location)
	startTime = startTime.AddDate(0, 0, offset)

	// Coordinates of LA are 34.05N 118.25W
	s.Around(35.994, -78.8986, startTime)
	for s.Sunrise().Before(startTime) {
		s.AddDays(1)
	}

	return s.Sunset(), s.Sunset().Sub(time.Now())
}

// Taken from https://github.com/keep94/sunrise
func GetSunrise(offset int) (time.Time, time.Duration) {
	var s Sunrise

	// Start time is June 1, 2013 PST
	location, _ := time.LoadLocation("America/New_York")
	startTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, location)
	startTime = startTime.AddDate(0, 0, offset)

	// Coordinates of LA are 34.05N 118.25W
	s.Around(35.994, -78.8986, startTime)
	for s.Sunrise().Before(startTime) {
		s.AddDays(1)
	}

	return s.Sunrise(), s.Sunrise().Sub(time.Now())
}

const (
	jepoch = float64(2451545.0)
	uepoch = int64(946728000.0)
)

// Sunrise gives sunrise and sunset times.
type Sunrise struct {
	location        *time.Location
	sinLat          float64
	cosLat          float64
	jstar           float64
	solarNoon       float64
	hourAngleInDays float64
}

// Around computes the sunrise and sunset times for latitude and longitude
// around currentTime. Generally, the computed sunrise will be no earlier
// than 24 hours before currentTime and the computed sunset will be no later
// than 24 hours after currentTime. However, these differences may exceed 24
// hours on days with more than 23 hours of daylight.
// The latitude is positive for north and negative for south. Longitude is
// positive for east and negative for west.
func (s *Sunrise) Around(latitude, longitude float64, currentTime time.Time) {
	s.location = currentTime.Location()
	s.sinLat = sin(latitude)
	s.cosLat = cos(latitude)
	s.jstar = math.Floor(
		julianDay(currentTime.Unix())-0.0009+longitude/360.0+0.5) + 0.0009 - longitude/360.0
	s.computeSolarNoonHourAngle()
}

// AddDays computes the sunrise and sunset numDays after
// (or before if numDays is negative) the current sunrise and sunset at the
// same latitude and longitude.
func (s *Sunrise) AddDays(numDays int) {
	s.jstar += float64(numDays)
	s.computeSolarNoonHourAngle()
}

// Sunrise returns the current computed sunrise. Returned sunrise has the same
// location as the time passed to Around.
func (s *Sunrise) Sunrise() time.Time {
	return goTime(s.solarNoon-s.hourAngleInDays, s.location)
}

// Sunset returns the current computed sunset. Returned sunset has the same
// location as the time passed to Around.
func (s *Sunrise) Sunset() time.Time {
	return goTime(s.solarNoon+s.hourAngleInDays, s.location)
}

func (s *Sunrise) computeSolarNoonHourAngle() {
	ma := mod360(357.5291 + 0.98560028*(s.jstar-jepoch))
	center := 1.9148*sin(ma) + 0.02*sin(2.0*ma) + 0.0003*sin(3.0*ma)
	el := mod360(ma + 102.9372 + center + 180.0)
	s.solarNoon = s.jstar + 0.0053*sin(ma) - 0.0069*sin(2.0*el)
	declination := asin(sin(el) * sin(23.45))
	s.hourAngleInDays = acos((sin(-0.83)-s.sinLat*sin(declination))/(s.cosLat*cos(declination))) / 360.0
}

func julianDay(unix int64) float64 {
	return float64(unix-uepoch)/86400.0 + jepoch
}

func goTime(julianDay float64, loc *time.Location) time.Time {
	unix := uepoch + int64((julianDay-jepoch)*86400.0)
	return time.Unix(unix, 0).In(loc)
}

func sin(degrees float64) float64 {
	return math.Sin(degrees * math.Pi / 180.0)
}

func cos(degrees float64) float64 {
	return math.Cos(degrees * math.Pi / 180.0)
}

func asin(x float64) float64 {
	return math.Asin(x) * 180.0 / math.Pi
}

func acos(x float64) float64 {
	if x >= 1.0 {
		return 0.0
	}
	if x <= -1.0 {
		return 180.0
	}
	return math.Acos(x) * 180.0 / math.Pi
}

func mod360(x float64) float64 {
	return x - 360.0*math.Floor(x/360.0)
}
