package backend

import (
	"testing"
)

func TestSolcBinEqual(t *testing.T) {
	type test = testSolcBinEqual
	t.Run("empty", test{
		a: "",
		b: "",
	}.run)
	t.Run("0x", test{
		a: "0x1",
		b: "1",
	}.run)
	t.Run("meta", test{
		a: "1234",
		b: "1234056fea165627a7a7230582086a6b3c7c6b942a6bcf7e8ea07686953fe563275da0b17537439e884054eede40029",
	}.run)
	t.Run("meta-empty", test{
		a: "1234",
		b: "1234056fea165627a7a723058200029",
	}.run)
	t.Run("suffix", test{
		a: "1234cdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdabcdefghi",
		b: "1234abababababababababababababababababababababababababababababab123456789",
	}.run)
	t.Run("0x-suffix-meta", test{
		a: "1234cdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdabcdefghi056fea165627a7a72305820ABCDE0029",
		b: "0x1234abababababababababababababababababababababababababababababab123456789",
	}.run)

	t.Run("meta-suffix", test{
		a:        "12345678",
		b:        "1234056fea165627a7a7230582000295678",
		mismatch: true,
	}.run)
	t.Run("meta-multi", test{
		a:        "1234567890",
		b:        "1234056fea165627a7a723058200029567056fea165627a7a723058200029890",
		mismatch: true,
	}.run)
}

type testSolcBinEqual struct {
	a, b     string
	mismatch bool
}

func (tt testSolcBinEqual) run(t *testing.T) {
	if SolcBinEqual(tt.a, tt.b) {
		if tt.mismatch {
			t.Error("expected mismatch")
		}
	} else {
		if !tt.mismatch {
			t.Error("expected match")
		}
	}
}
