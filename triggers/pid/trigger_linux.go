package pid

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/variadico/noti/runstat"
)

func (t *Trigger) Run(cmdErr chan error, stats chan runstat.Result) {
	fmt.Println("running pid!")

	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			// Check if PID exists.
			fmt.Println("check if exists")

			// os.FindProcess is unreliable. It won't return an error if the
			// pid doesn't exist.
			pid := filepath.Join("/proc", fmt.Sprint(t.pid))
			_, err := os.Stat(pid)
			if os.IsNotExist(err) {
				t.stats.Err = fmt.Errorf("pid %d does not exist", t.pid)
				stats <- t.stats
				return
			} else if err != nil {
				t.stats.Err = err
				stats <- t.stats
				return
			}

			data, err := ioutil.ReadFile(filepath.Join(pid, "status"))
			if err == nil {
				t.stats.Cmd = procName(string(data))
			}

			time.Sleep(2 * time.Second)
		}
	}
}

func procName(s string) string {
	re := regexp.MustCompile(`Name:\s*([\w- ]+)\s`)
	ans := re.FindStringSubmatch(s)

	if len(ans) != 2 {
		return ""
	}

	return ans[1]
}
