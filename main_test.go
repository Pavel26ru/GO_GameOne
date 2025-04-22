package main

import (
	"testing"
)

// Test cases for the game
func TestGame0(t *testing.T) {
	for caseNum, commands := range game0cases {
		game := initGame()
		for _, item := range commands {
			answer := game.HandleCommand(item.command)
			if answer != item.answer {
				t.Error("case:", caseNum, item.step,
					"\n\tcmd:", item.command,
					"\n\tresult:  ", answer,
					"\n\texpected:", item.answer)
			}

		}
	}
}
