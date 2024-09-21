package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ndum/m143_generator/generator"
	"github.com/ndum/m143_generator/utils"
	"github.com/spf13/cobra"
)

var (
	filesCount      int
	dirsCount       int
	totalSizeStr    string
	baseDir         string
	levels          int
	dirNamePattern  string
	fileNamePattern string
	fileExtension   string
	duplicates      int
	seed            int64
)

var rootCmd = &cobra.Command{
	Use:   "m143_generator",
	Short: "A tool to generate dummy files and directories for M143",
	Long: `m143_generator is a flexible and powerful tool written in Go for generating dummy files and directories.
It's ideal for simulating file system structures, testing backup and restore scenarios, or educational purposes.

Author: Nicolas Dumermuth (nd@nidum.org)
`,
	Run: func(cmd *cobra.Command, args []string) {
		if baseDir == "" {
			fmt.Println("Error: The output directory must be specified using the --dir flag.")
			cmd.Help()
			os.Exit(1)
		}

		settingsFilePath := filepath.Join(baseDir, "settings.json")

		if _, err := os.Stat(settingsFilePath); err == nil {
			fmt.Printf("A settings file was found in the target directory.\n")
			fmt.Printf("Do you want to recreate files based on these settings? [y/N]: ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))
			if input == "y" || input == "yes" {
				settings, err := utils.LoadSettings(settingsFilePath)
				if err != nil {
					fmt.Printf("Error loading settings: %v\n", err)
					os.Exit(1)
				}
				filesCount = settings.FilesCount
				dirsCount = settings.DirsCount
				totalSizeStr = settings.TotalSizeStr
				baseDir = settings.BaseDir
				levels = settings.Levels
				dirNamePattern = settings.DirNamePattern
				fileNamePattern = settings.FileNamePattern
				fileExtension = settings.FileExtension
				duplicates = settings.Duplicates
				seed = settings.Seed
				fmt.Println("Creating files based on saved settings...")
			} else {
				fmt.Println("Proceeding with provided command-line arguments.")
			}
		}

		totalSize, err := utils.ParseSize(totalSizeStr)
		if err != nil {
			fmt.Printf("Invalid size format: %v\n", err)
			os.Exit(1)
		}

		err = generator.GenerateDummyData(generator.Options{
			FilesCount:      filesCount,
			DirsCount:       dirsCount,
			TotalSize:       totalSize,
			BaseDir:         baseDir,
			Levels:          levels,
			DirNamePattern:  dirNamePattern,
			FileNamePattern: fileNamePattern,
			FileExtension:   fileExtension,
			Duplicates:      duplicates,
			Seed:            seed,
		})
		if err != nil {
			fmt.Printf("Error generating dummy data: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Dummy files and directories were successfully created.")

		err = utils.SaveSettings(settingsFilePath, utils.Settings{
			FilesCount:      filesCount,
			DirsCount:       dirsCount,
			TotalSizeStr:    totalSizeStr,
			BaseDir:         baseDir,
			Levels:          levels,
			DirNamePattern:  dirNamePattern,
			FileNamePattern: fileNamePattern,
			FileExtension:   fileExtension,
			Duplicates:      duplicates,
			Seed:            seed,
		})
		if err != nil {
			fmt.Printf("Error saving settings: %v\n", err)
		}
	},
}

func Execute() {
	rootCmd.Flags().IntVarP(&filesCount, "files", "f", 10, "Number of files to create")
	rootCmd.Flags().IntVarP(&dirsCount, "dirs", "d", 5, "Number of directories per level")
	rootCmd.Flags().StringVarP(&totalSizeStr, "size", "s", "500mb", "Total size of all files (e.g., 500mb)")
	rootCmd.Flags().StringVarP(&baseDir, "dir", "b", "", "Base directory for file and directory creation (required)")
	rootCmd.MarkFlagRequired("dir")
	rootCmd.Flags().IntVarP(&levels, "levels", "l", 1, "Number of subdirectory levels")
	rootCmd.Flags().StringVarP(&dirNamePattern, "dir-name-pattern", "D", "dir_{level}_{index}", "Pattern for directory names")
	rootCmd.Flags().StringVarP(&fileNamePattern, "file-name-pattern", "F", "file_{index}.{ext}", "Pattern for file names")
	rootCmd.Flags().StringVarP(&fileExtension, "file-extension", "e", "dat", "File extension for generated files")
	rootCmd.Flags().IntVarP(&duplicates, "duplicates", "u", 0, "Number of duplicate files to create")
	rootCmd.Flags().Int64VarP(&seed, "seed", "S", 0, "Seed value for random number generator")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
