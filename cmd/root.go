package cmd

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	pers "github.com/dhth/omm/internal/persistence"
	"github.com/dhth/omm/internal/ui"
	"github.com/dhth/omm/internal/utils"
	"github.com/spf13/cobra"
)

const (
	author              = "@dhth"
	repoIssuesUrl       = "https://github.com/dhth/omm/issues"
	defaultDataDir      = ".local/share"
	dbFileName          = "omm/omm.db"
	printTasksDefault   = 20
	taskListTitleMaxLen = 8
)

var (
	dbPath                string
	db                    *sql.DB
	taskListColor         string
	archivedTaskListColor string
	printTasksNum         uint8
	taskListTitle         string
	listDensityFlagInp    string
	editorFlagInp         string
	editorCmd             string
	showContextFlagInp    bool
)

func die(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		die("Something went wrong: %s", err)
	}
}

func setupDB() {

	if dbPath == "" {
		die("DB path cannot be empty")
	}

	dbPathFull := expandTilde(dbPath)

	var err error

	_, err = os.Stat(dbPathFull)
	if errors.Is(err, fs.ErrNotExist) {

		dir := filepath.Dir(dbPathFull)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			die(`Couldn't create directory for data files: %s
Error: %s`,
				dir,
				err)
		}

		db, err = getDB(dbPathFull)

		if err != nil {
			die(`Couldn't create omm's local database. This is a fatal error;
Let %s know about this via %s.

Error: %s`,
				author,
				repoIssuesUrl,
				err)
		}

		err = initDB(db)
		if err != nil {
			die(`Couldn't create omm's local database. This is a fatal error;
Let %s know about this via %s.

Error: %s`,
				author,
				repoIssuesUrl,
				err)
		}
		upgradeDB(db, 1)
	} else {
		db, err = getDB(dbPathFull)
		if err != nil {
			die(`Couldn't open omm's local database. This is a fatal error;
Let %s know about this via %s.

Error: %s`,
				author,
				repoIssuesUrl,
				err)
		}
		upgradeDBIfNeeded(db)
	}
}

var rootCmd = &cobra.Command{
	Use:   "omm",
	Short: "omm (\"on my mind\") is a keyboard-driven task manager for the command line",
	Long: `omm ("on my mind") is a keyboard-driven task manager for the command line.

It is intended to help you visualize and arrange the tasks you need to finish,
based on the priority you assign them. The higher a task is in omm's list, the
higher priority it takes.

Tip: Quickly add a task using 'omm "task summary goes here"'.
`,
	Args: cobra.MaximumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		if cmd.CalledAs() == "guide" {
			tempDir := os.TempDir()
			timestamp := time.Now().UnixNano()
			tempFileName := fmt.Sprintf("omm-%d.db", timestamp)
			tempFilePath := filepath.Join(tempDir, tempFileName)
			dbPath = tempFilePath
		}

		setupDB()
	},
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 0 {
			summary := utils.Trim(args[0], ui.TaskSummaryMaxLen)
			err := importTask(db, summary)
			if err != nil {
				die("There was an error adding the task: %s", err)
			}
			return
		}

		// config management
		if cmd.Flags().Lookup("editor").Changed {
			editorCmd = editorFlagInp
		} else {
			editorCmd = getUserConfiguredEditor(editorFlagInp)
		}

		var sc bool
		var scErr error
		if cmd.Flags().Lookup("show-context").Changed {
			sc = showContextFlagInp
		} else {
			sc, scErr = getUserConfiguredShowContext(showContextFlagInp)
			if scErr != nil {
				die("%s", scErr)
			}
		}

		var ld ui.ListDensityType
		var ldErr error
		if cmd.Flags().Lookup("list-density").Changed {
			switch listDensityFlagInp {
			case ui.CompactDensityVal:
				ld = ui.Compact
			case ui.SpaciousDensityVal:
				ld = ui.Spacious
			default:
				die("--list-density is incorrect")
			}
		} else {
			ld, ldErr = getUserConfiguredListDensity(listDensityFlagInp)
			if ldErr != nil {
				die("%s", ldErr)
			}
		}

		if len(taskListTitle) > taskListTitleMaxLen {
			taskListTitle = taskListTitle[:taskListTitleMaxLen]
		}

		config := ui.Config{
			DBPath:                dbPath,
			ListDensity:           ld,
			TaskListColor:         taskListColor,
			ArchivedTaskListColor: archivedTaskListColor,
			TaskListTitle:         taskListTitle,
			TextEditorCmd:         strings.Fields(editorCmd),
			ShowContext:           sc,
		}

		ui.RenderUI(db, config)
	},
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import tasks into omm from stdin",
	Run: func(cmd *cobra.Command, args []string) {

		var tasks []string
		taskCounter := 0

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if taskCounter > pers.TaskNumLimit {
				die("Max number of tasks that can be imported at a time: %d", pers.TaskNumLimit)
			}

			line := scanner.Text()
			line = strings.TrimSpace(line)
			if len(line) > ui.TaskSummaryMaxLen {
				line = utils.Trim(line, ui.TaskSummaryMaxLen)
			}

			if line != "" {
				tasks = append(tasks, line)
			}
			taskCounter++
		}

		if len(tasks) == 0 {
			die("Nothing to import")
		}

		err := importTasks(db, tasks)
		if err != nil {
			die("There was an error importing tasks: %s", err)
		}
	},
}

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Output tasks tracked by omm to stdout",
	Run: func(cmd *cobra.Command, args []string) {

		err := printTasks(db, printTasksNum, os.Stdout)
		if err != nil {
			die("There was an error importing tasks: %s", err)
		}
	},
}

