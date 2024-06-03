package sitemapper

import "encoding/xml"

type Sitemapindex struct {
	XMLName xml.Name `xml:"sitemapindex"`
	Sitemap []struct {
		Text string `xml:"chardata"`
		Loc  string `xml:"loc"`
	} `xml:"sitemap"`
}

// Sitemap with minimal data
type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	URL     []struct {
		Loc string `xml:"loc"`
	} `xml:"url"`
}
