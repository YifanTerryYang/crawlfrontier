package crawlfrontier

import (
	"net/http"
	"strconv"
	"errors"
	"log"
)

func New() frontier {
	return frontier{
		websitesList: make([]*website,0), 
		websitesIndex: make(map[string]int),
		logger: NewLogger(),
	}
}

func (fr *frontier) Start(port int) error {
	if port > 1024 && port <= 65535 {
		http.HandleFunc("/api/setseed", fr.setSeed)
		http.HandleFunc("/api/puturl", fr.putUrl)
		http.HandleFunc("/api/geturl", fr.getUrl)
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
		return nil
	} else {
		return errors.New("port out of range: " + strconv.Itoa(port))
	}	
}

func (fr *frontier) getClient(r *http.Request) (*website,error) {
	index, ok := r.URL.Query()["index"]
	if !ok {
		seed, ok := r.URL.Query()["seed"]  // get seed name
		if !ok {
			return nil, errors.New("need containing 'index' or 'seed' in request")
		}
		i, ok := fr.websitesIndex[seed[0]]    // get the index
		if !ok {
			return nil, errors.New("website not set up!")
		}
		if i >= len(fr.websitesList) {   // if index is less than length of websiteslist
			return nil, errors.New("website not set up!")
		}

		return fr.websitesList[i], nil
	}
	
	i, err := strconv.Atoi(index[0])
	if err != nil {  // if conversion failed
		return nil, errors.New("err to convert index")
	}
	if i >= len(fr.websitesList) {   // if index is less than length of websiteslist
		return nil, errors.New("website not set up!")
	}

	return fr.websitesList[i], nil
}