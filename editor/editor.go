package editor

import (
	"html/template"
	"net/http"
	"github.com/gorilla/csrf"
)



var templates *template.Template

func init() {

	templates = template.Must(template.ParseFiles("templates/editor.html","templates/index.html"))

}

func EditorSpace(resp http.ResponseWriter,req *http.Request) {

	templates.ExecuteTemplate(resp,"editorview",map[string]interface{}{
        csrf.TemplateTag: csrf.TemplateField(req),
    })


}
