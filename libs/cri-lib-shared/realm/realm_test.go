package realm

import "testing"

func TestValidate(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"cr-prod", true},
		{"sv-prod", true},
		{"demo", true},
		{"a", false},
		{"", false},
		{"CR-prod", false},
		{"cr_prod", false},
		{"cr-prod-cr-prod-cr-prod-cr-prod-cr", false},
	}
	for _, c := range cases {
		got := Validate(c.in) == nil
		if got != c.want {
			t.Errorf("Validate(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestSchemaName(t *testing.T) {
	if got := SchemaName("cr-prod"); got != "cr_prod" {
		t.Errorf("SchemaName(cr-prod) = %q", got)
	}
	if got := SchemaName("demo"); got != "demo" {
		t.Errorf("SchemaName(demo) = %q", got)
	}
}
