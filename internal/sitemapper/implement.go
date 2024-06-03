package sitemapper

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
)

type sitemapper struct {
}

func New() SitemapParser {
	return &sitemapper{}
}

func (s *sitemapper) Get(url string) []Sitemap {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	mult := parseMultipleSitemaps(body)
	single := parseUrlLoc(body)

	if len(mult) > 0 {
		return mult
	}

	if len(single.URL) > 0 {
		return []Sitemap{single}
	}

	return []Sitemap{}
}

// todo parse xml.gz, for example: https://ekb.cian.ru/sitemap.xml

func parseMultipleSitemaps(body []byte) []Sitemap {
	var parsedXml Sitemapindex
	xml.Unmarshal(body, &parsedXml)

	sitemaps := make([]Sitemap, len(parsedXml.Sitemap))

	for _, sm := range parsedXml.Sitemap {
		resp, err := http.Get(sm.Loc)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		sitemaps = append(sitemaps, parseUrlLoc(body))
	}

	return sitemaps
}

func parseUrlLoc(body []byte) Sitemap {
	var parsedXml Sitemap
	xml.Unmarshal(body, &parsedXml)

	return parsedXml
}
