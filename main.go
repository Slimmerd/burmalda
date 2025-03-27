package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type SlotMachine struct {
	symbols       []string
	weights       []int
	payouts       map[string]int
	playerBalance int
}

func NewSlotMachine() *SlotMachine {
	sm := &SlotMachine{
		symbols:       []string{"🍒", "🍋", "🍊", "🍇", " 7", "💎"},
		weights:       []int{20, 15, 10, 8, 5, 2}, // Lower numbers are rarer
		payouts:       make(map[string]int),
		playerBalance: 100,
	}

	// Set up payouts
	sm.payouts["🍒🍒🍒"] = 5
	sm.payouts["🍋🍋🍋"] = 10
	sm.payouts["🍊🍊🍊"] = 15
	sm.payouts["🍇🍇🍇"] = 20
	sm.payouts["777"] = 50
	sm.payouts["💎💎💎"] = 100
	sm.payouts["🍒🍒"] = 2
	sm.payouts["🍋🍋"] = 3
	sm.payouts["🍊🍊"] = 4
	sm.payouts["🍇🍇"] = 5
	sm.payouts["77"] = 10
	sm.payouts["💎💎"] = 20

	return sm
}

func (sm *SlotMachine) clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (sm *SlotMachine) printSlotDisplay(reels []string) {
	sm.clearScreen()
	fmt.Println("\n" + strings.Repeat("=", 40))
	fmt.Println(strings.Repeat(" ", 15) + "SLOT MACHINE")
	fmt.Println(strings.Repeat("=", 40))

	// Print reels
	fmt.Println("\n" + strings.Repeat(" ", 15) + "┌───┬───┬───┐")
	fmt.Printf("%s│%s │%s │%s │\n", strings.Repeat(" ", 15), reels[0], reels[1], reels[2])
	fmt.Println(strings.Repeat(" ", 15) + "└───┴───┴───┘\n")

	fmt.Printf("Your balance: $%d\n", sm.playerBalance)
	fmt.Println(strings.Repeat("=", 40))
}

func (sm *SlotMachine) weightedRandom() string {
	// Calculate the total weight
	totalWeight := 0
	for _, weight := range sm.weights {
		totalWeight += weight
	}

	// Generate a random number between 0 and totalWeight
	r := rand.Intn(totalWeight)

	// Find the symbol corresponding to this random value
	for i, weight := range sm.weights {
		r -= weight
		if r < 0 {
			return sm.symbols[i]
		}
	}

	// Should never reach here, but return first symbol as fallback
	return sm.symbols[0]
}

func (sm *SlotMachine) spin(betAmount int) string {
	if betAmount > sm.playerBalance {
		return "You don't have enough money for that bet!"
	}

	sm.playerBalance -= betAmount

	// Show spinning animation
	for i := 0; i < 3; i++ {
		tempReels := []string{sm.weightedRandom(), sm.weightedRandom(), sm.weightedRandom()}
		sm.printSlotDisplay(tempReels)
		time.Sleep(300 * time.Millisecond)
	}

	// Final result
	reels := []string{sm.weightedRandom(), sm.weightedRandom(), sm.weightedRandom()}
	sm.printSlotDisplay(reels)

	// Check for wins
	winAmount := sm.checkWin(reels, betAmount)

	if winAmount > 0 {
		sm.playerBalance += winAmount
		return fmt.Sprintf("You won $%d!", winAmount)
	} else {
		return "No win this time. Try again!"
	}
}

func (sm *SlotMachine) checkWin(reels []string, betAmount int) int {
	// Check for three of a kind
	if reels[0] == reels[1] && reels[1] == reels[2] {
		symbolCombo := reels[0] + reels[0] + reels[0]
		if multiplier, ok := sm.payouts[symbolCombo]; ok {
			return betAmount * multiplier
		}
	}

	// Check for two of a kind
	if reels[0] == reels[1] {
		symbolCombo := reels[0] + reels[0]
		if multiplier, ok := sm.payouts[symbolCombo]; ok {
			return betAmount * multiplier
		}
	}

	if reels[1] == reels[2] {
		symbolCombo := reels[1] + reels[1]
		if multiplier, ok := sm.payouts[symbolCombo]; ok {
			return betAmount * multiplier
		}
	}

	if reels[0] == reels[2] {
		symbolCombo := reels[0] + reels[0]
		if multiplier, ok := sm.payouts[symbolCombo]; ok {
			return betAmount * multiplier
		}
	}

	return 0
}

func (sm *SlotMachine) displayRules() {
	sm.clearScreen()
	fmt.Println("\n" + strings.Repeat("=", 40))
	fmt.Println(strings.Repeat(" ", 15) + "GAME RULES")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Println("Win multipliers:")

	for combo, multiplier := range sm.payouts {
		fmt.Printf("%s: %dx your bet\n", combo, multiplier)
	}

	fmt.Println("\nRarity (highest to lowest):")
	for i, symbol := range sm.symbols {
		fmt.Printf("%s: %d%%\n", symbol, 60/sm.weights[i])
	}

	fmt.Println("\nYou start with $100. Good luck!")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Print("\nPress Enter to return to the game...")
	fmt.Scanln()
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	slotMachine := NewSlotMachine()
	slotMachine.clearScreen()

	fmt.Println("\nWelcome to the CLI Slots Simulator!")
	fmt.Println("You start with $100.")

	for {
		fmt.Println("\nOptions:")
		fmt.Println("1. Spin the slots")
		fmt.Println("2. View rules and payouts")
		fmt.Println("3. Exit")

		fmt.Print("\nChoose an option (1-3): ")
		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			if slotMachine.playerBalance <= 0 {
				fmt.Println("You're out of money! Game over.")
				fmt.Print("Press Enter to exit...")
				fmt.Scanln()
				return
			}

			fmt.Printf("Enter your bet amount (1-%d): $", slotMachine.playerBalance)
			var betInput string
			fmt.Scanln(&betInput)
			bet, err := strconv.Atoi(betInput)

			if err != nil || bet < 1 {
				fmt.Println("Please enter a valid bet amount of at least $1.")
				continue
			}

			result := slotMachine.spin(bet)
			fmt.Println(result)

			if slotMachine.playerBalance <= 0 {
				fmt.Println("You're out of money! Game over.")
				fmt.Print("Press Enter to exit...")
				fmt.Scanln()
				return
			}

			fmt.Print("Press Enter to continue...")
			fmt.Scanln()

		case "2":
			slotMachine.displayRules()

		case "3":
			fmt.Printf("Thanks for playing! You're leaving with $%d.\n", slotMachine.playerBalance)
			return

		default:
			fmt.Println("Invalid choice. Please select 1, 2, or 3.")
		}
	}
}
