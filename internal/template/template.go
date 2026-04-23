package template

import (
    _ "embed"
    "fmt"
    "html/template"
    "strings"
)

// HTML adalah alias dari template.HTML supaya package build bisa menulisnya
// tanpa harus import "html/template" langsung.
type HTML = template.HTML

// Entry mewakili satu catatan jurnal yang sudah dikonversi.
type Entry struct {
    DateRaw string // "2026-01-02" — untuk sorting (tidak ditampilkan)
    Date    string // "2 January 2006" — ditampilkan di UI
    Slug    string // "20260102" — untuk ID elemen HTML
    Content HTML   // konten markdown yang sudah jadi HTML
    Preview string // satu baris preview untuk sidebar
}

//go:embed index.html
var rawTemplate string

// Render menyuntikkan data entries ke dalam template dan mengembalikan HTML lengkap.
func Render(entries []Entry) (string, error) {
    funcs := template.FuncMap{
        "delay": func(i int) string {
            return fmt.Sprintf("%dms", i*120)
        },
    }

    t, err := template.New("index").Funcs(funcs).Parse(rawTemplate)
    if err != nil {
        return "", fmt.Errorf("parse template gagal: %w", err)
    }

    data := struct {
        Title   string
        Entries []Entry
        Count   int
    }{
        Title:   "NurLaily — Daily Draft",
        Entries: entries,
        Count:   len(entries),
    }

    var buf strings.Builder
    if err := t.Execute(&buf, data); err != nil {
        return "", fmt.Errorf("eksekusi template gagal: %w", err)
    }

    return buf.String(), nil
}