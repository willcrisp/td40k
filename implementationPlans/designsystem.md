# Design System Specification

## 1. Overview & Creative North Star: The Machine Spirit
This design system is built to evoke the "Machine Spirit"—the intersection of ancient, gothic industrialism and cold, tactical data processing. Our Creative North Star is **"Gothic Brutalism."** We are moving away from the "clean" and "friendly" aesthetic of modern SaaS. Instead, the UI should feel like a heavy, physical terminal bolted into the bulkhead of a starship.

To break the "standard template" look, we utilize **Hard-Edged Asymmetry.** Layouts should feel intentional and structural, using heavy blocks of color, staggered container heights, and "mechanical" overlaps where one panel appears to be physically bolted over another. This is not just a digital interface; it is a tactical relic.

## 2. Colors: Tonal Warfare
Our palette is dominated by the darkness of the void, punctuated by the "hazard" warnings of the machine and the "parchment" of ancient records.

### The "No-Line" Rule
Traditional 1px borders are strictly prohibited for layout sectioning. In this system, boundaries are defined by **Industrial Tonal Shifts.** A sidebar does not have a line; it is simply `surface-container-lowest` sitting against a `surface` background. The eye should perceive depth through weight, not outlines.

### Surface Hierarchy & Nesting
Treat the UI as a series of physical plates.
*   **Base Layer:** `surface` (#131313) for the main viewport.
*   **Sub-Panels:** Use `surface-container-low` for secondary data.
*   **Active Modules:** Use `surface-container-high` for interactive terminal windows.
*   **Nesting:** To create a "recessed" look, place a `surface-container-lowest` field inside a `surface-container-highest` panel.

### The "Glass & Gradient" Rule
While the system is "gritty," it is not flat. Use **Anodized Gradients** for primary actions. A CTA should not be a flat red; it should transition from `primary-container` (#a31317) to `primary` (#ffb4ab) at a 45-degree angle to mimic the sheen of painted metal. Floating tactical overlays should utilize `backdrop-blur` (12px) paired with a 40% opacity `surface-container-high` to create a "targeting glass" effect.

## 3. Typography: The Edict and The Data
Typography is the bridge between the thematic "Gothic" and the functional "Tactical."

*   **Display & Headlines (Space Grotesk):** These are our "stencils." Used for unit names, phase titles, and sector designations. They must be set in ALL CAPS with a slightly tighter letter-spacing (-0.02em) to feel like stamped metal.
*   **Body & Titles (Inter/System-UI):** This is the "Logos." High-readability sans-serif for weapon profiles, rule text, and stat lines.
*   **Tonal Contrast:** Large `display-lg` headers should be high-contrast (`on-surface`), while supporting technical data should use `label-sm` in `tertiary` (#ffb952) to mimic glowing amber phosphor screens.

## 4. Elevation & Depth: Structural Layering
We do not use shadows to imply "lightness." We use them to imply "mass."

*   **The Layering Principle:** Depth is achieved by stacking. A `surface-container-highest` tactical map sits "above" the `surface` background. The contrast in value provides the elevation.
*   **Ambient Shadows:** For floating modals (e.g., a unit's detail card), use an extra-diffused shadow: `box-shadow: 0 20px 50px rgba(0, 0, 0, 0.6)`. The shadow color should be a tinted dark-red or dark-blue (based on the `on-surface` tone) to simulate atmospheric occlusion.
*   **The "Ghost Border" Fallback:** If a container requires definition against a similar background, use a **Ghost Border**: 1px solid `outline-variant` (#5a403d) at 15% opacity. It should be barely felt, only perceived.
*   **Riveted Corners:** To reinforce the industrial theme, use pseudo-elements (::before/::after) to place 2px "rivets" or L-shaped "brackets" in the corners of `surface-container-highest` elements.

## 5. Components

### Buttons: Tactical Actuators
*   **Primary:** `primary-container` background, no border, sharp 0px corners. Use a subtle inner-bevel effect (1px top-light) to make it feel like a physical key.
*   **Secondary:** `surface-variant` background with a `tertiary` (#ffb952) left-hand accent bar (4px width).
*   **States:** On hover, the background color should "flicker" (opacity shift from 100% to 90%) to mimic a CRT screen.

### Input Fields: Data Slates
*   **Styling:** Inputs must use `surface-container-lowest` as the background.
*   **Active State:** When focused, the Ghost Border increases to 40% opacity in `tertiary` (Hazard Orange), creating a "system active" glow.
*   **Labels:** All labels use `label-md` in `on-surface-variant` and must be set in ALL CAPS.

### Cards & Lists: Heavy Cargo
*   **Structural Rule:** **Forbid dividers.** To separate list items, use a 12px vertical gap. For intense data density, use alternating background tints: `surface-container-low` and `surface-container-high`.
*   **Headers:** Card headers should have a "Header Bar"—a solid block of `secondary-container` (#0a4c6a) behind the title text.

### Additional Component: The "Status Seal"
A custom component for this system. A circular or hexagonal "Seal" using the `primary` red, used to denote "Confirmed" actions or "Objective Secured" statuses, placed with a slight 5-degree rotation to break the grid.

## 6. Do's and Don'ts

### Do:
*   **Embrace the Grid-Break:** Allow tactical icons or "rivets" to bleed outside the edge of their containers by 4-8px.
*   **Use Mono-spacing for Stats:** Use a monospaced font variant for numbers in unit profiles to ensure they align perfectly in tables.
*   **Color for Intent:** Use `tertiary` (Yellow/Orange) strictly for warnings and technical readouts. Use `primary` (Red) for combat actions.

### Don't:
*   **No Rounded Corners:** `0px` is the absolute rule. Any radius over 0px destroys the "Industrial" feel.
*   **No Center Alignment:** Tactical interfaces are built for efficiency. Use Left-aligned or Full-bleed layouts. Avoid centered "hero" sections.
*   **No Pure White:** Use `on-surface` (#e5e2e1) for text. Pure white (#FFFFFF) is reserved only for the most intense "Glow" effects or flash-points.