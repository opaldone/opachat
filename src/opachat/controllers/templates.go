package controllers

import (
	"fmt"
	"net/http"
	"text/template"
)

func getSiteTemplates(filenames []string) (tmpl *template.Template) {
	var files []string

	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	tmpl = template.Must(template.New("").ParseFiles(files...))
	return
}

func GenerateHTMLEmp(w http.ResponseWriter, r *http.Request, data interface{}, filenames ...string) {
	filenames = append(filenames, "lays/layout")

	getSiteTemplates(filenames).ExecuteTemplate(w, "layout", data)
}
