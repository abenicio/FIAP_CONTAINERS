package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/doc_ocr/cmd"
	"github.com/doc_ocr/events"
	"github.com/gorilla/mux"
	"github.com/otiai10/gosseract"
)

var stream = events.CreateStream()

func GetFile(w http.ResponseWriter, r *http.Request) {

	return
}
func Health(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode("Still alive!")
}
func main() {

	router := mux.NewRouter()
	router.HandleFunc("/GETFILE", GetFile).Methods("GET")
	router.HandleFunc("/Health", Health).Methods("GET")
	router.HandleFunc("/RF", ReceiveFile).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func doOcr(data []byte) string {

	ocrCliente := gosseract.NewClient()
	defer ocrCliente.Close()
	ocrCliente.SetImageFromBytes(data)
	text, _ := ocrCliente.Text()
	return text
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
func ReceiveFile(w http.ResponseWriter, r *http.Request) {
	document := cmd.Document{}
	document.ID = document.GenerateNewUUID()
	var Buf bytes.Buffer
	file, header, err := r.FormFile("file")
	if err != nil {
		file, err := os.Create("result.txt")
		file.WriteString("ERRO")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		file.Close()
		panic(err)
	}
	defer file.Close()
	//name := strings.Split(header.Filename, ".")
	io.Copy(&Buf, file)
	extension := path.Ext(header.Filename)
	f, err := os.Create(document.ID + extension)

	d2 := Buf.Bytes()

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(d2)
	_, err = f.Write(d2)
	check(err)
	Buf.Reset()

	text := strings.TrimSpace(doOcr(d2))
	var buffer bytes.Buffer

	text = strings.TrimSuffix(text, "\n")
	for _, r := range text {
		if r != 10 {
			buffer.WriteString(string(r))
		}
	}

	c := []byte(fmt.Sprintf(`{"text":"%s"}`, buffer.String()))

	req, err := http.NewRequest("POST", "http://localhost:5000", bytes.NewBuffer(c))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseString := string(responseData)

	document.Creation = time.Now()
	document.Extension = header.Filename
	document.FileName = header.Filename
	document.TEXT = buffer.String()
	document.IMAGE = encoded
	document.CLASSIFICATION = responseString

	b, err := json.Marshal(document)
	if err != nil {
		fmt.Println(err)
		return
	}

	events.PutStream(*stream, string(b))

	fmt.Println(document.IMAGE)
	fmt.Println(responseString)
	return

}
