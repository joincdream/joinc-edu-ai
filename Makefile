.PHONY: all build serve test clean categories

BIN_DIR := bin
BIN := $(BIN_DIR)/site-cli
SRC := tools/site-cli/cmd/site-cli

all: build

# 1. 사이트 빌더 엔진 컴파일
$(BIN):
	@mkdir -p $(BIN_DIR)
	@echo "==> Compiling site-cli engine..."
	@cd tools/site-cli && go build -o ../../$(BIN) ./cmd/site-cli

# 2. 정적 웹사이트 일괄 컴파일 (dist/ 생성)
build: $(BIN)
	@echo "==> Building static site..."
	@$(BIN) build --source posts --pages pages --theme templates/default --output dist

# 3. 로컬 개발 서버 구동 (라이브 리로드 및 브라우저 감시)
serve: $(BIN)
	@echo "==> Starting local dev server on http://localhost:8080..."
	@$(BIN) serve --source posts --pages pages --theme templates/default

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
