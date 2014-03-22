package font

import (
	"fmt"
)

type Font struct {
	family	string
	size	int
}

func New(fm string, sz int) (font *Font) {
	fm = saneFamily("serif", fm)
	sz = saneSize(sz)
	return &Font{fm, sz}
}

func (font *Font) SetFamily(fm string) {
	font.family = saneFamily(font.family, fm)
}

func (font *Font) Family() string {
	return font.family
}

func (font *Font) SetSize(sz int) {
	font.size = saneSize(sz)
}

func (font *Font) Size() int {
	return font.size
}

func (font *Font) String() string {
	return fmt.Sprintf("{font-family: %q; font-size: %dpt;}", font.family, font.size) 
}

var saneFamily = func(old string, new_ string) string {
	if "" == new_ {
		return old
	}
	return new_
}

var saneSize = func(sz int) int {
	if sz < 5 {
		return 5
	} else if sz > 144 {
		return 144
	}
	return sz
}