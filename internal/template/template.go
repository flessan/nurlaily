package template

import (
    "embed"
    "html/template"
    "strings"

    "github.com/flessan/nurlaily/internal/model"
)

//go:embed index.html
var templateFS embed.FS

// safePageData wraps PageData to convert Content string → template.HTML
type safePageData struct {
    *model.PageData
    Days []safeDayData
}

type safeDayData struct {
    Date      string
    DateHuman string
    Entries   []safeEntryData
    Count     int
}

type safeEntryData struct {
    Time       string
    Content    template.HTML
    ContentRaw string
    Mood       string
    Tags       []string
    WordCount  int
    ReadTime   string
}

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

    // Wrap data for safe HTML
    safe := safePageData{PageData: &data}
    for _, d := range data.Days {
        sd := safeDayData{Date: d.Date, DateHuman: d.DateHuman, Count: d.Count}
        for _, e := range d.Entries {
            sd.Entries = append(sd.Entries, safeEntryData{
                Time:       e.Time,
                Content:    template.HTML(e.Content),
                ContentRaw: e.ContentRaw,
                Mood:       e.Mood,
                Tags:       e.Tags,
                WordCount:  e.WordCount,
                ReadTime:   e.ReadTime,
            })
        }
        safe.Days = append(safe.Days, sd)
    }

    var buf strings.Builder
    if err := t.Execute(&buf, safe); err != nil {
        return "", err
    }

    return buf.String(), nil
}
