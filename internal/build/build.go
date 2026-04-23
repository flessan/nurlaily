package build

import (
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "time"

    "github.com/yuin/goldmark"
    C "github.com/thio/nurlaily/cmd"
    tpl "github.com/thio/nurlaily/internal/template"
)

const distDir = "dist"

// Run membaca semua .md di drafts/, konversi ke HTML, suntik ke template, tulis ke dist/.
func Run() error {
    fmt.Printf("%s⚡%s Membangun website statis...\n", C.Purple(), C.Reset())

    entries, err := readAllDrafts()
    if err != nil {
        return err
    }

    if len(entries) == 0 {
        fmt.Printf("%s⚠%s Tidak ada draft di folder %s%s%s\n",
            C.Yellow(), C.Reset(), C.Cyan(), "drafts/", C.Reset())
        fmt.Printf("  Mulai menulis: %slaily draft \"catatanmu\"%s\n",
            C.Purple(), C.Reset())
        return nil
    }

    if err := os.MkdirAll(distDir, 0755); err != nil {
        return fmt.Errorf("gagal membuat folder %s: %w", distDir, err)
    }

    html, err := tpl.Render(entries)
    if err != nil {
        return fmt.Errorf("gagal render template: %w", err)
    }

    outPath := filepath.Join(distDir, "index.html")
    if err := os.WriteFile(outPath, []byte(html), 0644); err != nil {
        return fmt.Errorf("gagal menulis %s: %w", outPath, err)
    }

    fmt.Printf("%s✓%s %s%d%s draft diproses\n",
        C.Green(), C.Reset(), C.Cyan(), len(entries), C.Reset())
    fmt.Printf("%s✓%s Website siap di %s%s%s\n",
        C.Green(), C.Reset(), C.Cyan(), outPath, C.Reset())

    return nil
}

// readAllDrafts membaca semua file .md dan mengubahnya ke Entry.
func readAllDrafts() ([]tpl.Entry, error) {
    md := goldmark.New()
    var entries []tpl.Entry

    err := filepath.WalkDir("drafts", func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() || !strings.HasSuffix(path, ".md") {
            return nil
        }

        raw, err := os.ReadFile(path)
        if err != nil {
            return fmt.Errorf("gagal baca %s: %w", path, err)
        }

        // Konversi Markdown → HTML
        var buf strings.Builder
        if err := md.Convert(raw, &buf); err != nil {
            return fmt.Errorf("gagal konversi %s: %w", path, err)
        }

        // Ekstrak tanggal dari nama file (2026-04-23.md)
        basename := strings.TrimSuffix(filepath.Base(path), ".md")
        t, err := time.Parse("2006-01-02", basename)
        if err != nil {
            return nil // lewati file dengan nama tidak valid
        }

        slug := strings.ReplaceAll(basename, "-", "")

        entries = append(entries, tpl.Entry{
            DateRaw: basename,                    // "2026-04-23" — untuk sorting
            Date:    t.Format("2 January 2006"),   // "23 April 2026" — untuk tampilan
            Slug:    slug,                         // "20260423" — untuk anchor ID
            Content: tpl.HTML(buf.String()),       // template.HTML supaya tidak di-escape
            Preview: extractPreview(string(raw)),  // baris pertama yang bermakna
        })
        return nil
    })

    if err != nil && !os.IsNotExist(err) {
        return nil, err
    }

    // Urutkan dari terbaru ke terlama
    sort.Slice(entries, func(i, j int) bool {
        return entries[i].DateRaw > entries[j].DateRaw
    })

    return entries, nil
}

// extractPreview mengambil baris pertama yang bermakna dari source markdown.
func extractPreview(md string) string {
    lines := strings.Split(strings.TrimSpace(md), "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }
        // Lewati baris heading
        if strings.HasPrefix(line, "#") {
            continue
        }
        // Bersihkan markdown formatting ringan
        line = strings.NewReplacer("**", "", "*", "", "`", "", "## ", "").Replace(line)
        if len(line) > 65 {
            return line[:65] + "..."
        }
        return line
    }
    return "Catatan tanpa preview"
}