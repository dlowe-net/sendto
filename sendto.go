package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger"
)

const maxFileSize = 10e9

var (
	outputDirFlag = flag.String("dir", "", "Path to output directory")
	keyPathFlag   = flag.String("key", "", "Path to SSL Private Key")
	certPathFlag  = flag.String("cert", "", "Path to SSL certificate")
	portFlag      = flag.String("port", "80", "Port to use for HTTP")
	sPortFlag     = flag.String("sport", "443", "Port to use for HTTPS")
	hostnameFlag  = flag.String("hostname", "", "Hostname to use by default for SSL redirection")
	webDirFlag    = flag.String("web", "", "(debugging) path to web directory")
)

func ReportError(w http.ResponseWriter, err error) {
	log.Print(err)
	w.WriteHeader(http.StatusInternalServerError)
}

func SaveToFile(fname string, w http.ResponseWriter, r io.Reader) {
	// Make the filename safe
	fname = strings.ReplaceAll(fname, string(os.PathSeparator), "_")
	fname = strings.ReplaceAll(fname, "\0000", "?")
	file, err := os.Create(filepath.Join(*outputDirFlag, fname))
	if err != nil {
		ReportError(w, err)
		return
	}
	defer file.Close()

	_, err = io.CopyN(file, r, maxFileSize)
	if err != nil && err != io.EOF {
		ReportError(w, err)
		return
	}
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X_FILENAME") == "" {
		// Multipart
		reader, err := r.MultipartReader()
		if err != nil {
			ReportError(w, err)
			return
		}
		for {
			part, err := reader.NextPart()
			if err != nil {
				ReportError(w, err)
				return
			}
			if part.FormName() != "file" {
				ReportError(w, fmt.Errorf("validation error"))
				return
			}
			log.Printf("Receiving '%s' via mime", part.FileName())
			SaveToFile(part.FileName(), w, part)
		}
	} else {
		// Straight AJAX send
		fname := r.Header.Get("X_FILENAME")
		log.Printf("Receiving '%s' via ajax", fname)
		SaveToFile(fname, w, r.Body)
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	flag.Parse()
	if *outputDirFlag == "" {
		fmt.Fprintf(os.Stderr, "You must specify the output directory with -dir.\n")
		os.Exit(1)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var (
			inf io.ReadCloser
			err error
		)
		if *webDirFlag == "" {
			inf, err = pkger.Open("/web/index.html")
			if err != nil {
				ReportError(w, err)
			}
		} else {
			inf, err = os.Open(filepath.Join(*webDirFlag, "index.html"))
			if err != nil {
				ReportError(w, err)
			}
		}
		defer inf.Close()
		w.Header().Set("Content-Type", "text/html")
		_, err = io.Copy(w, inf)
		if err != nil {
			ReportError(w, err)
		}
	})
	http.HandleFunc("/upload", HandleUpload)

	if *keyPathFlag != "" && *certPathFlag != "" {
		go func() {
			log.Fatal(http.ListenAndServeTLS(":"+*sPortFlag, *certPathFlag, *keyPathFlag, nil))
		}()
		log.Fatal(http.ListenAndServe(":"+*portFlag, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			targetURL := *r.URL
			targetURL.Scheme = "https"
			if targetURL.Host == "" {
				targetURL.Host = *hostnameFlag
			}
			http.Redirect(w, r, targetURL.String(), http.StatusMovedPermanently)
		})))
	}

	log.Fatal(http.ListenAndServe(":"+*portFlag, nil))
}
