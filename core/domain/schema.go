package domain

type CrawlData interface {
	Fetch() error
	Render() (interface{}, error)
}
