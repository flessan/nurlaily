package cmd

import (
    "fmt"
    "log"
    "net/http"

    "github.com/flessan/nurlaily/internal/build"
    "github.com/spf13/cobra"
)

var outputDir string

var buildCmd = &cobra.Command{
    Use:   "build",
    Short: "Bangun website statis dari semua draft",
    RunE: func(cmd *cobra.Command, args []string) error {
        if err := build.BuildSite(outputDir); err != nil {
            return err
        }
        fmt.Printf("%s Website dibuat di %s/\n", green("✓"), bold(outputDir))
        return nil
    },
}

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Bangun dan preview website di localhost:3000",
    RunE: func(cmd *cobra.Command, args []string) error {
        if err := build.BuildSite(outputDir); err != nil {
            return err
        }
        fmt.Printf("%s Website dibuat. Preview di %s\n", green("✓"), bold("http://localhost:3000"))
        fmt.Println("Tekan Ctrl+C untuk berhenti.")
        fs := http.FileServer(http.Dir(outputDir))
        log.Fatal(http.ListenAndServe(":3000", fs))
        return nil
    },
}

func init() {
    buildCmd.Flags().StringVarP(&outputDir, "output", "o", "dist", "Folder output website")
    serveCmd.Flags().StringVarP(&outputDir, "output", "o", "dist", "Folder output website")
}
