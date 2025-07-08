package repo

import (
	"fmt"
	"testing"
)

func TestConfig_Init(t *testing.T) {
	config := GetConfig().Init()
	r, err := config.CurrentRepo()
	config.CheckErr(err)
	fmt.Println("repo url is: ", r.Url)
	
}
