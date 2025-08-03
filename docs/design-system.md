# Backup Service Design System

This design system provides a comprehensive guide for creating professional, consistent, and accessible UI components for modern dark-themed applications.

## Color Palette

### Core Colors
```css
/* Primary Backgrounds */
--bg-primary: #1a1a1a;           /* Main app background */
--bg-secondary: #252525;         /* Sidebar/header background */
--bg-tertiary: #2d2d2d;          /* Card/panel background */
--bg-quaternary: #363636;        /* Elevated elements */

/* Interactive Backgrounds */
--bg-hover: #3a3a3a;            /* Hover state background */
--bg-active: #404040;           /* Active state background */
--bg-selected: #2d4a6b;         /* Selected state background */

/* Text Colors */
--text-primary: #ffffff;         /* Primary text (headers, important content) */
--text-secondary: #b4b4b4;       /* Secondary text (descriptions, labels) */
--text-muted: #7a7a7a;          /* Muted text (timestamps, less important) */
--text-inverse: #1a1a1a;        /* Text on light backgrounds */

/* Border Colors */
--border-primary: #3a3a3a;      /* Main borders and dividers */
--border-secondary: #2d2d2d;    /* Subtle borders */
--border-accent: #4a4a4a;       /* Highlighted borders */

/* Status Colors */
--status-success: #00d4aa;      /* Success/positive states */
--status-success-bg: #0d2818;   /* Success background */
--status-success-border: #1a4731; /* Success border */

--status-warning: #ff9500;      /* Warning states */
--status-warning-bg: #2d1a00;   /* Warning background */
--status-warning-border: #4a2c00; /* Warning border */

--status-error: #ff4757;        /* Error/negative states */
--status-error-bg: #2d0a0d;     /* Error background */
--status-error-border: #4a141a; /* Error border */

/* Interactive Colors */
--accent-primary: #4a9eff;      /* Primary actions, links */
--accent-primary-hover: #3d8bef; /* Primary hover state */
--accent-primary-active: #2d7adf; /* Primary active state */

--accent-secondary: #6c757d;    /* Secondary actions */
--accent-secondary-hover: #5a6268; /* Secondary hover state */
```

## Typography

### Font Family
```css
font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
```

### Font Scale
```css
/* Headers */
--font-size-h1: 1.5rem;         /* 24px - Main page titles */
--font-size-h2: 1.25rem;        /* 20px - Section headers */
--font-size-h3: 1.125rem;       /* 18px - Subsection headers */

/* Body Text */
--font-size-base: 1rem;         /* 16px - Default body text */
--font-size-sm: 0.875rem;       /* 14px - Secondary text */
--font-size-xs: 0.75rem;        /* 12px - Small text, captions */

/* Font Weights */
--font-weight-normal: 400;
--font-weight-medium: 500;
--font-weight-semibold: 600;
--font-weight-bold: 700;

/* Line Heights */
--line-height-tight: 1.2;
--line-height-normal: 1.5;
--line-height-relaxed: 1.6;
```

## Spacing System

### Base Unit: 8px
```css
--spacing-xs: 0.25rem;          /* 4px */
--spacing-sm: 0.5rem;           /* 8px */
--spacing-md: 1rem;             /* 16px */
--spacing-lg: 1.5rem;           /* 24px */
--spacing-xl: 2rem;             /* 32px */
--spacing-2xl: 3rem;            /* 48px */
--spacing-3xl: 4rem;            /* 64px */
```

### Component Spacing
```css
--card-padding: var(--spacing-md);     /* 16px */
--card-gap: var(--spacing-sm);         /* 8px */
--section-gap: var(--spacing-lg);      /* 24px */
--element-gap: var(--spacing-sm);      /* 8px */
```

## Border Radius
```css
--radius-sm: 0.25rem;           /* 4px - Small elements */
--radius-md: 0.375rem;          /* 6px - Cards, buttons */
--radius-lg: 0.5rem;            /* 8px - Large components */
--radius-xl: 1rem;              /* 16px - Modal, overlays */
```

## Shadows
```css
--shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.1);
--shadow-md: 0 4px 6px rgba(0, 0, 0, 0.1);
--shadow-lg: 0 10px 15px rgba(0, 0, 0, 0.1);
--shadow-hover: 0 4px 12px rgba(0, 0, 0, 0.15);
```

## Transitions
```css
--transition-fast: 150ms ease-in-out;
--transition-normal: 250ms ease-in-out;
--transition-slow: 350ms ease-in-out;
```

## Component Guidelines

### Cards
- **Background**: `--bg-tertiary`
- **Border**: `1px solid --border-primary`
- **Border Radius**: `--radius-md`
- **Padding**: `--card-padding`
- **Hover**: Background transitions to `--bg-hover`

