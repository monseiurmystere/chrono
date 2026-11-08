# Chrono: Architecture & Implementation Plan

A modern, highly customizable screen locker for Sway and Wayland compositors supporting `ext-session-lock-v1`, featuring an HTML/XML and SCSS-driven DOM layout engine.

---

## 1. Problem Statement & Motivation

Wayland screen lockers today suffer from significant customizability and layout limitations:

- **`swaylock` / `swaylock-effects`**: Strictly constrained to a circular ring indicator and basic color/blur parameters. There is no support for dynamic layouts, responsive multi-monitor designs, or contextual widgets (battery, system status, date formats).
- **`gtklock`**: Dependent on GTK3 widgets and CSS, which makes non-standard layouts cumbersome and introduces heavyweight toolkit dependencies for a security-critical utility.
- **`hyprlock`**: Employs an ad-hoc configuration syntax requiring manual absolute coordinate offsets (`position = X, Y`), making responsive designs across varying display resolutions and aspect ratios brittle.

### Chrono's Solution
Chrono introduces a **DOM-like layout engine** configured using familiar web standards:
1. **HTML/XML Templates**: Declarative hierarchical structure with flexbox-driven flow and built-in semantic widgets.
2. **SCSS / CSS Styling**: Familiar styling syntax with variables, nesting, flexbox layout properties, border radii, colors, and selector states.
3. **Reactive State Providers**: Out-of-the-box system data providers (Battery via Linux sysfs, Clock/Date, Host/User info, PAM authentication state).
4. **Wayland Native**: First-class support for `ext-session-lock-v1`, ensuring security and multi-output resilience.

---

## 2. High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                 User Configuration & Theme                  │
│       ~/.config/chrono/index.html & style.scss              │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                        Chrono Engine                        │
│                                                             │
│  ┌────────────────────────┐      ┌───────────────────────┐  │
│  │    State Providers     │      │   Virtual DOM Tree    │  │
│  │  - Battery (sysfs)     │─────▶│  - Elements & Nodes   │  │
│  │  - Clock / Date Ticker │      │  - Custom Widgets     │  │
│  │  - PAM Auth State      │      │  - Template Hydration │  │
│  └────────────────────────┘      └───────────┬───────────┘  │
│                                              │              │
│  ┌────────────────────────┐                  │              │
│  │   SCSS / CSS Parser    │                  │              │
│  │  - SCSS Compilation   │                  │              │
│  │  - Selector Matching   │                  │              │
│  │  - Style Cascade       │                  │              │
│  └───────────┬────────────┘                  │              │
│              │                               │              │
│              └───────────────┬───────────────┘              │
│                              ▼                              │
│                 ┌─────────────────────────┐                 │
│                 │      Layout Engine      │                 │
│                 │  - Box Model (Geometry) │                 │
│                 │  - Flexbox Solver       │                 │
│                 └────────────┬────────────┘                 │
└──────────────────────────────┼──────────────────────────────┘
                               │ Calculated Layout Tree
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                  Renderer & Protocol Layer                  │
│  - 2D Canvas Painting (Colors, Borders, Fonts, Pixmaps)     │
│  - Shared Memory (wl_shm) Buffer Management                 │
│  - Wayland ext-session-lock-v1 Protocol Presentation        │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. Core Component Specifications

### 3.1. DOM Engine (`pkg/dom`)
- Parses HTML/XML templates into a virtual node tree (`ElementNode`, `TextNode`).
- Supports reactive template hydration via Go's standard templating engine for dynamic attributes and text values.
- Built-in custom semantic components:
  - `<chrono-clock format="15:04" />`: High-precision time and date rendering.
  - `<chrono-battery [device="BAT0"] />`: Dynamic battery state (level, charging status, time remaining).
  - `<chrono-auth-input placeholder="..." />`: Secure password entry field with configurable masking characters.
  - `<chrono-auth-feedback />`: Displays status messages and authentication error prompts from PAM.

### 3.2. Styling & SCSS Engine (`pkg/style`)
- SCSS parsing and compilation down to standard CSS rules.
- Selector matching:
  - Element tags (`header`, `div`, `span`, custom elements).
  - Classes (`.lock-card`, `.status-bar`).
  - IDs (`#auth-input`).
  - Attribute selectors (`[data-state="charging"]`, `[data-status="error"]`).
