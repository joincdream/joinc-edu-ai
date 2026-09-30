# 프로젝트 개발 에이전트 행동 규칙 (Agent Protocols)

당신은 이 저장소의 페어 프로그래밍 AI 에이전트입니다. 전체 코드를 임의로 탐색하며 컨텍스트를 낭비하지 않고, 지정된 스펙을 기반으로 정확하고 정밀하게 작업해야 합니다.

## 0. 작업 계층 분류 및 Content Fast Path (Task Layer Isolation)

요청을 수신하면 가장 먼저 다음 3가지 작업 계층 중 하나를 식별하고, 해당 계층의 허용 범위를 벗어난 도구 호출과 파일 조회를 엄격히 차단합니다:

| 작업 계층 | 대상 경로 | 허용 작업 및 엄격한 가드레일 (Hard Boundary) |
| :--- | :--- | :--- |
| **콘텐츠 계층 (Content Fast Path)**<br/>글 작성, 편집, 오타 교정, **발행(Publish)**, 드래프트 이동 | `posts/`, `docs/` | **[절대 금지] `tools/` 내부 Go 소스 코드 조회 절대 금지**<br/>**[절대 금지] `make build`, `go test` 등 빌드·테스트 명령 실행 금지**<br/>• 마크다운 본문, Frontmatter 수정 및 파일 이동(`git mv`)만 즉시 수행 (5초 이내 완료) |
| **빌더 엔진 계층 (Engine Path)**<br/>사이트 빌더 CLI 기능 개발, 버그 수정, 성능 최적화 | `tools/site-cli/` | • 제1조(OKF 스펙 참조) 및 제4조(단위 테스트/정적 빌드 검증) 필수 준수 |
| **템플릿/UI 계층 (Theme Path)**<br/>사이트 레이아웃, CSS 스타일, UI 안내 문구 수정 | `templates/` | • `templates/` 디렉터리 내 HTML, CSS, `messages.yaml`, `DESIGN.md`(Google Labs 규격) 파일만 최소 수정 |

---

## 1. OKF 스펙 최우선 참조 원칙 (Zero Wide Search)

- `tools/site-cli` 관련 작업(기능 개발, 버그 수정, 리팩토링, 성능 최적화 등)을 지시받았을 때, **전체 코드베이스나 디렉토리에 대한 임의의 광역 검색(`find`, `grep`, `ls -R` 등)을 절대 수행하지 않습니다.**
- 작업을 시작하기 전, 반드시 [`docs/okf/site-cli.md`](./docs/okf/site-cli.md)를 최우선으로 열람합니다.
- 사용자의 요구사항에 부합하는 서브시스템 ID(예: `cli_commands`, `parser_scanner`, `markdown_converter`, `template_engine`, `site_builder`, `server_watcher`, `deployer`)를 식별합니다.
- 식별된 컴포넌트에 명시된 **`target_codebase` 파일만 최소한으로 열람 및 수정**하여 컨텍스트 오염을 원천 차단합니다.

---

## 2. 작업 전 사전 설명 원칙 (Pre-Execution Self-Reflection)

코드를 수정하거나 명령을 실행하기 직전, 반드시 다음 3가지 사항을 사용자에게 텍스트로 먼저 명시해야 합니다:

1. **대상 서브시스템 및 타깃 파일**: `docs/okf/site-cli.md`에서 확인한 대상 소스 파일 경로
2. **작업 목적**: 요구사항 해결 목표 및 이유
3. **구현 세부 계획 및 부작용 방지 전략**: 수정할 함수/로직 및 기존 파이프라인 무결성 유지 방안

---

## 3. 핵심 개발 불변 원칙 (Global Constraints)

[`docs/okf/site-cli.md`](./docs/okf/site-cli.md) 및 [`docs/architecture/principles.md`](./docs/architecture/principles.md)의 원칙을 반드시 준수합니다:

- **3계층 엄격 격리**: 콘텐츠(`posts/`), 빌더 엔진(`tools/site-cli/`), 템플릿(`templates/`) 상호 오염 금지
- **100% Pure Static**: 런타임 DB나 백엔드 API 종속성을 배제하고 순수 정적 파일(HTML/CSS/JS)만 생성
- **UI 메시지 외재화**: UI 텍스트 및 안내 문구 하드코딩 금지 (`templates/<theme>/messages.yaml` 참조)
- **디자인 토큰 외재화**: UI 색상 및 스타일은 Google Labs `DESIGN.md` 스펙에 정의하고 WCAG AA 명도 대비 준수
- **Idiomatic Go**: 과도한 추상화 및 고수준 프레임워크 지양, 표준 라이브러리 및 간결하고 테스트 가능한 코드 작성

---

## 4. 결정론적 로컬 검증 및 완료 기준 (Definition of Done - Engine Path 전용)

**엔진 코드(`tools/site-cli/`) 수정을 마친 후**에는 AI의 자율적 추측에 의존하지 않고, 반드시 로컬 결정론적 검증을 직접 수행해야 합니다 (단, `posts/` 콘텐츠 작업 시에는 불필요하므로 절대 실행하지 않음):

1. **단위 테스트 실행**:
   ```bash
   cd tools/site-cli && go test -v ./...
   ```
2. **정적 빌드 검증**:
   ```bash
   make build
   ```
3. **변경 파일 최소화**: `git status`로 타깃 외 파일이 변경되지 않았는지 확인
