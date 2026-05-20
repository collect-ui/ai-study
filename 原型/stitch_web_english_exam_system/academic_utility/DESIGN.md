---
name: Academic Utility
colors:
  surface: '#f8f9fa'
  surface-dim: '#d9dadb'
  surface-bright: '#f8f9fa'
  surface-container-lowest: '#ffffff'
  surface-container-low: '#f3f4f5'
  surface-container: '#edeeef'
  surface-container-high: '#e7e8e9'
  surface-container-highest: '#e1e3e4'
  on-surface: '#191c1d'
  on-surface-variant: '#414754'
  inverse-surface: '#2e3132'
  inverse-on-surface: '#f0f1f2'
  outline: '#727785'
  outline-variant: '#c1c6d6'
  surface-tint: '#005bc0'
  primary: '#005bbf'
  on-primary: '#ffffff'
  primary-container: '#1a73e8'
  on-primary-container: '#ffffff'
  inverse-primary: '#adc7ff'
  secondary: '#005ac1'
  on-secondary: '#ffffff'
  secondary-container: '#4d8efe'
  on-secondary-container: '#00285c'
  tertiary: '#006d2c'
  on-tertiary: '#ffffff'
  tertiary-container: '#008939'
  on-tertiary-container: '#ffffff'
  error: '#ba1a1a'
  on-error: '#ffffff'
  error-container: '#ffdad6'
  on-error-container: '#93000a'
  primary-fixed: '#d8e2ff'
  primary-fixed-dim: '#adc7ff'
  on-primary-fixed: '#001a41'
  on-primary-fixed-variant: '#004493'
  secondary-fixed: '#d8e2ff'
  secondary-fixed-dim: '#adc6ff'
  on-secondary-fixed: '#001a41'
  on-secondary-fixed-variant: '#004494'
  tertiary-fixed: '#89fa9b'
  tertiary-fixed-dim: '#6ddd81'
  on-tertiary-fixed: '#002108'
  on-tertiary-fixed-variant: '#005320'
  background: '#f8f9fa'
  on-background: '#191c1d'
  surface-variant: '#e1e3e4'
typography:
  display-lg:
    fontFamily: Work Sans
    fontSize: 48px
    fontWeight: '700'
    lineHeight: 56px
    letterSpacing: -0.02em
  headline-lg:
    fontFamily: Work Sans
    fontSize: 32px
    fontWeight: '600'
    lineHeight: 40px
  headline-md:
    fontFamily: Work Sans
    fontSize: 24px
    fontWeight: '600'
    lineHeight: 32px
  headline-sm:
    fontFamily: Work Sans
    fontSize: 20px
    fontWeight: '600'
    lineHeight: 28px
  body-lg:
    fontFamily: Work Sans
    fontSize: 18px
    fontWeight: '400'
    lineHeight: 28px
  body-md:
    fontFamily: Work Sans
    fontSize: 16px
    fontWeight: '400'
    lineHeight: 24px
  body-sm:
    fontFamily: Work Sans
    fontSize: 14px
    fontWeight: '400'
    lineHeight: 20px
  label-lg:
    fontFamily: Work Sans
    fontSize: 14px
    fontWeight: '600'
    lineHeight: 20px
  label-md:
    fontFamily: Work Sans
    fontSize: 12px
    fontWeight: '500'
    lineHeight: 16px
  headline-lg-mobile:
    fontFamily: Work Sans
    fontSize: 28px
    fontWeight: '600'
    lineHeight: 36px
rounded:
  sm: 0.125rem
  DEFAULT: 0.25rem
  md: 0.375rem
  lg: 0.5rem
  xl: 0.75rem
  full: 9999px
spacing:
  base: 8px
  xs: 4px
  sm: 12px
  md: 24px
  lg: 40px
  xl: 64px
  gutter: 24px
  margin: 32px
---

## Brand & Style

The design system is engineered for the high-stakes environment of educational administration. It prioritizes cognitive ease and rapid information retrieval for teachers managing English examinations. The aesthetic is **Minimalist and Professional**, leaning into a modern "SaaS-for-Education" look that emphasizes utility over decoration.

