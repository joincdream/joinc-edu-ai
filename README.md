# joinc-ai.io

엔터프라이즈 AI 엔지니어링, 클라우드 아키텍처, AXP(AI Experience Platform) 전문 기술 블로그 및 지식 베이스 리포지토리입니다.

Go 기반 자체 경량 정적 사이트 생성기(SSG)를 통해 **런타임 서버 및 데이터베이스 제로(Pure Static)** 환경으로 GitHub Pages에 무중단 배포됩니다.

---

## 📁 디렉터리 아키텍처

```text
.
├── posts/                  # 기술 블로그 포스트 (Markdown)
│   ├── YYYY-MM-DD-slug.md
│   ├── drafts/             # 작성 중인 초안 포스트
│   └── assets/             # 포스트 내 이미지/다이어그램
│
├── pages/                  # 독립 단일 정적 페이지 (About, Contact 등)
│   ├── about.html          # 이력 및 비전 소개 페이지 (HTML 직접 관리)
│   └── assets/             # 페이지 전용 정적 에셋
│
├── templates/              # HTML/CSS 디자인 테마
│   └── default/
│       ├── base.html       # 공통 GNB, 푸터, 메타데이터 레이아웃
│       ├── index.html      # 홈 / 최신 포스트 목록
│       ├── detail.html     # 포스트 상세 보기
│       ├── category.html   # 카테고리별 포스트 모아보기
│       ├── page.html       # 마크다운 페이지 기본 템플릿
│       ├── messages.yaml   # 다국어/문구 리소스 번들
│       └── style.css       # Tailwind CSS 유틸리티 및 커스텀 스타일
│
├── tools/                  # 애플리케이션 및 CLI 도구 격리
│   └── site-cli/           # Go 기반 정적 사이트 생성기 & 로컬 개발 서버
│       ├── cmd/site-cli/   # CLI 엔트리포인트 (build, serve, deploy, categories)
│       ├── internal/       # 파서, 템플릿 엔진, 빌더, 라이브 서버
│       ├── go.mod
│       └── go.sum
│
├── docs/                   # 기술 문서 & 리서치 자료
│   ├── research/           # 딥다이브 리서치 (Market trends, Deep-dive)
│   └── architecture/       # 사이트 빌더 아키텍처 및 스펙
│
├── Makefile                # 통합 실행 명령어 모음
├── .gitignore              # 빌드 산출물(dist/) 및 바이너리 제외
└── .github/workflows/      # GitHub Actions 자동 배포 파이프라인
```

---

## 🚀 빠른 시작 (Local Development)

### 사전 요구사항
* [Go 1.23+](https://go.dev/) 설치 필요

### 1. 로컬 개발 서버 구동 (Hugo 스타일 라이브 리로드)
`posts/`, `pages/`, `templates/` 디렉터리의 변경 사항을 실시간 감지하여 브라우저를 자동 새로고침합니다.
```bash
make serve
```
* 접속 주소: `http://localhost:8080`

### 2. 정적 웹사이트 일괄 빌드
`dist/` 디렉터리에 배포 가능한 순수 HTML/CSS/이미지 정적 파일을 컴파일합니다.
```bash
make build
```

### 3. 카테고리 및 포스트 통계 확인
```bash
make categories
```

### 4. 단위 테스트 실행
```bash
make test
```

---

## ✍️ 기술 블로그 포스트 작성 가이드

`posts/` 디렉터리 하위에 마크다운(`*.md`) 파일로 작성하며, 상단에 다음과 같은 YAML Frontmatter를 선언합니다:

```markdown
---
title: "포스트 제목"
date: "2026-03-26"
category: "AI Agent"
tags: ["Agentic AI", "Orchestration", "LLM"]
summary: "포스트 요약 문구 (목록 카드 및 메타 태그 노출)"
status: "published" # 'published' 또는 'draft'
author: "윤상배"
---

# 본문 내용
...
```

* `status: draft`로 설정된 문서는 로컬 테스트(`make serve`) 시에는 노출되지 않으며, `site-cli build --drafts` 옵션으로만 빌드됩니다.
* 이미지는 `posts/assets/`에 저장하고, 마크다운에서는 `/assets/images/파일명.png` 경로로 참조합니다.

---

## 🌐 배포 (GitHub Pages)

`main` 브랜치에 코드가 푸시되면 `.github/workflows/deploy.yml` 워크플로우가 자동으로 실행되어 사이트를 빌드하고 GitHub Pages로 배포합니다.
