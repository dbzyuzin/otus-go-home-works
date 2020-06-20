package models

type Test struct {
	I, I2 int      `validate:"min:18"`
	S     string   `validate:s="a"`
	arrI  []int    `validate:"in:18,5"`
	arrS  []string `validate:"len:5"`
}
