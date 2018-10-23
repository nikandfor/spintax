// spintax is an text generator from Spintax template
package spintax

import (
	"fmt"
	"math/rand"
	"strings"
)

type (
	// Exp is an Spintax expression.
	// Spin returns all pars cat together
	Exp []Spintax
	// Alt is an alternative expression
	// Spin returns one of options
	Alt []Spintax
	// Str is an simple string
	Str string

	Spintax interface {
		// Returns one of options
		Spin() string
		// Counts total number of options
		Count() int
		// Returns template
		String() string
		// Returns all possible options
		All() []string
		// Iter gives options one by one
		Iter() <-chan string
	}
)

// Parse parses template
func Parse(exp string) (Spintax, error) {
	e, i, err := parseExp(exp, 0)
	if i != len(exp) {
		return nil, fmt.Errorf("brackets balance error: %v", i)
	}
	return e, err
}

func parseExp(e string, i int) (Spintax, int, error) {
	var r Exp
	var err error
	s := i
loop:
	for i < len(e) {
		c := e[i]
		switch c {
		case '{':
			if s != i {
				r = append(r, Str(e[s:i]))
			}
			var alt Spintax
			alt, i, err = parseAlt(e, i)
			if err != nil {
				return nil, i, err
			}
			i++
			s = i
			if alt != nil {
				r = append(r, alt)
			}
		case '|', '}':
			break loop
		default:
			i++
		}
	}
	if s != i {
		r = append(r, Str(e[s:i]))
	}
	if r == nil {
		return Str(""), i, nil
	}
	if len(r) == 1 {
		return r[0], i, nil
	}
	return r, i, nil
}

func parseAlt(e string, i int) (Spintax, int, error) {
	var r Alt
	var exp Spintax
	var err error
	d := 0
	for i < len(e) {
		if e[i] == '}' {
			d--
			break
		}
		if e[i] == '|' || e[i] == '{' {
			if e[i] == '{' {
				d++
			}
			i++
		}
		exp, i, err = parseExp(e, i)
		if err != nil {
			return nil, i, err
		}
		r = append(r, exp)
	}
	if d != 0 {
		return nil, i, fmt.Errorf("brackets balance error at pos %v", i)
	}
	if len(r) == 1 {
		return r[0], i, nil
	}
	return r, i, nil
}

func (e Exp) Spin() string {
	var b strings.Builder
	for _, e := range e {
		b.WriteString(e.Spin())
	}
	return b.String()
}

func (a Alt) Spin() string {
	e := a[rand.Intn(len(a))]
	return e.Spin()
}

func (s Str) Spin() string { return string(s) }

func (e Exp) Count() int {
	s := 1
	for _, e := range e {
		s *= e.Count()
	}
	return s
}

func (a Alt) Count() int {
	s := 0
	for _, e := range a {
		s += e.Count()
	}
	return s
}

func (s Str) Count() int { return 1 }

func (e Exp) String() string {
	var b strings.Builder
	for _, e := range e {
		b.WriteString(e.String())
	}
	return b.String()
}

func (a Alt) String() string {
	var b strings.Builder
	b.WriteString("{")
	for i, e := range a {
		if i != 0 {
			b.WriteString("|")
		}
		b.WriteString(e.String())
	}
	b.WriteString("}")
	return b.String()
}

func (s Str) String() string { return string(s) }

func (e Exp) All() []string {
	if len(e) == 1 {
		return e[0].All()
	}

	f := e[0].All()
	t := e[1:].All()

	var r []string
	for _, f := range f {
		for _, t := range t {
			r = append(r, f+t)
		}
	}

	return r
}

func (a Alt) All() []string {
	var r []string

	for _, e := range a {
		all := e.All()
		r = append(r, all...)
	}

	return r
}

func (s Str) All() []string { return []string{string(s)} }

func (e Exp) Iter() <-chan string {
	c := make(chan string, 1)
	go func() {
		l := len(e)
		t := make([]string, l)
		s := make([]<-chan string, l)
		var b strings.Builder
		var ok bool

		for j := 0; j < l; j++ {
			s[j] = e[j].Iter()
			t[j] = <-s[j]
		}

		for {
			b.Reset()

			for _, s := range t {
				b.WriteString(s)
			}
			c <- b.String()

			var j int
			for ; j < l; j++ {
				t[j], ok = <-s[j]
				if ok {
					break
				}
				s[j] = e[j].Iter()
				t[j], ok = <-s[j]
				if !ok {
					panic("no data at renewed Iter")
				}
			}
			if j == l {
				break
			}
		}

		close(c)
	}()
	return c
}

func (a Alt) Iter() <-chan string {
	c := make(chan string, 1)
	go func() {
		for _, e := range a {
			sub := e.Iter()
			for l := range sub {
				c <- l
			}
		}
		close(c)
	}()
	return c
}

func (s Str) Iter() <-chan string {
	c := make(chan string, 1)
	c <- string(s)
	close(c)
	return c
}
