# [TP-7] Go Backend & Svelte 템플릿 기반 하이브리드 SSG 마이그레이션 작업 계획서

> 💡 **핵심 가치 제안 (Executive Value Proposition)**  
> **"본 마이그레이션은 런타임 서버가 없는 100% 순수 정적(Pure Static) 호스팅의 장점을 온전히 유지하면서, Go 백엔드를 '초고속 데이터 추출기'로, Svelte를 '선언적 컴포넌트 프리렌더러'로 분리·결합하여 기존 바닐라 JS 스파게티 부채를 완전히 해소하고 향후 5년 이상의 프론트엔드 유지·보수·관리성을 확보합니다."**

---

## 0. 작업 배경 및 원칙 (Strict Guardrails)

### 0.1 왜 지금 마이그레이션인가?
1. **프론트엔드 상태 폭발 방지**: `base.html`에 산재하던 절차적 DOM 조작 코드(Mermaid 줌 모달, KaTeX 바인딩, i18n 배너, 다크모드 등)를 컴포넌트 내부로 완전 격리하여 부작용(Side-effect) 없는 UI 개발 환경을 마련합니다.
2. **Tailwind CDN 부채 청산**: 브라우저 런타임에 CSS를 파싱하던 개발용 Tailwind CDN 의존성을 제거하고, 빌드 타임 컴파일을 통해 100% FOUC 없는 최적화된 번들 CSS를 생성합니다.
3. **Go 템플릿 문법의 표현력 한계 극복**: 난해한 Go `html/template` 문법 대신 Svelte의 직관적인 컴포넌트(`Header.svelte`, `PostCard.svelte`, `TOC.svelte`)로 UI를 관리합니다.

### 0.2 작업 철칙 (Hard Boundaries)
* ⚠️ **[원칙 1] 마이그레이션에만 집중 (신규 기능 추가 절대 금지)**:
  - **Pagefind 검색창 모달 바인딩 등 신규 기능은 이번 작업에서 일체 개발하지 않습니다.**
  - 현재 운영 중인 `default-light` 테마의 시각적 디자인, 기능, UI 문구, 다국어 전환(`ko`/`en`), 모달 동작을 **Svelte 컴포넌트로 1:1 완벽하게 동일하게 이식(Parity)**하는 것에만 집중합니다.
* 🛡️ **[원칙 2] 기존 환경과의 병렬 공존 및 무중단 (Zero-Downtime)**:
  - 기존 `templates/default-light/` 테마를 삭제하거나 훼손하지 않습니다.
  - 신규 `templates/default-svelte/` 테마를 독립 디렉터리로 구성하여, `config.yaml`의 `theme` 설정에 따라 언제든 기존 Go 템플릿 엔진으로 즉시 롤백 가능한 상태를 유지합니다.
* 🌐 **[원칙 3] 100% Pure Static 산출물 보장**:
  - 클라이언트 런타임 SPA(CSR) 방식이 아니며, **빌드 타임에 모든 포스트와 카테고리가 100% 완성된 HTML 파일(`dist/**/*.html`)로 사전 렌더링(SSG)**되어야 합니다.

---

## 1. 기본 정보 및 매핑

