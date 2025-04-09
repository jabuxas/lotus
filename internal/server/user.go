package server

import (
	"math"
)

type User struct {
	id   int
	name string
	exp  int
}

const (
	base     = 100.0
	exponent = 1.5
)

func (u *User) calculateLevel() int {
	if u.exp < base {
		return 1
	}
	return int(math.Pow(float64(u.exp)/base, 1/exponent))
}
