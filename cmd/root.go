package cmd

import (
    "fmt"
    "os"
    "strings"

    "github.com/flessan/nurlaily/internal/draft"
    "github.com/spf13/cobra"
)

var (
    moodFlags []string
    tagFlags  []string
)

func isTerminal() bool {
    fi, err := os.Stdout.Stat()
    if err != nil {
        return false
    }
    return fi.Mode()&os.ModeCharDevice != 0
}

func c(code, s string) string {
    if !isTerminal() {
        return s
    }
    return code + s + "\033[0m"
}
func cyan(s string) string    { return c("\033[36m", s) }
func purple(s string) string  { return c("\033[35m", s) }
func green(s string) string   { return c("\033[32m", s) }
func yellow(s string) string  { return c("\033[33m", s) }
func gray(s string) string    { return c("\033[90m", s) }
func bold(s string) string    { return c("\033[1m", s) }
func red(s string) string     { return c("\033[31m", s) }

var rootCmd = &cobra.Command{
    Use:   "laily",
    Short: "NurLaily — Catat jurnal harian lewat terminal",
    Long:  "NurLaily adalah CLI untuk menulis jurnal harian dan mengubahnya menjadi website statis yang estetis.",
    RunE: func(cmd *cobra.Command, args []string) error {
        return cmd.Help()
    },
}

// ─── draft ───────────────────────────────────────────────
var draftCmd = &cobra.Command{
    Use:   "draft \"pesan\"",
    Short: "Tulis catatan baru",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        mood := ""
        if len(moodFlags) > 0 {
            mood = strings.Join(moodFlags, " ")
        }
        if err := draft.WriteEntry(args[0], mood, tagFlags); err != nil {
            return err
        }
        fmt.Printf("%s ✓ Catatan ditambahkan%s\n", green("✓"), gray(""))
        return nil
    },
}

// ─── list ────────────────────────────────────────────────
var listCmd = &cobra.Command{
    Use:   "list",
    Short: "Lihat daftar semua draft",
    RunE: func(cmd *cobra.Command, args []string) error {
        days, err := draft.ListDrafts()
        if err != nil {
            return err
        }
        if len(days) == 0 {
            fmt.Println(gray("Belum ada catatan. Tulis dengan ") + cyan("./laily draft \"pesan\""))
            return nil
        }
        for _, d := range days {
            fmt.Printf("  %s  %s  %s\n",
                cyan(d.Date),
                purple(d.DateHuman),
                gray(fmt.Sprintf("(%d catatan)", d.EntryCount)),
            )
        }
        return nil
    },
}

// ─── today ───────────────────────────────────────────────
var todayCmd = &cobra.Command{
    Use:   "today",
    Short: "Lihat catatan hari ini",
    RunE: func(cmd *cobra.Command, args []string) error {
        entries, err := draft.GetToday()
        if err != nil {
            return err
        }
        if len(entries) == 0 {
            fmt.Println(gray("Belum ada catatan hari ini."))
            return nil
        }
        for _, e := range entries {
            fmt.Printf("\n  %s", bold(cyan(e.Time)))
            if e.Mood != "" {
                fmt.Printf(" %s", e.Mood)
            }
            if len(e.Tags) > 0 {
                for _, t := range e.Tags {
                    fmt.Printf(" %s", green("#"+t))
                }
            }
            fmt.Printf("\n  %s\n", e.Content)
        }
        fmt.Println()
        return nil
    },
}

// ─── stats ───────────────────────────────────────────────
var statsCmd = &cobra.Command{
    Use:   "stats",
    Short: "Lihat statistik jurnal",
    RunE: func(cmd *cobra.Command, args []string) error {
        s, err := draft.GetStats()
        if err != nil {
            return err
        }
        fmt.Println()
        fmt.Printf("  %s %d hari menulis\n", purple("📅"), s.TotalDays)
        fmt.Printf("  %s %d catatan\n", purple("📝"), s.TotalEntries)
        fmt.Printf("  %s %d kata\n", purple("📊"), s.TotalWords)
        fmt.Printf("  %s %d hari (streak saat ini)\n", purple("🔥"), s.CurrentStreak)
        fmt.Printf("  %s %d hari (streak terpanjang)\n", purple("🏆"), s.LongestStreak)
        if len(s.TopTags) > 0 {
            fmt.Printf("  %s ", purple("🏷"))
            for i, t := range s.TopTags {
                if i > 0 {
                    fmt.Print(gray(", "))
                }
                fmt.Printf("%s(%d)", green("#"+t.Tag), t.Count)
            }
            fmt.Println()
        }
        fmt.Println()
        return nil
    },
}

// ─── delete ──────────────────────────────────────────────
var deleteCmd = &cobra.Command{
    Use:   "delete [tanggal] [indeks]",
    Short: "Hapus catatan (indeks mulai dari 0)",
    Long: `Hapus catatan berdasarkan tanggal dan indeks.
Contoh: ./laily delete 2026-04-23 0
Tanpa indeks: hapus seluruh catatan hari tersebut.`,
    Args: cobra.RangeArgs(1, 2),
    RunE: func(cmd *cobra.Command, args []string) error {
        date := args[0]
        if len(args) == 2 {
            var index int
            _, err := fmt.Sscanf(args[1], "%d", &index)
            if err != nil {
                return fmt.Errorf("indeks harus berupa angka")
            }
            if err := draft.DeleteEntry(date, index); err != nil {
                return err
            }
            fmt.Printf("%s Catatan #%d pada %s dihapus\n", green("✓"), index, cyan(date))
        } else {
            if err := draft.DeleteDay(date); err != nil {
                return err
            }
            fmt.Printf("%s Semua catatan %s dihapus\n", green("✓"), cyan(date))
        }
        return nil
    },
}

// ─── init ────────────────────────────────────────────────
var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Inisialisasi folder drafts dengan contoh catatan",
    RunE: func(cmd *cobra.Command, args []string) error {
        if err := draft.InitDrafts(); err != nil {
            return err
        }
        fmt.Printf("%s Folder drafts/ siap dengan catatan contoh\n", green("✓"))
        return nil
    },
}

func Execute() error {
    draftCmd.Flags().StringSliceVar(&moodFlags, "mood", nil, "Suasana hati (emoji, contoh: 😊)")
    draftCmd.Flags().StringSliceVar(&tagFlags, "tag", nil, "Tag kategori (contoh: go,belajar)")

    rootCmd.AddCommand(draftCmd, listCmd, todayCmd, statsCmd, deleteCmd, initCmd)
    rootCmd.AddCommand(buildCmd, serveCmd)
    return rootCmd.Execute()
}
