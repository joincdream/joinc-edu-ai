---
version: alpha
name: Default Dark
description: Joinc AI 엔터프라이즈 딥 테크 아키텍처 저널을 위한 장문 몰입형 다크 테마 디자인 시스템
colors:
  neutral: "#020617"
  on-neutral: "#f8fafc"
  surface: "#0f172a"
  surface-elevated: "#1e293b"
  border-subtle: "#334155"
  primary: "#1d4ed8"
  primary-hover: "#2563eb"
  link: "#60a5fa"
  text-primary: "#f8fafc"
  text-body: "#cbd5e1"
  text-muted: "#94a3b8"
  info: "#38bdf8"
  success: "#34d399"
  warning: "#fbbf24"
  danger: "#f87171"
typography:
  h1:
    fontFamily: Pretendard, -apple-system, BlinkMacSystemFont, system-ui, sans-serif
    fontSize: 2.25rem
    fontWeight: 700
    lineHeight: 1.2
    letterSpacing: -0.025em
  h2:
    fontFamily: Pretendard, -apple-system, BlinkMacSystemFont, system-ui, sans-serif
    fontSize: 1.5rem
    fontWeight: 700
    lineHeight: 1.3
    letterSpacing: -0.02em
  h3:
    fontFamily: Pretendard, -apple-system, BlinkMacSystemFont, system-ui, sans-serif
    fontSize: 1.25rem
    fontWeight: 600
    lineHeight: 1.4
  body-md:
    fontFamily: Pretendard, -apple-system, BlinkMacSystemFont, system-ui, sans-serif
    fontSize: 1.05rem
    fontWeight: 400
    lineHeight: 1.75
  code:
    fontFamily: JetBrains Mono, ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace
    fontSize: 0.875rem
    fontWeight: 500
    lineHeight: 1.5
rounded:
  sm: 4px
  md: 8px
  lg: 12px
  full: 9999px
spacing:
  xs: 4px
  sm: 8px
  md: 16px
  lg: 24px
  xl: 32px
components:
  header:
    backgroundColor: "{colors.surface}"
    height: 64px
  button-primary:
    backgroundColor: "{colors.primary}"
    textColor: "{colors.on-neutral}"
    rounded: "{rounded.md}"
    padding: 10px 18px
  button-primary-hover:
    backgroundColor: "{colors.primary-hover}"
    textColor: "{colors.on-neutral}"
    rounded: "{rounded.md}"
    padding: 10px 18px
  card-post:
    backgroundColor: "{colors.surface}"
    rounded: "{rounded.lg}"
    padding: 24px
  code-block:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.info}"
    rounded: "{rounded.md}"
    padding: 16px
  toc-navigation:
    textColor: "{colors.text-muted}"
    typography: "{typography.code}"
  tag-pill:
    backgroundColor: "{colors.surface-elevated}"
    textColor: "{colors.link}"
    rounded: "{rounded.full}"
    padding: 4px 12px
  heading-hero:
    textColor: "{colors.text-primary}"
    typography: "{typography.h1}"
  body-content:
    textColor: "{colors.text-body}"
    typography: "{typography.body-md}"
  alert-success:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.success}"
    rounded: "{rounded.md}"
  alert-warning:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.warning}"
    rounded: "{rounded.md}"
  alert-danger:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.danger}"
    rounded: "{rounded.md}"
  divider:
    backgroundColor: "{colors.border-subtle}"
    height: 1px
---

## Overview

Architectural Precision meets Journalistic Gravitas.

Joinc AI 블로그의 UI는 화려하거나 유치한 장식 요소를 배제하고, 장시간 고도의 기술 아티클과 시스템 설계를 분석하는 엔지니어가 눈의 피로 없이 콘텐츠에 완전히 몰입할 수 있도록 설계된 **엔터프라이즈 딥 테크 저널 캔버스**입니다. 

