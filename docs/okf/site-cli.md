---
version: "0.2.0"
type: "tool_specification"
title: "site-cli OKF Specification & Harness Master Index"
project: "site-cli"
description: "Go-Native Pure Static Site Generator & GitHub Pages Deployer 마스터 인덱스"
trust_signals:
  status: "verified"
  verified_by: "system-architect"
  staleness: "fresh"
  last_verified: "2026-09-28"

global_guardrails:
  three_tier_isolation: "콘텐츠(posts/) ↔ 엔진(tools/site-cli/) ↔ 템플릿(templates/) 상호 격리"
  pure_static_only: "런타임 DB/서버 배제, 100% 순수 정적 파일(HTML/CSS/JS) 빌드"
  message_externalization: "UI 텍스트 및 안내 문구의 코드 내 하드코딩 금지 (templates/<theme>/messages.yaml 활용)"
  idiomatic_go: "과도한 인터페이스 추상화 지양, Go 표준 라이브러리 및 에러 래핑(fmt.Errorf %w) 준수"
  verification_commands:
    unit_tests: "cd tools/site-cli && go test -v ./..."
    static_build: "make build"

components:
  - id: "cli_commands"
    name: "CLI 커맨드 및 옵션 플래그 인터페이스"
    description: "Cobra 기반 CLI 서브커맨드(build, serve, deploy, categories) 파싱 및 라이프사이클 제어"
    living_specs:
      - "docs/architecture/detailed-design.md#2-cli-커맨드-계층"
      - "docs/architecture/technical-requirements.md#cli-인터페이스-요구사항"
    target_codebase:
      entrypoint: "tools/site-cli/cmd/site-cli/main.go"
      commands:
        root: "tools/site-cli/internal/cli/root.go"
        build: "tools/site-cli/internal/cli/build.go"
        serve: "tools/site-cli/internal/cli/serve.go"
        deploy: "tools/site-cli/internal/cli/deploy.go"
        categories: "tools/site-cli/internal/cli/categories.go"
    verification:
      test_cmd: "cd tools/site-cli && go test ./internal/cli/..."

  - id: "parser_scanner"
    name: "Frontmatter 파서 및 포스트 스캐너"
    description: "마크다운 YAML Frontmatter 파싱, DTO 정규화, 유효성 검증(Validation)"
    living_specs:
      - "docs/architecture/detailed-design.md#3-파서-및-스캐너-계층"
    target_codebase:
      models:
        - "tools/site-cli/internal/model/types.go"
        - "tools/site-cli/internal/model/redirect.go"
      scanner: "tools/site-cli/internal/parser/scanner.go"
      frontmatter: "tools/site-cli/internal/parser/frontmatter.go"
    verification:
      tests:
        - "tools/site-cli/internal/parser/scanner_test.go"
        - "tools/site-cli/internal/parser/frontmatter_test.go"
        - "tools/site-cli/internal/model/types_test.go"

  - id: "markdown_converter"
    name: "Markdown-to-HTML 렌더링 파이프라인"
    description: "Goldmark 기반 GFM 파싱, 하이라이팅, 목차(TOC) 및 Mermaid 코드블록 보존"
    living_specs:
      - "docs/architecture/detailed-design.md#4-마크다운-변환-계층"
    target_codebase:
      converter: "tools/site-cli/internal/markdown/converter.go"
    verification:
      tests:
        - "tools/site-cli/internal/markdown/converter_test.go"

  - id: "template_engine"
    name: "템플릿 컴파일러 및 메시지 바인딩"
    description: "Go html/template 로딩, 레이아웃/컴포넌트 합성, messages.yaml UI 문구 주입"
    living_specs:
      - "docs/architecture/detailed-design.md#5-템플릿-엔진-계층"
    target_codebase:
      engine: "tools/site-cli/internal/template/engine.go"
    verification:
      tests:
        - "tools/site-cli/internal/template/engine_test.go"

  - id: "site_builder"
    name: "정적 사이트 빌더 및 에셋 파이프라인"
    description: "색인/상세/카테고리/정적페이지 일괄 렌더링, 리다이렉트 HTML 생성, 에셋 카피"
    living_specs:
      - "docs/architecture/detailed-design.md#6-빌더-및-파이프라인-계층"
    target_codebase:
      builder: "tools/site-cli/internal/builder/builder.go"
      redirect: "tools/site-cli/internal/builder/redirect.go"
    verification:
      tests:
        - "tools/site-cli/internal/builder/builder_test.go"

  - id: "server_watcher"
    name: "로컬 개발 서버 및 라이브 감시자"
    description: "fsnotify 기반 posts/templates 디렉토리 감시, 변경 시 인메모리 증분 리빌드 및 HTTP 서빙"
    living_specs:
      - "docs/architecture/detailed-design.md#7-개발-서버-및-워처-계층"
    target_codebase:
      server: "tools/site-cli/internal/server/server.go"
      watcher: "tools/site-cli/internal/server/watcher.go"

  - id: "deployer"
    name: "GitHub Pages 배포기"
    description: "빌드 산출물(dist/)을 gh-pages 브랜치에 안전하게 고립 커밋 및 원클릭 푸시"
    living_specs:
      - "docs/architecture/detailed-design.md#8-배포기-계층"
    target_codebase:
      git_deployer: "tools/site-cli/internal/deployer/git.go"
---

# `site-cli` OKF 스펙 및 하네스 마스터 인덱스

본 문서는 **Google OKF (Open Knowledge Format v0.2)** 표준을 준수하여 작성된 `tools/site-cli`의 시스템 마스터 인덱스입니다. AI 에이전트와 엔지니어 간의 결정론적 협업을 보장하고 컨텍스트 낭비 및 광역 코드 오염을 원천 차단하기 위한 하네스 가드레일 역할을 합니다.

