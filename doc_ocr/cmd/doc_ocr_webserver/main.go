package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/doc_ocr/cmd"
	"github.com/doc_ocr/events"
	"github.com/gorilla/mux"
	"github.com/otiai10/gosseract"
)

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
	f, err := os.Create(header.Filename)
	d2 := Buf.Bytes()
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
	document := cmd.Document{}
	document.ID = document.GenerateNewUUID()
	document.Creation = time.Now()
	document.Extension = header.Filename
	document.FileName = header.Filename
	document.TEXT = buffer.String()

	fmt.Println(document)
	fmt.Println(responseString)
	return

}
func main() {
	events.Push()
	router := mux.NewRouter()
	router.HandleFunc("/RF", ReceiveFile).Methods("POST")
	//router.HandleFunc("/ReceiveFile", doOcr).Methods("POST")
	//router.HandleFunc("/OCR", doOcr).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
