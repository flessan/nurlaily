package build

import (
    "bytes"
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "time"

    "github.com/flessan/nurlaily/internal/draft"
    "github.com/flessan/nurlaily/internal/model"
    tpl "github.com/flessan/nurlaily/internal/template"
    "github.com/yuin/goldmark"
)

func BuildSite(outputDir string) error {
    days, err := parseAllDrafts()
    if err != nil {
        return err
    }
    if len(days) == 0 {
        return fmt.Errorf("belum ada catatan. Tulis dengan `./laily draft \"pesan\"` terlebih dahulu")
    }

    tagMap := make(map[string]int)
    totalEntries, totalWords := 0, 0
    for _, day := range days {
        totalEntries += len(day.Entries)
        for _, e := range day.Entries {
            totalWords += e.WordCount
            for _, t := range e.Tags {
                tagMap[t]++
            }
        }
    }

    var allTags []model.TagInfo
    for name, count := range tagMap {
        allTags = append(allTags, model.TagInfo{Name: name, Count: count})
    }
    sort.Slice(allTags, func(i, j int) bool { return allTags[i].Count > allTags[j].Count })

    data := model.PageData{
        Days:         days,
        TotalDays:    len(days),
        TotalEntries: totalEntries,
        TotalWords:   totalWords,
        AllTags:      allTags,
        Title:        "NurLaily — Daily Draft",
        GeneratedAt:  time.Now().Format("2 January 2006, 15:04"),
    }

    html, err := tpl.Render(data)
    if err != nil {
        return fmt.Errorf("gagal render template: %w", err)
    }

    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return err
    }
    return os.WriteFile(filepath.Join(outputDir, "index.html"), []byte(html), 0644)
}

func parseAllDrafts() ([]model.DayData, error) {
    daysInfo, err := draft.ListDrafts()
    if err != nil {
        return nil, err
    }

    md := goldmark.New()
    var days []model.DayData

    for _, di := range daysInfo {
        entries, err := draft.GetEntries(di.Date)
        if err != nil {
            continue
        }

        var elist []model.EntryData
        for _, e := range entries {
            var buf bytes.Buffer
            if err := md.Convert([]byte(e.Content), &buf); err != nil {
                continue
            }
            elist = append(elist, model.EntryData{
                Time:       e.Time,
                Content:    buf.String(),
                ContentRaw: e.Content,
                Mood:       e.Mood,
                Tags:       e.Tags,
                WordCount:  e.WordCount,
                ReadTime:   readTime(e.WordCount),
            })
        }

        if len(elist) > 0 {
            days = append(days, model.DayData{
                Date:      di.Date,
                DateHuman: di.DateHuman,
                Entries:   elist,
                Count:     len(elist),
            })
        }
    }
    return days, nil
}

func readTime(words int) string {
    if words == 0 {
        return "< 1 min"
    }
    m := words / 200
    if m < 1 {
        return "< 1 min"
    }
    return fmt.Sprintf("%d min", m)
}
