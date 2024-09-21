package generator

import (
	"math/rand"
	"os"
)

func createRandomFile(path string, size int64, rng *rand.Rand) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 1024*1024)
	totalWritten := int64(0)
	for totalWritten < size {
		remaining := size - totalWritten
		if remaining < int64(len(buffer)) {
			buffer = buffer[:remaining]
		}
		_, err := rng.Read(buffer)
		if err != nil {
			return err
		}
		n, err := file.Write(buffer)
		if err != nil {
			return err
		}
		totalWritten += int64(n)
	}

	return nil
}

func createFileWithContent(path string, content []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	return err
}