var guideCmd = &cobra.Command{
	Use:   "guide",
	Short: "Starts a guided walkthrough of omm's features",
	PreRun: func(cmd *cobra.Command, args []string) {

		guideErr := insertGuideTasks(db)
		if guideErr != nil {
			die(`Failed to set up a guided walkthrough.
Let %s know about this via %s.

Error: %s`, author, repoIssuesUrl, guideErr)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		if cmd.Flags().Lookup("editor").Changed {
			editorCmd = editorFlagInp
		} else {
			editorCmd = getUserConfiguredEditor(editorFlagInp)
		}
		config := ui.Config{
			DBPath:                dbPath,
			ListDensity:           ui.Compact,
			TaskListColor:         taskListColor,
			ArchivedTaskListColor: archivedTaskListColor,
			TaskListTitle:         taskListTitle,
			TextEditorCmd:         strings.Fields(editorCmd),
			ShowContext:           true,
			Guide:                 true,
		}

		ui.RenderUI(db, config)
	},
}

func getUserHomeDir() string {
	currentUser, err := user.Current()

	if err != nil {
		die(`Couldn't get your home directory. This is a fatal error;
use --dbpath to specify database path manually
Let %s know about this via %s.

Error: %s`, author, repoIssuesUrl, err)
	}

	return currentUser.HomeDir
}

func init() {
	ros := runtime.GOOS
	var defaultDBPath string
	var dbPathAdditionalCxt string

	switch ros {
	case "linux":
		xdgDataHome := os.Getenv("XDG_DATA_HOME")
		if xdgDataHome != "" {
			defaultDBPath = filepath.Join(xdgDataHome, dbFileName)
		} else {
			defaultDBPath = filepath.Join(getUserHomeDir(), defaultDataDir, dbFileName)
		}
		dbPathAdditionalCxt = "; will use $XDG_DATA_HOME by default, if set"
	default:
		defaultDBPath = filepath.Join(getUserHomeDir(), defaultDataDir, dbFileName)
	}

	rootCmd.Flags().StringVarP(&dbPath, "db-path", "d", defaultDBPath, fmt.Sprintf("location of omm's database file%s", dbPathAdditionalCxt))
	rootCmd.Flags().StringVar(&taskListColor, "tl-color", ui.TaskListColor, "hex color used for the task list")
	rootCmd.Flags().StringVar(&archivedTaskListColor, "atl-color", ui.ArchivedTLColor, "hex color used for the archived tasks list")
	rootCmd.Flags().StringVar(&taskListTitle, "title", ui.TaskListDefaultTitle, fmt.Sprintf("title of the task list, will trim till %d chars", taskListTitleMaxLen))
	rootCmd.Flags().StringVar(&listDensityFlagInp, "list-density", ui.CompactDensityVal, fmt.Sprintf("type of density for the list; possible values: [%s, %s]; this config property can also be set via $OMM_LIST_DENSITY", ui.CompactDensityVal, ui.SpaciousDensityVal))
	rootCmd.Flags().StringVar(&editorFlagInp, "editor", "vi", "editor command to run when adding/editing context to a task; this config property can also be set via $OMM_EDITOR/$EDITOR/$VISUAL")
	rootCmd.Flags().BoolVar(&showContextFlagInp, "show-context", true, "whether to start omm with a visible task context pane or not; this can later be toggled on/off in the TUI; this config property can also be set via $OMM_SHOW_CONTEXT")

	tasksCmd.Flags().Uint8VarP(&printTasksNum, "num", "n", printTasksDefault, "number of tasks to print")
	tasksCmd.Flags().StringVarP(&dbPath, "db-path", "d", defaultDBPath, fmt.Sprintf("location of omm's database file%s", dbPathAdditionalCxt))

	importCmd.Flags().StringVarP(&dbPath, "db-path", "d", defaultDBPath, fmt.Sprintf("location of omm's database file%s", dbPathAdditionalCxt))

	guideCmd.Flags().StringVar(&editorFlagInp, "editor", "vi", "editor command to run when adding/editing context to a task; this config property can also be set via $OMM_EDITOR/$EDITOR/$VISUAL")

	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(tasksCmd)
	rootCmd.AddCommand(guideCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
