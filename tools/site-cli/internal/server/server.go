package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
)

// LiveReloadScript 클라이언트에 주입되는 SSE 이벤트 수신 스크립트
const LiveReloadScript = template.HTML(`
<script>
  (function() {
    const es = new EventSource('/livereload');
    es.onmessage = function(e) {
      if (e.data === 'reload') {
        console.log('[LiveReload] Change detected, reloading page...');
        window.location.reload();
      }
    };
    es.onerror = function() {
      // 서버 재시작 시 자동 재연결 대기
      setTimeout(() => new EventSource('/livereload'), 2000);
    };
  })();
</script>
`)

// SSEHub 활성화된 브라우저 SSE 연결을 관리하고 브로드캐스트를 전송합니다.
type SSEHub struct {
	mu      sync.Mutex
	clients map[chan string]bool
}

// NewSSEHub 새 SSE 허브를 생성합니다.
func NewSSEHub() *SSEHub {
	return &SSEHub{
		clients: make(map[chan string]bool),
	}
}

// Broadcast 모든 연결된 클라이언트에게 신호를 전송합니다.
func (h *SSEHub) Broadcast(msg string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.clients {
		select {
		case ch <- msg:
		default:
			// 채널이 꽉 찬 경우 건너뜀
		}
	}
}

// ServeHTTP SSE 엔드포인트 핸들러
func (h *SSEHub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	msgChan := make(chan string, 5)

	h.mu.Lock()
	h.clients[msgChan] = true
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, msgChan)
		close(msgChan)
		h.mu.Unlock()
	}()

	notify := r.Context().Done()

	for {
		select {
		case <-notify:
			return
		case msg := <-msgChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}

// StartLocalServer 로컬 정적 파일 서빙 및 Live Reload SSE 서버 구동
func StartLocalServer(addr string, distDir string, hub *SSEHub) error {
	mux := http.NewServeMux()

	// 1. Live Reload SSE 엔드포인트
	if hub != nil {
		mux.Handle("/livereload", hub)
	}

	// 2. 정적 파일 서빙 (distDir)
	fs := http.FileServer(http.Dir(distDir))
	mux.Handle("/", fs)

	log.Printf("[SERVER] Serving site at http://%s (Root: %s)", addr, distDir)
	return http.ListenAndServe(addr, mux)
}
