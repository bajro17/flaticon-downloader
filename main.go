package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"github.com/gocolly/colly/v2"
)

func main() {
	ime := os.Args[1]
	broj, _ := strconv.Atoi(os.Args[2])
	
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("www.flaticon.com"),
		colly.Async(true),
	)

	// On every a element which has href attribute call callback
	c.OnHTML(".icon--item", func(e *colly.HTMLElement) {
		link := e.Attr("data-icon_src")
		// Print link
		fmt.Printf("Link found:  %s\n", link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting1", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		path, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		filename := r.FileName() //Costum for youself
		// ext := filepath.Ext(r.Request.URL.String())
		// ceanExt := sanitize.BaseName(ext)
		if _, err := os.Stat(path + "/"+ ime ); os.IsNotExist(err) {
			err := os.Mkdir(path + "/"+ ime, 0755)
			if err != nil {
				fmt.Println(err)
			}
		}
		final_filename := path + "/"+ ime + fmt.Sprintf("/%s.%s", filename, "svg")

		err = ioutil.WriteFile(final_filename, r.Body, 0644)
		if err != nil {
			fmt.Println(err)
		}
	})

	for i := 1; i <= broj; i++ {
		fullURL := fmt.Sprintf("https://www.flaticon.com/search/%d?word=%s&type=icon&license=selection&order_by=4", i, ime)
		c.Visit(fullURL)
		c.Wait()
	  }

}
