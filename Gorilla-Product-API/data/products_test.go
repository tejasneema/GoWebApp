package data

import (
	"fmt"
	"testing"
)

func TestValidate(t *testing.T) {
	p := &Product{
		Name:  "Something",
		Price: 23323,
		SKU:   "acf-aqf-acf",
	}

	err := p.Validate()

	fmt.Println("Error log: ", err)
	if nil != err {
		t.Fatal(err)
	}
}
