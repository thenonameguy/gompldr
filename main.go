package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const ompldr string = "http://ompldr.org/"

func main() {
	argsCount := len(os.Args)
	if argsCount < 2 {
		log.Fatalln("Error: not enough arguments")
	}
	var clipboard string
	for i := 1; i < argsCount; i++ {
		url := Upload(os.Args[i])
		if i == argsCount-1 {
			clipboard += url
		} else {
			clipboard += url + "\n"
		}
		fmt.Printf("Uploaded %s to %s\n", os.Args[i], url)
	}
	write2clipboard(clipboard)
}

func Upload(file string) string {
	resp, err := postFile(file, ompldr+"upload")
	if err != nil {
		log.Fatalln(err)
	}
	content, _ := ioutil.ReadAll(resp.Body)
	regex := regexp.MustCompile("href=\"v(\\w+)\"")
	url := regex.FindString(string(content))
	if len(url) < 7 {
		log.Fatalln("Error: Ompldr limited your uploads for a while")
	}
	return ompldr + url[6:len(url)-1]
}

func write2clipboard(str string) {
	xclip := exec.Command("xclip", "-selection", "clipboard")
	xclip.Stdin = strings.NewReader(str)
	xclip.Start()
	xclip.Wait()
	fmt.Printf("\"%s\" copied to clipboard\n", str)
}

func postFile(filename string, target_url string) (*http.Response, error) {
	body_buf := new(bytes.Buffer)
	body_writer := multipart.NewWriter(body_buf)
	fb, err := body_writer.CreateFormFile("file1", filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer file.Close()
	io.Copy(fb, file)

	body_writer.Close()
	req, err := http.NewRequest("POST", target_url, body_buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", body_writer.FormDataContentType())

	return http.DefaultClient.Do(req)
}
