package draft

import (
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "sort"
    "strings"
    "time"
)

const DraftDir = "drafts"

var hariIndo = []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
var bulanIndo = []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni",
    "Juli", "Agustus", "September", "Oktober", "November", "Desember"}

type Entry struct {
    Time      string
    Content   string
    Mood      string
    Tags      []string
    WordCount int
}

type DayInfo struct {
    Date       string
    DateHuman  string
    EntryCount int
}

type Stats struct {
    TotalDays     int
    TotalEntries  int
    TotalWords    int
    CurrentStreak int
    LongestStreak int
    TopTags       []TagCount
}

type TagCount struct {
    Tag   string
    Count int
}

var tagRe = regexp.MustCompile(`#([a-zA-Z][\w]*)`)
var headerRe = regexp.MustCompile(`^## (\d{2}:\d{2})(?:\s+(.+))?$`)

// ─── WriteEntry ──────────────────────────────────────────
func WriteEntry(text string, mood string, tags []string) error {
    if err := ensureDir(); err != nil {
        return err
    }

    now := time.Now()
    filename := now.Format("2006-01-02") + ".md"
    fp := filepath.Join(DraftDir, filename)

    timestamp := now.Format("15:04")

    var header string
    if mood != "" {
        header = fmt.Sprintf("## %s %s", timestamp, mood)
    } else {
        header = fmt.Sprintf("## %s", timestamp)
    }

    if len(tags) > 0 {
        ts := make([]string, len(tags))
        for i, t := range tags {
            ts[i] = "#" + strings.TrimSpace(t)
        }
        text = text + " " + strings.Join(ts, " ")
    }

    f, err := os.OpenFile(fp, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer f.Close()

    info, _ := f.Stat()
    if info.Size() == 0 {
        f.WriteString(fmt.Sprintf("# %s\n", now.Format("2006-01-02")))
    }

    _, err = f.WriteString(fmt.Sprintf("\n%s\n\n%s\n", header, text))
    return err
}

// ─── ListDrafts ──────────────────────────────────────────
func ListDrafts() ([]DayInfo, error) {
    if err := ensureDir(); err != nil {
        return nil, err
    }

    entries, err := os.ReadDir(DraftDir)
    if err != nil {
        return nil, err
    }

    var days []DayInfo
    for _, e := range entries {
        if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || strings.HasPrefix(e.Name(), ".") {
            continue
        }
        dateStr := strings.TrimSuffix(e.Name(), ".md")
        content, err := os.ReadFile(filepath.Join(DraftDir, e.Name()))
        if err != nil {
            continue
        }
        parsed := parseEntries(string(content))
        days = append(days, DayInfo{
            Date:       dateStr,
            DateHuman:  formatDateHuman(dateStr),
            EntryCount: len(parsed),
        })
    }

    sort.Slice(days, func(i, j int) bool {
        return days[i].Date > days[j].Date
    })
    return days, nil
}

// ─── GetEntries ──────────────────────────────────────────
func GetEntries(date string) ([]Entry, error) {
    fp := filepath.Join(DraftDir, date+".md")
    content, err := os.ReadFile(fp)
    if err != nil {
        return nil, fmt.Errorf("file draft untuk %s tidak ditemukan", date)
    }
    return parseEntries(string(content)), nil
}

// ─── GetToday ────────────────────────────────────────────
func GetToday() ([]Entry, error) {
    today := time.Now().Format("2006-01-02")
    return GetEntries(today)
}

// ─── GetStats ────────────────────────────────────────────
func GetStats() (*Stats, error) {
    days, err := ListDrafts()
    if err != nil {
        return nil, err
    }

    stats := &Stats{TotalDays: len(days)}
    tagMap := make(map[string]int)
    var allDates []string

    for _, day := range days {
        allDates = append(allDates, day.Date)
        stats.TotalEntries += day.EntryCount
        entries, err := GetEntries(day.Date)
        if err != nil {
            continue
        }
        for _, e := range entries {
            stats.TotalWords += e.WordCount
            for _, t := range e.Tags {
                tagMap[t]++
            }
        }
    }

    stats.CurrentStreak = calcStreak(allDates, false)
    stats.LongestStreak = calcStreak(allDates, true)

    var tags []TagCount
    for tag, count := range tagMap {
        tags = append(tags, TagCount{Tag: tag, Count: count})
    }
    sort.Slice(tags, func(i, j int) bool { return tags[i].Count > tags[j].Count })
    if len(tags) > 10 {
        tags = tags[:10]
    }
    stats.TopTags = tags
    return stats, nil
}

// ─── DeleteEntry ─────────────────────────────────────────
func DeleteEntry(date string, index int) error {
    fp := filepath.Join(DraftDir, date+".md")
    content, err := os.ReadFile(fp)
    if err != nil {
        return fmt.Errorf("file draft untuk %s tidak ditemukan", date)
    }

    entries := parseEntries(string(content))
    if index < 0 || index >= len(entries) {
        return fmt.Errorf("indeks %d tidak valid (total: %d catatan)", index, len(entries))
    }

    entries = append(entries[:index], entries[index+1:]...)

    var b strings.Builder
    b.WriteString(fmt.Sprintf("# %s\n", date))
    for _, e := range entries {
        b.WriteString("\n")
        if e.Mood != "" {
            b.WriteString(fmt.Sprintf("## %s %s\n\n%s\n", e.Time, e.Mood, e.Content))
        } else {
            b.WriteString(fmt.Sprintf("## %s\n\n%s\n", e.Time, e.Content))
        }
    }

    return os.WriteFile(fp, []byte(b.String()), 0644)
}

// ─── DeleteDay ───────────────────────────────────────────
func DeleteDay(date string) error {
    fp := filepath.Join(DraftDir, date+".md")
    if _, err := os.Stat(fp); err != nil {
        return fmt.Errorf("file draft untuk %s tidak ditemukan", date)
    }
    return os.Remove(fp)
}

// ─── InitDrafts ──────────────────────────────────────────
func InitDrafts() error {
    if err := os.MkdirAll(DraftDir, 0755); err != nil {
        return err
    }

    now := time.Now()
    filename := now.Format("2006-01-02") + ".md"
    fp := filepath.Join(DraftDir, filename)

    if _, err := os.Stat(fp); err == nil {
        return fmt.Errorf("file draft hari ini sudah ada (%s)", filename)
    }

    content := fmt.Sprintf("# %s\n\n## %s \xF0\x9F\x8E\x89\n\nSelamat datang di **NurLaily**! Ini adalah catatan pertamamu.\n\n"+
        "Gunakan perintah `./laily draft \"pesan\"` untuk menulis catatan baru.\n"+
        "Tambahkan `--mood \xF0\x9F\x98\x8A` untuk menandai suasana hati.\n"+
        "Gunakan `#tag` di dalam pesan untuk mengkategorikan.\n\n"+
        "Jalankan `./laily build` untuk mengubah semua catatan menjadi website statis yang estetis.\n",
        now.Format("2006-01-02"), now.Format("15:04"))

    return os.WriteFile(fp, []byte(content), 0644)
}

// ─── Helpers ─────────────────────────────────────────────
func ensureDir() error {
    return os.MkdirAll(DraftDir, 0755)
}

func formatDateHuman(dateStr string) string {
    t, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        return dateStr
    }
    return fmt.Sprintf("%s, %d %s %d", hariIndo[t.Weekday()], t.Day(), bulanIndo[t.Month()-1], t.Year())
}

