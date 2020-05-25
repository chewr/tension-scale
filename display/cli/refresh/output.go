package refresh

import (
	"strings"
)

type CliOutput interface {
	NoColor() string
	WithColor() string
}

func FromString(s string) CliOutput {
	return FromStrings(s, s)
}

type basicImpl struct {
	noColor, withColor string
}

func (b basicImpl) NoColor() string   { return b.noColor }
func (b basicImpl) WithColor() string { return b.withColor }
func FromStrings(noColor, withColor string) CliOutput {
	return basicImpl{
		noColor:   noColor,
		withColor: withColor,
	}
}

type concatImpl struct {
	outputs []CliOutput
	delim   CliOutput
}

func (c *concatImpl) NoColor() string {
	elems := make([]string, len(c.outputs))
	for i, o := range c.outputs {
		elems[i] = o.NoColor()
	}
	return strings.Join(elems, c.delim.NoColor())
}

func (c *concatImpl) WithColor() string {
	elems := make([]string, len(c.outputs))
	for i, o := range c.outputs {
		elems[i] = o.WithColor()
	}
	return strings.Join(elems, c.delim.WithColor())
}

func Concat(delim CliOutput, outputs ...CliOutput) CliOutput {
	return &concatImpl{
		outputs: outputs,
		delim:   delim,
	}
}
