package generator

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

func ReplacePlaceholders(pattern string, values map[string]interface{}, rng *rand.Rand) string {
	result := pattern

	randstrRegex := regexp.MustCompile(`{randstr:(\d+)}`)
	result = randstrRegex.ReplaceAllStringFunc(result, func(match string) string {
		matches := randstrRegex.FindStringSubmatch(match)
		if len(matches) == 2 {
			length, err := strconv.Atoi(matches[1])
			if err == nil {
				return randomString(length, rng)
			}
		}
		return match
	})

	for key, value := range values {
		placeholder := fmt.Sprintf("{%s}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}

	return result
}

func randomString(length int, rng *rand.Rand) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, length)
	for i := range s {
		s[i] = letters[rng.Intn(len(letters))]
	}
	return string(s)
}
