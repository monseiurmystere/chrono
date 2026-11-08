# Chrono

> A modern, highly-customizable screen locker for Sway and Wayland compositors.

Chrono is designed to solve the customizability problem in Wayland screen lockers by introducing a **DOM-like layout engine** driven by familiar web technologies (**HTML/XML** and **SCSS**) paired with native Wayland protocol integration (`ext-session-lock-v1`).

---

## Features

- **DOM-Driven Layout**: Flexible, responsive element hierarchy supporting the CSS box model and Flexbox layout—no manual coordinate math across multi-monitor setups.
- **SCSS / CSS Theming**: Full styling control with nesting, variables, borders, border radii, colors, and pseudo-states.
- **Reactive State Providers**: Built-in, zero-bloat state monitors:
  - **Battery**: Real-time Linux sysfs (`/sys/class/power_supply`) tracking percentage, charging status, and power metrics.
  - **Clock / Date**: High-precision, low-wakeup tick provider with customizable date/time formatting.
  - **Authentication**: Responsive state feedback for PAM conversation prompts and errors.
- **Wayland Native**: First-class support for `ext-session-lock-v1`, ensuring secure session locking and robust multi-output handling.

---

## Repository Structure

```
.
├── .plans/          # Architecture plans and implementation roadmaps
├── build/           # Build scripts, packaging, and CI configuration
├── doc/             # Documentation, user guides, and theming references
└── src/             # Go source code for Chrono
```

---

## Architecture & Roadmap

For detailed architecture, component design, and project milestones, see the [Architecture and Implementation Plan](.plans/plan.md).
