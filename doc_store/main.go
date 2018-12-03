package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/doc_store/cmd"
	"github.com/gorilla/mux"
)

func ReceiveFile(w http.ResponseWriter, req *http.Request) {

	base64 := req.FormValue("iamge")
	fileName := req.FormValue("Extension")
	class := req.FormValue("classification")
	CreateDirIfNotExist("neg/")
	CreateDirIfNotExist("pos/")
	id := req.FormValue("id")
	ext := filepath.Ext(fileName)
	sDec, err := b64.StdEncoding.DecodeString(base64)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(class + "/" + id + ext)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(sDec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}
	return

}
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
func GetAllFiles(w http.ResponseWriter, req *http.Request) {
	var fls []cmd.Files
	log.Printf("CREATE DIRECTORY")
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(dir)
	CreateDirIfNotExist("neg/")
	CreateDirIfNotExist("pos/")
	files, err := ioutil.ReadDir("pos/")
	if err != nil {
		log.Fatal(err)
	}
	filesNegative, err := ioutil.ReadDir("neg/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fl := cmd.Files{CLASSIFICATION: "POSITIVE", FileName: f.Name()}
		fls = append(fls, fl)
	}
	for _, f := range filesNegative {
		fl := cmd.Files{CLASSIFICATION: "NEGATIVE", FileName: f.Name()}
		fls = append(fls, fl)
	}

	data, _ := json.Marshal(fls)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return

}
func Health(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode("Still alive!")
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/SF", ReceiveFile).Methods("POST")
	router.HandleFunc("/FILES", GetAllFiles).Methods("GET")
	router.HandleFunc("/HEALTH", Health).Methods("GET")
	log.Fatal(http.ListenAndServe(":8010", router))
}
