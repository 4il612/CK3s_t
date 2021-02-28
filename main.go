package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	http.ListenAndServe(getAddresFromFlags(), mux)
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		tr, err := template.ParseFiles("static/html/index.html")
		if err != nil {
			http.Error(w, "На сервере возникла проблема!", 500)
			return
		}
		err = tr.Execute(w, nil)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "На сервере возникла проблема!", 500)
			return
		}
		return
	}

	src, hdr, err := r.FormFile("main")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer src.Close()

	path, err := getPathFromFile()
	if err != nil {
		http.Error(w, "Ошибка пути сохранения", 500)
		return
	}

	dst, err := os.Create(filepath.Join(path, hdr.Filename))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer dst.Close()
	io.Copy(dst, src)
}

func getAddresFromFlags() string {
	hostPtr := flag.String("h", "localhost:", "host")
	portPtr := flag.String("p", "3000", "port")
	flag.Parse()
	addr := *hostPtr + *portPtr
	return addr
}

func getPathFromFile() (string, error) {
	file, err := os.Open("path.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var path string
	for scanner.Scan() {
		path = scanner.Text()
		break
	}
	return path, nil
}
