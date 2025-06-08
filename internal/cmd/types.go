package cmd

type taskOutputFormat uint8

const (
	taskOutputPlain taskOutputFormat = iota
	taskOutputJSON
)

func taskOutputFormats() []string {
	return []string{"plain", "json"}
}

func parseTaskOutputFormat(value string) (taskOutputFormat, bool) {
	switch value {
	case "plain":
		return taskOutputPlain, true
	case "json":
		return taskOutputJSON, true
	default:
		return taskOutputPlain, false
	}
}
