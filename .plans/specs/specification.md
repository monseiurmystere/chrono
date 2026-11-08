# Chrono Technical Specification

**Version:** 1.0.0-draft  
**Status:** In Review  
**Target Platform:** Linux / Wayland (`ext-session-lock-v1`)  

---

## 1. Executive Summary & Goals

Chrono is a secure, highly customizable screen locker for Wayland compositors (such as Sway and wlroots-based window managers). It replaces brittle coordinate-based locker configurations with an expressive, web-standard declarative DOM architecture driven by **HTML/XML templates** and **SCSS/CSS stylesheets**, paired with an embedded layout solver and 2D canvas painter.

### 1.1. Core Objectives
- **Web-Standard Declarative Authoring**: Users build themes with structured markup (`index.html`) and stylesheets (`style.scss`), eliminating manual coordinate calculations.
- **Embedded Flexbox & Box Model**: Native, fast layout calculation without the memory overhead, startup latency, or security attack surface of an embedded browser engine.
- **Reactive State Providers**: First-class system monitoring built directly into the template engine (Battery via sysfs, Clock/Date, PAM authentication state).
- **Security-First Wayland Implementation**: Strict compliance with the `ext-session-lock-v1` Wayland protocol, zero-allocation credential handling, and safe crash degradation.

### 1.2. Non-Goals
- Full browser DOM/CSS compliance (e.g., complex CSS grid, floating layouts, arbitrary JavaScript execution, WebGL).
- Complex MPRIS album art fetching pipelines across arbitrary network sources (focusing instead on standard local and sysfs metrics).

---

## 2. System Architecture & Data Flow

```
                      ┌───────────────────────────────┐
                      │    ~/.config/chrono/          │
                      │  - index.html (Markup)        │
                      │  - style.scss (Styles)        │
                      │  - chrono.toml (Settings)     │
                      └───────────────┬───────────────┘
                                      │
                                      ▼
                      ┌───────────────────────────────┐
                      │      Template & SCSS Engine   │
                      │  - SCSS Compiler -> CSS AST   │
                      │  - HTML Parser -> Raw DOM AST │
                      └───────────────┬───────────────┘
                                      │
                                      ▼
 ┌────────────────────────┐   ┌───────────────────────────────┐
 │   State Providers      │   │    Virtual DOM Tree (V-DOM)   │
 │ - Battery (sysfs)      ├──▶│ - Go Template Interpolation   │
 │ - Clock / Time Ticker  │   │ - Custom Semantic Elements    │
 │ - PAM Auth Controller  │   │ - Dynamic Data Attributes     │
 └────────────────────────┘   └───────────────┬───────────────┘
                                              │
                                              ▼
                              ┌───────────────────────────────┐
                              │      Style Resolver           │
                              │ - Selector Specificity Match  │
                              │ - Cascade & Inheritance       │
                              └───────────────┬───────────────┘
                                              │
                                              ▼
                              ┌───────────────────────────────┐
                              │    Layout Engine (Flexbox)    │
                              │ - Box Model (Margin/Padding)  │
                              │ - Axis Dimensions & Gap       │
                              │ - Alignment & Distribution    │
                              └───────────────┬───────────────┘
                                              │ Computed Geometry
                                              ▼
                              ┌───────────────────────────────┐
                              │      2D Canvas Painter        │
                              │ - Backgrounds & Borders       │
                              │ - Border Radius (Rounded Box) │
                              │ - Text & Typography           │
                              │ - Password Masking & Cursor   │
                              └───────────────┬───────────────┘
                                              │ Pixel Buffers (shm)
                                              ▼
                              ┌───────────────────────────────┐
                              │  Wayland Session Lock Layer   │
                              │ - ext-session-lock-v1         │
                              │ - Per-Output Scaling & Buffer │
                              └───────────────────────────────┘
```

---

## 3. Subsystem Specifications

### 3.1. Template & DOM Engine (`pkg/dom`)

The DOM engine parses standard HTML/XML documents into a hierarchical node tree and applies reactive template hydration.

