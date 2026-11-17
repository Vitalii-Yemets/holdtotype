package punct

import (
	"math"
	"strings"
	"unicode"
)

const Labels = 6

func Norm(w string) string {
	w = strings.ToLower(strings.ReplaceAll(w, "’", "'"))
	return strings.TrimFunc(w, func(r rune) bool {
		return r != '\'' && r != '_' && !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func Decide(logits []float32, last []int) ([]bool, []float64) {
	end := make([]bool, len(last))
	prob := make([]float64, len(last))
	for i, at := range last {
		row := logits[at*Labels : (at+1)*Labels]
		top, mx := 0, float64(row[0])
		for k := 1; k < Labels; k++ {
			if float64(row[k]) > mx {
				top, mx = k, float64(row[k])
			}
		}
		sum := 0.0
		p := make([]float64, Labels)
		for k := range row {
			p[k] = math.Exp(float64(row[k]) - mx)
			sum += p[k]
		}
		end[i] = top == 1 || top == 3
		prob[i] = (p[1] + p[3]) / sum
	}
	return end, prob
}
