package crawlfrontier

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
)

func (fr *frontier) setSeed(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == "GET" {
		// dealing with request
		requestURLparamlen := len(r.URL.Query())
		if requestURLparamlen != 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400-Wrong parameter count:" + strconv.Itoa(requestURLparamlen)))
			return
		}
		values, ok := r.URL.Query()["seed"]
		
		if !ok { // "seed" param not exists
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400-no seed received"))
			return
		}
		newSeed := values[0]
		newSeed = strings.ToLower(newSeed)
		if i, ok := fr.websitesIndex[newSeed]; ok {
			// if website exists already
			fmt.Fprint(w, strconv.Itoa(i))
			return
		}
		
		newwebsite := &website{
			seed:newSeed,
			urlchan:make(chan *URLEntity, 5000),
			urlstorage:&Queue{length:0},
			parser:getParser(newSeed),
		}
		go newwebsite.loadQueue()   // concurrent goroutine to load Queue from channel
		// Load websites list
		fr.websitesList = append(fr.websitesList, newwebsite)
		fr.websitesIndex[newSeed] = len(fr.websitesList) - 1

		//fmt.Println(len(fr.websitesList))
		//fmt.Println(len(fr.websitesIndex))
		fmt.Fprint(w, strconv.Itoa(fr.websitesIndex[newSeed]))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong method: " + method))
		return
	}
}

func (fr *frontier) putUrl(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == "POST" {
		client,err := fr.getClient(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))  // return error information if got error
			return
		}
		
		values, ok := r.URL.Query()["sc"]  // searchcriteria
		var searchcriteria string
		if !ok { // "seed" param not exists
			searchcriteria = ""
		} else {
			searchcriteria = values[0]
		}
		searchcriteria = strings.ToLower(searchcriteria)
		
		// get page from request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))  // return error information if got error
			return
		}
		page := string(body)
		go fr.logger.writeLog(client.seed, page)
		client.putUrl(searchcriteria, page)
		fmt.Fprint(w, "success")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong method: " + method))
		return
	}
}

func (fr *frontier) getUrl(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == "GET" {
		client,err := fr.getClient(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))  // return error information if got error
			return
		}
		url := client.readQueue()
		
		fmt.Fprint(w, url)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong method: " + method))
		return
	}
}