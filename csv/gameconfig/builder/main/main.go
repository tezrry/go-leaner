package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"

	"csv/gameconfig/builder/lexer"
	"csv/gameconfig/builder/reader"
	"csv/gameconfig/infra/metafile"
)

type Source struct {
	tbs map[string]*lexer.Table
}

func main() {
	var cfgName string
	flag.StringVar(&cfgName, "c", "test/config/config.json", "")
	f, err := os.Open(cfgName)
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		panic(err)
	}
	_ = f.Close()

	var outputDir string
	flag.StringVar(&outputDir, "o", "generated", "output directory")
	f, err = os.Open(outputDir)
	if err != nil {
		panic(err)
	}

	genFiles, err := f.Readdir(-1)
	if err != nil {
		panic(err)
	}
	_ = f.Close()
	genMap := make(map[string]os.FileInfo, len(genFiles))
	for _, fi := range genFiles {
		if fi.IsDir() {
			metaName := filepath.Join(outputDir, fi.Name(), "meta.json")
			f, err = os.OpenFile(metaName, os.O_RDWR|os.O_CREATE, 0666)
			genMap[fi.Name()] = fi
		}

	}

	var srcDir string
	flag.StringVar(&srcDir, "s", "test", "source directory")
	f, err = os.Open(srcDir)
	if err != nil {
		panic(err)
	}

	srcFiles, err := f.Readdir(-1)
	if err != nil {
		panic(err)
	}
	_ = f.Close()

	var src Source
	src.tbs = make(map[string]*lexer.Table, len(srcFiles))

	csv := new(reader.CSV)
	for _, fi := range srcFiles {
		if fi.IsDir() || filepath.Ext(fi.Name()) != csv.Suffix() {
			continue
		}

		name := strings.TrimSuffix(fi.Name(), csv.Suffix())
		meta, err2 := metafile.LoadTable(outputDir, name)
		if err2 != nil {
			if err2 == io.EOF {
				meta, err2 = metafile.CreateFile(outputDir, csv.Version(name), name)
			}
		}

		if csv.Version(name) == "" {

		}

		rows, err := csv.Read(name)
		if err != nil {
			panic(err)
		}

		src.tbs[name], err = lexer.InitTable(name, rows)
		if err != nil {
			panic(err)
		}

	}

	_ = f.Close()
}
