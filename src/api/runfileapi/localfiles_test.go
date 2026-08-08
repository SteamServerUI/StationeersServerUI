package runfileapi

import "testing"

func TestProtectedFileRecognizesOnlySSUITree(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"./SSUI/config.json", true},
		{"SSUI\\identity.json", true},
		{"./SSUI", true},
		{"./SSUI-mod/config.json", false},
		{"./server/config.json", false},
	}
	for _, test := range tests {
		if got := protectedFile(test.path); got != test.want {
			t.Fatalf("protectedFile(%q) = %v, want %v", test.path, got, test.want)
		}
	}
}
