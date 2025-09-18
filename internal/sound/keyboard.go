package sound

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/daiweiwei/lang-cli/internal/config"
)

var (
	// 用于控制声音播放的互斥锁和当前播放进程
	soundMutex    sync.Mutex
	currentCmd    *exec.Cmd
	lastPlayTime  time.Time
	playInterval  = 50 * time.Millisecond // 防抖间隔，避免过于频繁的声音播放
)

// getKeyboardSoundPath 获取键盘声音文件路径
func getKeyboardSoundPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// 如果获取用户主目录失败，回退到相对路径
		return filepath.Join("assets", "keyboard-sound.wav")
	}
	return filepath.Join(homeDir, ".lang-cli", "assets", "keyboard-sound.wav")
}

// PlayKeyboardSound 播放键盘按键声音
func PlayKeyboardSound() {
	// 检查是否启用键盘声音
	if !config.AppConfig.InputKeyboardSound {
		return
	}

	soundMutex.Lock()
	defer soundMutex.Unlock()

	// 防抖：如果距离上次播放时间太短，则跳过
	now := time.Now()
	if now.Sub(lastPlayTime) < playInterval {
		return
	}
	lastPlayTime = now

	// 停止当前正在播放的声音
	stopCurrentSound()

	soundPath := getKeyboardSoundPath()

	// 根据操作系统使用不同的播放命令
	switch runtime.GOOS {
	case "darwin": // macOS
		playSoundFile("afplay", soundPath)
	case "linux":
		playSoundFile("aplay", soundPath)
	case "windows":
		playSoundFile("powershell", "-c", fmt.Sprintf("(New-Object Media.SoundPlayer '%s').PlaySync()", soundPath))
	default:
		// 不支持的操作系统，静默忽略
		return
	}
}

// stopCurrentSound 停止当前正在播放的声音
func stopCurrentSound() {
	if currentCmd != nil && currentCmd.Process != nil {
		currentCmd.Process.Kill()
		currentCmd = nil
	}
}

// playSoundFile 播放指定的音频文件
func playSoundFile(command string, args ...string) {
	go func() {
		cmd := exec.Command(command, args...)
		
		// 更新当前播放的命令
		soundMutex.Lock()
		currentCmd = cmd
		soundMutex.Unlock()
		
		cmd.Run() // 在goroutine中异步执行，避免阻塞主线程
		
		// 播放完成后清理
		soundMutex.Lock()
		if currentCmd == cmd {
			currentCmd = nil
		}
		soundMutex.Unlock()
	}()
}

// PlayTypingSound 播放打字声音（统一使用键盘音频文件）
func PlayTypingSound(char rune) {
	// 所有字符都使用统一的键盘声音
	PlayKeyboardSound()
}

// StopAllSounds 停止所有正在播放的声音
func StopAllSounds() {
	soundMutex.Lock()
	defer soundMutex.Unlock()
	stopCurrentSound()
}



// TestSound 测试声音功能
func TestSound() error {
	fmt.Println("测试键盘声音...")
	PlayKeyboardSound()
	fmt.Println("测试打字声音...")
	PlayTypingSound('a')
	fmt.Println("声音测试完成")
	return nil
}