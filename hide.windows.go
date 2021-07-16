// +build windows

package higgs

import (
	"fmt"
	"os"
	"syscall"
)

func getFileAttrs(path string) (uint32, *uint16, error) {
	utf16PtrPath, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, nil, fmt.Errorf("something went wrong getting path's UTF16 pointer: \"%s\"", err)
	}
	attrs, err := syscall.GetFileAttributes(utf16PtrPath)
	return attrs, utf16PtrPath, err
}

// IsHidden checks whether "FileHide.Path" is hidden or not
func (h *FileHide) IsHidden() (bool, error) {
	attrs, _, err := getFileAttrs(h.Path)
	if err != nil {
		return false, fmt.Errorf("something went wrong getting the file attributes: \"%s\"", err)
	}
	if attrs&syscall.FILE_ATTRIBUTE_HIDDEN > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// Hide makes file or directory hidden
func (h *FileHide) Hide() (dstName string, err error) {
	return h.hide(true)
}

// Unhide makes file or directory unhidden
func (h *FileHide) Unhide() (dstName string, err error) {
	return h.hide(false)
}

func (h *FileHide) hide(hidden bool) (dstName string, err error) {
	_, err = os.Stat(h.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("\"%s\" is not exists", h.Path)
		}
		return "", fmt.Errorf("something went wrong getting file stat: \"%s\"", err)
	}

	attrs, utf16PtrPath, err := getFileAttrs(h.Path)
	if err != nil {
		return "", fmt.Errorf("something went wrong getting the file attributes: \"%s\"", err)
	}

	var newAttrs uint32
	if hidden {
		if attrs&syscall.FILE_ATTRIBUTE_HIDDEN > 0 {
			return h.Path, nil
		}
		// Add hidden attribute to file's current attributes
		newAttrs = attrs | syscall.FILE_ATTRIBUTE_HIDDEN
	} else {
		if attrs&syscall.FILE_ATTRIBUTE_HIDDEN == 0 {
			return h.Path, nil
		}
		newAttrs = attrs - (attrs & syscall.FILE_ATTRIBUTE_HIDDEN)
	}
	err = syscall.SetFileAttributes(utf16PtrPath, newAttrs)
	if err != nil {
		return "", fmt.Errorf("something went wrong setting the file attributes: \"%s\"", err)
	}

	return h.Path, nil
}
