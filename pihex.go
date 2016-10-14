package pihex

import (
	"math"
	"math/big"
)

const (
	PrecisionCutoff = 1.0e-17
	SubRanges       = 8
)

func Digit(n int64) byte {
	seriesResult1 := make(chan float64)
	seriesResult2 := make(chan float64)
	seriesResult3 := make(chan float64)
	seriesResult4 := make(chan float64)

	go series(n, 1, seriesResult1)
	go series(n, 4, seriesResult2)
	go series(n, 5, seriesResult3)
	go series(n, 6, seriesResult4)

	digit :=
		(4.0 * <-seriesResult1) -
			(2.0 * <-seriesResult2) -
			<-seriesResult3 -
			<-seriesResult4

	digit = digit - math.Floor(digit) + 1.0
	return byte(math.Floor(16.0*math.Remainder(digit, 1.0))) & 0x0f
}

func preNSubSeries(kStart, kEnd, n, magic_constant int64, result chan<- float64) {
	acc := 0.0
	temp := big.NewInt(0)
	sixteen := big.NewInt(16)
	nMinusK := big.NewInt(0)
	modulo := big.NewInt(0)

	for k := kStart; k <= kEnd; k++ {
		nMinusK.SetInt64(n - k)
		modulo.SetInt64((k << 3) + magic_constant)
		temp.Exp(sixteen, nMinusK, modulo)
		acc = acc + float64(temp.Int64())/float64(modulo.Int64())
		acc = acc - math.Floor(acc)
	}

	result <- acc
}

func preNSeries(n, magic_constant int64) float64 {
	acc := 0.0

	if n < 100 {
		// For low values of n, don't bother trying to split the
		// calculation across concurrent routines.
		subResult := make(chan float64)
		go preNSubSeries(0, n-1, n, magic_constant, subResult)
		acc = <-subResult
	} else {
		subResults := buildChannels(SubRanges)
		span := n / SubRanges
		var kStart, kEnd int64

		// Spawn concurrent calculations for each sub-range between 0 and n - 1
		for i := int64(0); i < SubRanges; i++ {
			kStart = span * i
			if i == SubRanges-1 {
				kEnd = n - 1
			} else {
				kEnd = (kStart + span) - 1
			}

			go preNSubSeries(kStart, kEnd, n, magic_constant, subResults[i])
		}

		// Sum results of concurrent calculations
		for i := 0; i < SubRanges; i++ {
			acc += <-subResults[i]
			acc = acc - math.Floor(acc)
		}
	}

	return acc
}

func postNSeries(n, magic_constant int64) float64 {
	var modulo, temp float64
	acc := 0.0

	// Calculate the range where k > n to the limits of precision.
	for k := n; k < n+100; k++ {
		modulo = float64((k << 3) + magic_constant)
		temp = math.Pow(16.0, float64(n-k)) / modulo
		if temp < PrecisionCutoff {
			break
		}
		acc = acc + temp
		acc = acc - math.Floor(acc)
	}

	return acc
}

func buildChannels(num int64) []chan float64 {
	channels := make([]chan float64, SubRanges)
	for i := 0; i < SubRanges; i++ {
		channels[i] = make(chan float64)
	}
	return channels
}

func series(n, magic_constant int64, result chan<- float64) {
	acc := preNSeries(n, magic_constant)
	acc += postNSeries(n, magic_constant)
	acc = acc - math.Floor(acc)
	result <- acc
}
