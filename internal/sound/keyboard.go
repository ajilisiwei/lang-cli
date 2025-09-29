package sound

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/daiweiwei/lang-cli/internal/config"
)

var (
	soundMutex  sync.Mutex
	runningCmds []*exec.Cmd
	maxRunning  = 2
)

func getKeyboardSoundPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("assets", "keyboard-sound.wav")
	}
	return filepath.Join(homeDir, ".lang-cli", "assets", "keyboard-sound.wav")
}

func PlayKeyboardSound() {
	if !config.AppConfig.InputKeyboardSound {
		return
	}

	soundPath := getKeyboardSoundPath()

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("afplay", "-q", "1", soundPath)
	case "linux":
		cmd = exec.Command("aplay", soundPath)
	case "windows":
		cmd = exec.Command("powershell", "-c", fmt.Sprintf("(New-Object Media.SoundPlayer '%s').PlaySync()", soundPath))
	default:
		return
	}

	if err := cmd.Start(); err != nil {
		return
	}

	soundMutex.Lock()
	runningCmds = append(runningCmds, cmd)
	if len(runningCmds) > maxRunning {
		stale := runningCmds[0]
		runningCmds = runningCmds[1:]
		if stale != nil && stale.Process != nil {
			_ = stale.Process.Kill()
		}
	}
	soundMutex.Unlock()

	go func() {
		_ = cmd.Wait()
		soundMutex.Lock()
		for i, running := range runningCmds {
			if running == cmd {
				runningCmds = append(runningCmds[:i], runningCmds[i+1:]...)
				break
			}
		}
		soundMutex.Unlock()
	}()
}

func PlayTypingSound(r rune) {
	PlayKeyboardSound()
}

func StopAllSounds() {
	soundMutex.Lock()
	cmds := runningCmds
	runningCmds = nil
	soundMutex.Unlock()

	for _, cmd := range cmds {
		if cmd != nil && cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	}
}

func TestSound() error {
	fmt.Println("测试键盘声音...")
	PlayKeyboardSound()
	fmt.Println("测试打字声音...")
	PlayTypingSound('a')
	fmt.Println("声音测试完成")
	return nil
}
