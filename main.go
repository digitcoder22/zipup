package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	calledName := filepath.Base(os.Args[0])

	switch calledName {
	case "zipup":
		runZipUp()
	case "unpack":
		runUnpack()
	default:
		if len(os.Args) < 2 {
			fmt.Println("🧰 zipup-suite: invoke via 'zipup' or 'unpack' command symlinks.")
			os.Exit(1)
		}
		
		subCmd := os.Args[1]
		os.Args = os.Args[1:]
		
		if subCmd == "zipup" {
			runZipUp()
		} else if subCmd == "unpack" {
			runUnpack()
		} else {
			fmt.Printf(">:C error: unknown utility action '%s'.\n", subCmd)
			os.Exit(1)
		}
	}
}

func runZipUp() {
	if len(os.Args) < 4 {
		fmt.Println(">:C error: missing required arguments.")
		fmt.Println("Usage: zipup [format] [level] [file/folder] <optional_destination>")
		fmt.Println("example: zipup zip 9 test_files /home/user/backups")
		os.Exit(1)
	}

	format := os.Args[1]
	levelStr := os.Args[2]
	sourceFile := os.Args[3]

	level, err := strconv.Atoi(levelStr)
	if err != nil || level < 1 || level > 9 {
		fmt.Println(">:C error: compression level must be an integer between 1 and 9.")
		os.Exit(1)
	}

	destinationDir := "."
	if len(os.Args) >= 5 {
		destinationDir = os.Args[4]
	}

	sourceFile = filepath.Clean(sourceFile)
	baseName := filepath.Base(sourceFile)
	outputFile := filepath.Join(filepath.Clean(destinationDir), baseName+"."+format)

	fmt.Printf("zipping up: %s\n", sourceFile)
	fmt.Printf("target format: %s (Level %d)\n", format, level)
	fmt.Printf("target output: %s\n", outputFile)

	switch format {
	case "zip":
		err := compressToZip(sourceFile, outputFile, level)
		if err != nil {
			fmt.Printf("aw man :(, compression failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("yay :), zipped up successfully!")

	case "7z":
		_, err := exec.LookPath("7z")
		if err != nil {
			fmt.Println(":O error: '7z' command not found on your system.")
			fmt.Println("Please install 'p7zip' using your package manager to use this format.")
			os.Exit(1)
		}

		levelFlag := fmt.Sprintf("-mx=%d", level)
		cmd := exec.Command("7z", "a", levelFlag, outputFile, sourceFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Printf("aw man :(, 7z zipup failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("yay :), 7z zipping up successfully!")

	case "tar.gz", "tar.xz":
		levelEnv := strconv.Itoa(level)
		var cmd *exec.Cmd
		
		if format == "tar.gz" {
			cmd = exec.Command("tar", "-I", fmt.Sprintf("gzip -%s", levelEnv), "-cf", outputFile, sourceFile)
		} else {
			cmd = exec.Command("tar", "-I", fmt.Sprintf("xz -%s", levelEnv), "-cf", outputFile, sourceFile)
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Printf("aw man :(, tar zipup failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("yay :), tar zipping up successfully!")

	default:
		fmt.Printf(":O error: format '%s' is not implemented yet. supported: 'zip', '7z', 'tar.gz', 'tar.xz'.\n", format)
		os.Exit(1)
	}
}

func runUnpack() {
	if len(os.Args) < 2 {
		fmt.Println(">:C error: missing target archive.")
		fmt.Println("Usage: unpack [archive_file] <optional_destination>")
		fmt.Println("example: unpack backup.zip /home/user/extracted")
		os.Exit(1)
	}

	archivePath := os.Args[1]
	
	destinationDir := "."
	if len(os.Args) >= 3 {
		destinationDir = os.Args[2]
	}

	var ext string
	if strings.HasSuffix(archivePath, ".tar.gz") {
		ext = ".tar.gz"
	} else if strings.HasSuffix(archivePath, ".tar.xz") {
		ext = ".tar.xz"
	} else {
		ext = filepath.Ext(archivePath)
	}

	fmt.Printf("unpacking: %s\n", archivePath)
	fmt.Printf("detected format extension: %s\n", ext)

	switch ext {
	case ".zip":
		err := extractZip(archivePath, destinationDir)
		if err != nil {
			fmt.Printf("aw man :(, extraction failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("yay :), unpacked successfully!")

	case ".7z":
		_, err := exec.LookPath("7z")
		if err != nil {
			fmt.Println(":O error: '7z' tool not found. Please install 'p7zip' to extract .7z files.")
			os.Exit(1)
		}

		destFlag := fmt.Sprintf("-o%s", destinationDir)
		cmd := exec.Command("7z", "x", archivePath, destFlag)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Printf("aw man :(, 7z unpack failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("yay :), 7z unpacked successfully!")

	case ".tar.gz", ".tar.xz":
		os.MkdirAll(destinationDir, 0755)

		cmd := exec.Command("tar", "-xf", archivePath, "-C", destinationDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Printf("aw man :(, tar unpack failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("yay :), tar unpacked successfully!")

	default:
		fmt.Printf(":O error: extension '%s' is not supported. supported: .zip, .7z, .tar.gz, .tar.xz\n", ext)
		os.Exit(1)
	}
}

func compressToZip(sourceDir string, outputFile string, level int) error {
	archiveFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	return filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath

		if info.IsDir() {
			if relPath == "." {
				return nil
			}
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			err = func() error {
				file, err := os.Open(filePath)
				if err != nil {
					return err
				}
				defer file.Close()
				_, err = io.Copy(writer, file)
				return err
			}()
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func extractZip(source string, target string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	err = os.MkdirAll(target, 0755)
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)

		// Fix: Force standard 0755 folder permissions so nested entries can be written safely
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, 0755)
			continue
		}

		// Ensure parent directories exist with correct permissions before writing a file
		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return err
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		// Ensure files are extracted with standard user readable/writable permissions
		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fileReader.Close()
			return err
		}

		_, err = io.Copy(targetFile, fileReader)
		targetFile.Close()
		fileReader.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
