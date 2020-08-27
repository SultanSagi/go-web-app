package views

import (
	"bytes"
	"fmt"
	"io"
	"html/template"
	"net/http"
	"path/filepath"
)

var LayoutDir string = "views/layouts"

func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

type Data struct {
	Yield interface{}
	Errors string
	SearchFirstName string
	SearchLastName string
	FilterOrderBy string
	FilterSort string
	PaginationPrevPage string
	PaginationNextPage string
	PaginationPage string
}

func (v *View) Render(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	switch data. (type) {
		case Data:
			//
		default:
			data = Data{
				Yield: data,
			}
	}
	var buf bytes.Buffer
	err := v.Template.ExecuteTemplate(&buf, v.Layout, data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "/*.gohtml")
	if err != nil {
		panic(err)
	}
	return files
}
