package utils

import (
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Settings struct {
	FilesCount      int    `json:"files_count"`
	DirsCount       int    `json:"dirs_count"`
	TotalSizeStr    string `json:"total_size_str"`
	BaseDir         string `json:"base_dir"`
	Levels          int    `json:"levels"`
	DirNamePattern  string `json:"dir_name_pattern"`
	FileNamePattern string `json:"file_name_pattern"`
	FileExtension   string `json:"file_extension"`
	Duplicates      int    `json:"duplicates"`
	Seed            int64  `json:"seed"`
}

func ParseSize(sizeStr string) (int64, error) {
	sizeStr = strings.TrimSpace(strings.ToLower(sizeStr))
	var multiplier int64 = 1

	switch {
	case strings.HasSuffix(sizeStr, "kb"):
		multiplier = 1024
		sizeStr = strings.TrimSuffix(sizeStr, "kb")
	case strings.HasSuffix(sizeStr, "mb"):
		multiplier = 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "mb")
	case strings.HasSuffix(sizeStr, "gb"):
		multiplier = 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "gb")
	case strings.HasSuffix(sizeStr, "b"):
		sizeStr = strings.TrimSuffix(sizeStr, "b")
	default:
		return 0, errors.New("unknown size suffix")
	}

	value, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return 0, err
	}

	return int64(value * float64(multiplier)), nil
}

func GenerateFileSizes(count int, totalSize int64, duplicates int, rng *rand.Rand) []int64 {
	sizes := make([]int64, count)
	var sum int64 = 0

	if duplicates > 0 {
		duplicateSize := totalSize / int64(count)
		for i := 0; i < duplicates; i++ {
			sizes[i] = duplicateSize
			sum += duplicateSize
		}
	}

	for i := duplicates; i < count-1; i++ {
		remaining := totalSize - sum - int64(count-i-1)*1024
		maxSize := int64(10 * 1024 * 1024)
		if remaining < maxSize {
			maxSize = remaining
		}
		size := rng.Int63n(maxSize-1024) + 1024
		sizes[i] = size
		sum += size
	}

	sizes[count-1] = totalSize - sum

	return sizes
}

func LoadSettings(path string) (Settings, error) {
	var settings Settings
	file, err := os.Open(path)
	if err != nil {
		return settings, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&settings)
	return settings, err
}

func SaveSettings(path string, settings Settings) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(settings)
}
