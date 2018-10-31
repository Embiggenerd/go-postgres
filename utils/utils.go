package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// RandHex simply creates random hex string
// To use to query session data
func RandHex(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func FileHash(filePath string) (string, error) {
	var fileHash string
	file, err := os.Open(filePath)
	if err != nil {
		return fileHash, err
	}
	defer file.Close()
	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return fileHash, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	fileHash = hex.EncodeToString(hashInBytes)
	return fileHash, nil
}

// func AmendFilename(oldPath, hash string) error {
// 	var b strings.Builder
// 	b.WriteString(oldPath)
// 	b.WriteString(".")
// 	b.WriteString(hash)
// 	err := os.Rename(oldPath, b.String())
// 	return err
// }

// // func GetCSSPath(filePath)
// func Visit(path string, f os.FileInfo, err error) error {
// 	if name := f.Name(); strings.HasPrefix(name, "main") {
// 		dir := filepath.Path(path)
// 		newname := strings.Replace(name, "name_", "name1_", 1)
// 		newpath := filepath.Join(dir, newname)
// 		fmt.Printf("mv %q %q\n", path, newpath)
// 		os.Rename(path, newpath)
// 	}
// 	return nil
// 	return nil
// }
// func Visit(path string, file os.FileInfo, err error) error {

// 	ok := strings.HasPrefix(file.Name(), "mainFloats.css")

// 	if ok {
// 		// var b strings.Builder
// 		// b.writeString(path)
// 		// b.WriteString("mainFloats.css")
// 		// b.WriteString("lalalala")
// 		//
// 		fmt.Println("zz", path, file)
// 		err := os.Rename(path, "assets/lalala")
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 	}
// 	return nil
// }
func copyFileToPublic(path string) error {
	source, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer source.Close()
	_, filename := filepath.Split(path)
	newPath := "./public/" + filename

	destination, err := os.Create(newPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source) // first var shows number of bytes
	if err != nil {
		fmt.Println(err)
		return err
	}

	// err = destination.Sync()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }
	return nil
}

func removeStaleFiles(prefix string) error {
	separated := strings.Split(prefix, ".")
	err := filepath.Walk("./public", func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ok := strings.HasPrefix(file.Name(), separated[0])
		if ok {
			err := os.Remove(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func BustaCache(filename, oldFile string) (string, error) {
	var filenamePlusHash string

	err := removeStaleFiles(filename)
	if err != nil {
		return filenamePlusHash, err
	}

	err = filepath.Walk("./assets", func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ok := strings.HasPrefix(file.Name(), filename)

		if ok {
			err := copyFileToPublic(path)
			if err != nil {
				fmt.Println(err)
				return err
			}
			newPath := "./public/" + file.Name()
			hash, err := FileHash(newPath)
			if err != nil {
				fmt.Println(err)
				return err
			}
			separated := strings.Split(file.Name(), ".")

			var b strings.Builder
			b.WriteString("./")
			b.WriteString(filepath.Dir(newPath))
			b.WriteString("/")
			b.WriteString(separated[0])
			b.WriteString(".")
			b.WriteString(hash)
			b.WriteString(".")
			b.WriteString(separated[1])

			_, filenamePlusHash = filepath.Split(b.String())

			err = os.Rename(newPath, b.String())
			if err != nil {
				fmt.Println(err)
				return err
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return filenamePlusHash, err
	}
	return filenamePlusHash, nil
}
