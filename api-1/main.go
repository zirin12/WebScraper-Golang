package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

// ProductURL : amazon product's webpage url
type ProductURL struct {
	URL string `json:"url"`
}

// ResponseStruct has the final response fields
type ResponseStruct struct {
	URL     string      `json:"url"`
	Product ProductInfo `json:"product"`
}

// ProductInfo has product related info
type ProductInfo struct {
	Title            string `json:"title"`
	ImageURL         string `json:"imageURL"`
	ShortDescription string `json:"description"`
	Rating           string `json:"rating"`
	Price            string `json:"price"`
	TotalReviews     string `json:"totalReviews"`
}

// Scrapes the actual data from the amazon web page
func scrapeData(document *goquery.Document) ProductInfo {
	var result ProductInfo

	var productTitle string
	var imageURL string
	var shortDescription string
	var rating string
	var price string
	var totalReviews string

	productTitle = strings.TrimSpace(document.Find("#productTitle").First().Text())

	document.Find("#landingImage").Each(func(index int, element *goquery.Selection) {
		imgSrc, _ := element.Attr("data-a-dynamic-image")
		i := strings.Index(imgSrc, "[")
		imageURL = imgSrc[2 : i-2]
	})

	document.Find("#featurebullets_feature_div li").EachWithBreak(func(index int, s *goquery.Selection) bool {
		if !s.HasClass("aok-hidden") {
			shortDescription = strings.TrimSpace(s.Text())
			return false
		}
		return true
	})

	document.Find("#centerCol span.a-color-price").EachWithBreak(func(index int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(), "$") {
			price = strings.TrimSpace(s.Text())
			return false
		}
		return true
	})

	totalReviews = document.Find("#acrCustomerReviewText").First().Text()
	rating = document.Find("i.a-star-4-5 span.a-icon-alt, i.a-star-5 span.a-icon-alt").First().Text()

	result.Title = productTitle
	result.ImageURL = imageURL
	result.ShortDescription = shortDescription
	result.Price = price
	result.TotalReviews = totalReviews
	result.Rating = rating

	return result
}

func productScraper(url string) ResponseStruct {

	var productData ProductInfo
	var result ResponseStruct

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux i686; rv:76.0) Gecko/20100101 Firefox/76.0")

	response, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	productData = scrapeData(document)
	result.URL = url
	result.Product = productData

	return result
}

// ScrapeProductHandler gets the url from body and calls another api
func ScrapeProductHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var url ProductURL
	var product ResponseStruct
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, "Invalid Request body", http.StatusInternalServerError)
		return
	}
	product = productScraper(url.URL)
	productJSON, err := json.Marshal(product)
	if err != nil {
		log.Fatalln(err)
	}
	// call the other api
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://api-2:8080/products", bytes.NewBuffer(productJSON))
	if err != nil {
		log.Fatalln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)

	}
	log.Println(response)
	defer response.Body.Close()

	// return the scraped product info as the response
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(productJSON)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/product", ScrapeProductHandler).Methods("POST")
	if err := http.ListenAndServe(":9080", r); err != nil {
		log.Fatal(err)
	}
}
