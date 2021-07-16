// +build windows

package higgs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

var tmpDir string

func touch(path, content string) {
	path = filepath.FromSlash(path)
	dir := filepath.Dir(path)
	if dir != "" {
		os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
	}
	ioutil.WriteFile(filepath.Join(tmpDir, path), []byte(content), 0644)
}

func hideFile(path string) error {
	attrs, utf16PtrPath, err := getFileAttrs(filepath.Join(tmpDir, path))
	if err != nil {
		return fmt.Errorf("something went wrong getting file attributes: \"%s\"\nError: \"%s\"", path, err)
	}
	if attrs&syscall.FILE_ATTRIBUTE_HIDDEN > 0 {
		return nil
	}
	return syscall.SetFileAttributes(utf16PtrPath, syscall.FILE_ATTRIBUTE_HIDDEN)
}

func isFileHidden(path string) (bool, error) {
	attrs, _, err := getFileAttrs(path)
	if err != nil {
		return false, fmt.Errorf("something went wrong getting file attributes: \"%s\"", err)
	}
	if attrs&syscall.FILE_ATTRIBUTE_HIDDEN > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func TestMain(t *testing.M) {
	tmpDir, _ = ioutil.TempDir("", "higgs*")
	touch("a", "a")
	touch("b", "b")
	touch("c/c.a", "c.a")
	touch("d/d.a", "d.a")
	err := hideFile("b")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = hideFile("d")
	if err != nil {
		fmt.Println(err)
		return
	}

	code := t.Run()
	defer func() {
		os.RemoveAll(tmpDir)
		os.Exit(code)
	}()
}

func TestIsHiddenWhenNotHidden(t *testing.T) {
	hidden, err := IsHidden(filepath.Join(tmpDir, "a"))

	if err != nil {
		t.Errorf(`Error: "%s"`, err)
	}
	if hidden == true {
		t.Errorf("wrong output, file is not hidden but the output says otherwise")
	}
}

func TestIsHiddenWhenHidden(t *testing.T) {
	hidden, err := IsHidden(filepath.Join(tmpDir, "b"))

	if err != nil {
		t.Errorf(`Error: "%s"`, err)
	}
	if hidden == false {
		t.Errorf("wrong output, file is hidden but the output says otherwise")
	}
}

func TestIsHiddenWhenNotExists(t *testing.T) {
	hidden, err := IsHidden(filepath.Join(tmpDir, "notexists"))

	if err == nil {
		t.Errorf("no error")
	}
	if hidden == true {
		t.Errorf("wrong output")
	}
}

func TestHideHidesWhenAlreadyHidden(t *testing.T) {
	path := filepath.Join(tmpDir, "b")
	_, err := Hide(path)

	if err != nil {
		t.Errorf("error: \"%s\"", err)
	}

	hidden, err := isFileHidden(path)

	if err != nil {
		t.Errorf("error: \"%s\"", err)
	}
	if !hidden {
		t.Errorf("file should be hidden but it is not: \"%s\"", err)
	}
}

func TestUnhideNotHidesWhenAlreadyNotHidden(t *testing.T) {
	path := filepath.Join(tmpDir, "a")
	_, err := Unhide(path)

	if err != nil {
		t.Errorf("error: \"%s\"", err)
	}

	hidden, err := isFileHidden(path)

	if err != nil {
		t.Errorf("error: \"%s\"", err)
	}
	if hidden {
		t.Errorf("file should not be hidden but it is: \"%s\"", path)
	}
}

func TestHideWhenNotExists(t *testing.T) {
	_, err := Hide(filepath.Join(tmpDir, "notexists"))

	if err == nil {
		t.Errorf("error: \"%s\"", err)
	}
}

func TestHideHidesFile(t *testing.T) {
	path := filepath.Join(tmpDir, "a")
	_, err := Hide(path)

	if err != nil {
		t.Errorf("error: \"%s\"", err)
	}

	hidden, err := isFileHidden(path)

	if err != nil {
		t.Errorf("error: \"%s\"", err)
	}
	if !hidden {
		t.Errorf("file should be hidden but it is not: \"%s\"", path)
	}
}

func TestUnhideUnhidesFile(t *testing.T) {
	path := filepath.Join(tmpDir, "b")
	_, err := Unhide(path)

	if err != nil {
		t.Errorf("error: \"%s\"", err)
	}

	hidden, err := isFileHidden(path)

	if err != nil {
		t.Errorf("error: \"%s\"", err)
	}
	if hidden {
		t.Errorf("file should not be hidden but it is: \"%s\"", path)
	}
}

func TestHideHidesDirectory(t *testing.T) {
	path := filepath.Join(tmpDir, "c")
	_, err := Hide(path)

	if err != nil {
		t.Errorf("error: \"%s\"", err)
	}

	hidden, err := isFileHidden(path)

	if err != nil {
		t.Errorf("error: \"%s\"", err)
	}
	if !hidden {
		t.Errorf("directory should be hidden but it is not: \"%s\"", path)
	}
}

func TestHideUnhidesDirectory(t *testing.T) {
	path := filepath.Join(tmpDir, "d")
	_, err := Unhide(path)

	if err != nil {
		t.Errorf("error: \"%s\"", err)
	}

	hidden, err := isFileHidden(path)

	if err != nil {
		t.Errorf("error: \"%s\"", err)
	}
	if hidden {
		t.Errorf("directory should not be hidden but it is: \"%s\"", path)
	}
}

func TestNewFileHideHiddenHidesFile(t *testing.T) {
	path := filepath.Join(tmpDir, "b")
	_, err := NewFileHide(path, false).Hide()

	if err != nil {
		t.Errorf("error: \"%s\"", err)
	}

	hidden, err := isFileHidden(path)

	if err != nil {
		t.Errorf("error: \"%s\"", err)
	}
	if !hidden {
		t.Errorf("file should be hidden but it is not: \"%s\"", path)
	}
}
