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

func AmendFilename(oldPath, hash string) error {
	var b strings.Builder
	b.WriteString(oldPath)
	b.WriteString(".")
	b.WriteString(hash)
	err := os.Rename(oldPath, b.String())
	return err
}

// func GetCSSPath(filePath)
func Visit(path string, f os.FileInfo, err error) error {
	if name := f.Name(); strings.HasPrefix(name, "main") {
		dir := filepath.Path(path)
		newname := strings.Replace(name, "name_", "name1_", 1)
		newpath := filepath.Join(dir, newname)
		fmt.Printf("mv %q %q\n", path, newpath)
		os.Rename(path, newpath)
	}
	return nil
	return nil
}
