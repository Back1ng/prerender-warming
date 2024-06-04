package sitemapper

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"errors"
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

func parseMultipleSitemaps(body []byte) []Sitemap {
	var parsedXml Sitemapindex
	xml.Unmarshal(body, &parsedXml)

	sitemaps := make([]Sitemap, len(parsedXml.Sitemap))

	for _, sm := range parsedXml.Sitemap {
		if sm.Loc[len(sm.Loc)-3:] == ".gz" {
			b := new(bytes.Buffer)
			resp, _ := http.Get(sm.Loc)
			defer resp.Body.Close()
			io.Copy(b, resp.Body)

			reader := bytes.NewReader(b.Bytes())
			gzreader, err := gzip.NewReader(reader)
			if err != nil {
				if errors.Is(err, gzip.ErrHeader) {
					sitemaps = append(sitemaps, parseUrlLoc(b.Bytes()))
					continue
				}
			}

			body, _ := io.ReadAll(gzreader)
			sitemaps = append(sitemaps, parseUrlLoc(body))
		} else {
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
	}

	return sitemaps
}

func parseUrlLoc(body []byte) Sitemap {
	var parsedXml Sitemap
	xml.Unmarshal(body, &parsedXml)

	return parsedXml
}
