package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
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

func getData(filterDate string) WebData {
	files, _ := ioutil.ReadDir("./static/data")
	wavs := make(map[string]bool)
	jpgs := make(map[string]bool)
	for _, f := range files {
		name := f.Name()
		if strings.Contains(f.Name(), ".mp3") {
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
	notes := make([]string, len(tosort))
	chicken := make([]string, len(tosort))
	egg := make([]string, len(tosort))
	availableDates := []string{}
	parseableDates := []string{}
	foundDate := make(map[string]bool)
	pictureCounts := make(map[string]int)
	i := 0
	for _, d := range tosort {
		if len(filterDate) > 0 {
			if filterDate != d.date.Format("01/02/2006") {
				continue
			}
		}
		sortedDates[i] = d.date.Format("3:04 PM")
		sortedNames[i] = chickenDateMap[d.date.String()]
		sortedHashes[i] = d.date.Format("20060102150405")
		if _, ok := foundDate[d.date.Format("01/02/2006")]; !ok {
			availableDates = append(availableDates, d.date.Format("January 02, 2006"))
			parseableDates = append(parseableDates, d.date.Format("01/02/2006"))
			foundDate[d.date.Format("01/02/2006")] = true
			pictureCounts[d.date.Format("01/02/2006")] = 0
		}
		pictureCounts[d.date.Format("01/02/2006")]++
		if _, err := os.Stat(path.Join("static", "data", sortedHashes[i]+".txt")); err == nil {
			b, _ := ioutil.ReadFile(path.Join("static", "data", sortedHashes[i]+".txt"))
			var chickenDat ChickenData
			json.Unmarshal(b, &chickenDat)
			notes[i] = chickenDat.Notes
			if chickenDat.Egg {
				egg[i] = "checked"
			}
			if chickenDat.Chicken {
				chicken[i] = "checked"
			}
		}
		i++
	}
	sortedDates = sortedDates[0:i]
	sortedNames = sortedNames[0:i]
	sortedHashes = sortedHashes[0:i]
	return WebData{SortedDates: sortedDates,
		SortedNames:    sortedNames,
		SortedHashes:   sortedHashes,
		AvailableDates: availableDates,
		ParseableDates: parseableDates,
		PictureCounts:  pictureCounts,
		RandomNumber:   rand.New(rand.NewSource(99)).Int31(),
		Notes:          notes,
		Chicken:        chicken,
		Egg:            egg,
	}
}

type WebData struct {
	Title          string
	SortedDates    []string
	SortedNames    []string
	SortedHashes   []string
	AvailableDates []string
	ParseableDates []string
	Notes          []string
	Chicken        []string
	Egg            []string
	PictureCounts  map[string]int
	Info           map[string]ChickenData
	RandomNumber   int32
}

type ChickenData struct {
	Notes   string
	Egg     bool
	Chicken bool
}

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/", func(c *gin.Context) {
		data := getData("")
		data.SortedDates = []string{}
		c.HTML(http.StatusOK, "base.tmpl", gin.H{
			"title": "Main website",
			"Data":  data,
		})
	})
	router.GET("/date/*date", func(c *gin.Context) {
		filterDate := c.Param("date")[1:]
		data := getData(filterDate)
		data.ParseableDates = []string{}
		c.HTML(http.StatusOK, "base.tmpl", gin.H{
			"title": "Main website",
			"Data":  data,
		})
	})
	router.POST("/update", func(c *gin.Context) {
		notes := c.PostForm("notes")
		egg := strings.Contains(c.DefaultPostForm("egg", ""), "on")
		chicken := strings.Contains(c.DefaultPostForm("chicken", ""), "on")
		id := c.PostForm("id")
		chickenData := ChickenData{Notes: notes, Egg: egg, Chicken: chicken}
		b, _ := json.Marshal(chickenData)
		ioutil.WriteFile(path.Join("static", "data", id+".txt"), b, 0644)
		log.Printf("Wrote JSON data to %s\n", path.Join("static", "data", id+".txt"))
		c.JSON(200, gin.H{
			"status":  "posted",
			"success": true,
			"egg":     egg,
			"chicken": chicken,
			"id":      id,
			"notes":   notes,
		})
	})
	router.Run(":8081")
}

// GetMD5Hash from http://stackoverflow.com/questions/2377881/how-to-get-a-md5-hash-from-a-string-in-golang
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))[:7]
}
