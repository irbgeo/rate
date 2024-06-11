package rater

import "strings"

const separator = "/"

func buildPairNames(tokens []string) map[string]struct{} {
	pairs := make(map[string]struct{}, len(tokens)*(len(tokens)-1)/2)

	for i, tokenIn := range tokens {
		if i <= len(tokens)-2 {
			for _, tokenOut := range tokens[i+1:] {
				pairName, _ := buildPairName(tokenIn, tokenOut)
				pairs[pairName] = struct{}{}
			}
		}
	}
	return pairs
}

func buildPairName(tokenIn, tokenOut string) (string, string) {
	return tokenIn + separator + tokenOut, tokenOut + separator + tokenIn
}

func parsePairName(pairName string) (string, string) {
	tokens := strings.Split(pairName, separator)
	return tokens[0], tokens[1]
}
