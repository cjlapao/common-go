package parser

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Tokenizer structure
type Tokenizer struct {
	TokenMatchers  []*TokenMatcher
	IgnoreMatchers []*TokenMatcher
}

// TokenMatcher token matcher structure
type TokenMatcher struct {
	Pattern string
	Regexp  *regexp.Regexp
	Token   int
}

// Token token structure
type Token struct {
	stringValue string
	Value       interface{}
	Type        int
}

// Tokenize tokenize string by converting it to bytes and passing the array to the tokeizeBytes function
func (t *Tokenizer) Tokenize(target string) ([]*Token, error) {
	return t.tokenizeBytes([]byte(target))
}

// Add adds token to the tokenizer
func (t *Tokenizer) Add(pattern string, token int) {
	rxp := regexp.MustCompile(pattern)
	matcher := &TokenMatcher{pattern, rxp, token}
	t.TokenMatchers = append(t.TokenMatchers, matcher)
}

// Ignore adds ignore case to the tokenizer
func (t *Tokenizer) Ignore(pattern string, token int) {
	rxp := regexp.MustCompile(pattern)
	matcher := &TokenMatcher{pattern, rxp, token}
	t.IgnoreMatchers = append(t.IgnoreMatchers, matcher)
}

// TokenizeBytes tokenizes the bytes
func (t *Tokenizer) tokenizeBytes(target []byte) ([]*Token, error) {
	result := make([]*Token, 0)
	match := true // false when no match is found
	for len(target) > 0 && match {
		match = false
		for _, m := range t.TokenMatchers {
			token := m.Regexp.Find(target)
			if len(token) > 0 {
				convValue, _ := convertValue(token, m.Token)
				parsed := Token{stringValue: strings.TrimSpace(string(token)), Value: convValue, Type: m.Token}
				result = append(result, &parsed)
				target = target[len(token):] // remove the token from the input
				match = true
				break
			}
		}
		for _, m := range t.IgnoreMatchers {
			token := m.Regexp.Find(target)
			if len(token) > 0 {
				match = true
				target = target[len(token):] // remove the token from the input
				break
			}
		}
	}

	if len(target) > 0 && !match {
		return result, errors.New("No matching token for " + string(target))
	}

	return result, nil
}

func convertValue(token []byte, tokenType int) (interface{}, error) {
	switch tokenType {
	case FilterTokenInteger:
		return strconv.Atoi(string(token))
	case FilterTokenBoolean:
		return strconv.ParseBool(string(token))
	case FilterTokenFloat:
		return strconv.ParseFloat(string(token), 10)
	case FilterTokenLiteral, FilterTokenString:
		return strings.TrimSpace(string(token)), nil
	case FilterTokenDateTime, FilterTokenDate, FilterTokenTime:
		return time.Parse(time.RFC1123, string(token))
	default:
		return strings.TrimSpace(string(token)), nil
	}
}
