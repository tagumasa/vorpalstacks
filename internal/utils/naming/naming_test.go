package naming

import "testing"

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"single word", "hello", "hello"},
		{"snake_case", "hello_world", "helloWorld"},
		{"multiple words", "my_variable_name", "myVariableName"},
		{"already camel", "helloWorld", "helloworld"},
		{"uppercase snake", "HELLO_WORLD", "helloWorld"},
		{"leading underscore", "_private_field", "privateField"},
		{"trailing underscore", "field_", "field"},
		{"double underscore", "hello__world", "helloWorld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToCamelCase(tt.input); got != tt.want {
				t.Errorf("ToCamelCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"single word", "hello", "Hello"},
		{"snake_case", "hello_world", "HelloWorld"},
		{"multiple words", "my_variable_name", "MyVariableName"},
		{"already pascal", "HelloWorld", "Helloworld"},
		{"uppercase snake", "HELLO_WORLD", "HelloWorld"},
		{"leading underscore", "_private_field", "PrivateField"},
		{"single char", "a", "A"},
		{"double underscore", "hello__world", "HelloWorld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToPascalCase(tt.input); got != tt.want {
				t.Errorf("ToPascalCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"pascal case", "HelloWorld", "hello_world"},
		{"camel case", "helloWorld", "hello_world"},
		{"single word", "Hello", "hello"},
		{"already snake", "hello_world", "hello_world"},
		{"uppercase", "HELLOWORLD", "h_e_l_l_o_w_o_r_l_d"},
		{"multiple capitals", "myXMLParser", "my_x_m_l_parser"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSnakeCase(tt.input); got != tt.want {
				t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"test.txt", "test.txt"},
		{"my/file:name.txt", "my_file_name.txt"},
		{"file*name?.txt", "file_name_.txt"},
		{"file|name.txt", "file_name.txt"},
		{"my<file>name.txt", "my_file_name.txt"},
		{"normal-file.txt", "normal-file.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := SanitizeFilename(tt.input); got != tt.want {
				t.Errorf("SanitizeFilename(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		{"hello_world", "hello_world"},
		{"my-name_123", "myname_123"},
		{"name@#$%", "name"},
		{"hello world!", "helloworld"},
		{"", ""},
		{"123", "123"},
		{"_underscore_", "_underscore_"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := SanitizeName(tt.input); got != tt.want {
				t.Errorf("SanitizeName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
