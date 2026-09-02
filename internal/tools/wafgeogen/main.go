// The wafgeogen tool builds the embedded IP-to-country and
// IP-to-origin-ASN tables the WAF inspection engine uses for
// GeoMatchStatement and AsnMatchStatement. Inputs are RIR
// delegated-extended statistics files (country table) and one RouteViews
// MRT RIB dump (ASN table); the raw inputs are downloaded at generation
// time and are never committed. The output is a compact derived blob
// committed to the repository and embedded with go:embed.
//
// Usage:
//
//	wafgeogen -delegated <dir> -rib <file.bz2> -date YYYYMMDD -out <file>
package main

import (
	"flag"
	"fmt"
	"net/netip"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	delegatedDir := flag.String("delegated", "", "directory of delegated-*-extended files")
	ribFile := flag.String("rib", "", "RouteViews MRT RIB dump (bz2)")
	date := flag.String("date", "", "generation date as YYYYMMDD")
	out := flag.String("out", "", "output blob path")
	flag.Parse()
	if *delegatedDir == "" || *ribFile == "" || *out == "" {
		flag.Usage()
		os.Exit(2)
	}
	if err := run(*delegatedDir, *ribFile, *date, *out); err != nil {
		fmt.Fprintln(os.Stderr, "wafgeogen:", err)
		os.Exit(1)
	}
}

func run(delegatedDir, ribFile, date, outPath string) error {
	// Country table: map every registry allocation to its country
	// code's index in the string table.
	entries, err := os.ReadDir(delegatedDir)
	if err != nil {
		return err
	}
	var records []delegatedRecord
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, "delegated-") || !strings.Contains(name, "-extended") {
			continue
		}
		f, err := os.Open(filepath.Join(delegatedDir, name))
		if err != nil {
			return err
		}
		records, err = parseDelegated(f, records)
		f.Close()
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	countryIndex := map[string]uint32{"": 0}
	var countries []string
	trie := &trieNode{}
	for _, rec := range records {
		idx, ok := countryIndex[rec.CountryCode]
		if !ok {
			countries = append(countries, rec.CountryCode)
			idx = uint32(len(countries))
			countryIndex[rec.CountryCode] = idx
		}
		raw, bits := prefix128(rec.Prefix)
		trie.insert(raw, bits, idx)
	}
	countryTable := encodeTable(trie.flatten())

	// ASN table: origin ASN of every announced prefix.
	rib, err := os.Open(ribFile)
	if err != nil {
		return err
	}
	defer rib.Close()
	asnByPrefix := map[netip.Prefix]uint32{}
	if err := parseMRT(rib, asnByPrefix); err != nil {
		return err
	}
	asnTrie := &trieNode{}
	for prefix, asn := range asnByPrefix {
		raw, bits := prefix128(prefix)
		asnTrie.insert(raw, bits, asn)
	}
	asnTable := encodeTable(asnTrie.flatten())

	blob := buildBlob(date, countries, countryTable, asnTable)
	if err := os.WriteFile(outPath, blob, 0o644); err != nil {
		return err
	}
	fmt.Printf("wafgeogen: wrote %s (%d bytes, %d country codes, %d asn prefixes)\n",
		outPath, len(blob), len(countries), len(asnByPrefix))
	return nil
}
