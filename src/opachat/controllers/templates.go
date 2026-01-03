package controllers

import (
	"fmt"
	"html/template"
	"net/http"

	"opachat/tools"
)

func getFm() (fm template.FuncMap) {
	fm = template.FuncMap{
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"dd": func(v any) template.HTML {
			return template.HTML(
				tools.ShowJSON(v, false),
			)
		},
	}

	return
}

func getSiteTemplates(filenames []string, fm template.FuncMap) (tmpl *template.Template) {
	var files []string

	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	if fm == nil {
		tmpl = template.Must(template.New("").ParseFiles(files...))
		return
	}

	tmpl = template.Must(template.New("").Funcs(fm).ParseFiles(files...))
	return
}

func GenerateHTMLEmp(w http.ResponseWriter, r *http.Request, data any, filenames ...string) {
	funcMap := getFm()
	filenames = append(filenames, "lays/layout")

	getSiteTemplates(filenames, funcMap).ExecuteTemplate(w, "layout", data)
}
