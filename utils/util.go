package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	regFileWithDateSuffix = regexp.MustCompile(`\d{2,}$`)
	regExcludeFileSuffix  = regexp.MustCompile(`(\.swp|\.gz|\.zip|\.tar|\.tmp|\.bz|\.bz2|\.z)$`)
	regLogFileSuffix      = regexp.MustCompile(`(\.log|\.json|\.line|\.ngx)$`)
	regFormatStr          = `((\.%v)|((\.%v)(\.|-)(\d{8}|\d{4}(-)\d{2}(-)\d{2})))$`
)

func MD5(s string) string {
	data := []byte(s)
	has := md5.Sum(data)
	md5 := fmt.Sprintf("%x", has)
	return md5
}

//func GetFileStat(p string) (os.FileInfo, *syscall.Stat_t, error) {
//	finfo, err := os.Stat(p)
//	if err != nil {
//		return nil, nil, err
//	}
//	fstat, ok := finfo.Sys().(*syscall.Stat_t)
//	if !ok {
//		return nil, nil, errors.New("*syscall.Stat_t invoke error")
//	}
//
//	return finfo, fstat, nil
//}

// ReadOpen opens a file for reading only
func ReadOpen(path string) (*os.File, error) {
	flag := os.O_RDONLY
	perm := os.FileMode(0)
	return os.OpenFile(path, flag, perm)
}

// as same as file size
func CurOffset(fs *os.File) (int64, error) {
	var offset int64
	offset, err := fs.Seek(0, os.SEEK_END)
	if err != nil {
		return -1, err
	}

	return offset, nil
}

func ExcludeFile(p string) bool {
	// exclude file
	return regExcludeFileSuffix.Match([]byte(p))
}

func ExcludeDir(p string) bool {
	name := filepath.Base(p)
	return strings.Index(name, ".") == 0
}

func IsLogFile(p string) bool {
	return regLogFileSuffix.Match([]byte(p))
}

func RemoveFileRotationSuffix(p string) string {
	idx := strings.LastIndex(p, ".")
	if idx > -1 {
		return p[0:idx]
	} else {
		return p
	}
}

func ContainsOne(s, substr string) bool {
	idx := strings.Index(s, substr)
	lIdx := strings.LastIndex(s, substr)
	return idx == lIdx
}

func CheckFileLogType(path string) (string, error) {
	// old code
	/*	if regFileWithDateSuffix.Match([]byte(path)) {
			idx := strings.LastIndex(path, ".")
			path = path[:idx]
		}

		var logType string
		if strings.HasSuffix(path, "."+LogTypeLogback) {
			logType = LogTypeLogback
		} else if strings.HasSuffix(path, "."+LogTypeJson) {
			logType = LogTypeJson
		} else if strings.HasSuffix(path, "."+LogTypeLine) {
			logType = LogTypeLine
		} else if strings.HasSuffix(path, "."+LogTypeNginx) {
			logType = LogTypeNginx
		} else {
			//do nothing
			return "", errors.New("nonsupport log type for this file!")
		}*/
	var logType string
	regStr := fmt.Sprintf(regFormatStr, LogTypeLogback, LogTypeLogback)
	reg := regexp.MustCompile(regStr)
	if reg.Match([]byte(path)) {
		logType = LogTypeLogback
	}
	if strings.EqualFold(logType, "") {
		regStr := fmt.Sprintf(regFormatStr, LogTypeJson, LogTypeJson)
		reg := regexp.MustCompile(regStr)
		if reg.Match([]byte(path)) {
			logType = LogTypeJson
		}
	}
	if strings.EqualFold(logType, "") {
		regStr := fmt.Sprintf(regFormatStr, LogTypeLine, LogTypeLine)
		reg := regexp.MustCompile(regStr)
		if reg.Match([]byte(path)) {
			logType = LogTypeLine
		}
	}
	if strings.EqualFold(logType, "") {
		regStr := fmt.Sprintf(regFormatStr, LogTypeNginx, LogTypeNginx)
		reg := regexp.MustCompile(regStr)
		if reg.Match([]byte(path)) {
			logType = LogTypeNginx
		}
	}
	if !strings.EqualFold(logType, "") {
		return logType, nil
	} else {
		return "", errors.New("nonsupport log type for this file!")
	}
}

// file name without rotation suffix
func GetFileRealName(path string) string {
	if regFileWithDateSuffix.Match([]byte(path)) {
		idx := strings.LastIndex(path, ".")
		path = path[:idx]
	}

	return path
}

// ${watch_root}/app/${APP_CODE}/${APP_CONTAINER_IP}_${APP_CONTAINER_ID}/"
func ExtractTaskKeyFromPath(p string, root, base, prefix, suffix string) (string, error) {
	wrm := filepath.Join(root, base)
	rp, err := filepath.Rel(wrm, p)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(rp, ".") {
		return "", errors.New("ExtractTaskKeyFromPath relative path is zero or negative!" + wrm)
	}

	k, err := FormatTaskKey(strings.Split(rp, string(os.PathSeparator))[0])
	if err != nil {
		return "", err
	}
	if !strings.EqualFold(k, "") && !strings.EqualFold(prefix, "") {
		k = prefix + "_" + k
	}
	if !strings.EqualFold(k, "") && !strings.EqualFold(suffix, "") {
		k = k + "_" + suffix
	}
	return k, nil
}

func FormatTaskKey(k string) (string, error) {
	regx := regexp.MustCompile(`^[\w\-.]{3,64}$`)
	if !regx.Match([]byte(k)) {
		return "", errors.New("FormatTaskKey() meet anomaly key!" + k)
	}

	k = strings.Replace(k, "-", "_", -1)
	k = strings.Replace(k, ".", "_", -1)

	return strings.ToLower(k), nil
}

func Md5(data []byte) (string, error) {
	if nil != data && len(data) > 0 {
		md5str := fmt.Sprintf("%x", md5.Sum(data))
		return md5str, nil
	}
	return "", errors.New("data is nil or len is zero")
}

func Md5ForFirst6Char(data []byte) (string, error) {
	if nil != data && len(data) > 0 {
		md5str := fmt.Sprintf("%x", md5.Sum(data))
		return string([]byte(md5str)[:6]), nil
	}
	return "", errors.New("data is nil or len is zero")
}
