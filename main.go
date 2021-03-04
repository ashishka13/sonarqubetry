package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

//go:generate swagger generate spec -m -o ./swagger.json

func main() {
	myController()
}

func myController() {
	r := mux.NewRouter()

	r.HandleFunc("/", sonarhome)
	r.HandleFunc("/getsonarqubelist/{sonarParameter}/{projectName}", requestsonarqube).Methods("GET")

	log.Println("running on 12345")

	http.ListenAndServe(":12345", r)
}

func sonarhome(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "sonarqube properties")
	fmt.Fprintf(w, "choose from "+
		"\n code_smells"+
		"\n vulnerabilities"+
		"\n new_vulnerabilities"+
		"\n security_hotspots"+
		"\n new_security_hotspots"+
		"\n security_review_rating"+
		"\n new_security_review_rating"+
		"\n security_hotspots_reviewed\n")
	fmt.Fprintf(w, "\ngo to http://localhost:12345/getsonarqubelist/your_choice_here/project_name or ")
}

func requestsonarqube(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sonarParameter := vars["sonarParameter"]
	projectName := vars["projectName"]

	_, bodybyte := sonarqubeapi(sonarParameter, projectName)

	w.Write(bodybyte)
}

func sonarqubeapi(sonarParameter, projectName string) (bodystring string, bodybyte []byte) {
	var username string = "admin"
	var passwd string = "adminadmin"
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:9000/api/measures/search?metricKeys="+sonarParameter+"&projectKeys="+projectName, nil)
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	bodybyte, err = ioutil.ReadAll(resp.Body)
	log.Printf("error occurred %v", err)

	bodystring = string(bodybyte)
	spew.Dump("response is ", bodystring)
	return
}