* **Jira 티켓**: [TP-7](https://joincdream.atlassian.net/browse/TP-7)
* **담당 서브시스템**:
  * `site_builder` (`tools/site-cli/internal/builder/`): 정규화된 `site-data.json` 내보내기 및 Svelte SSG 빌더 트리거 연동
  * `template_engine` (`templates/default-svelte/`): Svelte 5 + Vite 기반 신규 컴포넌트 템플릿 구축
  * `model` (`tools/site-cli/internal/model/`): 데이터 계약 DTO 검증 및 직렬화
* **타깃 파일 및 디렉터리**:
  * `tools/site-cli/internal/builder/export.go` (신규: `site-data.json` 직렬화 및 추출 로직)
  * `tools/site-cli/internal/builder/builder.go` (수정: Svelte 테마 감지 시 데이터 추출 및 Svelte 빌더 호출)
  * `templates/default-svelte/` (신규: Svelte 5 템플릿 패키지 전체)
    - `package.json`, `tsconfig.json`, `vite.config.ts`, `svelte.config.js`
    - `src/app.css` (Tailwind CSS v4 & 디자인 토큰)
    - `src/types/*.ts` (Go 데이터 계약 타입 매핑)
    - `src/state/*.svelte.ts` (Svelte 5 Runes 반응형 상태)
    - `src/components/{layout,post,content,ui}/*.svelte` (도메인/책임별 분리 컴포넌트)
    - `src/layouts/BaseLayout.svelte` (공통 HTML Shell)
    - `src/views/*.svelte` (라우트별 사전 렌더링 뷰)
    - `scripts/build-ssg.ts` (JSON 읽어 정적 HTML 파일들로 프리렌더링)
  * `config.yaml` (테마 스위치 설정)

---

## 1.1 프론트엔드 기술 스택 명세 (Frontend Tech Stack)

| 계층 (Layer) | 기술 / 라이브러리 | 버전 | 선정 이유 및 역할 |
| :--- | :--- | :--- | :--- |
| **UI Framework** | **Svelte** | `^5.0.0` (Svelte 5) | • 가상 DOM(Virtual DOM) 없는 고성능 정적 컴파일<br/>• Runes(`$state`, `$props`) 기반 직관적 반응형 상태 관리<br/>• SSR 사전 렌더링으로 100% 순수 정적 HTML 생성 |
| **Language** | **TypeScript** | `^5.5.0` | • Go `site-data.json` 계약과 1:1 매핑되는 엄격한 타입 안정성<br/>• 컴포넌트 Props 및 스키마 검증 |
| **Build Tool & Bundler** | **Vite** | `^6.0.0` | • 초고속 HMR(Hot Module Replacement) 로컬 개발 서버<br/>• SSR 및 SSG 빌드 파이프라인 공식 지원<br/>• Zero-config Rollup 기반 에셋 최적화 번들링 |
| **CSS Framework** | **Tailwind CSS** | `^4.0.0` | • `@tailwindcss/vite` 공식 플러그인 연동 (별도 PostCSS 불필요)<br/>• CSS 변수 기반 Google `DESIGN.md` 토큰 바인딩<br/>• 브라우저 런타임 CDN 제거 및 사용된 CSS만 100% 압축 빌드 |
| **Typography** | **Pretendard** | `v1.3.9` | • 한글/영문 가독성 최적화 시스템 폰트 |
| **Diagram Engine** | **Mermaid.js** | `^10.9.0` | • 본문 내 다이어그램 렌더링 및 모달 팝업 줌 인터랙션 캡슐화 |
| **Math Engine** | **KaTeX** | `^0.16.9` | • 수식(`$$`, `$`) 초고속 인라인/블록 수학 렌더링 |
| **Syntax Highlighter** | **PrismJS** | `^1.29.0` | • 마크다운 코드 블록 경량 구문 강조 스타일링 |
| **Package Manager / Runtime** | **npm** (또는 **Bun**) | Node `^20.0.0 LTS` | • 현대적 프론트엔드 패키지 의존성 관리 및 빌드 스크립트 실행 |
* **참조 문서**:
  * [docs/architecture/go-backend-svelte-architecture.md](../docs/architecture/go-backend-svelte-architecture.md) (상세 설계서)
  * [docs/architecture/principles.md](../docs/architecture/principles.md) (100% Pure Static, 3계층 격리)
  * [docs/okf/site-cli.md](../docs/okf/site-cli.md) (컴포넌트 인덱스)

---

## 2. 세부 구현 단계 (Implementation Phases)

### 📌 Phase 1: Go Backend - `site-data.json` 데이터 추출기 구현
* **목표**: 마크다운 파싱 및 역색인 결과를 정규화된 단일 JSON 계약 파일(`dist/.cache/site-data.json`)로 출력.
* **세부 작업**:
  1. `tools/site-cli/internal/builder/export.go` 구현:
     - `koPosts`, `enPosts`, `TaxonomyIndex`, `Categories`, `Messages`(`messages.yaml` 및 `messages_en.yaml`), `SiteMetadata`를 포함하는 `SiteDataBundle` 구조체 정의.
     - `ExportData(outputPath string) error` 함수 구현 (`json.MarshalIndent`).
  2. 단위 테스트 작성 (`export_test.go`):
     - 추출된 JSON이 포스트 본문(HTMLContent), TOC 트리, Frontmatter 메타데이터를 누락 없이 포함하는지 검증.

---

### 📌 Phase 2: Svelte 5 모범 사례 기반 디렉터리 및 아키텍처 설계 (`templates/default-svelte/`)

Svelte 5의 모범 사례(Runes, 단일 책임 원칙, 관심사의 분리, 타입 안전성)를 준수하여 다음과 같이 디렉터리 및 컴포넌트 구조를 설계합니다.

```text
templates/default-svelte/
├── package.json                         # Svelte 5, Vite, Tailwind CSS v4 의존성
├── tsconfig.json                        # TypeScript 엄격 모드 설정
├── svelte.config.js                     # Svelte 5 컴파일러 옵션 (Runes 활성화)
├── vite.config.ts                       # SSR 프리렌더링 및 에셋 번들러 설정
├── scripts/
│   └── build-ssg.ts                     # site-data.json 기반 정적 HTML 일괄 렌더러
│
├── src/
│   ├── app.css                          # Tailwind CSS v4 진입점 및 Pretendard 폰트
│   │
│   ├── types/                           # [계약 계층] Go site-data.json 1:1 타입 매핑
│   │   ├── site.ts                      # SiteDataBundle 전체 컨테이너 인터페이스
│   │   ├── post.ts                      # Post, Frontmatter, TOCItem 인터페이스
│   │   ├── taxonomy.ts                  # Category, Tag 인터페이스
│   │   └── messages.ts                  # messages.yaml i18n 리소스 인터페이스
│   │
│   ├── state/                           # [상태 계층] Svelte 5 Runes ($state, $derived) 모듈
│   │   ├── i18n.svelte.ts               # 현재 언어 및 번역 메시지 바인딩 헬퍼 ($t)
│   │   └── modal.svelte.ts              # 다이어그램 줌 모달 등 전역 인터랙션 상태
│   │
│   ├── components/                      # [컴포넌트 계층] 도메인/책임별 격리 컴포넌트
│   │   ├── layout/                      # 사이트 골격 및 전역 네비게이션
│   │   │   ├── Header.svelte            # 상단 GNB 및 [KO | EN] 언어 스위처
│   │   │   ├── Footer.svelte            # 하단 카피라이트 및 Joinc 링크
│   │   │   ├── LanguageBanner.svelte    # 브라우저 언어 감지 스마트 추천 배너
│   │   │   └── CategorySidebar.svelte   # 좌측 카테고리 트리 및 포스트 수 뱃지
│   │   │
│   │   ├── post/                        # 포스트 및 아티클 관련 컴포넌트
│   │   │   ├── PostCard.svelte          # 목록 그리드용 단일 포스트 카드
│   │   │   ├── PostHeader.svelte        # 상세 상단 제목, 카테고리, 작성일, 읽기시간
│   │   │   ├── PostFooter.svelte        # 상세 하단 목록 복귀 및 URL 복사 버튼
│   │   │   └── TableOfContents.svelte   # 우측 H2/H3 목차 앵커 트리
│   │   │
│   │   ├── content/                     # 마크다운 특수 렌더러 격리 컴포넌트
│   │   │   ├── MermaidViewer.svelte     # Mermaid 다이어그램 렌더링 및 줌 모달
│   │   │   └── KaTeXRenderer.svelte     # 수식 ($$, $) 렌더러
│   │   │
│   │   └── ui/                          # 재사용 가능한 마이크로 원자(Atomic) UI
│   │       ├── Badge.svelte             # 카테고리/드래프트 뱃지
│   │       ├── Modal.svelte             # 공통 팝업 모달 프레임 (ESC/배경 클릭 닫기)
│   │       └── CopyButton.svelte        # 클립보드 복사 인터랙션 버튼
│   │
│   ├── layouts/
│   │   └── BaseLayout.svelte            # 공통 HTML Shell (<head>, SEO hreflang, CSS 링크)
│   │
│   └── views/                           # [뷰 계층] 라우트별 사전 렌더링 페이지 엔트리
│       ├── HomeView.svelte              # 메인 홈 (전체 포스트 목록 그리드)
│       ├── DetailView.svelte            # 포스트 상세 본문 읽기 뷰
│       ├── CategoryView.svelte          # 카테고리별 아카이브 뷰
│       └── StaticPageView.svelte        # About 등 일반 정적 페이지 뷰
```

#### Svelte 5 핵심 설계 원칙 적용
1. **Svelte 5 Runes 도입**:
   - `export let` 대신 신규 문법인 `$props()`를 사용하여 타입 안전한 컴포넌트 인터페이스 정의.
   - 전역 상태는 `$state` 및 `$derived`를 활용한 `.svelte.ts` 모듈로 캡슐화하여, 불필요한 이벤트 버스나 스토어 보일러플레이트 제거.
2. **도메인 기반 컴포넌트 분리 (Separation of Concerns)**:
   - 거대한 단일 파일 생성을 지양하고, 레이아웃(`layout/`), 포스트(`post/`), 특수 본문 렌더러(`content/`), 원자 단위 UI(`ui/`)로 책임을 엄격히 분할.
3. **완전한 타입 안전성 (End-to-End Type Safety)**:
   - `src/types/`에 Go 백엔드의 `model.Post`, `model.Category`와 100% 일치하는 TypeScript 인터페이스를 선언하여 빌드 타임 프로퍼티 누락 방어.

---

### 📌 Phase 3: 기존 UI 레이아웃의 1:1 Svelte 컴포넌트화 (Parity Porting)
* **목표**: `default-light` 테마의 HTML/CSS/인터랙션을 단 1%의 기능 누락 없이 위 컴포넌트 구조로 분리 이식.
* **컴포넌트별 1:1 매핑 명세**:
  1. `BaseLayout.svelte`: 기존 `base.html`의 `<head>`, SEO `hreflang`, Pretendard 폰트 및 Google DESIGN.md CSS 토큰 주입.
  2. `Header.svelte`: 기존 상단 헤더, 모바일 토글 메뉴, `[KO | EN]` 언어 전환 버튼.
  3. `LanguageBanner.svelte`: 기존 인라인 JS의 브라우저 언어 감지 및 `localStorage` 닫기 상태 연동.
  4. `PostCard.svelte`: 기존 `index.html`의 포스트 카드 마크업(카테고리 뱃지, 제목, 요약, 작성일, Read More 링크).
  5. `CategorySidebar.svelte`: 좌측 카테고리 목록 네비게이션 및 포스트 수 집계 뱃지.
  6. `TableOfContents.svelte`: 기존 `detail.html`의 우측 목차 네비게이션.
  7. `MermaidViewer.svelte`: `mermaid.min.js` 렌더링 및 클릭 시 `Modal.svelte` 기반 확대 줌 모달.
  8. `KaTeXRenderer.svelte`: 본문 수식 자동 렌더링.
  9. `DetailView.svelte`: 본문, 작성자 정보, 공유/복사 버튼.
  10. `CategoryView.svelte`: 카테고리별 포스트 목록 그리드.
  11. `StaticPageView.svelte`: `pages/about.html` 및 `about.en.html` 렌더링.

---

### 📌 Phase 4: Svelte SSG 프리렌더러 스크립트 작성 (`scripts/build-ssg.ts`)
* **목표**: Go가 생성한 `site-data.json`을 읽어 디스크에 완성된 정적 HTML 파일들을 일괄 렌더링.
* **세부 작업**:
  1. `scripts/build-ssg.ts` 작성:
     - `site-data.json` 로드 및 타입 파싱.
     - 한국어 사이트 생성:
       - `dist/index.html` (HomeView.svelte SSR)
       - `dist/posts/{slug}/index.html` (DetailView.svelte SSR)
       - `dist/category/{slug}/index.html` (CategoryView.svelte SSR)
       - `dist/about/index.html` (StaticPageView.svelte SSR)
     - 영문 사이트 생성:
       - `dist/en/index.html`
       - `dist/en/posts/{slug}/index.html`
       - `dist/en/category/{slug}/index.html`
       - `dist/en/about/index.html`
     - 에셋 복사: 컴파일된 CSS/JS 번들을 `dist/assets/`로 배치.

---

### 📌 Phase 5: 파이프라인 통합 및 검증 (End-to-End Integration)
* **목표**: `make build` 한 번으로 Go 파싱 ➔ JSON 추출 ➔ Svelte 빌드가 완결되도록 통합.
* **세부 작업**:
  1. `Makefile`에 `templates/default-svelte` 의존성 설치 및 빌드 타깃 추가 (`build-svelte`).
  2. `config.yaml`에 `theme: default-svelte` 옵션 지원.
  3. 시각적 일치율(Visual Parity) 및 다국어 기능 전수 검증.

---

## 3. 검증 기준 및 완료 조건 (Definition of Done)

- [ ] **데이터 정합성**: `site-cli export-data`가 13개 한국어 포스트 및 13개 영문 포스트의 모든 메타데이터를 무결하게 JSON으로 추출할 것.
- [ ] **100% Pure Static 검증**: `dist/` 내에 생성된 모든 `.html` 파일이 자바스크립트 비활성화 환경에서도 본문 텍스트와 제목이 완벽하게 렌더링될 것.
- [ ] **기능 동일성 (1:1 Parity)**:
  - [ ] 한국어/영문 다국어 경로(`dist/` vs `dist/en/`) 유지
  - [ ] Mermaid 다이어그램 렌더링 및 클릭 확대 줌 모달 정상 동작
  - [ ] KaTeX 수식(`$$`, `$`) 렌더링 정상 동작
  - [ ] 코드 하이라이팅 및 클립보드 링크 복사 기능 정상 동작
  - [ ] 브라우저 언어 감지 추천 배너 노출 및 닫기 상태 저장 정상 동작
- [ ] **신규 기능 배제 확인**: Pagefind 등 스펙 외 신규 기능 코드가 포함되지 않았을 것.
- [ ] **롤백 무결성**: `config.yaml`에서 `theme: default-light`로 변경 시 기존 Go 템플릿으로 0.1초 만에 정상 빌드될 것.
- [ ] **단위 테스트 통과**: `cd tools/site-cli && go test ./...` All PASS.
