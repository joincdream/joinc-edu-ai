# [TP-6] Post 및 정적 페이지(About 등), 사이트 템플릿의 영문 배포 체계 구축 작업 계획서

> 💡 **핵심 가치 제안 (Executive Value Proposition)**  
> **"본 기능은 단일 저장소와 100% 순수 정적(Pure Static) 빌드 파이프라인의 간결함을 유지하면서, 기술 포스트, 독립 페이지(About 등), 사이트 테마 UI를 영문으로 서브패스(/en/)에 완벽히 동시 배포하여 글로벌 기술 생태계로의 독자층 확장을 실현합니다."**

---

## 0. 왜 이 기능이 필요한가? (핵심 가치 및 기대 효과)

1. **글로벌 독자 및 테크 커뮤니티 관점: "글로벌 접근성 확보 및 언어 장벽 해소"**
   - 딥 테크 아키텍처 및 엔터프라이즈 AI 관련 인사이트를 국내뿐 아니라 글로벌 엔지니어, 오픈소스 기여자들에게 제공합니다.
   - 단일 도메인(`www.joinc.co.kr`) 하에서 직관적인 언어 전환(Language Switcher)과 다국어 SEO(`hreflang`)를 통해 해외 검색 엔진 유입을 극대화합니다.

2. **개발자/운영자 관점: "단일 GitOps 파이프라인으로 무중단 멀티 언어 서빙"**
   - 서브도메인이나 별도 저장소 분리 없이 단일 GitHub Actions 워크플로우로 `dist/` 및 `dist/en/`을 일괄 컴파일합니다.
   - 기존의 한국어 단일 빌드 경로와 100% 하위 호환성을 유지하여 기존 콘텐츠 운영에 전혀 지장을 주지 않습니다.

3. **AI 페어 프로그래밍(에이전트) 관점: "결정론적 다국어 번역 및 배포 하네스 확립"**
   - 마크다운 Frontmatter와 다이어그램, 코드블록의 구조를 보존하는 일관된 영문 디렉터리(`posts/en/`, `pages/*.en.html`) 규칙을 확립합니다.
   - AI 에이전트가 새 글 작성 시 즉시 영문 버전을 생성하고 빌드 검증까지 원스톱으로 수행할 수 있는 기반을 제공합니다.

---

## 1. 기본 정보 및 매핑

