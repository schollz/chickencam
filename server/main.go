package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type chicken_data struct {
	date time.Time
}

type timeSlice []chicken_data

func (p timeSlice) Len() int {
	return len(p)
}

func (p timeSlice) Less(i, j int) bool {
	return p[i].date.After(p[j].date)
}

func (p timeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func getData() ([]string, []string) {
	files, _ := ioutil.ReadDir("./static/data")
	wavs := make(map[string]bool)
	jpgs := make(map[string]bool)
	for _, f := range files {
		name := f.Name()
		if strings.Contains(f.Name(), ".wav") {
			wavs[name[0:len(name)-4]] = true
		}
		if strings.Contains(f.Name(), ".jpg") {
			jpgs[name[0:len(name)-4]] = true
		}
	}

	// Sort them by date, descending
	var chickenMap = make(map[string]chicken_data)
	for key := range wavs {
		if _, ok := jpgs[key]; ok {
			t, err := time.Parse("20060102150405", key)
			if err == nil {
				chickenMap[key] = chicken_data{date: t}
			}
		}
	}
	tosort := make(timeSlice, 0, len(chickenMap))
	chickenDateMap := make(map[string]string)
	for k, d := range chickenMap {
		tosort = append(tosort, d)
		chickenDateMap[d.date.String()] = k
	}
	sort.Sort(tosort)
	sortedDates := []string{}
	sortedNames := []string{}
	for _, d := range tosort {
		sortedDates = append(sortedDates, d.date.Format("01/02/2006 3:04 PM"))
		sortedNames = append(sortedNames, chickenDateMap[d.date.String()])
	}

	return sortedDates, sortedNames
}
func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/", func(c *gin.Context) {
		sortedDates, sortedNames := getData()
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
			"Dates": sortedDates,
			"Names": sortedNames,
		})
	})
	router.Run(":8081")
	fmt.Println(getData())
}
