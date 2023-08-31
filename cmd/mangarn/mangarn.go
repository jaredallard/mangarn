// Copyright (C) 2023 Jared Allard <jaredallard@users.noreply.github.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jaredallard/mangarn/internal/parser"
)

// main runs the mangarn CLI.
func main() {
	ctx := context.Background()

	var exitCode uint = 0
	defer func() {
		os.Exit(int(exitCode))
	}()

	// fatal is a helper to set exitCode to 1 and log an error. Callers
	// must still call return after calling fatal.
	fatal := func(ctx context.Context, msg string, args ...any) {
		slog.ErrorContext(ctx, msg, args...)
		exitCode = 1
	}

	// read all the files in the directory and parse them.
	files, err := os.ReadDir(".")
	if err != nil {
		fatal(ctx, "failed to read directory", "error", err)
		return
	}

	pages := make([]parser.Page, 0, len(files))
	for _, file := range files {
		if file.IsDir() || file.Name() == ".DS_Store" {
			continue
		}

		p, err := parser.Parse(file.Name())
		if err != nil {
			fatal(ctx, "failed to parse file", "error", err)
			return
		}

		pages = append(pages, p)
	}
	if len(pages) == 0 {
		fatal(ctx, "no pages found")
		return
	}

	title := pages[0].Title
	volMap := make(map[uint][]parser.Page)

	for _, p := range pages {
		if p.Title != title {
			fatal(ctx, "title mismatch", "expected", title, "got", p.Title)
			return
		}

		if p.Volume == ^uint(0) {
			fatal(ctx, "unable to determine volume", "page", p)
			return
		}

		if p.PageNumber == ^uint(0) {
			fatal(ctx, "unable to determine page number", "page", p)
			return
		}

		// If we failed to determine the chapter number (0 is OK), then we
		// should fail.
		if p.Chapter == ^uint(0) {
			fatal(ctx, "unable to determine chapter", "page", p)
			return
		}

		volMap[p.Volume] = append(volMap[p.Volume], p)
	}

	// create cbz archives for each volume in a given title
	for vol, pages := range volMap {
		chapMap := make(map[uint][]*parser.Page)

		// sort the pages by chapter
		for i := range pages {
			p := &pages[i]
			chapMap[p.Chapter] = append(chapMap[p.Chapter], p)
		}

		// write the pages to a cbz archive per chapter
		for chap, pages := range chapMap {
			outputName := fmt.Sprintf("%s Vol.%d", title, vol)
			if chap != 0 {
				outputName += fmt.Sprintf(" Ch.%d", chap)
			}
			outputName += ".cbz"

			if err := writePages(ctx, outputName, pages); err != nil {
				fatal(ctx, "failed to write pages", "error", err)
				return
			}
		}
	}
}

// writePages writes the provided pages to a cbz archive at the provided
// output name. Pages are saved as a cbz archive with the pages being
// the PageNumber of the page and the extension of the FileName on the
// page.
func writePages(ctx context.Context, outputName string, pages []*parser.Page) error {
	if err := os.MkdirAll("output", 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	out, err := os.Create(filepath.Join("output", outputName))
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	zw := zip.NewWriter(out)
	defer zw.Close()

	for _, p := range pages {
		savePath := fmt.Sprintf("Image %d%s", p.PageNumber, filepath.Ext(p.FileName))

		slog.InfoContext(ctx, "saving file", "page", p, "path", savePath)
		if err := writeFileToZip(ctx, p.FileName, savePath, zw); err != nil {
			return fmt.Errorf("failed to write file to zip: %w", err)
		}
	}

	return nil
}

// writeFileToZip writes the provided file (src) to the provided
// (zw) *zip.Writer at the provided destination (dest).
func writeFileToZip(ctx context.Context, src string, dest string, zw *zip.Writer) error {
	f, err := zw.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create file in cbz: %w", err)
	}

	srcF, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open file to copy to cbz: %w", err)
	}
	defer srcF.Close()

	if _, err := io.Copy(f, srcF); err != nil {
		return fmt.Errorf("failed to copy file to cbz: %w", err)
	}

	return nil
}
