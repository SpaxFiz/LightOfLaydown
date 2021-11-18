package model

type CrawlData interface {
	Fetch() error
	Render() (interface{}, error)
}
