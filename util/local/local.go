package local

import (
	"path/filepath"
	"os"
	"bufio"
)

func GetHostFromLocalFile(name string) (result []string, err error) {
	filePath := filepath.Join(os.Getenv("HOME"), ".pssh_groups", name)
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return result, err
	}
	return result, err
}
