package db

type NovelRow struct {
	Title        string
	Content      string
	CrawlUrl     string
	NextCrawlUrl string
	OssId        int
}
type OssRow struct {
	Id       int
	Name     string
	CrawlUrl string
	Type     string
	Url      string
}

type PhotoRow struct {
	Title    string
	Url      string
	CrawlUrl string
	OssId    int
	Index    int
}