### Buttons
- **Primary**: Background `--accent-primary`, text `--text-inverse`
- **Secondary**: Background `--bg-tertiary`, text `--text-primary`, border `--border-primary`
- **Danger**: Background `--status-error`, text `--text-inverse`
- **Padding**: `8px 16px` (small), `12px 24px` (normal)
- **Border Radius**: `--radius-md`

### Status Indicators
- **Success**: Green theme colors for completed/healthy states
- **Warning**: Orange theme colors for attention-needed states  
- **Error**: Red theme colors for failed/critical states
- **Neutral**: Gray theme colors for pending/unknown states

### Status Pills
- **Background**: Status-specific background colors with borders
- **Text**: Uppercase, small font with letter spacing
- **Padding**: `4px 8px` with small border radius
- **States**: Success (green), Warning (orange), Error (red)

### Metric Cards
- **Layout**: Centered content with large numbers and small labels
- **Number**: Large font size (1.875rem) with bold weight
- **Label**: Small uppercase text with letter spacing
- **Colors**: Financial color coding for positive/negative values
- **Hover**: Lift effect with shadow and border color change

### Data Tables
- **Row Height**: 2.5rem for optimal density
- **Cell Padding**: 8px 16px for comfortable spacing
- **Typography**: Mixed font sizes for hierarchy (xs labels, sm values)
- **Colors**: Financial colors for positive/negative values

### Typography Hierarchy
1. **Primary Text**: Main content, hostnames, important labels
2. **Secondary Text**: Descriptions, timestamps, secondary information
3. **Muted Text**: Less important details, helper text

### Interactive States
- **Hover**: Subtle background color change, smooth transition
- **Focus**: Visible outline for accessibility
- **Active**: Slight press effect, darker background
- **Disabled**: Reduced opacity (0.5), no cursor interaction

## Accessibility Guidelines

- **Contrast Ratios**: Minimum 4.5:1 for normal text, 3:1 for large text
- **Focus Indicators**: Visible focus outlines on all interactive elements
- **Color Independence**: Information not conveyed by color alone
- **Touch Targets**: Minimum 44px for mobile interactions

## Grid System

- **Base Unit**: 8px grid system
- **Container Max Width**: No fixed max-width (fluid)
- **Grid Gaps**: Multiples of 8px (8px, 16px, 24px, 32px)
- **Breakpoints**:
  - Mobile: `< 768px`
  - Tablet: `768px - 1024px`
  - Desktop: `1024px - 1280px`
  - Large Desktop: `> 1280px`

## Implementation Notes

### CSS Custom Properties
All design tokens are implemented as CSS custom properties (CSS variables) for consistency and easy theming.

### Component Architecture
- Use design tokens consistently across all components
- Implement hover states and transitions for better UX
- Follow spacing system for consistent layouts
- Maintain accessibility standards

### Performance Considerations
- Use `transform` and `opacity` for animations
- Implement smooth transitions without janky effects
- Optimize for 60fps interactions

## Implementation Status

### ✅ Completed Features
- **Professional Color System**: Dark theme with 30+ color variables including financial colors
- **Typography Hierarchy**: 8 font sizes, 4 weights, 3 line heights
- **8px Grid System**: Consistent spacing using 8px base unit
- **Data Table Design**: Professional dense layouts inspired by financial interfaces
- **Status Pills**: Color-coded status badges with professional styling
- **Metric Dashboard Cards**: Large number displays with color-coded values
- **Financial Color Coding**: Green/red for positive/negative values
- **Micro-interactions**: Hover effects, transforms, and smooth animations
- **Status Indicators**: Color-coded status dots with hover animations
- **Professional Buttons**: 3D hover effects with slide animations
- **Responsive Grid**: Mobile-first grid system with breakpoints

### 🎨 Design System Features
- **30+ Color Variables**: Backgrounds, status, financial, and interactive colors
- **Typography Scale**: 8 font sizes + 4 font weights + 3 line heights  
- **Spacing System**: 7 spacing variables + 10 grid variables following 8px grid
- **Border Radius**: 4 radius options for different component sizes
- **Shadows**: 4 shadow levels for depth and hierarchy
- **Transitions**: 3 transition speeds for consistent animations
- **Table Components**: Row heights, cell padding, and border styling
- **Metric Cards**: Specialized variables for dashboard metrics

### 📱 Responsive Breakpoints
- Mobile: `< 768px`
- Tablet: `768px - 1024px` 
- Desktop: `1024px - 1280px`
- Large Desktop: `> 1280px`

### ✅ Validation Results
- **Build**: ✅ Successful compilation
- **TypeScript**: ✅ All type errors resolved
- **Linting**: ✅ Only style warnings (no errors)
- **Performance**: ✅ Optimized animations using transform/opacity