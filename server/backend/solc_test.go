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
		b: "1234a165627a7a7230582086a6b3c7c6b942a6bcf7e8ea07686953fe563275da0b17537439e884054eede40029",
	}.run)
	t.Run("meta-empty", test{
		a: "1234",
		b: "1234a165627a7a7230582023412341234123412341234123412341234123412341234123412341234123400029",
	}.run)
	t.Run("suffix", test{
		a: "1234cdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdabcdefghi",
		b: "1234abababababababababababababababababababababababababababababab123456789",
	}.run)
	t.Run("0x-suffix-meta", test{
		a: "1234cdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdabcdefghia165627a7a7230582012341234123412341234123412341234123412341234123412341234123412340029",
		b: "0x1234abababababababababababababababababababababababababababababab123456789",
	}.run)
	t.Run("meta-suffix-0029", test{
		a: `0x80518281529051600091600160a060020a03851691600080516020610dd28339815191529181900360200190a350505600ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa165627a7a723058206944ecd1b3fc415e254283b30bfe6e28fe0b1472f35ac60ccbe2cfc4780b16050029`,
		b: `80518281529051600091600160a060020a03851691600080516020610dd28339815191529181900360200190a350505600ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa165627a7a7230582086a6b3c7c6b942a6bcf7e8ea07686953fe563275da0b17537439e884054eede40029`,
	}.run)
	//TODO 30 and 31?
	t.Run("meta-suffix-0032", test{
		// bzzr1
		a: `0xabcd1234a265627a7a723158202769c06084c54e3ad332d8faf5bff6c45b1f4129a9aeefb1d07070898e40507e64736f6c634300050d0032`,
		// bzzr0
		b: `0xabcd1234a265627a7a72305820ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff64736f6c63431234560032`,
	}.run)
	t.Run("meta-suffix-0033", test{
		a: `0x6b656e54696d656c6f636b3a206e6f20746f6b656e7320746f2072656c65617365a2a26469706673582212200486dd6b5b50ae39b46b3721e657d9d19624e14d6d86ceb91a22192400441b5064736f6c634300060c0033`,
		b: `0x6b656e54696d656c6f636b3a206e6f20746f6b656e7320746f2072656c65617365a2a264697066735822122024e555269fe2ec67c5737c0c0328973111b22e56a19e7546fec699ff425fd57164736f6c634300060c0033`,
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
	t.Helper()
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
