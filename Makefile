.PHONY: all build serve dev test clean categories deploy install complexity complexity-cognit complexity-cyclo

BIN_DIR := bin
BIN := $(BIN_DIR)/site-cli
SRC := tools/site-cli/cmd/site-cli
INSTALL_DIR ?= $(HOME)/.local/bin

# 복잡도 임계치 기준 (모범 사례: gocognit=20, gocyclo=15)
COGNIT_THRESHOLD ?= 20
CYCLO_THRESHOLD  ?= 15

all: build

# 1. 사이트 빌더 엔진 컴파일
$(BIN):
	@mkdir -p $(BIN_DIR)
	@echo "==> Compiling site-cli engine..."
	@cd tools/site-cli && go build -o ../../$(BIN) ./cmd/site-cli

# 2. 정적 웹사이트 일괄 컴파일 (config.yaml 기반 dist/ 생성)
build: $(BIN)
	@echo "==> Building static site..."
	@$(BIN) build

# 3. 로컬 정적 사이트 서빙 (재컴파일 없이 dist/ 디렉터리 순수 서빙)
serve: $(BIN)
	@echo "==> Serving compiled static site from dist/ on http://localhost:8080..."
	@$(BIN) serve --dir dist

# 4. 실시간 감시 개발 서버 (파일 변경 시 자동 재컴파일 및 브라우저 새로고침)
dev: $(BIN)
	@echo "==> Starting dev server with live reload on http://localhost:8080..."
	@$(BIN) serve --dir dist --watch

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

# 7. GitHub Pages 배포 (테스트 및 빌드 후 원격 main 푸시)
deploy: test build
	@echo "==> Deploying to GitHub Pages (pushing to origin main)..."
	git push origin main

# 8. 로컬 환경에 site-cli 설치 (~/.local/bin)
install: $(BIN)
	@echo "==> Installing $(BIN) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@cp -f $(BIN) $(INSTALL_DIR)/
	@echo "[SUCCESS] Installed site-cli to $(INSTALL_DIR)/site-cli"

# 9. 인지 복잡도(Cognitive Complexity) 측정 (SonarSource 기준, 상위 10개 및 평균)
complexity-cognit:
	@echo "==> [Cognitive Complexity] gocognit (threshold: $(COGNIT_THRESHOLD), top 10 & avg)..."
	@cd tools/site-cli && gocognit -top 10 -avg -ignore "_test\.go" .

# 10. 순환 복잡도(Cyclomatic Complexity) 측정 (McCabe 기준, 상위 10개 및 평균)
complexity-cyclo:
	@echo "==> [Cyclomatic Complexity] gocyclo (threshold: $(CYCLO_THRESHOLD), top 10 & avg)..."
	@cd tools/site-cli && gocyclo -top 10 -avg -ignore "_test\.go" .

# 11. 복잡도 종합 측정 (Cognitive + Cyclomatic)
complexity: complexity-cognit complexity-cyclo


