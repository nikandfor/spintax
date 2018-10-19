package spintax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	e := Parse(`string`)
	assert.Equal(t, Str(`string`), e)

	e = Parse(`{str_a|str_b}`)
	assert.Equal(t, Alt{Str("str_a"), Str("str_b")}, e)

	e = Parse(`pref {str_a|str_b} suff`)
	assert.Equal(t, Exp{
		Str("pref "),
		Alt{Str("str_a"), Str("str_b")},
		Str(" suff"),
	}, e)

	e = Parse(`pref {str_a|subp {alt_a|alt_b|alt_c} subs} suff`)
	assert.Equal(t, Exp{
		Str("pref "),
		Alt{
			Str("str_a"),
			Exp{
				Str("subp "),
				Alt{Str("alt_a"), Str("alt_b"), Str("alt_c")},
				Str(" subs"),
			},
		},
		Str(" suff"),
	}, e)

	e = Parse(`pref {|str_a|str_b|} suff`)
	assert.Equal(t, Exp{
		Str("pref "),
		Alt{Str(""), Str("str_a"), Str("str_b"), Str("")},
		Str(" suff"),
	}, e)
}

func TestCount(t *testing.T) {
	e := Parse(`pref {str_a|subp {alt_a|alt_b|alt_c} subs} suff`)

	assert.Equal(t, 4, e.Count())
}

func TestGen(t *testing.T) {
	e := Parse(`pref {str_a|subp {alt_a|alt_b|alt_c} subs} suff`)
	vars := []string{
		"pref str_a suff",
		"pref subp alt_a subs suff",
		"pref subp alt_b subs suff",
		"pref subp alt_c subs suff",
	}

	for i := 0; i < 10; i++ {
		v := e.Spin()
		assert.Subset(t, vars, []string{v}, "%q is not result of %q", v, e)
	}
}

func TestParseRoll(t *testing.T) {
	s := `pref {str_a|subp {alt_a|alt_b|alt_c} subs} suff`
	assert.Equal(t, s, Parse(s).String())

	s = `pref {|str_a|str_b|} suff`
	assert.Equal(t, s, Parse(s).String())
}

func TestAll(t *testing.T) {
	all := Parse("string").All()
	assert.Equal(t, []string{"string"}, all)

	all = Parse("{a|b}").All()
	assert.Equal(t, []string{"a", "b"}, all)

	all = Parse("a {|b|c} d").All()
	assert.Equal(t, []string{"a  d", "a b d", "a c d"}, all)

	e := Parse(`pref {str_a|subp {alt_a|alt_b|alt_c} subs} suff`)
	vars := []string{
		"pref str_a suff",
		"pref subp alt_a subs suff",
		"pref subp alt_b subs suff",
		"pref subp alt_c subs suff",
	}
	assert.Equal(t, vars, e.All())
}
