# TUI Themes

## Design Intent

- Treat built-in themes as TUI-optimized interpretations, not exact copies of upstream palettes
- Accent colors form a coordinated palette; the TUI owns how they map to components and semantic roles
- Keep each accent readable against the background and visually distinct from the other accents
- Keep categorical colors distinct and readable against the theme background

## Changing Themes

- Theme values generate reusable TUI styles, list delegates, and Markdown styles; view logic consumes those derived styles and applies the terminal foreground and background
- When changing theme structure or semantics, trace all style generation, rendering, and runtime theme-switching consumers
- Contrast tests enforce role-specific readability floors and the muted/foreground hierarchy
