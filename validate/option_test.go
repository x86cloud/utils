package validate

import (
	"fmt"
	"testing"
	"unsafe"
)

type Struct struct {
	OrderBy string `form:"orderBy"`
	Name    string `form:"name"`
	Page    int32  `form:"page"`

	Descs    bool  `form:"desc"`
	PageSize int32 `form:"pageSize"`

	//Desc     bool   `form:"desc"`
}

//
func TestCapitalize(t *testing.T) {
	fmt.Println(unsafe.Alignof(Struct{}))
}
