# [TP-3] 테마 시스템에 Google DESIGN.md 규격 도입 및 시인성/토큰 일관성 확보 작업 계획서

> 💡 **핵심 가치 제안 (Executive Value Proposition)**  
> **"본 기능은 분산된 템플릿 스타일을 단 하나의 '디자인 진실 공급원(DESIGN.md)'으로 통합하여, 독자에게는 WCAG AAA 기준의 편안한 시인성을, 개발자에게는 단 한 파일 수정으로 테마 전체를 제어하는 생산성을, 그리고 AI 에이전트에게는 스타일 왜곡(Visual Drift) 없는 무결한 UI 생성 역량을 제공합니다."**

---

## 0. 왜 이 기능이 필요한가? (핵심 가치 및 기대 효과)

문서의 기술적 구현에 들어가기에 앞서, 본 작업이 완성되었을 때 사용자(독자), 개발자, AI에게 제공하는 **3대 핵심 가치**는 다음과 같습니다:

1. **독자(독서 경험) 관점: "눈이 편안한 장문 몰입형 시인성 보장"**
   - 딥 테크 아티클 특성상 15~20분 이상의 장시간 독서가 요구됩니다.
   - 배경과 본문 텍스트 간 명도 대비를 WCAG 권장 기준(7:1)을 훌쩍 넘는 **11.2:1 ~ 16.8:1**로 수학적으로 통제하여, 어떤 디스플레이 환경에서도 눈부심이나 가독성 저하 없이 콘텐츠에 깊이 몰입할 수 있는 최상의 독서 환경을 제공합니다.

2. **개발자/운영자 관점: "단 한 줄로 테마 전체를 제어하는 단일 진실 공급원(SSOT) 구축"**
   - 현재는 색상이나 폰트를 바꾸려면 수십 개의 HTML 파일과 인라인 Tailwind 클래스를 전수 탐색해 수정해야 합니다.
   - `DESIGN.md` 도입 후에는 단 하나의 파일 내 YAML 토큰 값만 변경하면 사이트 전역의 레이아웃, 카드, 헤더, 코드블록 스타일이 즉시 동기화되는 **'토큰 중심 정적 사이트 빌더'**로서의 극적인 유지보수 생산성을 얻게 됩니다.

3. **AI 페어 프로그래밍(에이전트) 관점: "디자인 환각과 스타일 표류(Visual Drift)의 영구적 종식"**
   - AI 에이전트에게 새 컴포넌트나 새 테마 제작을 지시할 때 발생하는 고질적인 문제(제멋대로의 hex 컬러, 불규칙한 패딩, 조화롭지 않은 모서리 둥글기 등)를 원천 차단합니다.
   - 에이전트는 기계적으로 검증된 토큰과 헌장(Do's & Don'ts)의 제약 안에서만 프론트엔드 코드를 작성하므로, **첫 번째 시도(First-Shot)부터 디자인 시스템과 완벽히 일치하는 프로덕션 품질의 UI를 생성**할 수 있습니다.

---

## 1. 기본 정보 및 매핑

