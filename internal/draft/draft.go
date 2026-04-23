package draft

import (
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "time"

    S "github.com/flessan/nurlaily/internal/style"
)

const draftsDir = "drafts"

func Write(text string) error {
    if err := os.MkdirAll(draftsDir, 0755); err != nil {
        return fmt.Errorf("gagal membuat folder %s: %w", draftsDir, err)
    }

    now := time.Now()
    filename := fmt.Sprintf("%s.md", now.Format("2006-01-02"))
    fp := filepath.Join(draftsDir, filename)

    f, err := os.OpenFile(fp, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("gagal membuka file: %w", err)
    }
    defer f.Close()

    waktu := now.Format("15:04")
    entry := fmt.Sprintf("\n## %s\n%s\n", waktu, text)

    if _, err := f.WriteString(entry); err != nil {
        return fmt.Errorf("gagal menulis ke file: %w", err)
    }

    fmt.Printf("%s✓%s Catatan disimpan di %s%s%s\n",
        S.Green, S.Reset, S.Cyan, fp, S.Reset)
    fmt.Printf("  %s%s%s\n",
        S.Gray, truncate(text, 55), S.Reset)

    return nil
}

func List() error {
    entries, err := os.ReadDir(draftsDir)
    if err != nil {
        if os.IsNotExist(err) {
            fmt.Printf("%s⚠%s Folder %s%s%s belum ada.\n",
                S.Yellow, S.Reset, S.Cyan, draftsDir, S.Reset)
            fmt.Printf("  Mulai menulis: %slaily draft \"catatanmu\"%s\n",
                S.Purple, S.Reset)
            return nil
        }
        return fmt.Errorf("gagal membaca folder: %w", err)
    }

    var files []string
    for _, e := range entries {
        if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
            files = append(files, e.Name())
        }
    }

    if len(files) == 0 {
        fmt.Printf("%s⚠%s Belum ada draft.\n", S.Yellow, S.Reset)
        return nil
    }

    sort.Sort(sort.Reverse(sort.StringSlice(files)))

    fmt.Printf("%s%sNurLaily — Daftar Draft%s\n", S.Purple, S.Gray, S.Reset)
    fmt.Printf("%s%s────────────────────%s\n\n", S.Purple, S.Gray, S.Reset)

    for _, f := range files {
        date := strings.TrimSuffix(f, ".md")
        info, _ := os.Stat(filepath.Join(draftsDir, f))
        size := "0 B"
        if info != nil {
            if info.Size() < 1024 {
                size = fmt.Sprintf("%d B", info.Size())
            } else {
                size = fmt.Sprintf("%.1f KB", float64(info.Size())/1024)
            }
        }
        fmt.Printf("  %s%s%s  %s%s%s\n",
            S.Cyan, date, S.Reset,
            S.Gray, size, S.Reset)
    }

    fmt.Printf("\n  Total: %s%d%s file\n", S.Cyan, len(files), S.Reset)
    return nil
}

func truncate(s string, max int) string {
    s = strings.ReplaceAll(s, "\n", " ")
    if len(s) <= max {
        return s
    }
    return s[:max] + "..."
}