func parseEntries(content string) []Entry {
    lines := strings.Split(content, "\n")
    var entries []Entry
    var cur *Entry
    var buf []string

    for _, line := range lines {
        if strings.HasPrefix(line, "# ") && !strings.HasPrefix(line, "## ") {
            continue
        }
        if m := headerRe.FindStringSubmatch(line); m != nil {
            if cur != nil {
                cur.Content = strings.TrimSpace(strings.Join(buf, "\n"))
                cur.Tags = extractTags(cur.Content)
                cur.WordCount = countWords(cur.Content)
                entries = append(entries, *cur)
            }
            cur = &Entry{Time: m[1], Mood: m[2]}
            buf = nil
        } else if cur != nil {
            buf = append(buf, line)
        }
    }
    if cur != nil {
        cur.Content = strings.TrimSpace(strings.Join(buf, "\n"))
        cur.Tags = extractTags(cur.Content)
        cur.WordCount = countWords(cur.Content)
        entries = append(entries, *cur)
    }
    return entries
}

func extractTags(text string) []string {
    matches := tagRe.FindAllStringSubmatchIndex(text, -1)
    var tags []string
    seen := make(map[string]bool)
    for _, m := range matches {
        pos := m[0]
        if pos == 0 || (pos > 0 && text[pos-1] == '\n') {
            continue
        }
        tag := text[m[2]:m[3]]
        if !seen[tag] {
            seen[tag] = true
            tags = append(tags, tag)
        }
    }
    return tags
}

func countWords(text string) int {
    return len(strings.Fields(text))
}

func calcStreak(dates []string, longest bool) int {
    if len(dates) == 0 {
        return 0
    }

    sorted := make([]string, len(dates))
    copy(sorted, dates)
    sort.Strings(sorted)

    if longest {
        best, cur := 1, 1
        for i := 1; i < len(sorted); i++ {
            prev, _ := time.Parse("2006-01-02", sorted[i-1])
            curr, _ := time.Parse("2006-01-02", sorted[i])
            if curr.Sub(prev).Hours()/24 == 1 {
                cur++
                if cur > best {
                    best = cur
                }
            } else {
                cur = 1
            }
        }
        return best
    }

    // Current streak (from today backwards)
    today := time.Now().Format("2006-01-02")
    streak := 0
    expected := today
    for _, d := range dates {
        if d == expected {
            streak++
            t, _ := time.Parse("2006-01-02", expected)
            expected = t.AddDate(0, 0, -1).Format("2006-01-02")
        } else if d < expected {
            break
        }
    }
    return streak
}