* **티켓 번호**: [TP-3](https://joincdream.atlassian.net/browse/TP-3)
* **담당 서브시스템**: `template_engine` (UI/UX 템플릿 계층 및 엔진 바인딩 계층)
* **타깃 파일**:
  * `templates/default/DESIGN.md` (신규: 기본 테마 디자인 토큰 및 시인성 가이드라인)
  * `AGENTS.md` (Theme Path 가이드라인에 DESIGN.md 토큰 준수 규칙 추가)
  * `tools/site-cli/internal/model/design.go` (신규: 디자인 토큰 DTO 모델)
  * `tools/site-cli/internal/template/engine.go` (DESIGN.md Frontmatter 파싱 및 Context 주입)
  * `templates/default/base.html` (DESIGN.md CSS Custom Properties 바인딩)
* **참조 문서**:
  * [docs/okf/site-cli.md](../docs/okf/site-cli.md) (`template_engine` 서브시스템)
  * [Google Labs Code DESIGN.md](https://github.com/google-labs-code/design.md) (포맷 규격 및 린터 사양)
  * [docs/architecture/principles.md](../docs/architecture/principles.md) (Pure Static, UI 외재화, 3계층 격리 원칙)

---

## 2. 배경 및 문제 정의

* **현상**:
  - `templates/default/base.html` 및 하위 페이지에 Tailwind CSS 클래스(`bg-slate-950`, `text-slate-100`, `border-slate-800`, `bg-blue-600`)가 하드코딩되어 있음.
  - 신규 컴포넌트나 테마 작성 시 AI 에이전트의 "시각적 표류(Visual Drift)" 및 임의의 스타일 생성 발생.
  - 배경과 텍스트의 명도 대비가 부족해 시인성(WCAG AA 명도 대비)이 훼손될 위험 존재.
* **원인**:
  - UI 텍스트는 `messages.yaml`로 외재화되어 있으나, 색상/폰트/간격 등 **디자인 토큰(Design Tokens)의 정식 규격 및 시각적 거버넌스 부재**.
* **목표**:
  1. Google Labs Code의 `DESIGN.md` 포맷을 도입하여 기본 테마(`templates/default/`)의 디자인 토큰을 정형화하고 WCAG AA 명도 대비율(4.5:1 이상)을 확보.
  2. 에이전트가 템플릿 수정 시 `DESIGN.md`를 최우선 참조하도록 `AGENTS.md` 가이드라인 연계.
  3. `tools/site-cli` 템플릿 엔진에서 `DESIGN.md`의 토큰을 파싱하여 템플릿에 동적 주입하는 Token-Driven 빌드 파이프라인 수립.

---

## 3. DESIGN.md에서 구체적으로 정의할 핵심 내용 (Design System Specification)

`templates/default/DESIGN.md`는 단순한 변수 목록이 아니라, **"엔터프라이즈 딥 테크 아키텍처 저널"**이라는 브랜드 아이덴티티와 **시인성(Accessibility)**을 보장하기 위해 다음 5대 영역을 명확히 정의합니다:

### 1) 브랜드 정체성 및 톤앤매너 (Brand Identity & Emotional Anchor)
- **지향점**: Architectural Precision & Journalistic Gravitas (건축적 정밀함과 기술 저널의 진중함).
- **톤앤매너**: 화려하거나 유치한 장식 요소를 배제하고, 장시간 고도의 기술 문서를 읽어도 눈의 피로가 없는 **매트한 다크 슬레이트 캔버스**와 **절제된 테크 블루 액센트**를 지향.

### 2) 색상 및 시인성 계층 (Colors & WCAG Contrast Hierarchy)
배경과 텍스트의 명도 대비를 철저히 통제하여 어떤 환경에서도 시인성을 확보합니다:

* **표면 및 배경 (Surfaces)**:
  - `neutral` (`#020617` / Slate 950): 완전한 블랙(#000000)의 눈부심과 대비 피로를 피하기 위한 깊은 슬레이트 캔버스.
  - `surface` (`#0f172a` / Slate 900): 헤더 바, 카드 배경, 사이드바 레이어.
  - `surface-elevated` (`#1e293b` / Slate 800): 코드 블록, 모달, 호버 하이라이트 영역.
  - `border-subtle` (`#334155` / Slate 700): 과도하지 않은 1px 구조적 분리선.
* **텍스트 시인성 계층 (Typography Contrast - WCAG AAA/AA 충족)**:
  - `text-primary` (`#f8fafc` / Slate 50): H1~H3 제목 및 핵심 수치. 대비율 **16.8:1** (WCAG AAA 기준 7:1 초과).
  - `text-body` (`#cbd5e1` / Slate 300): 긴 호흡의 본문 아티클 텍스트. 대비율 **11.2:1** (장시간 독서에 최적화된 부드러운 고대비).
  - `text-muted` (`#94a3b8` / Slate 400): 작성일, 카테고리 메타데이터, 보조 설명. 대비율 **6.4:1** (WCAG AA 기준 4.5:1 통과).
* **인터랙션 & 액센트 (Interactive & Brand)**:
  - `primary` (`#3b82f6` / Blue 500): 주 CTA 버튼, 진행 바, 활성 메뉴 인디케이터.
  - `primary-hover` (`#60a5fa` / Blue 400): 마우스 호버 및 포커스 링.
  - `link` (`#60a5fa`): 본문 하이퍼링크. 밑줄 오프셋(4px)을 필수로 동반하여 본문 텍스트와 시각적 명확 분리.
* **시맨틱 상태 (Status Alerts)**:
  - Info: `#38bdf8` | Success: `#34d399` | Warning: `#fbbf24` | Danger: `#f87171`

### 3) 타이포그래피 및 읽기 리듬 (Typography & Reading Immersion)
- **폰트 패밀리**:
  - UI 및 본문: 한글/영문 자폭 균형과 모니터 가독성이 입증된 **Pretendard**
  - 코드 및 수식: 고정폭 정렬과 식별성이 뛰어난 **JetBrains Mono / ui-monospace**
- **계층별 스케일 및 리듬**:
  - `h1`: 2.25rem(36px) | weight 700 | line-height 1.2 | letter-spacing -0.025em
  - `h2`: 1.5rem(24px) | weight 700 | line-height 1.3 | letter-spacing -0.02em | border-bottom 1px
  - `h3`: 1.25rem(20px) | weight 600 | line-height 1.4
  - `body-md`: 1.05rem(16.8px) | weight 400 | **line-height 1.75** (장문 기술 아티클의 최적 행간)
  - `code`: 0.875rem(14px) | weight 500 | line-height 1.5
- **가독 폭 제약**: 장문 독서 시 시선 이동 피로를 방지하기 위해 본문 최대 너비는 **`48rem (768px, prose-lg)`** 로 엄격히 제한.

### 4) 핵심 컴포넌트 토큰 (Component Tokens)
자주 재사용되는 UI 컴포넌트의 토큰 매핑 규칙을 명시:
- `header`: 높이 `4rem(64px)`, 배경 `surface/80` (backdrop-blur-md), 하단 보더 `border-subtle`
- `card-post`: 배경 `surface`, 보더 `border-subtle`, 라운딩 `md(12px)`, 호버 시 보더 `accent-primary/50` 강조
- `code-block`: 배경 `surface`, 보더 `border-subtle`, 라운딩 `sm(8px)`, 코드 컬러 `#38bdf8`
- `toc-navigation`: 우측 플로팅, 폰트 `0.875rem`, 활성 헤딩 링크에 `primary` 좌측 보더(2px) 표시
- `tag-pill`: 라운딩 `full(9999px)`, 배경 `primary/10`, 텍스트 `#93c5fd`, 패딩 `4px 12px`

### 5) 디자인 헌장 (Do's and Don'ts - AI 에이전트 행동 가드레일)
에이전트가 템플릿 코드를 작성할 때 자의적 판단을 원천 차단하는 불변 규칙:
- **DO**: 본문 텍스트와 배경은 최소 7:1 (WCAG AAA) 이상의 대비를 유지해야 한다.
- **DO**: 본문 하이퍼링크는 반드시 언더라인(또는 확실한 시각 피드백)을 동반하여 일반 텍스트와 구별해야 한다.
- **DO**: 모든 색상, 라운딩, 패딩은 `DESIGN.md`에 정의된 토큰 및 Tailwind 표준 스케일만 사용해야 한다.
- **DON'T**: 본문 일반 텍스트에 파란색(`text-blue-*`)을 임의로 적용하지 않는다 (링크로 오인 방지).
- **DON'T**: 배경에 순수 블랙(`black`, `#000000`)이나 텍스트에 순수 화이트를 100% 면적으로 직접 충돌시키지 않는다.
- **DON'T**: Tailwind 임의값(예: `bg-[#123456]`, `p-[13px]`)을 코드에 하드코딩하지 않는다.

---

## 4. 단계별 세부 구현 계획

### Phase 1. 기본 테마 `templates/default/DESIGN.md` 작성 및 가이드라인 확립
1. **`templates/default/DESIGN.md` 파일 생성**:
   - 위 3장의 상세 스펙에 맞춰 YAML Front Matter(토큰 그룹) 및 Markdown Prose(가이드라인) 완벽 작성.
2. **`npx @google/design.md lint` 자동 검증**:
   - 토큰 파싱 무결성 및 WCAG AA 명도 대비율 통과 여부 검증.
3. **`AGENTS.md` Theme Path 규칙 갱신**:
   - 템플릿 작업 시 `templates/<theme>/DESIGN.md`를 최우선 준수하도록 강제화.

### Phase 2. `tools/site-cli` 템플릿 엔진 토큰 바인딩 확장 (Token-Driven Theming)
1. **디자인 토큰 모델 DTO 정의 (`tools/site-cli/internal/model/design.go`)**:
   - YAML Front Matter를 파싱할 `DesignTokens`, `ColorTokens`, `TypographyTokens` 구조체 구현.
2. **`template.Engine` 로더 확장 (`tools/site-cli/internal/template/engine.go`)**:
   - `themeDir/DESIGN.md` 존재 시 YAML Frontmatter 파싱 및 `TemplateContext.Design`에 주입.
   - `DESIGN.md`가 없는 테마도 기존처럼 작동하도록 Fallback 처리.
3. **`templates/default/base.html` CSS 변수 연동**:
   - 템플릿 컨텍스트의 `Design` 값을 기반으로 `:root` CSS Custom Properties 자동 주입.

### Phase 3. 3계층 격리 및 검증
* **Pure Static 불변성**: 런타임 종속 없이 빌드 타임에만 파싱되어 순수 정적 파일(HTML/CSS)로 빌드.
* **Go 표준 라이브러리 준수**: `gopkg.in/yaml.v3`를 사용하여 엔진 무의존성 유지.

---

## 5. 코딩 가드레일 및 엄격한 격리 경계 (Hard Boundaries)

본 작업은 무분별한 코드베이스 수정을 차단하기 위해 **지정된 5개 파일 외의 어떤 코드도 수정하지 않는 엄격한 하드 바운더리**를 적용합니다:

### 1) 수정 및 생성 허용 파일 목록 (총 5개)
| 구분 | 파일 경로 | 작업 계층 | 수정/생성 상세 내용 |
| :---: | :--- | :--- | :--- |
| **신규** | `templates/default/DESIGN.md` | `Theme Path` | 3장에 정의된 색상/타이포/라운딩 YAML 토큰 및 Markdown 디자인 헌장(Do's & Don'ts) 작성 |
| **신규** | `tools/site-cli/internal/model/design.go` | `Engine Path` | `DESIGN.md` YAML Frontmatter를 파싱할 DTO 구조체 (`DesignTokens`, `ColorTokens` 등) 정의 |
| **수정** | `tools/site-cli/internal/model/types.go` | `Engine Path` | `TemplateContext` 구조체에 `Design *DesignTokens` 필드 1줄 추가 |
| **수정** | `tools/site-cli/internal/template/engine.go` | `Engine Path` | 테마 로딩 시 `DESIGN.md` 존재 확인 후 파싱하여 `ctx.Design`에 주입 (부재 시 Fallback 처리) |
| **수정** | `templates/default/base.html` | `Theme Path` | `<head>` 내에 `{{ if .Design }}` CSS Custom Properties(`:root { --color-... }`) 바인딩 블록 추가 |
| **수정** | `AGENTS.md` | `Root Rules` | Theme Path 가이드라인에 `templates/<theme>/DESIGN.md` 토큰 준수 규칙 1줄 추가 |

### 2) 함수 및 블록 단위 수정 지점 제한
* **`model/types.go`**: `TemplateContext` 구조체 내 필드 1개(`Design *DesignTokens`) 추가 외에 기존 코드 일체 불변.
* **`template/engine.go`**: 
  - `NewEngine`: `messages.yaml` 로드 직후 `DESIGN.md` 파일 존재 검사 및 언마샬링 로직만 추가.
  - `RenderPage` / `RenderCustomTemplate`: `ctx.Design = e.design` 1줄 바인딩만 추가.
  - 기존 `loadTemplate`, `funcMap`, 템플릿 컴파일 및 렌더링 로직 수정 금지.
* **`templates/default/base.html`**:
  - `<head>` 내부 CSS 변수 바인딩 블록 추가 외에 기존 HTML 마크업, body 구조, 외부 CDN 링크 수정 금지.

### 3) 절대 수정 금지 영역 (Hard Boundary)
* ❌ **콘텐츠 계층 (`posts/`)**: 기존 블로그 포스트 및 에셋 파일 일체 수정 금지.
* ❌ **마크다운 변환기 (`tools/site-cli/internal/markdown/`)**: Goldmark AST 변환 로직 수정 금지.
* ❌ **포스트 스캐너 (`tools/site-cli/internal/parser/`)**: Frontmatter 파서 및 디렉터리 워커 수정 금지.
* ❌ **빌더 파이프라인 (`tools/site-cli/internal/builder/`)**: 빌드 오케스트레이션 및 리다이렉트 생성 로직 수정 금지.
* ❌ **개별 페이지 템플릿 (`index.html`, `detail.html`, `category.html`, `page.html`)**: 이번 스코프에서는 `base.html`의 CSS 변수 주입만 다루며 개별 템플릿 마크업은 건드리지 않음.

---

## 6. 완료 기준 (Definition of Done)

- [ ] `task/TP-3_design_md_adoption_plan.md` 계획서 작성 및 저장
- [ ] `templates/default/DESIGN.md` 작성 및 `npx @google/design.md lint` 통과
- [ ] `AGENTS.md`에 `DESIGN.md` 참조 규칙 반영
- [ ] `tools/site-cli/internal/template/engine_test.go` 단위 테스트 추가 및 통과 (`cd tools/site-cli && go test -v ./...`)
- [ ] 정적 사이트 컴파일 검증 (`make build`)
- [ ] `git status`로 위 5개 대상 파일 외의 변경이 없음을 엄격 확인
- [ ] Jira 티켓 [TP-3](https://joincdream.atlassian.net/browse/TP-3) 코멘트 등록 및 완료 처리

