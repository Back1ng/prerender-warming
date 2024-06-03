package sitemapper

type SitemapParser interface {
	Get(url string) []Sitemap
}
