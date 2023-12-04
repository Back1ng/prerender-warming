package sitemapper

import "encoding/xml"

// Sitemap with minimal data
type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	URL     []struct {
		Loc string `xml:"loc"`
	} `xml:"url"`
}
