package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "    hello  world    ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "   foo   bar   baz   ",
			expected: []string{"foo", "bar", "baz"},
		},
		{
			input:    "singleword",
			expected: []string{"singleword"},
		},
		{
			input:    "   spaced   ",
			expected: []string{"spaced"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		for i, _ := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("cleanInput = %v; want %v", word, expectedWord)
			}
		}
	}
}
