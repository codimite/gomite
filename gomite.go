package gomite

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type Gomite struct {
	Port    string
	Handler http.Handler
}

type Opts struct {
	Addr           string
	Handler        http.Handler
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

type OptFunc func(*Opts)

func defaultOpts(server Gomite) Opts {
	return Opts{
		Addr:           ":" + server.Port,
		Handler:        server.Handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func WithReadTimeout(t time.Duration) OptFunc {
	return func(opts *Opts) {
		opts.ReadTimeout = t
	}
}

func WithWriteTimeout(t time.Duration) OptFunc {
	return func(opts *Opts) {
		opts.WriteTimeout = t
	}
}

type GomiteHandlerFunc http.HandlerFunc
type Handler http.Handler

var Templates *template.Template

var templateDirs = []string{"templates", "templates/partials"}

func getTemplates(tempDirs []string) (templates *template.Template, err error) {
	var allFiles []string
	for _, dir := range tempDirs {
		files2, _ := ioutil.ReadDir(dir)
		for _, file := range files2 {
			filename := file.Name()
			if strings.HasSuffix(filename, ".gohtml") {
				filePath := filepath.Join(dir, filename)
				allFiles = append(allFiles, filePath)
			}
		}
	}

	templates, err = template.New("").ParseFiles(allFiles...)
	return
}

func init() {
	Templates, _ = getTemplates(templateDirs)
}

func InitTemplates(tempDirs []string) {
	Templates, _ = getTemplates(tempDirs)
}

func (server Gomite) Start(opts ...OptFunc) {
	o := defaultOpts(server)
	for _, fn := range opts {
		fn(&o) // changing default options
	}
	s := &http.Server{
		Addr:           o.Addr,
		Handler:        o.Handler,
		ReadTimeout:    o.ReadTimeout,
		WriteTimeout:   o.WriteTimeout,
		MaxHeaderBytes: o.MaxHeaderBytes,
	}
	log.Fatal(s.ListenAndServe())
}

func (server Gomite) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, handler)
}

func (server Gomite) Handle(pattern string, handler http.Handler) {
	http.Handle(pattern, handler)
}