#### 3.1.1. Node Data Model
```go
type NodeType int

const (
    ElementNode NodeType = iota
    TextNode
    CommentNode
)

type Node struct {
    Type       NodeType
    TagName    string
    Attributes map[string]string
    Children   []*Node
    Parent     *Node
    Text       string

    // Layout and Styling metadata attached during pipeline passes
    Style      *ComputedStyle
    Box        *BoundingBox
}
```

#### 3.1.2. Semantic Custom Elements
Chrono defines specific custom elements that interface directly with internal providers:
- `<chrono-clock>`:
  - Attributes:
    - `format`: Go reference time layout string (e.g., `"15:04"`, `"Monday, January 2"`, `"3:04 PM"`). Default: `"15:04"`.
    - `timezone`: Optional IANA timezone identifier (e.g., `"Local"`, `"UTC"`).
- `<chrono-battery>`:
  - Attributes:
    - `device`: Battery name in sysfs (e.g., `"BAT0"`). Default: auto-detect primary battery.
  - Automatically injects data attributes: `data-state` (`"charging"`, `"discharging"`, `"full"`), `data-level` (integer percentage).
- `<chrono-auth-input>`:
  - Attributes:
    - `placeholder`: Text to display when no input is entered.
    - `mask`: Character used to mask password inputs (e.g., `"•"`, `"*"`, or `""` for hidden). Default: `"•"`.
    - `autofocus`: Boolean indicating whether to direct keyboard input to this element.
- `<chrono-auth-feedback>`:
  - Renders authentication status and error messages from PAM.
  - Automatically updates `data-status`: `"idle"`, `"authenticating"`, `"error"`, `"locked-out"`.

#### 3.1.3. Template Context Hydration
Templates are evaluated against a reactive context object:
```go
type RenderContext struct {
    User    UserContext
    Clock   ClockContext
    Battery BatteryContext
    Auth    AuthContext
}
```

---

### 3.2. Styling & SCSS Engine (`pkg/style`)

#### 3.2.1. SCSS Preprocessing
- Chrono embeds an SCSS preprocessor that compiles SCSS files into flat CSS ASTs at runtime.
- Supports variable interpolation (`$primary: #58a6ff`), nesting (`.container { .child { ... } }`), parent selector reference (`&`), and mixins.

#### 3.2.2. CSS Selector Matching & Specificity
Selectors are evaluated according to W3C specificity rules:
1. **ID Selectors**: `#password-field` (Weight: 1-0-0)
2. **Class & Attribute Selectors**: `.auth-card`, `[data-state="charging"]` (Weight: 0-1-0)
3. **Element Selectors**: `div`, `chrono-clock`, `span` (Weight: 0-0-1)

#### 3.2.3. Supported Property Palette
| Category | Property | Supported Values |
|---|---|---|
| **Display** | `display` | `flex`, `block`, `none` |
| **Flexbox** | `flex-direction` | `row`, `column`, `row-reverse`, `column-reverse` |
| | `justify-content` | `flex-start`, `flex-end`, `center`, `space-between`, `space-around`, `space-evenly` |
| | `align-items` | `flex-start`, `flex-end`, `center`, `stretch`, `baseline` |
| | `gap` | `<length>` (px, em, rem) |
| | `flex-grow` | `<number>` |
| | `flex-shrink` | `<number>` |
| | `flex-basis` | `<length>`, `auto` |
| **Box Model** | `width`, `height` | `<length>`, `<percentage>`, `auto` |
| | `min-width`, `max-width` | `<length>`, `<percentage>` |
| | `min-height`, `max-height`| `<length>`, `<percentage>` |
| | `margin` | `<length>` (top, right, bottom, left) |
| | `padding` | `<length>` (top, right, bottom, left) |
| | `box-sizing` | `border-box`, `content-box` |
| **Decoration** | `background-color` | Hex, RGB, RGBA, HSL, HSLA, named colors |
| | `border` | `<width> <style> <color>` |
| | `border-radius` | `<length>` (uniform or per-corner) |
| | `opacity` | `<number>` (0.0 to 1.0) |
| **Typography** | `color` | Hex, RGB, RGBA, named colors |
| | `font-family` | System font family name |
| | `font-size` | `<length>` (px, pt, rem) |
| | `font-weight` | `normal`, `bold`, `100`–`900` |
| | `text-align` | `left`, `center`, `right` |

