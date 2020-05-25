package refresh

import (
	"strings"

	"github.com/fatih/color"
)

type Renderable interface {
	Render() bool
}

type CliOutput interface {
	Renderable
	NoColor() string
	WithColor() string
}

type noShow struct{}

func (noShow) NoColor() string   { return "" }
func (noShow) WithColor() string { return "" }
func (noShow) Render() bool      { return false }

func NoShow() CliOutput {
	return noShow{}
}

func FromString(s string) CliOutput {
	return FromStrings(s, s)
}

func WithColors(s string, colors ...color.Attribute) CliOutput {
	return FromStrings(s, color.New(colors...).Sprint(s))
}

type basicImpl struct {
	noColor, withColor string
}

func (b basicImpl) NoColor() string   { return b.noColor }
func (b basicImpl) WithColor() string { return b.withColor }
func (b basicImpl) Render() bool      { return true }
func FromStrings(noColor, withColor string) CliOutput {
	return basicImpl{
		noColor:   noColor,
		withColor: withColor,
	}
}

type concatImpl struct {
	outputs []CliOutput
	// TODO(rchew) handle delim with Render() == false? Or just make delim a string instead?
	delim CliOutput
}

func (c *concatImpl) NoColor() string {
	return join(CliOutput.NoColor, c.delim, c.outputs...)
}

func (c *concatImpl) WithColor() string {
	return join(CliOutput.WithColor, c.delim, c.outputs...)
}

func join(fn func(CliOutput) string, delim CliOutput, elems ...CliOutput) string {
	outputs := make([]string, 0, len(elems))
	for _, o := range elems {
		if o.Render() {
			outputs = append(outputs, fn(o))
		}
	}
	return strings.Join(outputs, fn(delim))
}

func (c *concatImpl) Render() bool { return true }

func Concat(delim CliOutput, outputs ...CliOutput) CliOutput {
	return &concatImpl{
		outputs: outputs,
		delim:   delim,
	}
}
