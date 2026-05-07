// driveScanner — find photos that crash unifi-drive 4.2.4's photo backup scanner.
//
// Usage:
//   driveScanner [-quarantine DIR] [-v] <photo-backup-dir> [more-dirs...]
//
// Walks each given directory and runs the same EXIF parser drive uses
// (imagemeta v0.3.1). Files that cause a panic are listed.
// With -quarantine, panicking files are moved to that directory.
//
// Safe by default: just scans and reports. Add -quarantine only after review.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/evanoberholster/imagemeta"
)

var (
	quarantineDir = flag.String("quarantine", "", "if set, move panicking files here")
	verbose       = flag.Bool("v", false, "verbose: show each file checked")
	showStack     = flag.Bool("stack", false, "show panic stack trace for each bad file")
)

var imageExts = map[string]bool{
	".heic": true, ".heif": true,
	".jpg": true, ".jpeg": true,
	".png":  true,
	".tif":  true, ".tiff": true,
}

func tryDecode(data []byte) (panicked bool, panicMsg string, stack string) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicMsg = fmt.Sprintf("%v", r)
			stack = string(debug.Stack())
		}
	}()
	_, _ = imagemeta.Decode(bytes.NewReader(data))
	return
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: driveScanner [-quarantine DIR] [-v] [-stack] <photo-backup-dir> [more...]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Scans photo backup folders for files that crash drive 4.2.4's EXIF parser.")
		fmt.Fprintln(os.Stderr, "By default, just lists problem files. Use -quarantine to move them.")
		os.Exit(2)
	}

	var bad []string
	checked := 0

	for _, root := range flag.Args() {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			if !imageExts[strings.ToLower(filepath.Ext(path))] {
				return nil
			}
			checked++
			if *verbose {
				fmt.Fprintf(os.Stderr, "[%d] %s\n", checked, path)
			}

			data, rerr := os.ReadFile(path)
			if rerr != nil {
				return nil
			}

			panicked, msg, stack := tryDecode(data)
			if panicked {
				fmt.Printf("BAD: %s\n", path)
				fmt.Printf("     panic: %s\n", msg)
				if *showStack {
					fmt.Println(stack)
				}
				bad = append(bad, path)
			}
			return nil
		})
	}

	fmt.Printf("\n=== summary ===\n")
	fmt.Printf("checked:  %d files\n", checked)
	fmt.Printf("panicked: %d files\n", len(bad))

	if *quarantineDir != "" && len(bad) > 0 {
		if err := os.MkdirAll(*quarantineDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create quarantine dir: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\n=== quarantining %d files to %s ===\n", len(bad), *quarantineDir)
		for _, f := range bad {
			dest := filepath.Join(*quarantineDir, filepath.Base(f))
			if _, statErr := os.Stat(dest); statErr == nil {
				dest = dest + ".dup"
			}
			if err := os.Rename(f, dest); err != nil {
				fmt.Fprintf(os.Stderr, "  ERR move %s: %v\n", f, err)
			} else {
				fmt.Printf("  moved: %s\n", f)
			}
		}
	}

	if len(bad) > 0 {
		os.Exit(1)
	}
}
