package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
	"github.com/lmtani/learning-go-loadtest/internal/executor"
)

// Progress text style
var progressTextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FAFAFA")).
	Bold(true)

// Global progress bar instance to maintain state
var progressBar = progress.New(
	progress.WithGradient("#7D56F4", "#2D3748"),
	progress.WithWidth(50),
	progress.WithoutPercentage(),
)

// RenderProgressBar creates a styled progress bar based on the current progress
func RenderProgressBar(update executor.ProgressUpdate) string {
	// Calculate the progress percentage (0.0 to 1.0)
	percent := float64(update.CompletedRequests) / float64(update.TotalRequests)

	// Render the progress bar
	bar := progressBar.ViewAs(percent)

	// Create the progress text
	progressText := fmt.Sprintf(
		"%d/%d requests completed (%.1f%%) - Elapsed: %.1fs",
		update.CompletedRequests,
		update.TotalRequests,
		percent*100,
		update.ElapsedTime.Seconds(),
	)

	// Style the text
	styledText := progressTextStyle.Render(progressText)

	// Combine the bar and text
	return strings.Join([]string{bar, styledText}, "\n")
}

// ClearProgressBar clears the progress bar by moving the cursor up and clearing lines
func ClearProgressBar() {
	// Move cursor up 2 lines (bar + text) and clear them
	fmt.Print("\033[2A\033[K\033[K")
}
