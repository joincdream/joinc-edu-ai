# [JC-12] KaTeX 기반 포스트 제목 및 본문 수식(LaTeX) 렌더링 지원 작업 계획서

## 1. 기본 정보 및 매핑

* **티켓 번호**: [JC-12](https://joincdream.atlassian.net/browse/JC-12)
* **담당 서브시스템**: `template_engine` (UI/UX 렌더링 계층)
* **타깃 파일**:
  * `templates/default/base.html` (KaTeX CDN 및 auto-render 스크립트 추가)
* **참조 문서**:
  * [docs/okf/site-cli.md](../docs/okf/site-cli.md) (`template_engine` 서브시스템)
  * [docs/architecture/principles.md](../docs/architecture/principles.md) (Pure Static & 3계층 격리 원칙)

---

## 2. 배경 및 문제 정의

* **현상**: 포스트 제목(`$0.95^{10}$`) 및 본문의 블록 수식(`$$P(\text{Success}) = \prod_{i=1}^{n} p_i = (0.95)^{10} \approx 0.598 \quad (59.8\%)$$`)이 컴파일되지 않고 텍스트 그대로 화면에 노출됨.
* **원인**: 템플릿에 Mermaid 및 PrismJS만 로드되어 있고, LaTeX 수식 렌더러(KaTeX)가 포함되어 있지 않음.
* **목표**: 
  1. 포스트 제목(`h1`) 및 본문(`.prose`)의 인라인 수식(`$...$`)과 블록 수식(`$$...$$`)이 표준 수학 기호로 올바르게 렌더링되도록 구현.
  2. 코드 블록(`pre`, `code`) 및 Mermaid 다이어그램과의 스크립트 충돌 방지.

---

## 3. 단계별 세부 구현 계획

### Step 1. KaTeX 라이브러리 및 확장 기능 로드
`templates/default/base.html`의 `<head>` 태그 내에 KaTeX 최신 안정 버전 CDN 추가:
```html
<!-- KaTeX CSS & JS -->
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.css" />
<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.js"></script>
<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/contrib/auto-render.min.js"></script>
```

### Step 2. 수식 Auto-render 바인딩 및 충돌 방지 옵션 설정
`templates/default/base.html` 하단 스크립트 영역에 DOMContentLoaded 핸들러 구성:
```javascript
// 3. KaTeX 수식 자동 렌더링 (제목 및 본문 적용)
document.addEventListener("DOMContentLoaded", () => {
  if (typeof renderMathInElement === "function") {
    renderMathInElement(document.body, {
      delimiters: [
        { left: "$$", right: "$$", display: true },   // 디스플레이 수식
        { left: "$", right: "$", display: false },     // 인라인 수식
        { left: "\\(", right: "\\)", display: false },
        { left: "\\[", right: "\\]", display: true }
      ],
      ignoredTags: ["script", "noscript", "style", "textarea", "pre", "code"],
      ignoredClasses: ["mermaid-raw", "mermaid-container"],
      throwOnError: false
    });
  }
});
```

### Step 3. 3계층 격리 및 사이드 이펙트 방지 전략
* **코어 엔진(`tools/site-cli/`) 불변 유지**: 빌더 엔진 코드 수정 없이 템플릿 계층(`templates/default/`)만 수정하여 Go 파이프라인의 무결성 보장.
* **Mermaid 및 PrismJS 충돌 격리**: `ignoredTags`에 `pre`, `code`를 지정하고, `ignoredClasses`에 Mermaid 컨테이너를 지정하여 코드 블록 내의 `$` 기호나 다이어그램 구문이 수식으로 오인 파싱되는 현상을 원천 방지.

---

## 4. 완료 기준 (Definition of Done)

- [ ] `templates/default/base.html`에 KaTeX CDN 및 auto-render 스크립트 반영
- [ ] 정적 사이트 정상 컴파일 검증 (`make build`)
- [ ] Go 엔진 단위 테스트 무결성 확인 (`cd tools/site-cli && go test -v ./...`)
- [ ] 수식 포함 포스트의 HTML 산출물에서 KaTeX 로드 및 정상 렌더링 확인
- [ ] Jira 티켓 [JC-12](https://joincdream.atlassian.net/browse/JC-12) 상태 업데이트 및 코멘트 기록
