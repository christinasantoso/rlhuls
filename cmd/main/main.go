package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/naoina/denco"
	"github.com/naoina/miyabi"
)

func main() {
	var pages []string
	fis, err := ioutil.ReadDir("./static")
	if err != nil {
		log.Fatal(err)
	}
	for _, fi := range fis {
		pages = append(pages, fi.Name())
	}

	mux := denco.NewMux()
	handler, err := mux.Build([]denco.Handler{
		mux.GET("/", func(w http.ResponseWriter, r *http.Request, params denco.Params) {
			tmpl, err := template.New("index").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>index</title>
	</head>
	<body>
		{{ range .Pages }}
			<a href="{{ . }}">{{ . }}</a></br>
		{{end}}
	</body>
</html>
			`)
			if err != nil {
				w.WriteHeader(http.StatusPaymentRequired)
				w.Write([]byte(err.Error()))
				return
			}
			err = tmpl.Execute(w, struct {
				Pages []string
			}{
				Pages: pages,
			})
			if err != nil {
				w.WriteHeader(http.StatusPaymentRequired)
				w.Write([]byte(err.Error()))
				return
			}
			return
		}),
		mux.GET("/:page", func(w http.ResponseWriter, r *http.Request, params denco.Params) {
			page := params.Get("page")
			file := fmt.Sprintf("./static/%s", page)
			f, err := os.Open(file)
			if err != nil {
				w.WriteHeader(http.StatusPaymentRequired)
				w.Write([]byte(err.Error()))
				return
			}
			defer f.Close()
			http.ServeContent(w, r, page, time.Now(), f)
			return
		}),
		mux.GET("/health", func(w http.ResponseWriter, r *http.Request, params denco.Params) {
			w.WriteHeader(http.StatusOK)
			return
		}),
	})
	if err != nil {
		panic(err)
	}
	log.Fatal(miyabi.ListenAndServe(":6969", handler))
}