- Cascade and inheritance computation:
  - Colors (`color`, `background-color`).
  - Typography (`font-family`, `font-size`, `font-weight`, `text-align`).
  - Box dimensions (`width`, `height`, `min-width`, `max-width`, `margin`, `padding`, `border`, `border-radius`).
  - Flexbox properties (`display: flex`, `flex-direction`, `justify-content`, `align-items`, `gap`, `flex-grow`).

### 3.3. Layout Solver (`pkg/layout`)
- Computes concrete pixel coordinates and dimensions (`x, y, width, height`) for all nodes across varying display resolutions.
- Box Model calculation: handles content-box vs border-box, margins, paddings, and borders.
- Flexbox Algorithm: implements row/column layouts, spacing distribution (`space-between`, `center`, `flex-start`, `flex-end`), alignment, and gap calculations without requiring an external heavy browser runtime.

### 3.4. 2D Painting & Rendering (`pkg/render`)
- Renders the computed layout nodes to an RGBA pixel buffer:
  - Solid fills, gradients, and background image rendering.
  - Crisp anti-aliased text rendering with TrueType/OpenType font loading.
  - Borders, outlines, and rounded corners (`border-radius`).
  - Input field cursor rendering and password character masking.

### 3.5. State Providers (`pkg/providers`)
- **Battery Provider**:
  - Reads `/sys/class/power_supply/` for lightweight, zero-dependency monitoring of battery level, status (Charging, Discharging, Full), and energy metrics.
- **Clock Provider**:
  - Ticker-driven updates aligned to second or minute boundaries to minimize CPU wakeups while ensuring precise time display.
- **System Provider**:
  - Hostname, username, display name, and current uptime.

### 3.6. Authentication & PAM Integration (`pkg/auth`)
- Asynchronous PAM worker running in an isolated execution thread to ensure UI animations and clock updates remain responsive while PAM blocks.
- Conversation function implementation supporting `PAM_PROMPT_ECHO_OFF`, `PAM_PROMPT_ECHO_ON`, `PAM_ERROR_MSG`, and `PAM_TEXT_INFO`.
- Emits structured state events (`AuthIdle`, `AuthAuthenticating`, `AuthSuccess`, `AuthFailed`) to the reactive DOM engine.

### 3.7. Wayland Client & Session Lock Protocol (`pkg/wayland`)
- Connects to the Wayland display socket.
- Binds to `ext-session-lock-v1` (wlroots/Sway standard protocol) and `wl_shm`.
- Creates lock surfaces for each active Wayland output (multi-monitor support with per-output scaling and resolution handling).
- Coordinates buffer commits, frame callbacks, and secure unlocking upon successful authentication.

---

## 4. Phased Implementation Roadmap

### Phase 1: Project Skeleton & Domain Models
- [x] Initialize Git repository with `base` branch.
- [x] Create comprehensive architecture and roadmap plan in `.plans/plan.md`.
- [ ] Initialize Go module `io.tylr.chrono`.
- [ ] Setup core directory layout and foundational package interfaces.

### Phase 2: State Providers & DOM Tree
- [ ] Implement Battery provider (`/sys/class/power_supply/`).
- [ ] Implement Clock provider with configurable time formats.
- [ ] Implement DOM parser (`pkg/dom`) supporting HTML/XML structures and template hydration.
- [ ] Unit tests for template hydration and DOM node creation.

### Phase 3: Styling Engine & Flexbox Layout
- [ ] Implement CSS stylesheet parser and selector matcher.
- [ ] Implement SCSS compilation pipeline.
- [ ] Implement Box Model and Flexbox layout solver (`pkg/layout`).
- [ ] Test layout solver against nested flexbox layouts.

### Phase 4: 2D Graphics Painting
- [ ] Implement software buffer canvas renderer (`pkg/render`).
- [ ] Implement font loading and text layout/rasterization.
- [ ] Implement background, border, and border-radius rendering.

### Phase 5: Wayland `ext-session-lock-v1` Integration
- [ ] Implement Wayland client connection and output discovery.
- [ ] Integrate `ext-session-lock-v1` protocol handling and surface allocation.
- [ ] Multi-monitor buffer presentation and frame loop.

### Phase 6: Authentication Engine & Event Loop
- [ ] Implement PAM worker with non-blocking conversation channel.
- [ ] Bind keyboard input handling from Wayland to password input node.
- [ ] Implement unlock teardown sequence upon `AuthSuccess`.

### Phase 7: Theming, CLI, and Default Assets
- [ ] Implement CLI flags, configuration loading (`~/.config/chrono/chrono.toml`).
- [ ] Create sleek, modern default theme (`index.html` + `style.scss`).
- [ ] End-to-end testing on Sway/Wayland session.
