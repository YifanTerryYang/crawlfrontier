package crawlfrontier

import (
	"fmt"
	"regexp"
	"strings"
)

func FindAllURL(s, pattern string) []string {
	re := regexp.MustCompile(pattern)
	return re.FindAllString(s, -1)
}

func URLcleanup(url string) (result string) {
	result = strings.Replace(url, "href=", "", -1)
	result = strings.Replace(result, "src=", "", -1)
	result = strings.Replace(strings.Replace(result, `"`, "", -1), " ", "", -1)
	if result[0:1] == "/" {
		result = result[1:]
	}
	return
}

func URLcleanup_singlequote(url string) (result string) {
	result = strings.Replace(url, "href=", "", -1)
	result = strings.Replace(strings.Replace(result, `'`, "", -1), " ", "", -1)
	if result[0:1] == "/" {
		result = result[1:]
	}
	return
}

func (web *website) putUrl(searchcriteria, page string) {
	fmt.Println(web.seed, ":", searchcriteria, ":", web.urlstorage.Size())
	go web.parser(searchcriteria, page, web.urlchan)
}

func (web *website) loadQueue() { // run only once at initialization stage
	var s *URLEntity
	for {
		s = <-web.urlchan
		web.urlstorage.Add(s)
	}
}

func (web *website) readQueue() (url string) {
	entity := web.urlstorage.Get()
	if entity != nil {
		url = `<searchcriteria>` + entity.searchcriteria + `</searchcriteria>` + `<url>` + entity.url + `</url>`
	} else {
		url = ""
	}
	if web.urlstorage.Size() == 0 {
		fmt.Println(web.seed, " EMPTY now!")
	}
	return
}

/* When there is a new website, update parser function here */
func getParser(seedname string) (f func(searchcriteria, page string, outpipe chan<- *URLEntity)) {
	switch seedname {
	case "www.example1.com":
		f = func(searchcriteria, page string, outpipe chan<- *URLEntity) {
			fullseed := "http://" + seedname + "/"
			// get full urls
			pattern := `<\s*a[^>]*rel="noopener nofollow"[^>]*>`
			hrefpattern := `href="?(\/[%=()'+,\w\.\?\&-]{0,})+"?`
			//titlepattern := `title="[^"]*"`
			// get sub urls with href
			hrefurls := FindAllURL(page, pattern)

			for _, val := range hrefurls {
				urls := FindAllURL(val, hrefpattern)
				if len(urls) == 0 {
					continue
				}
				url := urls[0]
				url = URLcleanup(url)
				url = fullseed + url

				newentity := &URLEntity{
					url:            url,
					searchcriteria: searchcriteria,
				}
				outpipe <- newentity
			}
			return
		}
	case "www.example2.com":
		f = func(searchcriteria, page string, outpipe chan<- *URLEntity) {
			fullseed := "http://" + seedname + "/"

			// get full urls
			pattern := `<\s*a[^>]*rel='nofollow'[^>]*>`
			hrefpattern := `href='?(\/[%=()+,\w\.!\?\&-]{0,})+'?`
			//hrefpattern := `href="?(.*)"?`
			// get sub urls with href
			hrefurls := FindAllURL(page, pattern)

			for _, val := range hrefurls {
				urls := FindAllURL(val, hrefpattern)
				if len(urls) == 0 {
					continue
				}
				url := urls[0]
				url = URLcleanup_singlequote(url)
				url = fullseed + url

				newentity := &URLEntity{
					url:            url,
					searchcriteria: searchcriteria,
				}
				outpipe <- newentity
			}
			return
		}
	case "www.example3.com":
		f = func(searchcriteria, page string, outpipe chan<- *URLEntity) {
			fullseed := "http://" + seedname + "/"

			// get full urls
			pattern := `<\s*a[^>]*class="dice-btn-link loggedInVisited[^>]*>`
			hrefpattern := `href="?(\/[%=()'+,*\w\.!\?\&-]{0,})+"?`
			//hrefpattern := `href="?(.*)"?`
			// get sub urls with href
			hrefurls := FindAllURL(page, pattern)

			for _, val := range hrefurls {
				urls := FindAllURL(val, hrefpattern)
				if len(urls) == 0 {
					continue
				}
				url := urls[0]
				url = URLcleanup(url)
				url = fullseed + url

				newentity := &URLEntity{
					url:            url,
					searchcriteria: searchcriteria,
				}
				outpipe <- newentity
			}
			return
		}
	case "www.example4.com":
		f = func(searchcriteria, page string, outpipe chan<- *URLEntity) {
			fullseed := "http://" + seedname + "/"

			// get full urls
			pattern := `<\s*a[^>]*data-gtm="jrp-job-list\|job-title-click\|[^>]*>`
			hrefpattern := `href="?(\/[%=()'+,*\w\.!\?\&-]{0,})+"?`
			//hrefpattern := `href="?(.*)"?`
			// get sub urls with href
			hrefurls := FindAllURL(page, pattern)

			for _, val := range hrefurls {
				urls := FindAllURL(val, hrefpattern)
				if len(urls) == 0 {
					continue
				}
				url := urls[0]
				url = URLcleanup(url)
				url = fullseed + url

				newentity := &URLEntity{
					url:            url,
					searchcriteria: searchcriteria,
				}
				outpipe <- newentity
			}
			return
		}
	case "www.example5.com":
		f = func(searchcriteria, page string, outpipe chan<- *URLEntity) {
			// get full urls
			pattern := `(?:H|h)(?:T|t)(?:T|t)(?:P|p)(?:S|s):\/\/(?:W|w)(?:W|w)(?:W|w)\.(?:L|l)(?:I|i)(?:N|n)(?:K|k)(?:E|e)(?:D|d)(?:I|i)(?:N|n)\.(?:C|c)(?:O|o)(?:M|m)\/(?:J|j)(?:O|o)(?:B|b)(?:S|s)\/(?:V|v)(?:I|i)(?:E|e)(?:W|w)(?:[^"]*)`
			// get sub urls with href
			hrefurls := FindAllURL(page, pattern)

			for _, url := range hrefurls {
				newentity := &URLEntity{
					url:            url,
					searchcriteria: searchcriteria,
				}
				outpipe <- newentity
			}
			return
		}
	default:
		f = func(searchcriteria, page string, outpipe chan<- *URLEntity) {
			fmt.Println("default function")
		}
	}
	return
}
