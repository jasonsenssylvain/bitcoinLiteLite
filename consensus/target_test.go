package consensus

import (
	//	"fmt"
	"fmt"
	"math/big"
	"testing"
)

type onevec struct {
	b uint32
	e string
	d float64
}

var testvecs = []onevec{
	{b: 0x1D00FFFF, e: "00000000FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"},
	{b: 0x1b0404cb, e: "00000000000404CB000000000000000000000000000000000000000000000000"},
	{b: 0x1d00ffff, e: "00000000FFFF0000000000000000000000000000000000000000000000000000"},
	{b: 436330132, d: 8974296.01488785},
	{b: 436543292, d: 3275464.59},
	{b: 436591499, d: 2864140.51},
	{b: 436841986, d: 1733207.51},
	{b: 437155514, d: 1159929.50},
	{b: 436789733, d: 1888786.71},
	{b: 453031340, d: 92347.59},
	{b: 453281356, d: 14484.16},
	{b: 470771548, d: 16.62},
	{b: 486604799, d: 1.00},
}

func TestTarget(t *testing.T) {
	yd, _ := new(big.Int).SetString(testvecs[0].e, 16)
	cs := GetCompact(yd)
	fmt.Println(cs)

	fmt.Println("=====================")
	// for i := range testvecs {
	// 	fmt.Println(testvecs[i].b)
	// 	x := SetCompact(testvecs[i].b)
	// 	fmt.Println(x)
	// 	fmt.Println(x.Bits())
	// 	d := GetDifficulty(testvecs[i].b)
	// 	fmt.Println(d)
	// 	c := GetCompact(x)
	// 	fmt.Println(c)
	// 	//fmt.Printf("%d. %d/%d -> %.8f / %.8f\n", i, testvecs[i].b, c, d, testvecs[i].d)
	// 	if testvecs[i].b != c {
	// 		t.Error("Set/GetCompact mismatch at alement", i)
	// 	}

	// 	if testvecs[i].e != "" {
	// 		y, _ := new(big.Int).SetString(testvecs[i].e, 16)
	// 		if x.Cmp(y) != 0 {
	// 			t.Error("Target mismatch at alement", i)
	// 		}
	// 	}

	// 	if testvecs[i].d != 0 && math.Abs(d-testvecs[i].d) > 0.1 {
	// 		t.Error("Difficulty mismatch at alement", i)
	// 	}
	// }
}
