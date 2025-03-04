package reader

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"os"

	"csv/gameconfig/infra/metafile"
	"go-learner/slice"
)

type CSV struct{}

func (inst *CSV) Read(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = file.Close()
	}()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("empty file %s", fileName)
	}

	return rows, nil
}

func (inst *CSV) Version(fileName string) string {
	//cmd := exec.Command("git", "rev-parse", "HEAD:", fileName)
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//if err := cmd.Run(); err != nil {
	//	return true
	//}
	//
	//// If output is empty, the file is not modified
	//status := strings.TrimSpace(out.String())
	//return status != ""
	b := make([]byte, metafile.VersionLen)
	_, _ = rand.Read(b)
	return slice.ByteSlice2String(b)
}

func (inst *CSV) Suffix() string {
	return ".csv"
}