---

### 3.3. Layout Engine (`pkg/layout`)

The layout engine executes a two-pass algorithm to compute exact screen coordinates:

1. **Measurement Pass (Bottom-Up)**:
   - Calculate intrinsic content sizes for text, inputs, and custom widgets.
   - Compute constrained box boundaries including padding and border dimensions.
2. **Distribution Pass (Top-Down)**:
   - Root container initializes with the Wayland output surface dimensions (`width, height`).
   - For `display: flex` containers:
     - Distribute remaining main-axis space according to `justify-content` and `gap`.
     - Allocate flex expansion based on `flex-grow` and `flex-shrink`.
     - Align items along the cross-axis based on `align-items`.
3. **Pixel Snapping**:
   - Round computed floating-point coordinates to integer device pixels to prevent visual blurring.

---

### 3.4. State Providers (`pkg/providers`)

#### 3.4.1. Battery Provider (`pkg/providers/battery`)
- Discovers batteries under `/sys/class/power_supply/BAT*`.
- Reads:
  - `capacity`: Integer 0–100.
  - `status`: `"Charging"`, `"Discharging"`, `"Full"`, `"Not charging"`.
  - `current_now` / `voltage_now` / `power_now`: Power draw statistics.
- Implements a dual-trigger update model:
  - Periodic polling (default: 5 seconds, configurable).
  - Linux `netlink` / `uevent` listener for instant reaction to AC adapter plug/unplug events.

#### 3.4.2. Clock Provider (`pkg/providers/clock`)
- Employs a self-synchronizing ticker:
  - Aligns ticks precisely with second or minute boundaries (`time.Now().Truncate(...)`).
  - Automatically adjusts sleep duration to prevent drift and avoid unnecessary wakeups.

#### 3.4.3. PAM Authentication Provider (`pkg/providers/auth`)
- Isolates PAM Cgo calls in an independent OS thread (`runtime.LockOSThread`).
- Communicates with the main UI loop over a non-blocking Go channel.
- State Machine:
  - `IDLE`: Awaiting user input.
  - `AUTHENTICATING`: Worker executing PAM verification.
  - `SUCCESS`: Authentication approved; signals Wayland unlock.
  - `FAILURE`: Authentication failed; displays error feedback, clears password buffer.
  - `LOCKED_OUT`: Rate-limit or lockout period enforced by PAM.

---

### 3.5. Wayland Integration (`pkg/wayland`)

#### 3.5.1. Session Lock Protocol (`ext-session-lock-v1`)
- Binds to `ext_session_lock_manager_v1`.
- Calls `ext_session_lock_manager_v1.lock` to request exclusive session lock.
- Awaits `ext_session_lock_v1.locked` before confirming session security.
- Allocates an `ext_session_lock_surface_v1` for each active `wl_output`.
- Commits initial rendered frames to all outputs before acknowledging completion.

#### 3.5.2. Multi-Monitor Handling
- Outputs are identified via `wl_output` global registry events.
- Handles dynamic hotplugging: creates new surfaces when displays are connected; cleans up buffers when displays are disconnected.
- Supports independent resolution scaling (`fractional-scale-v1` or `wl_output.scale`).

---

## 4. Security & Robustness Requirements

1. **Credential Safety**:
   - Password strings must never be stored as immutable Go strings where GC might defer cleanup.
   - Passwords must be held in byte slices (`[]byte`) and explicitly zeroed out immediately following PAM conversation completion.
2. **Crash Resilience**:
   - If Chrono crashes while holding an `ext-session-lock-v1` session, the Wayland compositor locks down the session permanently until another locker binds or the compositor is restarted.
3. **Signal Trapping**:
   - Ignores `SIGINT`, `SIGTERM`, and `SIGQUIT` until unlocking has completed via successful authentication.
