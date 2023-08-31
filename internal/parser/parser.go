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

// Package parser implements a parser for parsing file names and
// returning metadata about them. This parser is purely designed for
// usage with books, primarily Manga.
package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Contains various regexes used for parsing.
var (
	// volRegexp is used to parse the volume number from a file name.
	//
	// https://regex101.com/r/rPo3zK/1
	volRegexp = regexp.MustCompile(`(?:Vol|Volume|v)(?:\.\s)?(?P<uint>\d+)`)

	// chapRegexp is used to parse the chapter number from a file name.
	//
	// https://regex101.com/r/vSGk9D/1
	chapRegexp = regexp.MustCompile(`(?:Ch|Chapter|c)(?:\.\s)?(?P<uint>\d+)`)

	// absRegexp is used to parse the absolute page number from a file
	// name.
	absRegexp = regexp.MustCompile(`(?P<uint>^\d+)`)

	// pnumRegexp is used to parse the page number from a file name.
	pnumRegexp = regexp.MustCompile(`(?:p(\d+))|[^.0-9]{2}(?:(\d+)\.[a-zA-Z]+$)`)

	// titleRegexp is used to parse the title from a file name.
	//
	//
	// https://regex101.com/r/gjmMaR/1
	titleRegexp = regexp.MustCompile(`[_ ]?(?P<title>[a-zA-Z'_ ]+?)[_ ]?(?:[cvp]|(?:Vol|Volume|Chapter|Ch)\. )\d+`)
)

// Page is a struct that contains metadata about a page in a manga.
type Page struct {
	// Title is the title of the manga. Used as a sorting key for
	// organizing manga.
	Title,

	// FileName is the original file name that contains this page.
	FileName string

	// Volume is the volume number that this series is apart of. When this
	// is unable to be determined, this will be 0.
	Volume,

	// PageNumber is the page number of this page in the manga. This is
	// the page number relative to the chapter. If this is unable to be
	// determined (e.g., contains multiple pages), this will be ^uint(0).
	PageNumber,

	// AbsolutePageNumber is the page number of this page in the manga.
	// This is the page number relative to the entire series. If this is
	// unable to be determined (e.g., contains multiple pages), this will
	// be ^uint(0).
	AbsolutePageNumber,

	// Chapter is the chapter number that this series is apart of. When
	// this is unable to be determined, this will be 0.
	Chapter uint
}

// String implements fmt.Stringer for Page.
func (p Page) String() string {
	return fmt.Sprintf(
		"%s: abs=%d, vol=%d, chap=%d, page=%d (source: %s)",
		p.Title, p.AbsolutePageNumber, p.Volume, p.Chapter, p.PageNumber, p.FileName,
	)
}

// Parse parses the provided file name and returns a Page.
func Parse(name string) (Page, error) {
	p := parse(name)
	if p.Title == "" {
		return p, fmt.Errorf("unable to parse title from filename")
	}

	return p, nil
}

// parseUintFromRegex returns a uint of the result of the capture group
// named 'uint' from the provided regexp on the provided name. Returns
// the value of defaultValue if no matches are found.
func parseUintFromRegex(name string, re *regexp.Regexp, defaultValue uint) uint {
	match := re.FindStringSubmatch(name)
	if len(match) == 0 {
		return defaultValue
	}

	for i, name := range re.SubexpNames() {
		if name == "uint" {
			r, err := strconv.ParseUint(match[i], 10, 64)
			if err != nil {
				panic(fmt.Errorf("parseUintFromRegex: unable to parse uint from string: %w", err))
			}
			return uint(r)
		}
	}

	panic("parseUintFromRegex: unable to find capture group 'int' in regexp provided")
}

// parseString returns a string from the result of the first capture
// group found in the regex string. Returns an empty string if no
// matches are found.
func parseString(name string, re *regexp.Regexp) string {
	match := re.FindStringSubmatch(name)
	if len(match) == 0 || len(match) < 2 {
		return ""
	}

	return match[1]
}

// parse is the underlying implementation for Parse
func parse(name string) Page {
	vol := parseUintFromRegex(name, volRegexp, 0)
	chap := parseUintFromRegex(name, chapRegexp, 0)
	abs := parseUintFromRegex(name, absRegexp, ^uint(0))
	title := strings.ReplaceAll(parseString(name, titleRegexp), "_", " ")

	pnum := ^uint(0)
	for i, v := range pnumRegexp.FindStringSubmatch(name) {
		if v == "" || i == 0 {
			continue
		}

		pnumU64, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			panic(fmt.Errorf("parse: unable to parse page number from string: %v", err))
		}
		pnum = uint(pnumU64)
	}

	return Page{
		FileName:           name,
		Title:              title,
		AbsolutePageNumber: abs,
		PageNumber:         uint(pnum),
		Volume:             vol,
		Chapter:            chap,
	}
}
