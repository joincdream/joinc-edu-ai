package server

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher 파일시스템 변경 사항을 감시하고 디바운스된 콜백을 실행합니다.
type Watcher struct {
	fsWatcher *fsnotify.Watcher
	dirs      []string
	onChange  func()
	debounce  time.Duration
	stopChan  chan struct{}
}

// NewWatcher 새 파일 감시자를 생성합니다.
func NewWatcher(dirs []string, debounce time.Duration, onChange func()) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("watcher: failed to create fsnotify watcher: %w", err)
	}

	if debounce <= 0 {
		debounce = 200 * time.Millisecond
	}

	return &Watcher{
		fsWatcher: fsWatcher,
		dirs:      dirs,
		onChange:  onChange,
		debounce:  debounce,
		stopChan:  make(chan struct{}),
	}, nil
}

// Start 감시 루프를 백그라운드 고루틴으로 시작합니다.
func (w *Watcher) Start() error {
	// 감시 대상 디렉터리 재귀 등록
	for _, dir := range w.dirs {
		if err := w.addRecursive(dir); err != nil {
			log.Printf("[WATCHER] Warning: failed to watch directory %s: %v", dir, err)
		}
	}

	go w.eventLoop()
	return nil
}

// Stop 감시자를 종료합니다.
func (w *Watcher) Stop() error {
	close(w.stopChan)
	return w.fsWatcher.Close()
}

// addRecursive 지정된 디렉터리와 하위 모든 디렉터리를 감시자에 등록합니다.
func (w *Watcher) addRecursive(rootDir string) error {
	return filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "dist" {
				return filepath.SkipDir
			}
			return w.fsWatcher.Add(path)
		}
		return nil
	})
}

// eventLoop 파일 이벤트 수신 및 디바운스 처리 루프
func (w *Watcher) eventLoop() {
	var (
		timer *time.Timer
		mu    sync.Mutex
	)

	for {
		select {
		case <-w.stopChan:
			return

		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}

			// 무시할 이벤트 필터링 (Chmod 등)
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) == 0 {
				continue
			}

			// 새 디렉터리 생성 시 감시 등록
			if event.Op&fsnotify.Create != 0 {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					_ = w.addRecursive(event.Name)
				}
			}

			// 디바운스 타이머 리셋
			mu.Lock()
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(w.debounce, func() {
				log.Printf("[WATCHER] Change detected: %s (%s) -> Triggering Live Reload", event.Name, event.Op)
				if w.onChange != nil {
					w.onChange()
				}
			})
			mu.Unlock()

		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
			log.Printf("[WATCHER] fsnotify error: %v", err)
		}
	}
}