---

## 1. 아키텍처 개요 및 3계층 격리

`site-cli`는 콘텐츠와 디자인이 완벽히 분리된 **헤드리스(Headless) 정적 사이트 생성기**입니다.

```mermaid
flowchart LR
    classDef default font-family:Pretendard,sans-serif,font-size:13px;
    classDef content fill:#f8fafc,stroke:#64748b,color:#0f172a,rx:8px;
    classDef engine fill:#eff6ff,stroke:#2563eb,color:#1e3a8a,font-weight:bold,rx:8px;
    classDef tmpl fill:#fffbeb,stroke:#d97706,color:#78350f,rx:8px;
    classDef out fill:#f0fdf4,stroke:#059669,color:#064e3b,font-weight:bold,rx:8px;

    C["📄 <b>입력 콘텐츠</b><br/><small>posts/**/*.md<br/>pages/**/*.html</small>"]:::content
    E["⚙️ <b>코어 엔진</b><br/><small>tools/site-cli (Go)</small>"]:::engine
    T["🎨 <b>템플릿 / 메시지</b><br/><small>templates/default/<br/>messages.yaml</small>"]:::tmpl
    O["🌐 <b>100% 순수 정적 산출물</b><br/><small>dist/ (HTML, CSS, JS)</small>"]:::out

    C --> E
    T --> E
    E --> O
```

### 상위 설계 문서 링크 (Living Specs)
- 📋 [기획 사양서 (Planning)](../architecture/planning.md)
- 📐 [기술 요구사항 명세서 (TRD)](../architecture/technical-requirements.md)
- 🏗️ [상세 설계서 (Detailed Design)](../architecture/detailed-design.md)
- 🧭 [핵심 개발 원칙 (Principles)](../architecture/principles.md)
- 📊 [시장 비교 분석 및 전략적 포지셔닝 (Competitive Analysis & Strategy)](../architecture/competitive-analysis-and-strategy.md)
- ⚡ [Go Backend & Svelte 하이브리드 SSG 아키텍처 (Go & Svelte Architecture)](../architecture/go-backend-svelte-architecture.md)

---

## 2. AI 에이전트 작업 가드레일 (Agent Protocols)

AI 에이전트가 `site-cli` 관련 작업(기능 추가, 리팩토링, 버그 수정)을 수행할 때는 아래 프로토콜을 반드시 준수해야 합니다:

1. **타깃 파일 격리 원칙 (Zero Wide Search)**:
   - 전체 코드베이스를 대상으로 광역 검색(`find`, `grep`)을 무차별적으로 수행하지 않습니다.
   - 요청된 과업의 서브시스템 ID(예: `parser_scanner`, `template_engine`)를 상단 Frontmatter에서 식별하고, 지정된 **`target_codebase` 파일만 열람/수정**합니다.
2. **사전 설명 원칙 (Pre-execution Self-reflection)**:
   - 코드를 수정하기 전 반드시 다음 3가지를 명시적으로 서술합니다:
     - 대상 파일 (Target File)
     - 수정 목적 (Objective)
     - 핵심 로직 및 사이드 이펙트 방지 전략 (Implementation Details)
3. **결정론적 로컬 검증 (Deterministic Self-Validation)**:
   - 코드 수정 직후 반드시 `cd tools/site-cli && go test -v ./...`를 실행하여 기존 테스트의 무결성을 검증합니다.
   - 빌드 영향을 주는 변경일 경우 `make build`를 실행하여 `dist/` 산출물 생성 여부를 확인합니다.
4. **UI 하드코딩 엄격 금지**:
   - 새롭게 노출되는 UI 텍스트는 Go 코드나 HTML 템플릿에 하드코딩하지 않고, 반드시 `templates/<theme>/messages.yaml`에 키를 추가한 뒤 바인딩합니다.

---

## 3. 서브시스템별 책임 및 타깃 매핑 요약

| 서브시스템 ID | 책임 영역 | 핵심 타깃 소스 | 주 검증 테스트 |
| :--- | :--- | :--- | :--- |
| `cli_commands` | Cobra CLI 플래그 및 커맨드 생명주기 | `internal/cli/`, `cmd/site-cli/main.go` | `internal/cli/...` |
| `parser_scanner` | YAML Frontmatter 파싱, DTO 정규화 | `internal/parser/`, `internal/model/` | `scanner_test.go`, `frontmatter_test.go` |
| `markdown_converter` | Goldmark GFM 파싱, 코드블록 변환 | `internal/markdown/converter.go` | `converter_test.go` |
| `template_engine` | HTML 템플릿 컴파일, messages.yaml 주입 | `internal/template/engine.go` | `engine_test.go` |
| `site_builder` | 정적 사이트 일괄 렌더링, 리다이렉트 생성 | `internal/builder/` | `builder_test.go` |
| `server_watcher` | fsnotify 로컬 감시 및 라이브 서빙 | `internal/server/` | 로컬 `make serve` 검증 |
| `deployer` | gh-pages 브랜치 격리 GitOps 배포 | `internal/deployer/git.go` | GitOps dry-run 검증 |

---

## 4. 완료 기준 (Definition of Done)

모든 작업은 아래 기준을 만족해야 최종 완료로 간주됩니다:
- [ ] 관련 단위 테스트 통과 (`cd tools/site-cli && go test ./...`)
- [ ] 정적 사이트 빌드 성공 (`make build`)
- [ ] `git status` 상에 지정된 컴포넌트 외의 엉뚱한 파일 변경이 없을 것