The system aims to evoke a sense of **competence, reliability, and calm**. By utilizing generous whitespace and a restricted color palette, it reduces the visual noise often associated with complex data management. The interface acts as a silent partner to the educator, ensuring that exam parameters—difficulty levels, unit coverage, and grading metrics—are the primary focus of the user's attention.

## Colors

The palette is anchored by **Academic Blue (#1A73E8)**, a color synonymous with trust and intellectual focus. This primary shade is reserved for actionable elements, progress indicators, and primary branding moments. 

- **Backgrounds:** A tiered system of Soft Grays (#F8F9FA to #F1F3F4) is used to differentiate the canvas from individual modules or cards.
- **Typography:** Deep Charcoal (#202124) is used for body text to maintain a high contrast ratio (conforming to WCAG AA standards), while a lighter gray (#5F6368) handles metadata and labels.
- **Semantic Colors:** Green is utilized for "Validated" or "Published" exam states, while subtle amber is used for "Draft" modes.

## Typography

The typography strategy employs **Work Sans** across all levels. This choice provides a grounded, professional feel with excellent legibility at small sizes—crucial for reading exam questions and data tables. 

The hierarchy is structured to support "scanning." Large, bold headlines identify the current module (e.g., "Exam Bank"), while high-contrast labels clearly demarcate form fields. For educational content like reading comprehension passages, `body-lg` is recommended to ensure maximum readability for the reviewer.

## Layout & Spacing

The design system utilizes a **Fixed Grid** philosophy for its primary dashboard to ensure a stable, predictable management environment. 

- **Grid:** A 12-column grid system with 24px gutters.
- **Breakpoints:**
    - Mobile (< 600px): Fluid 4-column layout, 16px margins.
    - Tablet (600px - 1024px): 8-column layout, 24px margins.
    - Desktop (> 1024px): Fixed 1200px max-width container, centered.
- **Rhythm:** An 8px baseline grid dictates all vertical spacing. Elements like input fields and buttons are 48px or 56px high to provide comfortable "hit targets" for educators using touch-enabled laptops.

## Elevation & Depth

To maintain a clean, academic aesthetic, this design system avoids heavy shadows. Instead, it uses **Tonal Layers** and **Low-Contrast Outlines**.

- **Cards:** White surfaces (#FFFFFF) sit atop the Soft Gray background. They are defined by a 1px border (#E0E0E0) rather than a shadow.
- **Hover States:** When a card or list item is hovered, it gains a very soft, ambient shadow (0px 4px 12px rgba(0, 0, 0, 0.05)) and the border color shifts to the Primary Blue.
- **Modals:** Only high-level interruptions (like "Delete Exam") use a backdrop blur (4px) and a medium-diffused shadow to pull them to the foreground.

## Shapes

The shape language is **Soft**. A 4px (0.25rem) corner radius is applied to buttons, input fields, and tags. This subtle rounding suggests modern software without losing the "official" feel required for an academic tool.

- **Buttons & Inputs:** 4px radius.
- **Cards:** 8px radius (`rounded-lg`) to provide a clear container for grouped content.
- **Interactive Tags:** Pill-shaped (fully rounded) only for status indicators (e.g., "Hard", "Grade 7") to distinguish them from actionable buttons.

## Components

### Buttons
Primary buttons are solid Academic Blue with white text. Secondary buttons use an outline style with 1px Academic Blue borders. High-contrast focus states are mandatory, featuring a 2px offset ring.

### Cards
The primary container for exam questions and metadata. Cards must include a clear header area for the question type (e.g., "Multiple Choice") and a footer for difficulty/unit tags.

### Input Fields
Standardized height of 48px. Labels are always visible above the field (never just placeholders) to aid memory during complex form entry. Active states use a 2px Academic Blue border.

### Chips & Selection
Selection components for "Grade" or "Difficulty" use a "Toggle Group" style—horizontal blocks that highlight the selected option in blue. This provides a clear, high-contrast visual of current filters.

### Lists
Data tables for exam lists use alternating row colors (Zebra striping) with #F8F9FA to help educators track information across wide screens. Row height is generous (64px) to avoid visual cramping.