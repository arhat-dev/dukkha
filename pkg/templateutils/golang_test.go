package templateutils

import "testing"

func TestJSEscaping(t *testing.T) {
	testCases := []struct {
		in, exp string
	}{
		{`a`, `a`},
		{`'foo`, `\'foo`},
		{`Go "jump" \`, `Go \"jump\" \\`},
		{`Yukihiro says "今日は世界"`, `Yukihiro says \"今日は世界\"`},
		{"unprintable \uFDFF", `unprintable \uFDFF`},
		{`<html>`, `\u003Chtml\u003E`},
		{`no = in attributes`, `no \u003D in attributes`},
		{`&#x27; does not become HTML entity`, `\u0026#x27; does not become HTML entity`},
	}
	for _, tc := range testCases {
		s := golangNS{}.JSEscapeString(tc.in)
		if s != tc.exp {
			t.Errorf("JS escaping [%s] got [%s] want [%s]", tc.in, s, tc.exp)
		}
	}
}
