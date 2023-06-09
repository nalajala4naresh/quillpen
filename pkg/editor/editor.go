package editor

import (
	"github.com/gorilla/csrf"
	
	"net/http"
)



func EditorSpace(resp http.ResponseWriter, req *http.Request) {

	templates.ExecuteTemplate(resp, "editorview", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(req),
	})

}
