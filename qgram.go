package qgram

import (
	"io"
	"sort"
	"unicode"

	"github.com/wmentor/tokens"
)

// calc qgram frequency map
func CalcMap(in io.Reader) map[string]int {
	result := make(map[string]int)

	pushRune := makeCollector(result)

	pushWord := func(word string) {
		for _, r := range word {
			pushRune(r)
		}
	}

	pushRune('_')

	wCnt := 0

	tokenizer := tokens.New(in)

	for {
		token, err := tokenizer.Token()
		if err != nil {
			break
		}

		if !checkWord(token) {
			continue
		}

		if wCnt > 0 {
			pushRune('_')
		}

		wCnt++

		pushWord(token)
	}

	pushRune('_')
	pushRune('_')

	return result
}

// all qgrams in lexial order
func QGrams(in io.Reader) []string {
	hash := CalcMap(in)
	list := make([]string, 0, len(hash))
	for ng := range hash {
		list = append(list, ng)
	}
	sort.Strings(list)
	return list
}

// first limit ordered by frequency
func Popular(in io.Reader, limit int) []string {
	hash := CalcMap(in)
	list := make([]string, 0, len(hash))
	for ng := range hash {
		list = append(list, ng)
	}

	sort.Slice(list, func(i, j int) bool {
		t1 := hash[list[i]]
		t2 := hash[list[j]]

		return t1 > t2 || t1 == t2 && list[i] < list[j]
	})

	if len(list) > limit {
		list = list[:limit]
	}

	return list
}

func mapSimilarity(m1, m2 map[string]int) float64 {
	all := map[string]bool{}
	both := map[string]bool{}

	for k := range m1 {
		all[k] = true
		if _, has := m2[k]; has {
			both[k] = true
		}
	}

	for k := range m2 {
		all[k] = true
	}

	if allS := len(all); allS > 0 {
		return float64(len(both)) / float64(allS)
	}

	return 1
}

// calc text similarity
func Similarity(in1 io.Reader, in2 io.Reader) float64 {
	return mapSimilarity(CalcMap(in1), CalcMap(in2))
}

func checkWord(word string) bool {
	if word == "-" {
		return false
	}

	for _, r := range word {
		if !unicode.IsLetter(r) && r != '-' && r != '\'' {
			return false
		}
	}

	return true
}

func makeCollector(result map[string]int) func(rune) {
	var strs [3]string

	i := 0

	pushRune := func(r rune) {
		str := string(r)

		i++

		if i >= 4 {
			result[strs[0]+str]++
		}

		strs[0] = strs[1] + str
		strs[1] = strs[2] + str
		strs[2] = str
	}

	return pushRune
}
