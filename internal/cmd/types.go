package cmd

type showTaskOutputFormat uint8

const (
	taskOutputPlain showTaskOutputFormat = iota
	taskOutputJSON
)

func showTaskOutputFormats() []string {
	return []string{"plain", "json"}
}

func parseShowTaskOutputFormat(value string) (showTaskOutputFormat, bool) {
	switch value {
	case "plain":
		return taskOutputPlain, true
	case "json":
		return taskOutputJSON, true
	default:
		return taskOutputPlain, false
	}
}
