package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/thio/nurlaily/internal/build"
    "github.com/thio/nurlaily/internal/draft"
)

// Warna ANSI untuk output terminal yang enak dilihat
const (
    reset  = "\033[0m"
    green  = "\033[32m"
    cyan   = "\033[36m"
    purple = "\033[35m"
    gray   = "\033[90m"
    yellow = "\033[33m"
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

// Perintah: laily draft "isi catatan"
var draftCmd = &cobra.Command{
    Use:   "draft [teks]",
    Short: "Tulis catatan baru ke jurnal hari ini",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        text := args[0]
        if err := draft.Write(text); err != nil {
            fmt.Fprintf(os.Stderr, "%s✗%s %s\n", red(), reset, err)
            os.Exit(1)
        }
    },
}

// Perintah: laily build
var buildCmd = &cobra.Command{
    Use:   "build",
    Short: "Bangun website statis dari semua draft",
    Run: func(cmd *cobra.Command, args []string) {
        if err := build.Run(); err != nil {
            fmt.Fprintf(os.Stderr, "%s✗%s %s\n", red(), reset, err)
            os.Exit(1)
        }
    },
}

// Perintah: laily list
var listCmd = &cobra.Command{
    Use:   "list",
    Short: "Tampilkan semua draft yang tersimpan",
    Run: func(cmd *cobra.Command, args []string) {
        if err := draft.List(); err != nil {
            fmt.Fprintf(os.Stderr, "%s✗%s %s\n", red(), reset, err)
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

// red mengembalikan kode warna merah ANSI
func red() string { return "\033[31m" }

// Export warna agar bisa dipakai package lain
func Green() string  { return green }
func Cyan() string   { return cyan }
func Purple() string { return purple }
func Gray() string   { return gray }
func Yellow() string { return yellow }
func Reset() string  { return reset }
func Red() string    { return red() }