package util

import (
	"path/filepath"
	"runtime"
)

func ContainString(ss []string, s string) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}

func Unique(ss []string) (result []string) {
	smap := make(map[string]struct{})
	for _, s := range ss {
		if _, ok := smap[s]; ok {
			continue
		}
		result = append(result, s)
		smap[s] = struct{}{}
	}

	return result
}

func GetLastDir(path string) string {
	// 规范化路径（处理冗余分隔符、'.'和'..'）
	cleanPath := filepath.Clean(path)

	// 判断是否为根目录（不同平台有差异）
	if isRoot(cleanPath) {
		return ""
	}

	// 直接获取路径最后一部分
	return filepath.Base(cleanPath)
}

func isRoot(path string) bool {
	if runtime.GOOS == "windows" {
		vol := filepath.VolumeName(path)
		rest := path[len(vol):]

		if vol != "" {
			// 处理网络路径（如 \\server\share）
			if len(vol) > 2 && vol[0] == '\\' && vol[1] == '\\' {
				return rest == "" || rest == string(filepath.Separator)
			}
			// 处理本地磁盘（如 C:\）
			return rest == string(filepath.Separator)
		}
		return false
	}
	// UNIX 系统：根目录为单斜杠 /
	return path == string(filepath.Separator)
}

func Default(s, ss string) string {
	if s != "" {
		return s
	}
	return ss
}
