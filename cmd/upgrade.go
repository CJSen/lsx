package cmd

import (
	"fmt"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/CJSen/lsx/config"
	"github.com/CJSen/lsx/utils"
	"github.com/spf13/cobra"
)

const (
	maxConcurrency = 6
)

func NewUpgradeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade all commands from remote.",
		Run: func(cmd *cobra.Command, args []string) {
			upgradeCmd()
		},
	}

	return cmd
}

func upgradeCmd() {
	var num, failed int64
	l := len(commands)
	cfg := config.GlobalConfig

	ch := make(chan string, maxConcurrency)
	wg := sync.WaitGroup{}
	failedCmds := make([]string, 0)

	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for cmd := range ch {
				url := cfg.RemoteBaseUrl + fmt.Sprintf("/command/%s.md", cmd)
				path := filepath.Join(cfg.DataDir, fmt.Sprintf("%s.md", cmd))
				if err := utils.RetryDownloadFile(url, path, cmd); err != nil {
					atomic.AddInt64(&failed, 1)
					failedCmds = append(failedCmds, cmd)
				}
				atomic.AddInt64(&num, 1)
				fmt.Printf("[busy working]: upgrade command:<%d/%d> => %s\n", atomic.LoadInt64(&num), l, cmd)
			}
		}()
	}

	for _, c := range commands {
		ch <- c
	}
	close(ch)
	wg.Wait()
	fmt.Printf("[clap]: all commands are upgraded. All: %d, Failed: %d\n", num, failed)
	if failed > 0 {
		fmt.Printf("[warn]: failed commands: %v\n", failedCmds)
	}
}
