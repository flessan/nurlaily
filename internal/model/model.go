package model

type EntryData struct {
    Time       string
    Content    string // akan di-pass sebagai template.HTML di render
    ContentRaw string
    Mood       string
    Tags       []string
    WordCount  int
    ReadTime   string
}

type DayData struct {
    Date      string
    DateHuman string
    Entries   []EntryData
    Count     int
}

type TagInfo struct {
    Name  string
    Count int
}

type PageData struct {
    Days         []DayData
    TotalDays    int
    TotalEntries int
    TotalWords   int
    AllTags      []TagInfo
    Title        string
    GeneratedAt  string
}
