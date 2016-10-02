package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
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

func getData() ([]string, []string, []string) {
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
	sortedDates := make([]string, len(tosort))
	sortedNames := make([]string, len(tosort))
	sortedHashes := make([]string, len(tosort))
	for i, d := range tosort {
		sortedDates[i] = d.date.Format("01/02/2006 3:04 PM")
		sortedNames[i] = chickenDateMap[d.date.String()]
		sortedHashes[i] = GetMD5Hash(sortedDates[i])
	}

	return sortedDates, sortedNames, sortedHashes
}

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/", func(c *gin.Context) {
		sortedDates, sortedNames, sortedHashes := getData()
		r := rand.New(rand.NewSource(99))
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":        "Main website",
			"Dates":        sortedDates,
			"Names":        sortedNames,
			"Hashes":       sortedHashes,
			"RandomNumber": r.Int31(),
		})
	})
	router.Run(":8081")
	fmt.Println(getData())
}

// GetMD5Hash from http://stackoverflow.com/questions/2377881/how-to-get-a-md5-hash-from-a-string-in-golang
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))[:7]
}
