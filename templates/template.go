package templates

import (
	"bytes"
	"embed"
	"html/template"

	"github.com/gin-gonic/gin"
)

// TemplateName is an enumeration of template name that are allowed.
// ENUM(
// success.tmpl
// verify.tmpl
// )
type TemplateName string

// String implements the Stringer interface.
func (x TemplateName) String() string {
	return string(x)
}

var _TemplateNameNames = []string{
	string(TemplateNameVerifyTmpl),
	string(TemplateNameSuccessTmpl),
}

// TemplateNameNames returns a list of possible string values of TemplateName.
func TemplateNameNames() []string {
	tmp := make([]string, len(_TemplateNameNames))
	copy(tmp, _TemplateNameNames)
	return tmp
}

const (
	TemplateNameVerifyTmpl  TemplateName = "verify.tmpl"
	TemplateNameSuccessTmpl TemplateName = "success.tmpl"
)

//go:embed *
var f embed.FS

var templates = template.Must(template.New("").ParseFS(f, TemplateNameNames()...))

// SetHTMLTemplate set templates into gin engine.
func SetHTMLTemplate(r *gin.Engine) {
	r.SetHTMLTemplate(templates)
}

// GenerateHTML returns html with filler.
func GenerateHTML(n TemplateName, filler any) (string, error) {
	buf := new(bytes.Buffer)
	if err := templates.ExecuteTemplate(buf, n.String(), filler); err != nil {
		return "", err
	}

	return buf.String(), nil
}