* **티켓 번호**: [TP-6](https://joincdream.atlassian.net/browse/TP-6)
* **담당 서브시스템**:
  * `template_engine` (템플릿 및 messages.yaml UI 다국어화, 언어 스위처 바인딩)
  * `site_builder` (멀티 랭귀지 다국어 정적 빌드 오케스트레이션)
  * `parser_scanner` (언어별 디렉터리 스캔 및 메타데이터 정규화)
* **타깃 파일**:
  * `config.yaml` (다국어 언어 목록 및 매핑 스키마 추가)
  * `templates/default-light/messages_en.yaml` (신규: 영문 UI 텍스트 리소스 번들)
  * `templates/default-light/base.html` (수정: `html lang`, `BaseURL` 연동, 언어 스위처 UI, `hreflang` 태그)
  * `pages/about.en.html` (신규: 영문 About 페이지)
  * `posts/*.en.md` (신규: 영문 번역 포스트 - 원문과 동일 디렉터리에 `.en.md` 접미사로 나란히 배치)
  * `tools/site-cli/internal/model/types.go` (수정: `TemplateContext`에 다국어 메타데이터 필드 추가)
  * `tools/site-cli/internal/template/engine.go` (수정: 언어별 메시지 번들 다중 로딩 및 전환 지원)
  * `tools/site-cli/internal/parser/frontmatter.go` (수정: `ExtractSlug`에서 `.en.md` 확장자 지원)
  * `tools/site-cli/internal/parser/scanner.go` (수정: 언어별 포스트 필터링 지원)
  * `tools/site-cli/internal/builder/builder.go` (수정: 다국어 순차 빌드 및 `/en/` 서브패스 산출물 생성)
* **참조 문서**:
  * [docs/okf/site-cli.md](../docs/okf/site-cli.md) (`site_builder`, `template_engine`, `parser_scanner` 서브시스템)
  * [docs/architecture/principles.md](../docs/architecture/principles.md) (100% Pure Static, UI 메시지 외재화, 3계층 엄격 격리)
  * [AGENTS.md](../AGENTS.md) (Agent Protocols 및 완료 기준)

---

## 2. 배경 및 문제 정의

* **현상**:
  - 현재 모든 블로그 포스트와 `pages/about.html`은 한국어로만 작성되고 배포됨.
  - `templates/default-light/base.html`에 `<html lang="ko">` 및 홈 링크(`<a href="/">`), 메뉴 링크(`<a href="/about/">`)가 한국어 기준으로 고정되어 있음.
  - UI 텍스트 리소스가 `templates/default-light/messages.yaml` 단일 파일로만 존재함.
  - `site-cli`는 단일 `posts` 디렉터리를 스캔하여 `dist/` 루트에만 HTML을 생성함.
* **원인**:
  - 초기 사이트 빌더 엔진이 단일 언어(한국어) 환경만을 전제로 설계되어 다국어 컨텍스트 주입 및 서브패스 빌드 파이프라인 부재.
* **목표**:
  1. 단일 빌드 명령(`make build`)으로 한국어 기본 사이트(`dist/`)와 영문 사이트(`dist/en/`)를 동시에 생성하는 정적 빌드 파이프라인 구축.
  2. 포스트 및 About 페이지의 영문화 지원 (`posts/*.en.md`, `pages/about.en.html` -> `dist/en/about/index.html`).
  3. UI 템플릿(헤더 네비게이션, 푸터, TOC, 카드 등)에 영문 메시지(`messages_en.yaml`) 주입 및 상호 언어 전환 스위처(`[ KO | EN ]`) 구현.
  4. 구글 검색엔진 최적화를 위한 `<link rel="alternate" hreflang="...">` 메타 태그 자동 주입.

---

## 3. 구체적 아키텍처 및 상세 사양 (Detailed Architecture Specification)

### 1) URL 및 서브패스 라우팅 사양 (GitHub Pages 최적화)
단일 도메인(`www.joinc.co.kr`) 하에서 서브패스 기반의 깔끔한 URL 체계를 운영합니다:

| 페이지 유형 | 한국어 (기본) URL | 영문 (EN) URL | 산출물 파일 경로 |
| :--- | :--- | :--- | :--- |
| **메인 홈** | `www.joinc.co.kr/` | `www.joinc.co.kr/en/` | `dist/index.html`<br/>`dist/en/index.html` |
| **소개 페이지** | `www.joinc.co.kr/about/` | `www.joinc.co.kr/en/about/` | `dist/about/index.html`<br/>`dist/en/about/index.html` |
| **포스트 상세** | `www.joinc.co.kr/posts/{slug}/` | `www.joinc.co.kr/en/posts/{slug}/` | `dist/posts/{slug}/index.html`<br/>`dist/en/posts/{slug}/index.html` |
| **카테고리** | `www.joinc.co.kr/category/{slug}/` | `www.joinc.co.kr/en/category/{slug}/` | `dist/category/{slug}/index.html`<br/>`dist/en/category/{slug}/index.html` |
| **공통 에셋** | `www.joinc.co.kr/assets/...` | `www.joinc.co.kr/assets/...` (공유) | `dist/assets/...` |

### 2) 디렉터리 및 콘텐츠 관리 구조 (방안 1: 동일 디렉터리 접미사 방식 채택)
원문 작성 후 영문 번역을 진행하는 워크플로우에 최적화하여, **원문과 번역본을 동일 디렉터리에 나란히 배치**합니다:

```text
posts/
├── 2026-09-26-redefining-ai-coding-agent-autonomy.md     # 한글 원문
├── 2026-09-26-redefining-ai-coding-agent-autonomy.en.md  # 영문 번역본 (동일 슬러그)
├── 2026-09-29-ai-market-trends-h2.md                     # 한글 원문 (미번역 상태 즉시 식별)
└── assets/                                               # 공통 이미지/다이어그램 에셋
```

* **포스트 (Posts)**:
  - 한국어: `posts/<slug>.md` (기존과 100% 동일)
  - 영문: `posts/<slug>.en.md` (원문 바로 옆에 `.en.md` 접미사로 생성)
  - **슬러그 1:1 일치**: `site-cli` 파서가 `.en.md`를 제거하여 양쪽 모두 `redefining-ai-coding-agent-autonomy`라는 동일한 고유 슬러그를 공유 -> 상호 언어 스위칭이 100% 무결하게 연결됨.
  - **에셋 경로 무결성**: 디렉터리 깊이가 동일하므로 본문 내 `![arch](./assets/...)` 상대 경로를 번역 시 단 1자도 수정하지 않고 그대로 재사용 가능.
  - **번역 추적성**: 디렉터리 조회 시 `.en.md` 존재 여부만으로 번역 진행 상태를 한눈에 즉시 파악.
* **독립 단일 페이지 (Pages)**:
  - 한국어: `pages/about.html` (또는 `about.md`)
  - 영문: `pages/about.en.html` (또는 `about.en.md`) -> 빌더에서 `about.en` 감지 시 `/en/about/index.html`로 라우팅
* **UI 메시지 번들**:
  - 한국어: `templates/default-light/messages.yaml` (기존 유지)
  - 영문: `templates/default-light/messages_en.yaml` (신규)

### 3) `config.yaml` 스키마 확장안 (하위 호환성 100% 보장)
기존 설정을 그대로 유지하면서, 다국어 설정을 옵션으로 추가합니다:

```yaml
title: "AI Info"
subtitle: "Enterprise AI & Software Engineering Tech Blog"
baseURL: "/"
cname: "www.joinc.co.kr"
theme: "default-light"

# 기본 경로
source: "posts"
pages: "pages"
output: "dist"
redirects: "redirect.yaml"

# 다국어(i18n) 설정 (선택적 활성화)
i18n:
  default_lang: "ko"
  languages:
    - code: "ko"
      name: "한국어"
      base_url: "/"
      source_dir: "posts"
      post_pattern: "*.md"      # .en.md 제외한 순수 .md
      page_suffix: ""
      messages_file: "messages.yaml"
      output_prefix: ""
    - code: "en"
      name: "English"
      base_url: "/en/"
      source_dir: "posts"       # 동일 posts 디렉터리 참조
      post_pattern: "*.en.md"   # 영문 번역본만 스캔
      page_suffix: ".en"
      messages_file: "messages_en.yaml"
      output_prefix: "en"
```

### 4) `TemplateContext` 확장 (`model/types.go`)
템플릿 렌더링 시 언어별 컨텍스트를 주입하기 위해 필드를 확장합니다:

```go
type AlternateLink struct {
    Lang string
    URL  string
}

type TemplateContext struct {
    SiteTitle      string
    SiteSubtitle   string
    BaseURL        string
    CurrentLang    string           // "ko" 또는 "en"
    AlternateLangs []AlternateLink  // hreflang 및 언어 스위처용
    Messages       MessageBundle
    Categories     []*Category
    ActiveCategory *Category
    Posts          []*Post
    Post           *Post
    TOC            []TOCItem
    CurrentPath    string
    ExtraHead      template.HTML
    Design         *DesignTokens
}
```

### 5) 언어 전환 스위처(Language Switcher) 및 SEO (`base.html`)
* **헤더 네비게이션**:
  ```html
  <div class="flex items-center gap-1 bg-slate-100 p-1 rounded-lg text-xs font-semibold">
    <a href="{{ .SwitchURL.ko }}" class="px-2 py-1 rounded {{ if eq .CurrentLang "ko" }}bg-white text-blue-600 shadow-sm{{ else }}text-slate-500 hover:text-slate-900{{ end }}">KO</a>
    <a href="{{ .SwitchURL.en }}" class="px-2 py-1 rounded {{ if eq .CurrentLang "en" }}bg-white text-blue-600 shadow-sm{{ else }}text-slate-500 hover:text-slate-900{{ end }}">EN</a>
  </div>
  ```
* **다국어 SEO (`hreflang`)**:
  ```html
  {{ range .AlternateLangs }}
  <link rel="alternate" hreflang="{{ .Lang }}" href="{{ .URL }}" />
  {{ end }}
  ```

### 6) 브라우저 언어 감지 및 스마트 제안 배너 모듈 (Industry Best Practice)
Google Search Central 및 W3C i18n 가이드라인을 준수하여, **강제 자동 리디렉션을 배제**하고 **"브라우저 언어 감지 + 비침습적 제안 배너 + 사용자 선택 기억(localStorage)"** 방식을 적용합니다:

* **동작 시나리오**:
  1. **첫 방문 감지**: 사용자가 한국어 페이지(`/` 또는 `/posts/...`)에 접속했을 때, `navigator.language`가 영어(`en`)로 시작하고 이전에 안내를 닫은 이력(`localStorage.getItem('i18n_banner_dismissed')`)이 없는지 검사합니다.
  2. **상단 스마트 배너 노출**: 조건 충족 시 헤더 상단에 영문 전환 추천 배너를 노출합니다.
     - 배너 문구: *"🌐 English version is available for this content."*
     - 액션 버튼: **[Switch to English]** (해당 페이지의 영문 URL 또는 `/en/`으로 이동) / **[Dismiss]** (닫기)
  3. **사용자 선택 영구 기억**:
     - **[Switch to English]** 클릭 시: `localStorage.setItem('preferred_lang', 'en')`을 저장하고 영문 페이지로 이동.
     - **[Dismiss]** 클릭 시: `localStorage.setItem('i18n_banner_dismissed', 'true')`를 저장하여 다음 방문 시 배너를 다시 띄우지 않고 한국어 페이지에 머물도록 존중.
  4. **SEO 안전성 보장**: 정적 HTML에는 배너가 `hidden` 상태로 존재하며 클라이언트 자바스크립트로만 제어되므로, Googlebot 크롤러는 강제 리디렉션 없이 한국어/영문 원본을 완벽히 색인합니다.

---

## 4. 단계별 세부 구현 계획

### Phase 1. 템플릿 UI 리소스 및 독립 페이지 영문화 (Theme & Content)
1. **`templates/default-light/messages_en.yaml` 생성**:
   - Navigation, Search, Reading Time, TOC, Card, Detail, Empty, Footer 등 전 항목의 영문 텍스트 정의.
2. **`pages/about.en.html` 작성**:
   - 기존 `pages/about.html`의 레이아웃과 디자인 토큰을 100% 계승하면서 영문 소개 콘텐츠로 작성.
3. **영문 샘플 포스트 작성 (`posts/*.en.md`)**:
   - 한글 포스트 원본 바로 옆에 `posts/2026-09-26-redefining-ai-coding-agent-autonomy.en.md` 샘플 번역본을 작성하여 슬러그 1:1 매핑 및 빌드 검증 준비.

### Phase 2. 빌더 엔진(`tools/site-cli`) 다국어 파이프라인 확장 (Engine Path)
1. **파서 및 슬러그 추출 확장 (`internal/parser/frontmatter.go`, `scanner.go`)**:
   - `ExtractSlug`: `.en.md` 확장자를 감지하여 정상적으로 스트립, 원문과 동일한 고유 슬러그(`redefining-ai-coding-agent-autonomy`) 추출 보장.
   - `ScanPosts`: 언어별 포스트 필터링(한국어는 `.en.md` 제외, 영문은 `.en.md`만 수집) 지원.
2. **모델 확장 (`internal/model/types.go`)**:
   - `TemplateContext`에 `CurrentLang`, `AlternateLangs` 필드 추가.
3. **템플릿 엔진 확장 (`internal/template/engine.go`)**:
   - 지정된 메시지 파일(`messages_en.yaml`)을 로드하여 독립된 `MessageBundle`을 반환하는 메서드(`GetMessagesFor(filename string)`) 구현.
4. **빌더 오케스트레이션 확장 (`internal/builder/builder.go`)**:
   - `i18n` 설정이 활성화된 경우:
     - 1차: 기본 언어(`ko`) 빌드 수행 (`dist/`)
     - 2차: 영문 언어(`en`) 빌드 수행 (`dist/en/`)
     - `pages/` 디렉터리 내 `about.en.html`을 감지하여 `dist/en/about/index.html`로 렌더링.
     - `posts/*.en.md` 포스트를 스캔하여 `dist/en/posts/{slug}/`로 렌더링.
5. **엔진 단위 테스트 작성 및 통과**:
   - `builder_test.go`, `frontmatter_test.go`, `scanner_test.go`에 다국어 및 `.en.md` 슬러그 검증 테스트 케이스 추가.

### Phase 3. 템플릿 네비게이션 및 SEO 연동 (Theme Path)
1. **`templates/default-light/base.html` 수정**:
   - `<html lang="{{ .CurrentLang }}">` 동적 바인딩.
   - 홈 링크(`<a href="{{ .BaseURL }}">`) 및 소개 링크(`<a href="{{ .BaseURL }}about/">`)를 동적 BaseURL로 바인딩.
   - 헤더 우측에 `[ KO | EN ]` 언어 전환 스위처 컴포넌트 추가.
   - `<head>` 영역에 `hreflang` 메타 태그 자동 렌더링.
   - **브라우저 언어 스마트 감지 배너 모듈 구현**:
     - `navigator.language`가 영어(`en`)일 때 상단에 비침습적 영문 전환 제안 배너 노출.
     - [Switch to English] / [Dismiss] 클릭 시 `localStorage`에 상태를 영구 저장하여 사용자 선호 존중 및 재노출 방지.

### Phase 4. E2E 정적 빌드 및 배포 파이프라인 검증
1. **로컬 컴파일 및 정적 산출물 검증**:
   - `make build` 실행 후 `dist/index.html`, `dist/about/index.html`, `dist/en/index.html`, `dist/en/about/index.html`, `dist/en/posts/...` 생성 확인.
2. **로컬 개발 서버 검증**:
   - `make serve` 실행 후 브라우저에서 `http://localhost:8080/` 및 `http://localhost:8080/en/` 접근 및 상호 언어 스위칭 확인.
   - 영문 브라우저 환경에서 상단 스마트 제안 배너 정상 노출 및 `Dismiss` 클릭 시 `localStorage` 차단 동작 확인.
3. **CI/CD 워크플로우 무결성 확인**:
   - [`.github/workflows/deploy.yml`](../.github/workflows/deploy.yml)이 기존과 동일하게 `dist/` 산출물을 GitHub Pages로 완벽히 배포하는지 확인.

---

## 5. 코딩 가드레일 및 엄격한 격리 경계 (Hard Boundaries)

본 작업은 기존 블로그 운영 무결성을 보장하기 위해 **지정된 타깃 파일 외의 어떤 코드도 임의 수정하지 않는 엄격한 하드 바운더리**를 적용합니다:

### 1) 수정 및 생성 허용 파일 목록
| 구분 | 파일 경로 | 작업 계층 | 수정/생성 상세 내용 |
| :---: | :--- | :--- | :--- |
| **신규** | `templates/default-light/messages_en.yaml` | `Theme Path` | UI 텍스트 전반의 영문 번들 리소스 파일 작성 |
| **신규** | `pages/about.en.html` | `Content Fast Path` | 영문 소개 페이지 작성 |
| **신규** | `posts/2026-09-26-redefining-ai-coding-agent-autonomy.en.md` | `Content Fast Path` | 영문 빌드 검증용 첫 샘플 번역 포스트 작성 |
| **수정** | `config.yaml` | `Config Path` | `i18n` 다국어 설정 블록 추가 |
| **수정** | `templates/default-light/base.html` | `Theme Path` | `lang`, `BaseURL`, 언어 스위처, `hreflang`, 언어 감지 스마트 배너 & localStorage 스크립트 바인딩 |
| **수정** | `tools/site-cli/internal/model/types.go` | `Engine Path` | `TemplateContext` 내 `CurrentLang`, `AlternateLangs` 필드 추가 |
| **수정** | `tools/site-cli/internal/template/engine.go` | `Engine Path` | 다국어 메시지 파일 동적 로드 로직 추가 |
| **수정** | `tools/site-cli/internal/parser/frontmatter.go` | `Engine Path` | `ExtractSlug` 내 `.en.md` 확장자 스트립 및 동일 슬러그 유지 |
| **수정** | `tools/site-cli/internal/parser/scanner.go` | `Engine Path` | 언어별(`*.md` vs `*.en.md`) 포스트 스캔 필터링 |
| **수정** | `tools/site-cli/internal/builder/builder.go` | `Engine Path` | 다국어 순차 빌드 및 `/en/` 서브패스 산출물 생성 오케스트레이션 |

### 2) 절대 수정 금지 영역 (Hard Boundary)
* ❌ **기존 한국어 포스트 (`posts/*.md`)**: 기존 한글 포스트 원본 및 날짜/메타데이터 일체 불변.
* ❌ **마크다운 파서 코어 (`tools/site-cli/internal/markdown/`)**: Goldmark 변환 로직 및 수식/다이어그램 렌더링 파이프라인 수정 금지.
* ❌ **GitHub Actions 워크플로우 (`.github/workflows/deploy.yml`)**: 기존의 `make build` -> `dist/` 배포 구조 변경 금지.

---

## 6. 완료 기준 (Definition of Done)

- [x] `task/TP-6_multilingual_english_deployment_plan.md` 작업 계획서 수립 완료
- [x] `templates/default-light/messages_en.yaml` 작성 완료
- [x] `pages/about.en.html` 및 `posts/*.en.md` 샘플 번역 포스트 작성 완료
- [x] `config.yaml`에 `i18n` 설정 스키마 반영 완료
- [x] `tools/site-cli` 파서(`frontmatter.go`, `scanner.go`) 및 빌더 다국어 로직 구현 및 단위 테스트 통과 (`cd tools/site-cli && go test -v ./...`)
- [x] `make build` 실행 시 `dist/` 및 `dist/en/` 산출물 정상 생성 확인 (동일 슬러그로 한국어/영문 포스트 생성 확인)
- [x] 브라우저 로컬 서빙(`make serve`) 환경에서 헤더 언어 스위처(`KO <-> EN`) 동작 확인
- [x] 브라우저 언어(영문) 설정 시 상단 스마트 안내 배너 노출 및 `localStorage` 억제/전환 기능 검증
- [x] `git status`로 지정된 파일 외의 무분별한 파일 변경이 없음을 엄격 확인
- [x] Jira 티켓 [TP-6](https://joincdream.atlassian.net/browse/TP-6) 코멘트 등록 및 상태 업데이트
