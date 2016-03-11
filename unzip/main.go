package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	DIRPREFIX = "xxx"
	TMPDIR    = "ttt"
)

var (
	ErrNotMatch = fmt.Errorf("dir name not match case name")
	ErrParse    = fmt.Errorf("parse cfg.json failed")
)

func UnZip(zipFile string) (dirName string, err error) {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return
	}

	for _, k := range r.Reader.File {
		if k.FileInfo().IsDir() {
			name := filepath.Join(TMPDIR, k.Name)
			err := os.MkdirAll(name, 0664)
			if err != nil {
				return "", err
			}
			dirName = k.Name
			continue
		}
		r, err := k.Open()
		if err != nil {
			return "", err
		}
		defer r.Close()
		NewFile, err := os.OpenFile(filepath.Join(TMPDIR, k.Name), os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0664)
		if err != nil {
			return "", err
		}
		io.Copy(NewFile, r)
		NewFile.Close()
	}
	return
}

func CheckDir(dirName string) (err error) {
	data, err := ioutil.ReadFile(filepath.Join(dirName, "cfg.json"))
	if err != nil {
		return err
	}

	ret := make(map[string]interface{})
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return
	}

	name, ok := ret["name"].(string)
	if !ok {
		err = ErrParse
		return
	}
	if name != dirName {
		err = ErrNotMatch
		return
	}
	return moveDir(dirName, filepath.Join(DIRPREFIX, name))
}

func moveDir(src, dst string) (err error) {
	os.RemoveAll(dst)
	cp := exec.Command("cp", "-a", src, dst)
	fmt.Printf("cp -a %s %s\n", src, dst)
	return cp.Run()
}

func main() {
	name, err := UnZip("xxx.zip")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = CheckDir(name)
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		fmt.Printf("unzip and check failed! Err: %v\n", err)
	} else {
		println("passed")
	}

	err = MoveDir("/Upload/demo", "/Users/demo")
	if err != nil {
		fmt.Printf("unzip and check failed! Err: %v\n", err)
	}
}
