package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lmtani/learning-go-loadtest/internal/entities"
)

var (
	// Define styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Width(50).
			Align(lipgloss.Center)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#2D3748")).
			Padding(0, 1).
			Width(50).
			Align(lipgloss.Center)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Width(50)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981"))

	failureStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3B82F6"))

	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#4B5563")).
				Padding(0, 1)

	tableCellStyle = lipgloss.NewStyle().
			Padding(0, 1)
)

// RenderReport creates a styled output for the load test report
func RenderReport(report *entities.Report) string {
	var output strings.Builder

	// Calculate success rate
	successRate := float64(report.SuccessfulRequests) / float64(report.TotalRequests) * 100

	// Build summary box
	summaryContent := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n\n%s\n%s\n%s\n%s\n%s\n%s",
		fmt.Sprintf("Total Requests:    %s", infoStyle.Render(fmt.Sprintf("%d", report.TotalRequests))),
		fmt.Sprintf("Successful:        %s  (%0.1f%%)", successStyle.Render(fmt.Sprintf("%d", report.SuccessfulRequests)), successRate),
		fmt.Sprintf("Failed:            %s  (%0.1f%%)", failureStyle.Render(fmt.Sprintf("%d", report.FailedRequests)), 100-successRate),
		fmt.Sprintf("Total Time:        %s", infoStyle.Render(fmt.Sprintf("%.3fs", report.TotalTime.Seconds()))),

		"Time Statistics:",
		fmt.Sprintf("Min Time:          %s", infoStyle.Render(fmt.Sprintf("%.3fs", report.MinResponseTime.Seconds()))),
		fmt.Sprintf("Avg Time:          %s", infoStyle.Render(fmt.Sprintf("%.3fs", report.AvgResponseTime.Seconds()))),
		fmt.Sprintf("Max Time:          %s", infoStyle.Render(fmt.Sprintf("%.3fs", report.MaxResponseTime.Seconds()))),
		fmt.Sprintf("P50:               %s", infoStyle.Render(fmt.Sprintf("%.3fs", report.P50ResponseTime.Seconds()))),
		fmt.Sprintf("P95:               %s", infoStyle.Render(fmt.Sprintf("%.3fs", report.P95ResponseTime.Seconds()))),
	)

	summaryBox := boxStyle.Render(summaryContent)
	summaryTitle := titleStyle.Render("LOAD TEST RESULTS")
	summary := lipgloss.JoinVertical(lipgloss.Left, summaryTitle, summaryBox)

	output.WriteString(summary)
	output.WriteString("\n\n")

	// Build status code distribution table
	if len(report.StatusCodeDistribution) > 0 {
		// Table header
		statusHeader := tableHeaderStyle.Render("Status")
		countHeader := tableHeaderStyle.Render("Count")
		percentHeader := tableHeaderStyle.Render("Percentage")
		tableHeader := lipgloss.JoinHorizontal(lipgloss.Top, statusHeader, countHeader, percentHeader)

		// Table rows
		var rows []string
		statusColumnWidth := 20
		countColumnWidth := 10
		percentColumnWidth := 15

		// Apply widths to header cells
		statusHeader = tableHeaderStyle.Width(statusColumnWidth).Render("Status")
		countHeader = tableHeaderStyle.Width(countColumnWidth).Render("Count")
		percentHeader = tableHeaderStyle.Width(percentColumnWidth).Render("Percentage")
		tableHeader = lipgloss.JoinHorizontal(lipgloss.Top, statusHeader, countHeader, percentHeader)

		for code, count := range report.StatusCodeDistribution {
			percentage := float64(count) / float64(report.TotalRequests) * 100

			var statusText string
			switch code {
			case 200:
				statusText = fmt.Sprintf("%d OK", code)
			case 400:
				statusText = fmt.Sprintf("%d Bad Request", code)
			case 404:
				statusText = fmt.Sprintf("%d Not Found", code)
			case 500:
				statusText = fmt.Sprintf("%d Server Error", code)
			default:
				statusText = fmt.Sprintf("%d", code)
			}

			statusCell := tableCellStyle.Width(statusColumnWidth).Render(statusText)
			countCell := tableCellStyle.Width(countColumnWidth).Render(fmt.Sprintf("%d", count))
			percentCell := tableCellStyle.Width(percentColumnWidth).Render(fmt.Sprintf("%.1f%%", percentage))

			row := lipgloss.JoinHorizontal(lipgloss.Top, statusCell, countCell, percentCell)
			rows = append(rows, row)
		}

		tableContent := lipgloss.JoinVertical(lipgloss.Left, rows...)
		tableTitle := subtitleStyle.Render("STATUS CODE DISTRIBUTION")
		table := boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, tableHeader, tableContent))

		output.WriteString(lipgloss.JoinVertical(lipgloss.Left, tableTitle, table))
	}

	return output.String()
}
