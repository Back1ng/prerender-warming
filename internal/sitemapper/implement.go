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

func (s *sitemapper) Get(url string) Sitemap {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	var parsedXml Sitemap
	xml.Unmarshal(body, &parsedXml)

	return parsedXml
}
