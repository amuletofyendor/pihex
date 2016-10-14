package pihex

import (
	"math"
	"math/big"
)

const (
	PrecisionCutoff = 1.0e-17
)

func Digit(n int64) byte {
	digit :=
		(4.0 * series(n, 1)) -
			(2.0 * series(n, 4)) -
			series(n, 5) -
			series(n, 6)

	digit = digit - math.Floor(digit) + 1.0
	return byte(math.Floor(16.0*math.Remainder(digit, 1.0))) & 0x0f
}

func series(n, magic_constant int64) float64 {
	acc := 0.0

	{
		temp := new(big.Int)
		sixteen := big.NewInt(16)
		var modulo int64

		for k := int64(0); k < n; k++ {
			modulo = 8*k + magic_constant
			temp.Exp(sixteen, big.NewInt(n-k), big.NewInt(modulo))
			acc = acc + float64(temp.Int64())/float64(modulo)
			acc = acc - math.Floor(acc)
		}
	}

	{
		var modulo, temp float64

		for k := n; k < n+100; k++ {
			modulo = float64(8*k + magic_constant)
			temp = math.Pow(16.0, float64(n-k)) / modulo
			if temp < PrecisionCutoff {
				break
			}
			acc = acc + temp
			acc = acc - math.Floor(acc)
		}
	}

	return acc
}
