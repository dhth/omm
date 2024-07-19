package cmd

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	pers "github.com/dhth/omm/internal/persistence"
	"github.com/dhth/omm/internal/types"
	"github.com/dhth/omm/internal/ui"
	"github.com/dhth/omm/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultConfigFilename   = "omm"
	envPrefix               = "OMM"
	author                  = "@dhth"
	repoIssuesUrl           = "https://github.com/dhth/omm/issues"
	defaultConfigDir        = ".config"
	defaultDataDir          = ".local/share"
	defaultConfigDirWindows = "AppData/Roaming"
	defaultDataDirWindows   = "AppData/Local"
	configFileName          = "omm/omm.toml"
	dbFileName              = "omm/omm.db"
	printTasksDefault       = 20
	taskListTitleMaxLen     = 8
)

var (
	configFileExtIncorrectErr = errors.New("config file must be a TOML file")
	configFileDoesntExistErr  = errors.New("config file does not exist")
	dbFileExtIncorrectErr     = errors.New("db file needs to end with .db")

	maxImportLimitExceededErr = fmt.Errorf("Max number of tasks that can be imported at a time: %d", pers.TaskNumLimit)
	nothingToImportErr        = errors.New("Nothing to import")

	listDensityIncorrectErr = errors.New("List density is incorrect; valid values: compact/spacious")

	taskSummaryTooLongErr = fmt.Errorf("Task summary is too long; max length allowed: %d", types.TaskSummaryMaxLen)
)

func Execute(version string) {
	rootCmd, err := NewRootCommand()

	rootCmd.Version = version
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	_ = rootCmd.Execute()
}

