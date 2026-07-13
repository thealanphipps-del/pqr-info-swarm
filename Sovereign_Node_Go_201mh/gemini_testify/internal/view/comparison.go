package view

import (
"html/template"
"io"
)

type ComparisonView struct {
Templates *template.Template
}

func (v *ComparisonView) RenderComparison(w io.Writer, data interface{}) {
v.Templates.ExecuteTemplate(w, "index.html", data)
}
