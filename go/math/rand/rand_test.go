package rand_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	rand_ "github.com/kaydxh/golang/go/math/rand"

	"github.com/stretchr/testify/assert"
)

func TestRand(t *testing.T) {
	s := fmt.Sprintf("%08v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000))
	fmt.Printf("s: %v\n", s)

	ss := "abc_dvg_123"
	fmt.Printf("ss %v\n", ss)
	ts := strings.TrimPrefix(ss, "abc_dvg")
	fmt.Printf("ts: %v\n", ts)

	ns := "_abc"
	nns := strings.Split(ns, "_")
	fmt.Printf("nns: %v\n", nns)

}

func TestRangeInt(t *testing.T) {
	testCases := []struct {
		min int
		max int
	}{
		{
			min: 10,
			max: 12,
		},
		{
			min: 10000000,
			max: 100000000,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			r, err := rand_.RangeInt(testCase.min, testCase.max)
			if err != nil {
				t.Fatalf("failed to rand int, err: %v", err)
			}
			t.Logf("random: %v", r)

			assert.GreaterOrEqual(t, r, testCase.min)
			assert.LessOrEqual(t, r, testCase.max)

		})
	}
}

func TestRead(t *testing.T) {
	testCases := []struct {
		p []byte
	}{
		{
			p: make([]byte, 10),
		},
		{
			p: make([]byte, 20),
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			n, err := rand_.Read(testCase.p)
			if err != nil {
				t.Fatalf("failed to rand int, err: %v", err)
			}
			t.Logf("read n: %v, p: %v", n, testCase.p)

			assert.Equal(t, len(testCase.p), n)

		})
	}
}

func TestRangeString(t *testing.T) {
	testCases := []struct {
		n int
	}{
		{
			n: 0,
		},
		{
			n: 5,
		},
		{
			n: 8,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			str := rand_.RangeString(testCase.n)
			t.Logf("str: %v", str)
			assert.Equal(t, len(str), testCase.n)

		})
	}

}
