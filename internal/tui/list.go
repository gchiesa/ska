package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const (
	listItemMaxWidth = 30
	listMaxItems     = 5
)

// listItem implements list.Item interface for compact list display
type listItem struct {
	title string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return "" }
func (i listItem) FilterValue() string { return i.title }

// createListWidget creates a compact list model without descriptions
func createListWidget(input InteractiveInput) list.Model {
	var items []list.Item

	// Get items from static Items first, then fall back to ItemsFunc
	if len(input.Items) > 0 {
		for _, item := range input.Items {
			items = append(items, listItem{title: item})
		}
	} else if input.ItemsFunc != nil {
		for _, item := range input.ItemsFunc() {
			items = append(items, listItem{title: item})
		}
	}

	// Create a compact delegate with single-line items
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1) // Single line per item
	delegate.SetSpacing(0)
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#83ADF4")).
		BorderForeground(lipgloss.Color("#83ADF4")).
		Bold(true)
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(lipgloss.Color("#416767"))

	// Calculate visible items
	visibleItems := len(items)
	if visibleItems > listMaxItems {
		visibleItems = listMaxItems
	}
	if visibleItems < 1 {
		visibleItems = 1
	}

	// Find max item width for proper sizing
	maxItemWidth := listItemMaxWidth
	for _, item := range items {
		if len(item.(listItem).title) > maxItemWidth {
			maxItemWidth = len(item.(listItem).title)
		}
	}

	// List height = visible items only (no title, status bar, or help)
	listWidget := list.New(items, delegate, maxItemWidth+6, visibleItems)
	listWidget.SetShowTitle(false)
	listWidget.SetShowStatusBar(false)
	listWidget.SetShowHelp(false)
	listWidget.SetShowPagination(len(items) > listMaxItems) // Only show pagination if needed
	listWidget.SetFilteringEnabled(true)

	// Select the default item if provided
	if input.Default != "" {
		for i, item := range items {
			if item.(listItem).title == input.Default {
				listWidget.Select(i)
				break
			}
		}
	}

	return listWidget
}
