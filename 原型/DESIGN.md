---
name: Intelligent Learning System
colors:
  surface: '#f7f9fc'
  surface-dim: '#d8dadd'
  surface-bright: '#f7f9fc'
  surface-container-lowest: '#ffffff'
  surface-container-low: '#f2f4f7'
  surface-container: '#eceef1'
  surface-container-high: '#e6e8eb'
  surface-container-highest: '#e0e3e6'
  on-surface: '#191c1e'
  on-surface-variant: '#414751'
  inverse-surface: '#2d3133'
  inverse-on-surface: '#eff1f4'
  outline: '#717783'
  outline-variant: '#c1c7d3'
  surface-tint: '#0060ac'
  primary: '#005da7'
  on-primary: '#ffffff'
  primary-container: '#2976c7'
  on-primary-container: '#fdfcff'
  inverse-primary: '#a4c9ff'
  secondary: '#006b58'
  on-secondary: '#ffffff'
  secondary-container: '#67f7d5'
  on-secondary-container: '#00705c'
  tertiary: '#7f5300'
  on-tertiary: '#ffffff'
  tertiary-container: '#a06900'
  on-tertiary-container: '#fffbff'
  error: '#ba1a1a'
  on-error: '#ffffff'
  error-container: '#ffdad6'
  on-error-container: '#93000a'
  primary-fixed: '#d4e3ff'
  primary-fixed-dim: '#a4c9ff'
  on-primary-fixed: '#001c39'
  on-primary-fixed-variant: '#004883'
  secondary-fixed: '#6bfad8'
  secondary-fixed-dim: '#48ddbc'
  on-secondary-fixed: '#002019'
  on-secondary-fixed-variant: '#005142'
  tertiary-fixed: '#ffddb4'
  tertiary-fixed-dim: '#ffb953'
  on-tertiary-fixed: '#291800'
  on-tertiary-fixed-variant: '#633f00'
  background: '#f7f9fc'
  on-background: '#191c1e'
  surface-variant: '#e0e3e6'
typography:
  display-lg:
    fontFamily: Plus Jakarta Sans
    fontSize: 32px
    fontWeight: '700'
    lineHeight: 40px
    letterSpacing: -0.02em
  headline-lg:
    fontFamily: Plus Jakarta Sans
    fontSize: 24px
    fontWeight: '600'
    lineHeight: 32px
  headline-md:
    fontFamily: Plus Jakarta Sans
    fontSize: 20px
    fontWeight: '600'
    lineHeight: 28px
  body-lg:
    fontFamily: Be Vietnam Pro
    fontSize: 16px
    fontWeight: '400'
    lineHeight: 24px
  body-md:
    fontFamily: Be Vietnam Pro
    fontSize: 14px
    fontWeight: '400'
    lineHeight: 20px
  label-md:
    fontFamily: Be Vietnam Pro
    fontSize: 12px
    fontWeight: '500'
    lineHeight: 16px
    letterSpacing: 0.05em
rounded:
  sm: 0.25rem
  DEFAULT: 0.5rem
  md: 0.75rem
  lg: 1rem
  xl: 1.5rem
  full: 9999px
spacing:
  container-padding: 20px
  stack-gap: 16px
  inline-gap: 12px
  section-margin: 32px
---

## Brand & Style
The design system is engineered for a K12 AI-driven educational environment, specifically optimized for the WeChat Mini Program ecosystem. The brand personality balances **institutional trust** with **youthful vitality**, ensuring the interface feels like a reliable mentor rather than a cold machine.

The style is a fusion of **Corporate Modern** and **Soft Minimalism**. It utilizes heavy whitespace to reduce cognitive load—essential for students—and employs "Fresh Green" accents to signify growth and the iterative nature of AI learning. The emotional response is one of clarity, safety, and encouragement.

## Colors
The palette is anchored by **Safe Blue**, providing a professional foundation that evokes wisdom and focus. **Fresh Green** is used strategically for interactive elements, success states, and growth indicators, creating a vibrant contrast that feels approachable.

**Neutral Tones** are intentionally cool-leaning (#F5F7FA) to maintain a clean, clinical (but not cold) environment that allows educational content to take center stage. Status colors are calibrated for high legibility against white backgrounds to ensure students can quickly identify areas requiring attention.

## Typography
The typography system prioritizes readability across various mobile screen densities. **Plus Jakarta Sans** is used for headings to provide a friendly, rounded aesthetic that feels modern and optimistic. For body text, **Be Vietnam Pro** offers exceptional clarity and a contemporary rhythm that reduces eye strain during long reading sessions.

On WeChat Mini Program screens, we use a slightly larger base font size (16px) for primary content to ensure accessibility for younger students and parents. Information hierarchy is established through weight shifts rather than aggressive size changes to maintain a compact, mobile-friendly layout.

## Layout & Spacing
This design system utilizes a **fluid grid** based on an 8px rhythmic scale. Given the constraints of the WeChat Mini Program environment, the layout relies on a standard **20px side margin** for containers to prevent content from touching the edge of the device frame.

Vertical spacing is generous to create a sense of "breath" between different learning modules. Components are grouped using **16px gaps**, while internal element spacing (like icon-to-text) uses **12px**. The layout adapts to larger tablet screens by centering the core content column and increasing side padding, rather than stretching elements excessively.

## Elevation & Depth
To maintain a "friendly" and "soft" atmosphere, the design system avoids harsh shadows. Depth is primarily achieved through **Tonal Layers**, where the primary background is `#F5F7FA` and interactive cards are pure white (#FFFFFF).

When physical elevation is required (e.g., for "Floating Action Buttons" or "Active Task Cards"), we use **Ambient Shadows**. These shadows are extremely diffused, using the Primary Blue color at 8% opacity rather than pure black, ensuring the UI feels light and integrated rather than heavy or disconnected.

## Shapes
The shape language is defined by significant roundedness to evoke a sense of safety and approachability. A **12px base radius** is applied to standard buttons and input fields. For larger structural elements, such as **Content Cards** or **Modal Sheets**, a **16px to 24px radius** is used. This softness removes the "clinical" edge often found in enterprise software, making the AI tools feel like a natural part of the student's learning journey.

## Components

### Buttons
Primary buttons utilize a full-color gradient from Safe Blue to a slightly lighter tint, with 12px rounded corners. The text is always centered and semi-bold.

### Cards
Educational cards (for lessons or quizzes) are the workhorse of the system. They feature a white background, a subtle 1px border (#E4E7ED), and the "Rounded-LG" (16px) corner radius. On tap, they should show a slight scale-down effect (98%) to provide tactile feedback.

### Input Fields
Search bars and quiz inputs use a background of `#F5F7FA` to contrast against white cards. They feature a 12px radius and use the Safe Blue for the cursor and active border state.

### Progress Bars
A signature AI component. Progress bars are "Pill-shaped" (32px radius) and use a Fresh Green fill to represent completed work, providing a satisfying visual reward for the student.

### Chips
Used for subject tags (e.g., "Math", "AI Tutor"). Chips use a low-opacity version of the Primary color (10% Safe Blue) with high-contrast blue text to ensure they are legible but secondary to main actions.