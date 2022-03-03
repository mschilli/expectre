package expectre

import (
	"log"
	"regexp"
	"strings"
	"testing"
)

func TestTTyinScript(t *testing.T) {
	exp := New()
	err := exp.Spawn("scripts/ttyin.sh")
	if err != nil {
		panic(err)
	}
	exp.Debug = false

	if exp.Debug {
		log.Printf("Testing ttyin")
	}

	// startup message
	_ = <-exp.Stdout

	text := <-exp.Stdout
	if text != "Input: " {
		t.Fatalf("Unexpected text: '%s'\n", text)
	}

	exp.Stdin <- "blah blah\n"

	// throw away echo
	_ = <-exp.Stdout

	text = strings.TrimRight(<-exp.Stdout, " \n\r")
	if text != "Input was: blah blah" {
		t.Fatalf("Unexpected text: '%s'\n", text)
	}

	exp.Cancel()
	<-exp.Released
}

func TestWhoScript(t *testing.T) {
	exp := New()
	err := exp.Spawn("scripts/who.sh")
	if err != nil {
		panic(err)
	}
	exp.Debug = false

	if exp.Debug {
		log.Printf("Testing who")
	}

	exp.ExpectString("Who are you?")
	exp.Send("Fred\n")
	match, err := exp.ExpectString("Hello,")
	if err != nil {
		panic(err)
	}

	out := strings.TrimRight(match, " \n\r")
	if out != "Hello, Fred." {
		t.Fatalf("Unexpected text: '%s'\n", match)
	}

	exp.ExpectEOF()
	exp.Cancel()
}

func TestRegex(t *testing.T) {
	exp := New()
	err := exp.Spawn("scripts/who.sh")
	if err != nil {
		panic(err)
	}

	rex := regexp.MustCompile("(.)(.)(.) are you\\?")

	res, err := exp.ExpectRegexp(rex)

	if err != nil {
		t.Fatalf("Unexpected error on regexp match: %v\n", err)
	}

	if len(res) != 1 ||
		res[0][0] != "Who are you?" ||
		res[0][1] != "W" ||
		res[0][2] != "h" ||
		res[0][3] != "o" ||
		false {
		t.Fatalf("Unexpected res: %v\n", res)
	}

	exp.Cancel()
}
