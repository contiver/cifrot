package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"unicode/utf8"
)

const alphsize = 27
const es_index = 0.07464494

var mindiff = math.Inf(1)
var freqcount [27]float64
var alph = []rune("ABCDEFGHIJKLMNÑOPQRSTUVWXYZ")
var es_letterprob = []float64{
	/* a    b     c     d     e     f     g     h     i */
	12.50, 1.27, 4.43, 5.14, 13.24, 0.79, 1.17, 0.81, 6.91,
	/* j    k     l     m     n     ñ     o     p     q */
	0.45, 0.08, 5.84, 2.61, 7.09, 0.22, 8.98, 2.75, 0.83,
	/* r    s     t     u     v     w     x     y     z */
	6.62, 7.44, 4.42, 4.00, 0.98, 0.03, 0.19, 0.79, 0.42,
}

func encrypt(key int, message string) string {
	s := make([]rune, utf8.RuneCountInString(message))
	for i, r := range message {
		s[i] = shiftforward(r, key)
	}
	return string(s)
}

func decrypt(message string) string {
	letterfreq(message)
	probablekey := analyzefrequencies()
	return dec(probablekey, message)
}

func letterfreq(m string) {
	msgsize := 0.0
	for _, r := range m {
		i := letterindex(r, 0)
		freqcount[i]++
		msgsize++
	}
	for i := 0; i < alphsize; i++ {
		freqcount[i] /= msgsize
	}
}

func dec(key int, message string) string {
	s := make([]rune, utf8.RuneCountInString(message))
	for i, r := range message {
		s[i] = shiftbackward(r, key)
	}
	return string(s)
}

func analyzefrequencies() int {
	probablekey := 0
	for k := 1; k < alphsize; k++ {
		ic := 0.0 /* Index of coincidence */
		for i := 0; i < alphsize; i++ {
			ic += es_letterprob[i] * freqcount[(i+k)%alphsize]
		}
		diff := es_index - ic
		if diff < mindiff {
			mindiff = diff
			probablekey = k
		}
	}
	return probablekey
}

func validatekey(key int) {
	if key >= alphsize {
		fmt.Fprintln(os.Stderr, "error: llave %d inválida", key)
		os.Exit(1)
	}
}

func shiftforward(r rune, key int) rune {
	return alph[letterindex(r, key)]
}

func letterindex(r rune, key int) int {
	if r == 'Ñ' {
		return (14 + key) % alphsize
	}
	if r >= 'o' && r <= 'z' || r >= 'O' && r <= 'Z' {
		r = r - 'a' + 1 + rune(key)
	} else if 'a' <= r && r <= 'n' || 'A' <= r && r <= 'N' {
		r = r - 'a' + rune(key)
	} else {
		fmt.Fprintln(os.Stderr, "error: rune %d inválido", r)
		os.Exit(1)
	}
	return int(r) % alphsize
}

func die(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a)
	os.Exit(1)
}

func shiftbackward(r rune, key int) rune {
	return shiftforward(r, alphsize-key)
}

func usage() {
	/* TODO print usage */
	os.Exit(1)
}

func getinput() string {
	bio := bufio.NewReader(os.Stdin)
	s, _ := bio.ReadString('\n')
	return s[:len(s)-1]
}

func main() {
	var k int
	if len(os.Args) < 2 {
		usage()
	}
	if len(os.Args) > 2 {
		k, _ = strconv.Atoi(os.Args[2])
		validatekey(k)
	}

	switch os.Args[1] {
	case "cif":
		if len(os.Args) < 3 {
			usage()
		}
		m := getinput()
		encrypt(k, m)
		fmt.Println(m)
	case "dec":
		m := getinput()
		if len(os.Args) < 3 {
			dec(k, m)
		} else {
			decrypt(m)
		}
		fmt.Println(m)
	default:
		usage()
	}
}
