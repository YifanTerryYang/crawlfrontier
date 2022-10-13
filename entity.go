package crawlfrontier

type frontier struct {
	websitesList []*website
	websitesIndex map[string]int   // website name to index of websitesList
	logger *Logger
}

type website struct {
	seed string
	urlchan chan *URLEntity
	urlstorage *Queue
	parser func(searchcriteria, page string, outpipe chan<- *URLEntity)
}