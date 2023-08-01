package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
)

type item struct {
	Url    string
	Name   string
	Price  string
	ImgUrl string
}

func processHTMLElement(h *colly.HTMLElement) item {
	return item{
		Url:    h.ChildAttr("a", "href"),
		Name:   h.ChildText("h2"),
		Price:  h.ChildText(".price"),
		ImgUrl: h.ChildAttr("img", "src"),
	}
}

func main() {
	//Create a new collector
	c := colly.NewCollector(
		colly.AllowedDomains("scrapeme.live"),
	)

	//get the items and put it into a slice
	var items []item
	c.OnHTML("li.product", func(h *colly.HTMLElement) {
		item := processHTMLElement(h)
		items = append(items, item)
	})

	//go to the next page
	c.OnHTML("a.next.page-numbers", func(h *colly.HTMLElement) {
		next_page := h.Request.AbsoluteURL(h.Attr("href"))
		err := c.Visit(next_page)
		if err != nil {
			log.Println("error on visit next page")
			return
		}
	})

	//print the url of request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})

	//visit the site
	err := c.Visit("https://scrapeme.live/shop")
	if err != nil {
		log.Println("error on visit page")
		return
	}
	fmt.Println(items)

	//create and close the csv file
	file, err := os.Create("products.csv")
	if err != nil {
		log.Fatalln("failed to create csv file", err)
	}
	defer file.Close()

	//create a new writer
	writer := csv.NewWriter(file)

	//create the headers
	headers := []string{
		"url",
		"name",
		"price",
		"image",
	}

	//write the headers into csv file
	err = writer.Write(headers)
	if err != nil {
		log.Println("error on write headers")
		return
	}

	//a loop to write all the products
	for _, product := range items {
		record := []string{
			product.Url,
			product.Name,
			product.Price,
			product.ImgUrl,
		}
		err = writer.Write(record)
		if err != nil {
			log.Println("error on write record")
			return
		}
	}
	defer writer.Flush()

}
