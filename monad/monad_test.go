package monad_test

import (
	"strconv"
	"testing"

	"go.osspkg.com/casecheck"
	. "go.osspkg.com/do/monad"
)

func str2int(arg string) (int64, error) {
	return strconv.ParseInt(arg, 10, 64)
}
func int2str(arg int64) (string, error) {
	return strconv.FormatInt(arg, 10), nil
}

func TestUnit_Monad(t *testing.T) {
	result, err := Bind(
		Bind(
			Some("123"),
			str2int,
		),
		int2str,
	).
		Return()

	casecheck.NoError(t, err)
	casecheck.Equal(t, "123", result)
}

func BenchmarkMonad(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b1 := Some("123")
			b2 := Bind(b1, str2int)
			b3 := Bind(b2, int2str)
			Bind(b3, str2int)
		}
	})
}
