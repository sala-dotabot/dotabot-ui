package tests

import (
	"regexp"
	"testing"
)

func TestRegexp(t *testing.T) {
	message := "subscribe 70766996"
	re := regexp.MustCompile("subscribe (\\d+)")
	matches := re.FindStringSubmatch(message)

	if (len(matches) != 2) {
		t.Error("matches length is not 2")
	}
	accountId := matches[1]
	if (accountId != "70766996") {
		t.Error("accountId was not parsed")
	}
}