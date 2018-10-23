package spintax

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	e, err := Parse(`{}`)
	assert.NoError(t, err)
	assert.Equal(t, Str(""), e)

	e, err = Parse(`string`)
	assert.NoError(t, err)
	assert.Equal(t, Str(`string`), e)

	e, err = Parse(`{str_a|str_b}`)
	assert.NoError(t, err)
	assert.Equal(t, Alt{Str("str_a"), Str("str_b")}, e)

	e, err = Parse(`{str}`)
	assert.NoError(t, err)
	assert.Equal(t, Str("str"), e)

	e, err = Parse(`pref {str_a|str_b} suff`)
	assert.NoError(t, err)
	assert.Equal(t, Exp{
		Str("pref "),
		Alt{Str("str_a"), Str("str_b")},
		Str(" suff"),
	}, e)

	e, err = Parse(`pref {str_a|subp {alt_a|alt_b|alt_c} subs} suff`)
	assert.NoError(t, err)
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

	e, err = Parse(`pref {|str_a|str_b|} suff`)
	assert.NoError(t, err)
	assert.Equal(t, Exp{
		Str("pref "),
		Alt{Str(""), Str("str_a"), Str("str_b"), Str("")},
		Str(" suff"),
	}, e)
}

func TestError(t *testing.T) {
	_, err := Parse(`{string`)
	assert.Error(t, err)

	_, err = Parse(`string}`)
	assert.Error(t, err)

	_, err = Parse(`aa{bb{cc|dd}ee`)
	assert.Error(t, err)

	_, err = Parse(`aa{bb|cc}dd}ee`)
	assert.Error(t, err)

	_, err = Parse(`{{str{{{{ing}}`)
	assert.Error(t, err)
}

func TestCount(t *testing.T) {
	e, err := Parse(`pref {str_a|subp {alt_a|alt_b|alt_c} subs} suff`)
	assert.NoError(t, err)

	assert.Equal(t, 4, e.Count())
}

func TestSpin(t *testing.T) {
	e, err := Parse(`pref {str_a|subp {alt_a|alt_b|alt_c} subs} suff`)
	assert.NoError(t, err)
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
	e, err := Parse(s)
	assert.NoError(t, err)
	assert.Equal(t, s, e.String())

	s = `pref {|str_a|str_b|} suff`
	e, err = Parse(s)
	assert.NoError(t, err)
	assert.Equal(t, s, e.String())
}

func TestAll(t *testing.T) {
	e, err := Parse("string")
	assert.NoError(t, err)
	assert.Equal(t, []string{"string"}, e.All())

	e, err = Parse("{a|b}")
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, e.All())

	e, err = Parse("a {|b|c} d")
	assert.NoError(t, err)
	assert.Equal(t, []string{"a  d", "a b d", "a c d"}, e.All())

	e, err = Parse(`pref {str_a|subp {alt_a|alt_b|alt_c} subs} suff`)
	assert.NoError(t, err)
	vars := []string{
		"pref str_a suff",
		"pref subp alt_a subs suff",
		"pref subp alt_b subs suff",
		"pref subp alt_c subs suff",
	}
	assert.Equal(t, vars, e.All())
}

func TestIter(t *testing.T) {
	e, err := Parse("string")
	assert.NoError(t, err)
	all := iter_all(e)
	assert.Equal(t, []string{"string"}, all)

	e, err = Parse("{a|b}")
	assert.NoError(t, err)
	all = iter_all(e)
	assert.Equal(t, []string{"a", "b"}, all)

	e, err = Parse("a {|b|c} d")
	assert.NoError(t, err)
	all = iter_all(e)
	assert.Equal(t, []string{"a  d", "a b d", "a c d"}, all)

	e, err = Parse(`pref {str_a|subp {alt_a|alt_b|alt_c} subs} suff`)
	assert.NoError(t, err)
	all = iter_all(e)
	vars := []string{
		"pref str_a suff",
		"pref subp alt_a subs suff",
		"pref subp alt_b subs suff",
		"pref subp alt_c subs suff",
	}
	assert.Equal(t, vars, all)
}

func iter_all(e Spintax) []string {
	c := e.Iter()
	var res []string
	for s := range c {
		res = append(res, s)
	}
	return res
}

func ExampleSpintax() {
	e, err := Parse(`first {single|second {mid_a|mid_b|mid_c} before_last} last`)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	for _, s := range e.All() {
		fmt.Printf("%v\n", s)
	}
	// Output:
	// first single last
	// first second mid_a before_last last
	// first second mid_b before_last last
	// first second mid_c before_last last
}
