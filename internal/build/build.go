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
    S "github.com/flessan/nurlaily/internal/style"
    tpl "github.com/flessan/nurlaily/internal/template"
)

const distDir = "dist"

func Run() error {
    fmt.Printf("%s⚡%s Membangun website statis...\n", S.Purple, S.Reset)

    entries, err := readAllDrafts()
    if err != nil {
        return err
    }

    if len(entries) == 0 {
        fmt.Printf("%s⚠%s Tidak ada draft di folder %s%s%s\n",
            S.Yellow, S.Reset, S.Cyan, "drafts/", S.Reset)
        fmt.Printf("  Mulai menulis: %slaily draft \"catatanmu\"%s\n",
            S.Purple, S.Reset)
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
        S.Green, S.Reset, S.Cyan, len(entries), S.Reset)
    fmt.Printf("%s✓%s Website siap di %s%s%s\n",
        S.Green, S.Reset, S.Cyan, outPath, S.Reset)

    return nil
}

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

        var buf strings.Builder
        if err := md.Convert(raw, &buf); err != nil {
            return fmt.Errorf("gagal konversi %s: %w", path, err)
        }

        basename := strings.TrimSuffix(filepath.Base(path), ".md")
        t, err := time.Parse("2006-01-02", basename)
        if err != nil {
            return nil
        }

        slug := strings.ReplaceAll(basename, "-", "")

        entries = append(entries, tpl.Entry{
            DateRaw: basename,
            Date:    t.Format("2 January 2006"),
            Slug:    slug,
            Content: tpl.HTML(buf.String()),
            Preview: extractPreview(string(raw)),
        })
        return nil
    })

    if err != nil && !os.IsNotExist(err) {
        return nil, err
    }

    sort.Slice(entries, func(i, j int) bool {
        return entries[i].DateRaw > entries[j].DateRaw
    })

    return entries, nil
}

func extractPreview(md string) string {
    lines := strings.Split(strings.TrimSpace(md), "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }
        if strings.HasPrefix(line, "#") {
            continue
        }
        line = strings.NewReplacer("**", "", "*", "", "`", "", "## ", "").Replace(line)
        if len(line) > 65 {
            return line[:65] + "..."
        }
        return line
    }
    return "Catatan tanpa preview"
}