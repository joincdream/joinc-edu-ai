---
version: alpha
name: Default Light
description: Joinc AI 엔터프라이즈 딥 테크 아키텍처 저널을 위한 가독성 중심의 라이트 테마 디자인 시스템
colors:
  neutral: "#ffffff"
  on-neutral: "#0f172a"
  surface: "#ffffff"
  surface-elevated: "#f1f5f9"
  border-subtle: "#e2e8f0"
  primary: "#1d4ed8"
  primary-hover: "#1e40af"
  link: "#1d4ed8"
  text-primary: "#0f172a"
  text-body: "#334155"
  text-muted: "#64748b"
  info: "#0369a1"
  success: "#15803d"
  warning: "#b45309"
  danger: "#b91c1c"
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
    textColor: "{colors.neutral}"
    rounded: "{rounded.md}"
    padding: 10px 18px
  button-primary-hover:
    backgroundColor: "{colors.primary-hover}"
    textColor: "{colors.neutral}"
    rounded: "{rounded.md}"
    padding: 10px 18px
  card-post:
    backgroundColor: "{colors.surface}"
    rounded: "{rounded.lg}"
    padding: 24px
  code-block:
    backgroundColor: "{colors.surface-elevated}"
    textColor: "{colors.text-primary}"
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
    backgroundColor: "{colors.surface-elevated}"
    textColor: "{colors.success}"
    rounded: "{rounded.md}"
  alert-warning:
    backgroundColor: "{colors.surface-elevated}"
    textColor: "{colors.warning}"
    rounded: "{rounded.md}"
  alert-danger:
    backgroundColor: "{colors.surface-elevated}"
    textColor: "{colors.danger}"
    rounded: "{rounded.md}"
  alert-info:
    backgroundColor: "{colors.surface-elevated}"
    textColor: "{colors.info}"
    rounded: "{rounded.md}"
  divider:
    backgroundColor: "{colors.border-subtle}"
    height: 1px
---

## Overview

Clarity, Readability, and Architectural Rigor in Daylight.

Joinc AI 라이트 테마는 인쇄된 테크니컬 리포트와 고품질 엔지니어링 서적의 조판 원칙을 웹 인터페이스로 재해석한 **모던 라이트 저널 캔버스**입니다. 

완전한 순백색(#ffffff)과 부드러운 슬레이트(#f1f5f9) 표면 레이어를 기반으로 하여, 밝은 조명 환경에서도 시각적 피로 없이 수식, 다이어그램, 장문 아키텍처 분석글에 집중할 수 있도록 돕습니다.

## Colors

모든 색상은 배경과 텍스트 간의 수학적 명도 대비(Contrast Ratio)를 엄격히 통제하여 시인성을 확보합니다:

- **표면 및 배경 (Surfaces)**:
  - `neutral (#ffffff)`: 캔버스 전체의 기본 배경색 (Pure White).
  - `surface (#ffffff)`: 포스트 카드 및 메인 헤더의 기본 표면.
  - `surface-elevated (#f1f5f9)`: 사이드바, 태그 칩, 인라인 코드, 호버 상태의 표면 레이어.
  - `border-subtle (#e2e8f0)`: 부드럽고 명확한 구조적 분리선 (Slate 200).
- **텍스트 시인성 계층 (WCAG AAA/AA 엄수)**:
  - `text-primary (#0f172a)`: 메인 헤드라인 및 주요 타이틀. 명도 대비 **18.7:1** (WCAG AAA 초과).
  - `text-body (#334155)`: 장문 기술 본문. 명도 대비 **9.4:1** (WCAG AAA 초과)로 장시간 읽기 최적화.
  - `text-muted (#64748b)`: 메타데이터, 작성일, 목차 링크. 명도 대비 **4.6:1** (WCAG AA 4.5:1 초과).
- **액센트 및 상호작용 (Interactive)**:
  - `primary (#1d4ed8)`: 주 CTA 버튼, 탭 활성 인디케이터 (명도 대비 **5.8:1**).
  - `primary-hover (#1e40af)`: 마우스 오버 시 시각적 피드백 (명도 대비 **8.0:1**).
  - `link (#1d4ed8)`: 본문 하이퍼링크 및 태그 텍스트.

## Typography

가독성이 입증된 Pretendard와 JetBrains Mono 조합을 유지하여 한국어 기술 문서 및 영문 코드 블록의 시각적 리듬을 통일합니다:

- `h1`: 2.25rem (36px), Bold 700, -0.025em 자간.
- `h2`: 1.5rem (24px), Bold 700, -0.02em 자간.
- `h3`: 1.25rem (20px), Semi-Bold 600.
- `body-md`: 1.05rem (16.8px), Regular 400, 1.75 줄간격.
- `code`: 0.875rem (14px), Medium 500, JetBrains Mono.

## Components Specification

- **Header**: 배경 `#ffffff`, 높이 `64px`, 하단 1px 경계선 `#e2e8f0`.
- **Card Post**: 카드 배경 `#ffffff`, 모서리 `12px(rounded.lg)`, 내부 패딩 `24px`.
- **Code Block**: 배경 `#f1f5f9`, 코드 텍스트 `#0f172a`, 모서리 `8px(rounded.md)`.
- **Tag Pill**: 배경 `#f1f5f9`, 텍스트 `#1d4ed8`, 모서리 `9999px(rounded.full)`.
