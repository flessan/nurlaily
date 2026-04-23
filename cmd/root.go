package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/flessan/nurlaily/internal/build"
    "github.com/flessan/nurlaily/internal/draft"
    S "github.com/flessan/nurlaily/internal/style"
)

var rootCmd = &cobra.Command{
    Use:   "laily",
    Short: "NurLaily — Daily Draft",
    Long: `NurLaily adalah CLI untuk mencatat jurnal harian lewat terminal
lalu mengubahnya menjadi website statis yang estetis secara instan.

Cara pakai:
  laily draft "catatanmu hari ini"
  laily build
  laily list`,
    SilenceUsage: true,
}

var draftCmd = &cobra.Command{
    Use:   "draft [teks]",
    Short: "Tulis catatan baru ke jurnal hari ini",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        text := args[0]
        if err := draft.Write(text); err != nil {
            fmt.Fprintf(os.Stderr, "%s✗%s %s\n", S.Red, S.Reset, err)
            os.Exit(1)
        }
    },
}

var buildCmd = &cobra.Command{
    Use:   "build",
    Short: "Bangun website statis dari semua draft",
    Run: func(cmd *cobra.Command, args []string) {
        if err := build.Run(); err != nil {
            fmt.Fprintf(os.Stderr, "%s✗%s %s\n", S.Red, S.Reset, err)
            os.Exit(1)
        }
    },
}

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "Tampilkan semua draft yang tersimpan",
    Run: func(cmd *cobra.Command, args []string) {
        if err := draft.List(); err != nil {
            fmt.Fprintf(os.Stderr, "%s✗%s %s\n", S.Red, S.Reset, err)
            os.Exit(1)
        }
    },
}

func init() {
    rootCmd.AddCommand(draftCmd)
    rootCmd.AddCommand(buildCmd)
    rootCmd.AddCommand(listCmd)
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}