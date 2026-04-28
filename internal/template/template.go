package template

import (
    "embed"
    "strings"
    "text/template"

    "github.com/flessan/nurlaily/internal/model"
)

//go:embed index.html
var templateFS embed.FS

func Render(data model.PageData) (string, error) {
    content, err := templateFS.ReadFile("index.html")
    if err != nil {
        return "", err
    }

    funcMap := template.FuncMap{
        "join": func(sep string, tags []string) string {
            return strings.Join(tags, sep)
        },
    }

    t, err := template.New("index").Funcs(funcMap).Parse(string(content))
    if err != nil {
        return "", err
    }

    var buf strings.Builder
    if err := t.Execute(&buf, data); err != nil {
        return "", err
    }

    return buf.String(), nil
}