func setupDB(dbPathFull string) (*sql.DB, error) {

	var db *sql.DB
	var err error

	_, err = os.Stat(dbPathFull)
	if errors.Is(err, fs.ErrNotExist) {

		dir := filepath.Dir(dbPathFull)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, fmt.Errorf(`Couldn't create directory for data files: %s
Error: %s`,
				dir,
				err)
		}

		db, err = getDB(dbPathFull)

		if err != nil {
			return nil, fmt.Errorf(`Couldn't create omm's local database. This is a fatal error;
Let %s know about this via %s.

Error: %s`,
				author,
				repoIssuesUrl,
				err)
		}

		err = initDB(db)
		if err != nil {
			return nil, fmt.Errorf(`Couldn't create omm's local database. This is a fatal error;
Let %s know about this via %s.

Error: %s`,
				author,
				repoIssuesUrl,
				err)
		}
		err = upgradeDB(db, 1)
		if err != nil {
			return nil, err
		}
	} else {
		db, err = getDB(dbPathFull)
		if err != nil {
			return nil, fmt.Errorf(`Couldn't open omm's local database. This is a fatal error;
Let %s know about this via %s.

Error: %s`,
				author,
				repoIssuesUrl,
				err)
		}
		err = upgradeDBIfNeeded(db)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func NewRootCommand() (*cobra.Command, error) {

	var (
		configPath            string
		configPathFull        string
		dbPath                string
		dbPathFull            string
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

	rootCmd := &cobra.Command{
		Use:   "omm",
		Short: "omm (\"on my mind\") is a keyboard-driven task manager for the command line",
		Long: `omm ("on my mind") is a keyboard-driven task manager for the command line.

It is intended to help you visualize and arrange the tasks you need to finish,
based on the priority you assign them. The higher a task is in omm's list, the
higher priority it takes.

Tip: Quickly add a task using 'omm "task summary goes here"'.
`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			configPathFull = expandTilde(configPath)

			if filepath.Ext(configPathFull) != ".toml" {
				return configFileExtIncorrectErr
			}
			_, err := os.Stat(configPathFull)

			fl := cmd.Flags()
			if fl != nil {
				cf := fl.Lookup("config-path")
				if cf != nil && cf.Changed && errors.Is(err, fs.ErrNotExist) {
					return configFileDoesntExistErr
				}
			}

			err = initializeConfig(cmd, configPathFull)
			if err != nil {
				return err
			}

			if cmd.CalledAs() == "guide" {
				tempDir := os.TempDir()
				timestamp := time.Now().UnixNano()
				tempFileName := fmt.Sprintf("omm-%d.db", timestamp)
				tempFilePath := filepath.Join(tempDir, tempFileName)
				dbPath = tempFilePath
			}

			dbPathFull = expandTilde(dbPath)
			if filepath.Ext(dbPathFull) != ".db" {
				return dbFileExtIncorrectErr
			}

			db, err = setupDB(dbPathFull)
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) != 0 {
				if len(args[0]) > types.TaskSummaryMaxLen {
					return taskSummaryTooLongErr
				}

				err := importTask(db, args[0])
				if err != nil {
					return err
				}
				return nil
			}

			// config management
			if cmd.Flags().Lookup("editor").Changed {
				editorCmd = editorFlagInp
			} else {
				editorCmd = getUserConfiguredEditor(editorFlagInp)
			}

			var ld ui.ListDensityType
			switch listDensityFlagInp {
			case ui.CompactDensityVal:
				ld = ui.Compact
			case ui.SpaciousDensityVal:
				ld = ui.Spacious
			default:
				return listDensityIncorrectErr
			}

			if len(taskListTitle) > taskListTitleMaxLen {
				taskListTitle = taskListTitle[:taskListTitleMaxLen]
			}

			config := ui.Config{
				DBPath:                dbPathFull,
				ListDensity:           ld,
				TaskListColor:         taskListColor,
				ArchivedTaskListColor: archivedTaskListColor,
				TaskListTitle:         taskListTitle,
				TextEditorCmd:         strings.Fields(editorCmd),
				ShowContext:           showContextFlagInp,
			}

			ui.RenderUI(db, config)

			return nil
		},
	}
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "%s" .Version}}
`)

	importCmd := &cobra.Command{
		Use:   "import",
		Short: "Import tasks into omm from stdin",
		RunE: func(cmd *cobra.Command, args []string) error {

			var tasks []string
			taskCounter := 0

			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				if taskCounter > pers.TaskNumLimit {
					return maxImportLimitExceededErr
				}

				line := scanner.Text()
				line = strings.TrimSpace(line)
				if len(line) > types.TaskSummaryMaxLen {
					line = utils.Trim(line, types.TaskSummaryMaxLen)
				}

				if line != "" {
					tasks = append(tasks, line)
				}
				taskCounter++
			}

			if len(tasks) == 0 {
				return nothingToImportErr
			}

			err := importTasks(db, tasks)
			if err != nil {
				return err
			}

			return nil
		},
	}

	tasksCmd := &cobra.Command{
		Use:   "tasks",
		Short: "Output tasks tracked by omm to stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printTasks(db, printTasksNum, os.Stdout)
		},
	}

	guideCmd := &cobra.Command{
		Use:   "guide",
		Short: "Starts a guided walkthrough of omm's features",
		PreRunE: func(cmd *cobra.Command, args []string) error {

			guideErr := insertGuideTasks(db)
			if guideErr != nil {
				return fmt.Errorf(`Failed to set up a guided walkthrough.
Let %s know about this via %s.

Error: %s`, author, repoIssuesUrl, guideErr)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			if cmd.Flags().Lookup("editor").Changed {
				editorCmd = editorFlagInp
			} else {
				editorCmd = getUserConfiguredEditor(editorFlagInp)
			}
			config := ui.Config{
				DBPath:                dbPathFull,
				ListDensity:           ui.Compact,
				TaskListColor:         taskListColor,
				ArchivedTaskListColor: archivedTaskListColor,
				TaskListTitle:         taskListTitle,
				TextEditorCmd:         strings.Fields(editorCmd),
				ShowContext:           true,
				Guide:                 true,
			}

			ui.RenderUI(db, config)

			return nil
		},
	}

	ros := runtime.GOOS
	var defaultConfigPath, defaultDBPath string
	var configPathAdditionalCxt, dbPathAdditionalCxt string
	hd, err := os.UserHomeDir()

	if err != nil {
		return nil, fmt.Errorf(`Couldn't get your home directory. This is a fatal error;
use --dbpath to specify database path manually
Let %s know about this via %s.

