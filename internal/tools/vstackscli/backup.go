package vstackscli

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// RunBackup dispatches data backup commands (create, list, restore, info).
// Backup operations work directly on the filesystem and do not require a
// running server.
//
// Parameters:
//   - args: Command-line arguments (sub-command and flags)
func RunBackup(args []string) {
	if len(args) == 0 {
		printBackupUsage()
		os.Exit(1)
	}

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "create":
		backupCreateCmd(cmdArgs)
	case "list":
		backupListCmd()
	case "restore":
		backupRestoreCmd(cmdArgs)
	case "info":
		backupInfoCmd(cmdArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printBackupUsage()
		os.Exit(1)
	}
}

func printBackupUsage() {
	fmt.Println("Backup Commands:")
	fmt.Println("  vstacks backup create   [-output <path>]")
	fmt.Println("  vstacks backup list")
	fmt.Println("  vstacks backup restore  <file>")
	fmt.Println("  vstacks backup info     <file>")
}

func backupCreateCmd(args []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	output := fs.String("output", "", "Output zip path (default: data/backup_YYYYMMDD_HHmmss.zip)")
	dataPath := fs.String("data", "./data", "Data directory path")
	_ = fs.Parse(args)

	if *dataPath == "" {
		fmt.Fprintln(os.Stderr, "Error: -data is required")
		os.Exit(1)
	}

	zipPath := *output
	if zipPath == "" {
		zipPath = filepath.Join(*dataPath, fmt.Sprintf("backup_%s.zip", time.Now().Format("20060102_150405")))
	}

	src, err := os.Stat(*dataPath)
	if err != nil || !src.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: data directory not found: %s\n", *dataPath)
		os.Exit(1)
	}

	zipFile, err := os.Create(zipPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create zip: %v\n", err)
		os.Exit(1)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)

	err = filepath.Walk(*dataPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(*dataPath, filePath)
		if err != nil {
			return err
		}

		w, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(w, f)
		return err
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create backup: %v\n", err)
		os.Exit(1)
	}

	if err := zipWriter.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to close zip: %v\n", err)
		os.Exit(1)
	}

	info, _ := os.Stat(zipPath)
	fmt.Printf("Backup created: %s (%s)\n", zipPath, formatSize(info.Size()))
}

func backupListCmd() {
	dataPath := "./data"

	entries, err := os.ReadDir(dataPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot read data directory: %v\n", err)
		os.Exit(1)
	}

	var backups []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), "backup_") && strings.HasSuffix(e.Name(), ".zip") {
			info, _ := e.Info()
			backups = append(backups, fmt.Sprintf("  %-45s %s", e.Name(), formatSize(info.Size())))
		}
	}

	sort.Strings(backups)

	if len(backups) == 0 {
		fmt.Println("No backups found")
		return
	}

	fmt.Printf("Backups in %s/ (%d):\n", dataPath, len(backups))
	for _, b := range backups {
		fmt.Println(b)
	}
}

func backupRestoreCmd(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: backup file is required")
		os.Exit(1)
	}

	zipPath := args[0]
	dataPath := "./data"

	if _, err := os.Stat(zipPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: backup file not found: %s\n", zipPath)
		os.Exit(1)
	}

	archiveZip := filepath.Join(dataPath, fmt.Sprintf("archive_%s.zip", time.Now().Format("20060102_150405")))

	if err := os.Rename(dataPath, archiveZip); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to archive current data: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(dataPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create data directory: %v\n", err)
		os.Exit(1)
	}

	if err := unzip(zipPath, dataPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to restore backup: %v\n", err)
		fmt.Fprintf(os.Stderr, "Current data archived to: %s\n", archiveZip)
		os.Exit(1)
	}

	fmt.Printf("Backup restored from %s\n", zipPath)
	fmt.Printf("Previous data archived to: %s\n", archiveZip)
}

func backupInfoCmd(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: backup file is required")
		os.Exit(1)
	}

	zipPath := args[0]

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to open zip: %v\n", err)
		os.Exit(1)
	}
	defer r.Close()

	var totalSize int64
	for _, f := range r.File {
		totalSize += int64(f.UncompressedSize64)
	}

	fmt.Printf("Backup: %s (%s, %d files)\n", zipPath, formatSize(totalSize), len(r.File))

	for _, f := range r.File {
		if !f.FileInfo().IsDir() {
			fmt.Printf("  %s (%s)\n", f.Name, formatSize(f.FileInfo().Size()))
		}
	}
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
