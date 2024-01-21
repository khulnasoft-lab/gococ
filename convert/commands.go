package convert

/*
Conversion commands.

Convert one type to another.

This package provides conversions that can be used in chains of commands.
*/

import (
	"strconv"

	"github.com/Khulnasoft-lab/gococ"
)

// Convert a string to an integer.
//
// Params:
//   - str (string): A string that contains a number.
//
// Returns:
//   - An integer.
func Atoi(c gococ.Context, p *gococ.Params) (interface{}, gococ.Interrupt) {
	src := p.Get("str", "0").(string)
	return strconv.Atoi(src)
}
