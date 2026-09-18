package vstackscli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

	"vorpalstacks/internal/core/storage"
)

// RunStorage dispatches offline storage maintenance commands. These
// commands operate directly on the region databases and require the server
// to be stopped.
func RunStorage(args []string) {
	if len(args) == 0 {
		printStorageUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "vacuum":
		vacuumCmd(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
		printStorageUsage()
		os.Exit(1)
	}
}

func printStorageUsage() {
	fmt.Println("Storage Commands:")
	fmt.Println("  vstacks storage vacuum [-data-path <dir>]   Reclaim space: compact every region database (server must be stopped)")
}

func vacuumCmd(args []string) {
	flags := flag.NewFlagSet("vacuum", flag.ExitOnError)
	dataPath := flags.String("data-path", "./data", "Data directory path")
	_ = flags.Parse(args)

	if err := runVacuum(*dataPath, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// runVacuum full-range compacts every Pebble database directory under
// dataPath, printing per-database before/after sizes. Directories without
// a Pebble database marker (blob stores, chunk files) are left untouched.
func runVacuum(dataPath string, out io.Writer) error {
	entries, err := os.ReadDir(dataPath)
	if err != nil {
		return fmt.Errorf("data directory not readable: %w", err)
	}

	var dbDirs []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dir := filepath.Join(dataPath, entry.Name())
		if isPebbleDBDir(dir) {
			dbDirs = append(dbDirs, dir)
		}
	}
	if len(dbDirs) == 0 {
		fmt.Fprintln(out, "no region databases found — nothing to vacuum")
		return nil
	}

	for _, dir := range dbDirs {
		name := filepath.Base(dir)
		before, err := dirSize(dir)
		if err != nil {
			return fmt.Errorf("%s: size measurement failed: %w", name, err)
		}
		if err := requireFreeSpace(dataPath, before); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}

		store, err := storage.NewPebbleStorage(&storage.Config{Path: dir})
		if err != nil {
			if _, statErr := os.Stat(filepath.Join(dir, "LOCK")); statErr == nil {
				return fmt.Errorf("%s: cannot open database (%w); the server appears to be running — stop it before vacuuming", name, err)
			}
			return fmt.Errorf("%s: cannot open database: %w", name, err)
		}

		if err := vacuumStepError(store.Compact(), store.Close()); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}

		after, err := dirSize(dir)
		if err != nil {
			return fmt.Errorf("%s: size measurement failed: %w", name, err)
		}
		fmt.Fprintf(out, "%s: %s → %s (reclaimed %s)\n", name, humanSize(before), humanSize(after), humanSize(before-after))
	}
	return nil
}

// vacuumStepError labels and joins the two failure modes of one
// database's vacuum step. Both failures report together: returning on the
// compaction error alone would swallow a close failure behind it, and a
// dropped close leaves the database lock held with nothing reported.
func vacuumStepError(compactErr, closeErr error) error {
	if compactErr != nil {
		compactErr = fmt.Errorf("compaction failed: %w", compactErr)
	}
	if closeErr != nil {
		closeErr = fmt.Errorf("close failed: %w", closeErr)
	}
	return errors.Join(compactErr, closeErr)
}

// isPebbleDBDir reports whether dir holds a Pebble database, recognised by
// the LOCK and MANIFEST files every Pebble database directory carries. It
// never opens or creates the directory, so non-database directories are
// left untouched.
func isPebbleDBDir(dir string) bool {
	if _, err := os.Stat(filepath.Join(dir, "LOCK")); err != nil {
		return false
	}
	matches, err := filepath.Glob(filepath.Join(dir, "MANIFEST-*"))
	return err == nil && len(matches) > 0
}

// dirSize sums the byte sizes of every regular file under dir.
func dirSize(dir string) (int64, error) {
	var total int64
	err := filepath.WalkDir(dir, func(_ string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		total += info.Size()
		return nil
	})
	return total, err
}

// requireFreeSpace aborts before a compaction that would rewrite the whole
// live set: at least twice the database size must remain free on the
// filesystem holding dataPath.
func requireFreeSpace(dataPath string, dbSize int64) error {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(dataPath, &stat); err != nil {
		return fmt.Errorf("free-space check failed: %w", err)
	}
	free := int64(stat.Bavail) * int64(stat.Bsize)
	if free < 2*dbSize {
		return fmt.Errorf("insufficient free space: compaction rewrites the live set (needs up to %s, only %s free)",
			humanSize(2*dbSize), humanSize(free))
	}
	return nil
}

func humanSize(n int64) string {
	switch {
	case n >= 1024*1024*1024:
		return fmt.Sprintf("%.1f GiB", float64(n)/(1024*1024*1024))
	case n >= 1024*1024:
		return fmt.Sprintf("%.1f MiB", float64(n)/(1024*1024))
	case n >= 1024:
		return fmt.Sprintf("%.1f KiB", float64(n)/1024)
	default:
		return fmt.Sprintf("%d B", n)
	}
}
