
# m143_generator

**m143_generator** is a flexible and powerful tool written in Go for generating dummy files and directories. It's ideal for simulating file system structures, testing backup and restore scenarios, or educational purposes. This tool is specifically designed for use in **Module 143** of IT vocational education.

## Features

- **Customizable Directory and File Structures**: Define patterns for directory and file names using placeholders.
- **Multiple Placeholder Support**: Use placeholders like `{index}`, `{level}`, `{random}`, `{timestamp}`, `{date}`, `{time}`, `{uuid}`, `{ext}`, and `{randstr:N}`.
- **Nested Directories**: Generate directories with multiple levels of subdirectories.
- **Random File Content**: Creates files with random data of specified total size.
- **Duplicate Files Generation**: Create a specified number of files with identical content and hash values.
- **Seeded Randomness for Reproducibility**: Use a seed value to reproduce the exact same file and directory structure with identical content.
- **Settings Persistence**: Save and reuse settings for consistent data generation.
- **Cross-Platform Compatibility**: Works on Windows and Linux.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
    - [Command-Line Arguments](#command-line-arguments)
    - [Placeholders](#placeholders)
    - [Examples](#examples)
- [macOS Users](#macos-users)
- [Settings File](#settings-file)
- [Contributing](#contributing)
- [License](#license)

## Installation

### Prerequisites

- Go 1.23 or higher

### Steps

1. **Clone the Repository**

   ```bash
   git clone https://github.com/ndum/m143_generator.git
   cd m143_generator
   ```

2. **Install Dependencies**

   ```bash
   go mod tidy
   ```


3. **Build the Application**

   ```bash
   go build -o m143_generator
   ```

## Usage

Run the application using the built binary and provide the desired command-line arguments.

**Note**: The `--dir` parameter is required to specify the output directory.

```bash
./m143_generator --dir /path/to/output [options]
```

### Command-Line Arguments

- `-b, --dir`: Base directory for file and directory creation (**required**)
- `-f, --files`: Number of files to create (default: `10`)
- `-d, --dirs`: Number of directories per level (default: `5`)
- `-l, --levels`: Number of subdirectory levels (default: `1`)
- `-s, --size`: Total size of all files (e.g., `500mb`, default: `500mb`)
- `-D, --dir-name-pattern`: Pattern for directory names (default: `dir_{level}_{index}`)
- `-F, --file-name-pattern`: Pattern for file names (default: `file_{index}.{ext}`)
- `-e, --file-extension`: File extension for generated files (default: `dat`)
- `-u, --duplicates`: Number of duplicate files to create (default: `0`)
- `-S, --seed`: Seed value for random number generator (default: current time)
- `-h, --help`: Display help information

### Placeholders

You can use the following placeholders in your directory and file name patterns:

- `{level}`: Current level in the directory tree (zero-based).
- `{index}`: Index of the directory or file (running number starting from 0).
- `{random}`: Random number (derived from the seed).
- `{timestamp}`: Fixed Unix timestamp based on the seed.
- `{date}`: Fixed date based on the seed in format `YYYY-MM-DD`.
- `{time}`: Fixed time based on the seed in format `HH-MM-SS`.
- `{uuid}`: Deterministic UUID based on the seed.
- `{ext}`: File extension specified with `--file-extension`.
- `{randstr:N}`: Random string of letters with length `N` (derived from the seed).

**Note**: To ensure that duplicate files have the same names, avoid using time-based or random placeholders like `{timestamp}`, `{date}`, `{time}`, `{uuid}`, or `{random}` in your `--file-name-pattern` when generating duplicates. Instead, use deterministic placeholders like `{index}` or `{randstr:N}` with a fixed seed.

### Examples

#### Example 1: Basic Usage

Generate 10 files in 5 directories with default settings:

```bash
./m143_generator --dir /path/to/output
```

This will create:

- **Directories**: 5 directories named `dir_0_0` to `dir_0_4` in the specified output directory.
- **Files**: 10 files named `file_0.dat` to `file_9.dat`, distributed randomly among the directories.

#### Example 2: Create Duplicate Files

Generate files with duplicates for testing deduplication:

```bash
./m143_generator   --dir /path/to/output   -f 100   -d 2   -l 3   -s 500mb   -u 5   -D "dir_{level}_{index}"   -F "file_{index}.{ext}"   -e "dat"
```

Explanation:

- **Files**: Creates 100 files, with 5 duplicates.
- **Duplicates**: The first 5 files will have identical content and hash values.
- **Use Case**: Useful for testing deduplication features in backup software.

#### Example 3: Use Seed for Reproducibility

Generate files with a seed value to ensure reproducibility:

```bash
./m143_generator   --dir /path/to/output   -f 50   -d 3   -l 2   -s 200mb   -S 12345   -D "dir_{level}_{index}"   -F "file_{index}.{ext}"   -e "txt"
```

Explanation:

- **Seed**: Uses seed `12345` for the random number generator.
- **Reproducibility**: Running this command multiple times with the same seed will produce identical files, including their content and hash values.
- **Use Case**: Useful for consistent testing and educational purposes.

#### Example 4: Combine Duplicates and Seed

Generate duplicate files with reproducible content:

```bash
./m143_generator   --dir /path/to/output   -f 20   -d 2   -l 2   -s 100mb   -u 3   -S 67890   -D "folder_{index}"   -F "data_{index}.{ext}"   -e "bak"
```

Explanation:

- **Files**: Creates 20 files, with 3 duplicates.
- **Seed**: Uses seed `67890` for consistent results.
- **Duplicate Files**: The first 3 files will be identical across runs with the same seed.

## macOS Users

**Note**: Pre-built binaries for macOS are not provided due to code signing requirements. If you're using macOS, please build the application from source using the following steps:

1. **Install Go**: Ensure you have Go installed on your system. You can download it from the [official website](https://golang.org/dl/).

2. **Clone the Repository**:

   ```bash
   git clone https://github.com/yourusername/m143_generator.git
   cd m143_generator
   ```

3. **Install Dependencies**:

   ```bash
   go mod tidy
   ```

4. **Build the Application**:

   ```bash
   go build -o m143_generator
   ```

You can now run the application using:

```bash
./m143_generator --dir /path/to/output [options]
```

## Settings File

After the first run, a `settings.json` file is saved in the base directory. This file contains all the parameters used during that run, including any seed and duplicates specified.

### Reusing Settings

On subsequent runs, the application checks for the existence of `settings.json` in the base directory. If found, it prompts:

```
A settings file was found in the target directory.
Do you want to recreate files based on these settings? [y/N]:
```

- **`y` or `yes`**: The application will use the settings from `settings.json`, recreating the files exactly as before.
- **`n` or `no`** (or any other input): The application will use the command-line arguments provided.

This feature is especially useful for:

- **Backup and Restore Scenarios**: Recreate the exact file structure and content for testing backups, including duplicate files.
- **Educational Purposes**: Consistently reproduce file structures for classroom demonstrations.

## Contributing

Contributions are welcome! Please follow these steps:

1. **Fork the Repository**

   Click the "Fork" button at the top right of the repository page.

2. **Clone Your Fork**

   ```bash
   git clone https://github.com/ndum/m143_generator.git
   cd m143_generator
   ```

3. **Create a Feature Branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Commit Your Changes**

   ```bash
   git commit -am 'Add new feature'
   ```


**Important Note for macOS Users**: Pre-built binaries for macOS are not available due to Apple's code signing requirements. You will need to manually build the application by following the build steps under the macOS Users section.
