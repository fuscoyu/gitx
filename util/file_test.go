package util

import (
	"fmt"
	"testing"
)

func TestFindOriginURL(t *testing.T) {
	fmt.Println(FindOriginURL("/Users/yu/code/git.internal.yunify.com/chenyu/doc"))
	fmt.Println(FindOriginURL("/Users/yu/code/git.internal.yunify.com/bi/pitrix-wh-daemon"))
}
