package controllers

import (
	"fmt"
	"html/template"
	"net/http"
)

func getFm() (fm template.FuncMap) {
	fm = template.FuncMap{
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"ro": ro,
	}

	return
}

func getSiteTemplates(filenames []string, fm template.FuncMap) (tmpl *template.Template) {
	var files []string

	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/site/%s.html", file))
	}

	if fm == nil {
		tmpl = template.Must(template.New("").ParseFiles(files...))
		return
	}

	tmpl = template.Must(template.New("").Funcs(fm).ParseFiles(files...))

	return
}

func GenerateHTMLEmp(writer http.ResponseWriter, data interface{}, filenames []string) {
	funcMap := getFm()

	filenames = append(filenames, "layout_emp")

	getSiteTemplates(filenames, funcMap).ExecuteTemplate(writer, "layout_emp", data)
}
