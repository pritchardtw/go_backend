package main

import "testing"

func TestHashPassword(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{
			"angryMonkey",
			"ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==",
		},
	}
	for _, c := range cases {
		got := hashPassword(c.in)
		if got != c.want {
			t.Errorf("hashPassword(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
