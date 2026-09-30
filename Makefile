.PHONY: all build serve dev test clean categories

BIN_DIR := bin
BIN := $(BIN_DIR)/site-cli
SRC := tools/site-cli/cmd/site-cli

# 테마 설정 (기본값: default, 예: make build THEME=default-light)
THEME ?= default
THEME_DIR := $(if $(filter templates/%,$(THEME)),$(THEME),templates/$(THEME))

all: build

# 1. 사이트 빌더 엔진 컴파일
$(BIN):
	@mkdir -p $(BIN_DIR)
	@echo "==> Compiling site-cli engine..."
	@cd tools/site-cli && go build -o ../../$(BIN) ./cmd/site-cli

# 2. 정적 웹사이트 일괄 컴파일 (dist/ 생성)
build: $(BIN)
	@echo "==> Building static site with theme: $(THEME_DIR)..."
	@$(BIN) build --source posts --pages pages --theme $(THEME_DIR) --output dist

# 3. 로컬 정적 사이트 서빙 (재컴파일 없이 dist/ 디렉터리 순수 서빙)
serve: $(BIN)
	@echo "==> Serving compiled static site from dist/ on http://localhost:8080..."
	@$(BIN) serve --dir dist

# 4. 실시간 감시 개발 서버 (파일 변경 시 자동 재컴파일 및 브라우저 새로고침)
dev: $(BIN)
	@echo "==> Starting dev server with live reload on http://localhost:8080 (theme: $(THEME_DIR))..."
	@$(BIN) serve --dir dist --watch --source posts --pages pages --theme $(THEME_DIR)

# 4. 카테고리별 통계 및 글 목록 조회
categories: $(BIN)
	@$(BIN) categories --source posts

# 5. 엔진 단위 테스트 실행
test:
	@echo "==> Running engine unit tests..."
	cd tools/site-cli && go test -v ./...

# 6. 빌드 산출물 정리
clean:
	@echo "==> Cleaning build artifacts..."
	rm -rf dist $(BIN_DIR)
