package util

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
