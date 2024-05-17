package anagram

import (
	"testing"
)

func TestAnagram(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		name string
		arg1 string
		arg2 string
		want bool
	}{
		{
			name: "empty",
			want: false,
		},
		{
			name: "one_char",
			arg1: "a",
			arg2: "a",
			want: true,
		},
		{
			name: "anagram_1",
			arg1: "foo",
			arg2: "ofo",
			want: true,
		},
		{
			name: "anagram_2",
			arg1: "foobaranagram",
			arg2: "anarobramfoag",
			want: true,
		},
		{
			name: "anagram_cyr",
			arg1: "инфографика_ABC",
			arg2: "ииоAгBнраCфф_ка",
			want: true,
		},
		{
			name: "not_anagram_wrong_len",
			arg1: "bar",
			arg2: "baar",
			want: false,
		},
		{
			name: "not_anagram_correct_len",
			arg1: "bar",
			arg2: "baa",
			want: false,
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := anagram(test.arg1, test.arg2)
			if test.want != got {
				t.Errorf("got: %v, want: %v", got, test.want)
			}
		})
	}
}