Error: %s`, author, repoIssuesUrl, err)
	}

	switch ros {
	case "linux":
		xdgDataHome := os.Getenv("XDG_DATA_HOME")
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome != "" {
			defaultConfigPath = filepath.Join(xdgConfigHome, configFileName)
		} else {
			defaultConfigPath = filepath.Join(hd, defaultConfigDir, configFileName)
		}
		if xdgDataHome != "" {
			defaultDBPath = filepath.Join(xdgDataHome, dbFileName)
		} else {
			defaultDBPath = filepath.Join(hd, defaultDataDir, dbFileName)
		}
		configPathAdditionalCxt = "; will use $XDG_CONFIG_HOME by default, if set"
		dbPathAdditionalCxt = "; will use $XDG_DATA_HOME by default, if set"
	case "windows":
		defaultConfigPath = filepath.Join(hd, defaultConfigDirWindows, configFileName)
		defaultDBPath = filepath.Join(hd, defaultDataDirWindows, dbFileName)
	default:
		defaultConfigPath = filepath.Join(hd, defaultConfigDir, configFileName)
		defaultDBPath = filepath.Join(hd, defaultDataDir, dbFileName)
	}

	rootCmd.Flags().StringVarP(&configPath, "config-path", "c", defaultConfigPath, fmt.Sprintf("location of omm's TOML config file%s", configPathAdditionalCxt))
	rootCmd.Flags().StringVarP(&dbPath, "db-path", "d", defaultDBPath, fmt.Sprintf("location of omm's database file%s", dbPathAdditionalCxt))
	rootCmd.Flags().StringVar(&taskListColor, "tl-color", ui.TaskListColor, "hex color used for the task list")
	rootCmd.Flags().StringVar(&archivedTaskListColor, "atl-color", ui.ArchivedTLColor, "hex color used for the archived tasks list")
	rootCmd.Flags().StringVar(&taskListTitle, "title", ui.TaskListDefaultTitle, fmt.Sprintf("title of the task list, will trim till %d chars", taskListTitleMaxLen))
	rootCmd.Flags().StringVar(&listDensityFlagInp, "list-density", ui.CompactDensityVal, fmt.Sprintf("type of density for the list; possible values: [%s, %s]", ui.CompactDensityVal, ui.SpaciousDensityVal))
	rootCmd.Flags().StringVar(&editorFlagInp, "editor", "vi", "editor command to run when adding/editing context to a task")
	rootCmd.Flags().BoolVar(&showContextFlagInp, "show-context", true, "whether to start omm with a visible task context pane or not; this can later be toggled on/off in the TUI")

	tasksCmd.Flags().Uint8VarP(&printTasksNum, "num", "n", printTasksDefault, "number of tasks to print")
	tasksCmd.Flags().StringVarP(&configPath, "config-path", "c", defaultConfigPath, fmt.Sprintf("location of omm's TOML config file%s", configPathAdditionalCxt))
	tasksCmd.Flags().StringVarP(&dbPath, "db-path", "d", defaultDBPath, fmt.Sprintf("location of omm's database file%s", dbPathAdditionalCxt))

	importCmd.Flags().StringVarP(&configPath, "config-path", "c", defaultConfigPath, fmt.Sprintf("location of omm's TOML config file%s", configPathAdditionalCxt))
	importCmd.Flags().StringVarP(&dbPath, "db-path", "d", defaultDBPath, fmt.Sprintf("location of omm's database file%s", dbPathAdditionalCxt))

	guideCmd.Flags().StringVar(&editorFlagInp, "editor", "vi", "editor command to run when adding/editing context to a task")

	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(tasksCmd)
	rootCmd.AddCommand(guideCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd, nil
}

func initializeConfig(cmd *cobra.Command, configFile string) error {
	v := viper.New()

	v.SetConfigName(filepath.Base(configFile))
	v.SetConfigType("toml")
	v.AddConfigPath(filepath.Dir(configFile))

	var err error
	if err = v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	err = bindFlags(cmd, v)
	if err != nil {
		return err
	}

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) error {
	var err error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := strings.ReplaceAll(f.Name, "-", "_")

		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			fErr := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if fErr != nil {
				err = fErr
				return
			}
		}
	})
	return err
}
