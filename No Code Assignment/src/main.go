package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

type Flashcard struct {
    Front     string
    Back      string
    LearnRank int // Learning value (0: unlearned, 1: slightly learned, 2: nearly learned, 3: learned)
}

func main() {
    myApp := app.New()
    myWindow := myApp.NewWindow("Flashcards App")

    flashcards := loadFlashcards("flashcards.txt")

    if flashcards == nil {
        widget.NewLabel("Failed to load flashcards from file.")
        return
    }

    cardIndex := 0
    frontLabel := widget.NewLabel("")
    backLabel := widget.NewLabel("")
    backLabel.Hide()

    // Sidebar labels for card counts
    unlearnedCountLabel := widget.NewLabel(fmt.Sprintf("Unlearned: %d", countCardsByRank(flashcards, 0)))
    slightlyLearnedCountLabel := widget.NewLabel(fmt.Sprintf("Slightly Learned: %d", countCardsByRank(flashcards, 1)))
    nearlyLearnedCountLabel := widget.NewLabel(fmt.Sprintf("Nearly Learned: %d", countCardsByRank(flashcards, 2)))
    learnedCountLabel := widget.NewLabel(fmt.Sprintf("Learned: %d", countCardsByRank(flashcards, 3)))

    updateCountLabels := func() {
        unlearnedCountLabel.SetText(fmt.Sprintf("Unlearned: %d", countCardsByRank(flashcards, 0)))
        slightlyLearnedCountLabel.SetText(fmt.Sprintf("Slightly Learned: %d", countCardsByRank(flashcards, 1)))
        nearlyLearnedCountLabel.SetText(fmt.Sprintf("Nearly Learned: %d", countCardsByRank(flashcards, 2)))
        learnedCountLabel.SetText(fmt.Sprintf("Learned: %d", countCardsByRank(flashcards, 3)))

        // Check if all cards have been learned and display the "Congratulations" screen
        if countCardsByRank(flashcards, 3) == len(flashcards) {
            myWindow.SetContent(congratulationsScreen())
        }
    }

    updateLabels := func() {
        if cardIndex >= 0 && cardIndex < len(flashcards) {
            card := flashcards[cardIndex]
            frontLabel.SetText("Front: " + card.Front)
            backLabelText := "Back: " + card.Back 
            backLabel.SetText(backLabelText)
        }
    }

    // Function to update the learning value and advance to the next card
    updateLearningValue := func(inc bool) {
        if cardIndex >= 0 && cardIndex < len(flashcards) {
            card := &flashcards[cardIndex]
            if inc && card.LearnRank < 4 { // Increase learning value with an upper limit of 4
                card.LearnRank++
            } else if !inc && card.LearnRank > 0 { // Decrease learning value with a lower limit of 0
                card.LearnRank--
            }

            fmt.Printf("Card %s has a learning level of %d\n", card.Front, card.LearnRank)

            // Update the labels and counts
            updateLabels()
            updateCountLabels()

            // Advance to the next card
            cardIndex++
            if cardIndex >= len(flashcards) {
                cardIndex = 0
            }

            // Hide the back label after updating
            backLabel.Hide()
        }
    }

    goodButton := widget.NewButton("Good", func() {
        updateLearningValue(true)
    })

    badButton := widget.NewButton("Bad", func() {
        updateLearningValue(false)
    })

    revealButton := widget.NewButton("Toggle Reveal", func() {
        // Toggle the visibility of the back label
        if backLabel.Hidden {
            backLabel.Show()
        } else {
            backLabel.Hide()
        }
    })

    updateLabels()
    content := container.NewVBox(
        container.NewCenter(frontLabel),
        container.NewCenter(backLabel),
        container.NewHBox(badButton, goodButton, revealButton),
    )

    // Create the sidebar container with card counts on the left
    sidebar := container.NewVBox(
        unlearnedCountLabel,
        slightlyLearnedCountLabel,
        nearlyLearnedCountLabel,
        learnedCountLabel,
    )

    // Create a grid layout for the sidebar and content
    grid := container.NewGridWithColumns(2,
        sidebar,
        content,
    )

    myWindow.SetContent(grid)
    myWindow.ShowAndRun()
}

func loadFlashcards(filename string) []Flashcard {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return nil
    }
    defer file.Close()

    var flashcards []Flashcard
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        parts := strings.Split(line, "|::|")
        if len(parts) == 2 {
            flashcard := Flashcard{
                Front: strings.TrimSpace(parts[0]),
                Back:  strings.TrimSpace(parts[1]),
            }
            flashcards = append(flashcards, flashcard)
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading file:", err)
        return nil
    }

    return flashcards
}

func countCardsByRank(flashcards []Flashcard, rank int) int {
    count := 0
    for _, card := range flashcards {
        if card.LearnRank == rank {
            count++
        }
    }
    return count
}

func congratulationsScreen() fyne.CanvasObject {
    return container.NewCenter(
        container.NewVBox(
            widget.NewLabel("Congratulations!"),
            widget.NewLabel("You have learned all the flashcards."),
        ),
    )
}
