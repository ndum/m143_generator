package generator

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/ndum/m143_generator/utils"
)

type Options struct {
	FilesCount      int
	DirsCount       int
	TotalSize       int64
	BaseDir         string
	Levels          int
	DirNamePattern  string
	FileNamePattern string
	FileExtension   string
	Duplicates      int
	Seed            int64
}

func GenerateDummyData(opts Options) error {
	var rng *rand.Rand
	if opts.Seed != 0 {
		rng = rand.New(rand.NewSource(opts.Seed))
	} else {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	dirPaths, err := createDirectories(opts.BaseDir, opts.DirsCount, opts.Levels, opts.DirNamePattern, rng, opts.Seed)
	if err != nil {
		return err
	}

	fileSizes := utils.GenerateFileSizes(opts.FilesCount, opts.TotalSize, opts.Duplicates, rng)

	var duplicateContent []byte
	if opts.Duplicates > 0 {
		duplicateSize := fileSizes[0]

		duplicateContent = make([]byte, duplicateSize)
		_, err := rng.Read(duplicateContent)
		if err != nil {
			return fmt.Errorf("error generating content for duplicate files: %v", err)
		}
	}

	var fixedTime time.Time
	if opts.Seed != 0 {
		fixedTime = time.Unix(opts.Seed, 0)
	} else {
		fixedTime = time.Now()
	}
	fixedTimestamp := fixedTime.Unix()
	fixedDate := fixedTime.Format("2006-01-02")
	fixedTimeStr := fixedTime.Format("15-04-05")

	for i := 0; i < opts.FilesCount; i++ {
		dir := dirPaths[rng.Intn(len(dirPaths))]

		placeholderValues := map[string]interface{}{
			"index":     i,
			"random":    rng.Intn(100000),
			"timestamp": fixedTimestamp,
			"date":      fixedDate,
			"time":      fixedTimeStr,
			"uuid":      GenerateDeterministicUUID(rng),
			"ext":       opts.FileExtension,
		}

		fileName := ReplacePlaceholders(opts.FileNamePattern, placeholderValues, rng)

		filePath := filepath.Join(dir, fileName)

		if i < opts.Duplicates {
			err := createFileWithContent(filePath, duplicateContent)
			if err != nil {
				return fmt.Errorf("error creating duplicate file %s: %v", filePath, err)
			}
		} else {
			err := createRandomFile(filePath, fileSizes[i], rng)
			if err != nil {
				return fmt.Errorf("error creating file %s: %v", filePath, err)
			}
		}
	}

	return nil
}

func createDirectories(baseDir string, dirsPerLevel, levels int, dirNamePattern string, rng *rand.Rand, seed int64) ([]string, error) {
	var dirPaths []string

	currentDirs := []string{baseDir}

	var fixedTime time.Time
	if seed != 0 {
		fixedTime = time.Unix(seed, 0)
	} else {
		fixedTime = time.Now()
	}
	fixedTimestamp := fixedTime.Unix()
	fixedDate := fixedTime.Format("2006-01-02")
	fixedTimeStr := fixedTime.Format("15-04-05")

	for level := 0; level < levels; level++ {
		var newDirs []string
		for _, parentDir := range currentDirs {
			for i := 0; i < dirsPerLevel; i++ {
				placeholderValues := map[string]interface{}{
					"level":     level,
					"index":     i,
					"random":    rng.Intn(100000),
					"timestamp": fixedTimestamp,
					"date":      fixedDate,
					"time":      fixedTimeStr,
					"uuid":      GenerateDeterministicUUID(rng),
				}
				dirName := ReplacePlaceholders(dirNamePattern, placeholderValues, rng)
				dirPath := filepath.Join(parentDir, dirName)
				err := os.MkdirAll(dirPath, os.ModePerm)
				if err != nil {
					return nil, fmt.Errorf("error creating directory %s: %v", dirPath, err)
				}
				newDirs = append(newDirs, dirPath)
			}
		}
		dirPaths = append(dirPaths, newDirs...)
		currentDirs = newDirs
	}

	return dirPaths, nil
}

func GenerateDeterministicUUID(rng *rand.Rand) string {
	var uuid [16]byte
	rng.Read(uuid[:])

	uuid[6] = (uuid[6] & 0x0F) | 0x40
	uuid[8] = (uuid[8] & 0x3F) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16],
	)
}
