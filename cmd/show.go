package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

const (
	commandDir       = ".commands"
	commandCfg       = "config.json"
	linuxCommandJson = "linux-command.json"
	commandUrl       = commandUrlBase + "/command/%s.md"
)

type config struct {
	Dir string `json:"dir"`
}

var (
	commandPath = filepath.Join(getHomedir(), commandDir)
	configPath  = filepath.Join(getHomedir(), commandDir, commandCfg)
	defaultCfg  = config{Dir: commandPath}
)

func NewShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <command>",
		Short: "Show the specified command usage.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("[sorry]: the show command does not accept any arguments")
				return
			}
			force, _ := cmd.Flags().GetBool("force")
			showCmd(args[0], force)
		},
	}
	cmd.Flags().BoolP("force", "f", false, "force to refresh command usage from remote.")
	return cmd
}

func showCmd(cmd string, force bool) {
	cmd = strings.ToLower(cmd)

	if force {
		if err := retryDownloadCmd(cmd); err != nil {
			if errors.Is(err, ErrCommandNotFound) {
				fmt.Printf("[sorry]: could not found command <%s>\n", cmd)
				return
			}
			fmt.Printf("[sorry]: failed to download command <%s>\n", cmd)
		}
	}

	cfg, err := getConfigContent()
	if err != nil {
		fmt.Println("[sorry]: failed to get config content")
		return
	}

	// 减少文件操作
	//found := false
	//for _, str := range commands {
	//	if str == cmd {
	//		found = true
	//		break
	//	}
	//}
	//if !found {
	//	fmt.Printf("[sorry-for]: could not found command <%s>\n", cmd)
	//	return
	//}

	p := path.Join(cfg.Dir, fmt.Sprintf("%s.md", cmd))
	if !isFileExist(p) {
		err := retryDownloadCmd(cmd)
		if err != nil {
			fmt.Printf("[sorry]: failed to retrieve command <%s>\n", cmd)
			return
		}
		if errors.Is(err, ErrCommandNotFound) {
			fmt.Printf("[sorry]: could not found command <%s>\n", cmd)
			return
		}
	}

	source, err := os.ReadFile(p)
	if err != nil {
		fmt.Printf("[sorry]: failed to open file <%s>\n", p)
		return
	}
	markdown.BlueBgItalic = color.New(color.FgBlue).SprintFunc()
	result := markdown.Render(string(source), 80, 6)
	fmt.Println(string(result))
}

func getHomedir() string {
	home, _ := homedir.Expand("~")
	return home
}

func makeCmdDir(dir string) error {
	if _, err := os.Stat(dir); err != nil && !os.IsExist(err) {
		return os.Mkdir(dir, 0755)
	}
	return nil
}

func isFileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func genConfigFile() error {
	if err := makeCmdDir(commandPath); err != nil {
		return err
	}

	if !isFileExist(configPath) {
		bs, _ := json.Marshal(defaultCfg)
		return os.WriteFile(configPath, bs, 0666)
	}

	return nil
}

func getConfigContent() (*config, error) {
	if err := genConfigFile(); err != nil {
		return nil, err
	}

	bs, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cfgpath := config{}
	if err := json.Unmarshal(bs, &cfgpath); err != nil {
		return nil, err
	}

	return &cfgpath, nil
}

var (
	ErrCommandNotFound = errors.New("command not found")
)

func retryDownloadCmd(cmd string) error {
	for j := 0; j < maxRetry; j++ {
		if err := downloadCmd(cmd); err != nil {
			continue
		}
		break
	}

	return nil
}

func downloadCmd(cmd string) error {
	c, err := getConfigContent()
	if err != nil {
		return err
	}

	if err := makeCmdDir(c.Dir); err != nil {
		return err
	}

	resp, err := http.Get(fmt.Sprintf(commandUrl, cmd))
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		return ErrCommandNotFound
	}

	defer resp.Body.Close()

	content := make([]byte, 0)
	reader := bufio.NewReader(resp.Body)
	for {
		line, _, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		content = append(content, line...)
		content = append(content, []byte("\n")...)
	}

	fp := path.Join(c.Dir, fmt.Sprintf("%s.md", cmd))
	err = os.WriteFile(fp, content, 0666)
	if err != nil {
		return err
	}
	return nil
}
