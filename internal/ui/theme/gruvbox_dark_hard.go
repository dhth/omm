package theme

const themeNameGruvboxDarkHard = "gruvbox-dark-hard"

func gruvboxDarkHard() Theme {
	theme := gruvboxDark()
	theme.Name = themeNameGruvboxDarkHard
	theme.Background = "#1d2021"

	return theme
}
