package main

import (
	"encoding/json"
	// "os"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type animal struct {
	Name string `json:"title"`
	Url  string `json:"permalink"`
	Tags struct {
		Species string `json:"species"`
	} `json:"tags"`
}

type catDetail struct {
	name string
	url  string
	lbs  int
	oz   int
}

type expectedResponse struct {
	Items []animal `json:"items"`
}

const listUrl = "https://www.sfspca.org/wp-json/sfspca/v1/filtered-posts/get-adoptions?per_page=100"

var (
	reLbs *regexp.Regexp
	reOz  *regexp.Regexp
)

func init() {
	reLbs = regexp.MustCompile("([0-9]+)\\slb")
	reOz = regexp.MustCompile("([0-9]+)\\soz")
}

func main() {
	fattestCats, err := getFattestCats()
	if err != nil {
		log.Fatal("error cat", err)
	}
	if len(fattestCats) == 0 {
		log.Fatal("not cat", err)
	}
	log.Printf("fat cat: %+v", fattestCats)
}

func saveFattestCats([]catDetail) error {
	// TODO
	return nil
}

func getFattestCats() ([]catDetail, error) {
	catItems, err := fetchCatItems()
	if err != nil {
		return nil, fmt.Errorf("error fetching cats :( %v", err)
	}

	fattestCats := make([]catDetail, 0, len(catItems))
	var highestWeight int

	for _, c := range catItems {
		cat, err := fetchCatDetails(c)
		if err != nil {
			log.Println("error fetching cat detail", err)
		} else {
			fmt.Printf("Weighing cat: %s\n", cat.name)
			// fmt.Printf("Details: %+v\n", cat)
			if weight := cat.getWeight(); weight > highestWeight {
				highestWeight = weight
				fattestCats = []catDetail{cat}
			} else if weight == highestWeight {
				fattestCats = append(fattestCats, cat)
			}
		}
	}

	return fattestCats, nil
}

func interpretResults(fattestCats []catDetail) {
	// TODO
}

func (c catDetail) getWeight() int {
	return c.lbs*16 + c.oz
}

func fetchCatItems() ([]animal, error) {
	response, err := http.Get(listUrl)
	if err != nil {
		return nil, err
	}

	// TODO check for error response

	var items expectedResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&items)
	if err != nil {
		return nil, err
	}

	cats := make([]animal, 0, len(items.Items))
	for _, i := range items.Items {
		if i.Tags.Species == "Cat" {
			cats = append(cats, i)
		}
	}

	return cats, nil
}

func fetchCatDetails(c animal) (catDetail, error) {
	response, err := http.Get(c.Url)
	if err != nil {
		return catDetail{}, err
	}

	detailB, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return catDetail{}, err
	}
	detail := string(detailB)

	var (
		lbs int
		oz  int
	)

	lbMatch := reLbs.FindStringSubmatch(detail)
	ozMatch := reOz.FindStringSubmatch(detail)

	if len(lbMatch) > 1 {
		parsed, err := strconv.ParseInt(lbMatch[1], 10, 32)
		if err != nil {
			return catDetail{}, err
		}
		lbs = int(parsed)
	}

	if len(ozMatch) > 1 {
		parsed, err := strconv.ParseInt(ozMatch[1], 10, 32)
		if err != nil {
			return catDetail{}, err
		}
		oz = int(parsed)
	}

	return catDetail{name: c.Name, url: c.Url, lbs: lbs, oz: oz}, nil
}
