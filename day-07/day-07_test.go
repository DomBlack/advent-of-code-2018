package main

import (
	"testing"
)

var exampleInput = []string{
	"Step C must be finished before step A can begin.",
	"Step C must be finished before step F can begin.",
	"Step A must be finished before step B can begin.",
	"Step A must be finished before step D can begin.",
	"Step B must be finished before step E can begin.",
	"Step D must be finished before step E can begin.",
	"Step F must be finished before step E can begin.",
}

func Test_part1(t *testing.T) {
	want := "CABDFE"
	if got := part1(GraphFromStrings(exampleInput)); got != want {
		t.Errorf("part1() = %v, want %v", got, want)
	}
}

func Test_part2(t *testing.T) {
	wantWord := "CABFDE"
	wantClock := 15
	gotWord, gotClock := part2(GraphFromStrings(exampleInput), 0, 2)

	if gotClock != wantClock {
		t.Errorf("part2() = %v, want %v", gotClock, wantClock)
	}

	if gotWord != wantWord {
		t.Errorf("part2() = %v, want %v", gotWord, wantWord)
	}
}
