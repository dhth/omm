package ui

type ListDensityType uint8

const (
	Compact ListDensityType = iota
	Spacious
)

const (
	CompactDensityVal  = "compact"
	SpaciousDensityVal = "spacious"
)

type Config struct {
	ListDensity           ListDensityType
	TaskListColor         string
	ArchivedTaskListColor string
	ContextPaneColor      string
	TaskListTitle         string
	TextEditorCmd         []string
	Guide                 bool
	DBPath                string
	ShowContext           bool
}
