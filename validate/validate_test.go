package validate

import (
	"testing"
)

type Test struct {
	Eq      int     `validate:"Eq=10"`
	EqFloat float32 `validate:"Eq=10"`
	Name    string  `validate:"Length=10"`
	//Email   string  `validate:"email"`
	TestSub []TestSub
}

type TestSub struct {
	Eq      int     `validate:"Eq=10"`
	EqFloat float32 `validate:"Eq=12"`
	Name    string  `validate:"Length=3"`
}

func BenchmarkValidate(b *testing.B) {
	t := Test{
		Eq:      10,
		EqFloat: 10,
		Name:    "1234567890",
		TestSub: []TestSub{
			{
				Eq:      10,
				EqFloat: 10,
				Name:    "111",
			},
		},
	}
	for i := 0; i < b.N; i++ {
		err := Validate(t)
		if err != nil {
			return
		}
	}
}