화면의 기본 골격은 눈부심을 일으키는 완전 블랙(#000000) 대신 차분한 딥 슬레이트(#020617)를 기초로 하며, 절제된 테크 블루(#3b82f6)를 통해 중요한 상호작용 지점을 명확히 안내합니다.

## Colors

모든 색상은 배경과 텍스트 간의 수학적 명도 대비(Contrast Ratio)를 엄격히 통제하여 시인성을 확보합니다:

- **표면 및 배경 (Surfaces)**:
  - `neutral (#020617)`: 캔버스 전체의 기본 배경색 (Slate 950).
  - `surface (#0f172a)`: 헤더 바, 포스트 카드, 사이드바를 구분하는 베이스 표면 레이어.
  - `surface-elevated (#1e293b)`: 호버 하이라이트, 모달, 인라인 코드 배경.
  - `border-subtle (#334155)`: 과도한 눈길을 끌지 않는 1px 구조적 분리선.
- **텍스트 시인성 계층 (WCAG AAA/AA 엄수)**:
  - `text-primary (#f8fafc)`: 메인 헤드라인 및 주요 타이틀. 명도 대비 **16.8:1** (WCAG AAA 초과).
  - `text-body (#cbd5e1)`: 장문 기술 본문. 명도 대비 **11.2:1**로 눈부심 없는 최상의 가독성 보장.
  - `text-muted (#94a3b8)`: 메타데이터, 작성일, 각주. 명도 대비 **6.4:1** (WCAG AA 4.5:1 초과).
- **액센트 및 링크 (Interactive)**:
  - `primary (#3b82f6)`: 주 CTA 버튼, 탭 인디케이터, 강조 태그.
  - `primary-hover / link (#60a5fa)`: 마우스 오버 피드백 및 본문 하이퍼링크.

## Typography

가독성이 입증된 Pretendard와 개발자에게 친숙한 모노스페이스 폰트를 결합합니다:

- **폰트 패밀리**:
  - UI 및 아티클 본문: `Pretendard`, `-apple-system`, `sans-serif`
  - 코드 및 수식: `JetBrains Mono`, `ui-monospace`, `monospace`
- **본문 읽기 리듬**:
  - 기술 문서의 긴 호흡을 유지하기 위해 본문(`.prose p`)은 **`font-size: 1.05rem`**과 **`line-height: 1.75`**를 엄격히 준수합니다.
  - 가독 폭 제약: 한 줄당 최적 글자 수(65~75자) 유지를 위해 본문 최대 너비는 **`48rem (768px, prose-lg)`** 로 제한됩니다.

## Layout & Spacing

8px 기반의 일관된 공간 배분을 적용합니다:

- `spacing.sm (8px)`: 인라인 요소 간 간격, 배지 내부 여백.
- `spacing.md (16px)`: 카드 내부 패딩, 헤더 좌우 여백.
- `spacing.lg (24px)`: 섹션 간 기본 분리 여백, 카드 간격.
- `spacing.xl (32px)`: H2 대주제 간 여백, 주요 레이아웃 블록 간격.

## Elevation & Depth

과도한 그림자(Drop Shadow)를 지양하고 1px 보더(`border-subtle`)와 미세한 명도 차이(`surface` ➔ `surface-elevated`)를 통해 계층을 구분합니다.

## Shapes

- 카드 및 모달: 부드럽고 현대적인 **`rounded.lg (12px)`**
- 버튼 및 인라인 코드: 절제된 **`rounded.md (8px)`**
- 태그 배지: 완전한 원형인 **`rounded.full (9999px, Pill)`**

## Components

- **Header**: 상단에 고정(`sticky`)되며, `surface/80`의 반투명 배경과 `backdrop-blur-md`를 결합하여 스크롤 시 콘텐츠와 자연스럽게 겹칩니다.
- **Card Post**: 마우스 호버 시 그림자를 키우지 않고, 보더 색상을 `primary/50`으로 미세하게 밝혀 인터랙션을 암시합니다.
- **TOC Navigation**: 우측에 플로팅되며, 현재 읽고 있는 활성 헤딩에는 `primary` 색상의 2px 레프트 보더가 표시됩니다.
- **Code Block**: `surface` 배경과 `border-subtle`로 감싸 본문 텍스트와 시각적으로 명확히 분리합니다.

## Do's and Don'ts

AI 코딩 에이전트와 엔지니어는 템플릿 코드 작성 시 다음 불변 규칙을 따라야 합니다:

- **DO**: 본문 텍스트와 배경의 명도 대비는 항상 7:1 (WCAG AAA) 이상을 유지할 것.
- **DO**: 본문 하이퍼링크는 4px 오프셋의 밑줄을 동반하여 일반 텍스트와 시각적으로 명확히 구분할 것.
- **DO**: 모든 색상, 마진, 패딩은 `DESIGN.md` 토큰 또는 Tailwind 표준 클래스만 사용할 것.
- **DON'T**: 본문 일반 텍스트에 파란색(`text-blue-*`)을 임의로 칠하지 말 것 (링크로 오인 유발).
- **DON'T**: 완전 블랙(`#000000`) 배경에 완전 화이트(`#ffffff`) 글자를 100% 면적으로 충돌시키지 말 것 (눈의 피로 초래).
- **DON'T**: Tailwind 임의값(Arbitrary Values: `bg-[#123456]`, `p-[13px]`)을 코드에 하드코딩하지 말 것.
