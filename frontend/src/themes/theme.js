// src/theme/theme.js

// Available Vars
//"--bg-primary"
//"--bg-secondary"
//"--bg-tertiary"
//"--bg-hover"
//"--bg-active"
//"--text-primary"
//"--text-secondary"
//"--text-accent"
//"--text-warning"
//"--border-color"
//"--accent-primary"
//"--accent-secondary"
//"--accent-tertiary"
//"--shadow-light"
//"--shadow-medium"
//"--transition-speed"
//"--top-nav-height"


// Define your themes
const themes = {
    forest: {
      name: "Forest Dark",
      properties: {
        "--bg-primary": "#1e1e1e",
        "--bg-secondary": "#252526",
        "--bg-tertiary": "#2d2d2d",
        "--bg-hover": "#3c3c3c",
        "--bg-active": "#3e4033",
        "--text-primary": "#d4d4d4",
        "--text-secondary": "#a9a9a9",
        "--text-accent": "#6a9955",
        "--text-warning": "#ce9178",
        "--border-color": "#3e3e3e",
        "--accent-primary": "#6a9955",
        "--accent-secondary": "#4d7240",
        "--accent-tertiary": "#5f7e52",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.3)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.4)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  
    vaxholmDark: {
      name: "Vaxholm Dark",
      properties: {
        "--bg-primary": "#121a12",
        "--bg-secondary": "#1b2a1b",
        "--bg-tertiary": "#243224",
        "--bg-hover": "#2e3a2e",
        "--bg-active": "#3a4a3a",
        "--text-primary": "#d9e6d9",
        "--text-secondary": "#a3b3a3",
        "--text-accent": "#7a9a7a",
        "--text-warning": "#c9a67a",
        "--border-color": "#2a3a2a",
        "--accent-primary": "#7a9a7a",
        "--accent-secondary": "#5f7a5f",
        "--accent-tertiary": "#6a8a6a",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.5)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.6)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  
    archipelagoPastel: {
      name: "Archipelago Pastel",
      properties: {
        "--bg-primary": "#2e3b3e",
        "--bg-secondary": "#3e4b4e",
        "--bg-tertiary": "#4e5b5e",
        "--bg-hover": "#5e6b6e",
        "--bg-active": "#6e7b7e",
        "--text-primary": "#dce7e7",
        "--text-secondary": "#b0c0c0",
        "--text-accent": "#a3c1ad",
        "--text-warning": "#d9bba3",
        "--border-color": "#4a5a5a",
        "--accent-primary": "#a3c1ad",
        "--accent-secondary": "#8bb394",
        "--accent-tertiary": "#7aa383",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.25)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.35)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  
    colorblindFriendly: {
      name: "Colorblind Friendly",
      properties: {
        "--bg-primary": "#121212",
        "--bg-secondary": "#1e1e1e",
        "--bg-tertiary": "#2a2a2a",
        "--bg-hover": "#383838",
        "--bg-active": "#454545",
        "--text-primary": "#ffffff",
        "--text-secondary": "#bfbfbf",
        "--text-accent": "#ffb300", // bright yellow for visibility
        "--text-warning": "#ff3b3b", // bright red
        "--border-color": "#666666",
        "--accent-primary": "#ffb300",
        "--accent-secondary": "#ffaa00",
        "--accent-tertiary": "#cc8800",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.7)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.8)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  
    cyberpunkGlow: {
      name: "Cyberpunk Glow",
      properties: {
        "--bg-primary": "#0a0a23",
        "--bg-secondary": "#1a1a3a",
        "--bg-tertiary": "#2a2a5a",
        "--bg-hover": "#3a3a7a",
        "--bg-active": "#4a4a9a",
        "--text-primary": "#e0e0ff",
        "--text-secondary": "#a0a0ff",
        "--text-accent": "#ff00ff",
        "--text-warning": "#ff4d4d",
        "--border-color": "#660066",
        "--accent-primary": "#ff00ff",
        "--accent-secondary": "#cc00cc",
        "--accent-tertiary": "#990099",
        "--shadow-light": "0 0 8px #ff00ff",
        "--shadow-medium": "0 0 16px #ff00ff",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },

    v1Classic: {
        name: "Stationeers Server UI (Classic)",
        properties: {
            "--bg-primary": "#0a0a14",
            "--bg-secondary": "#1b1b2f8f",
            "--bg-tertiary": "#2a2a5a",
            "--bg-hover": "#2a2a5a",
            "--bg-active": "#4a4a9a",
            "--text-primary": "#00fca9",
            "--text-secondary": "#00fca9",
            "--text-accent": "#00fca9",
            "--text-warning": "#ff4d4d",
            "--border-color": "#660066",
            "--accent-primary": "#0eefa9",
            "--accent-secondary": "#cc00cc",
            "--accent-tertiary": "#990099",
            "--shadow-light": "0 0 8px #0df2aa",
            "--shadow-medium": "0 0 16px #0df2aa",
            "--transition-speed": "250ms",
            "--top-nav-height": "3rem",
            "--sidebar-width": "150px",
            "--sidebar-collapsed-width": "60px",
        },
      },
  
    lightArchipelago: {
      name: "Light Archipelago",
      properties: {
        "--bg-primary": "#f0f4f3",
        "--bg-secondary": "#d9e4e1",
        "--bg-tertiary": "#c0d1cd",
        "--bg-hover": "#b0c4bf",
        "--bg-active": "#a0b4af",
        "--text-primary": "#2a3a33",
        "--text-secondary": "#4a5a53",
        "--text-accent": "#7a9a7a",
        "--text-warning": "#c97a5a",
        "--border-color": "#a0b0a8",
        "--accent-primary": "#7a9a7a",
        "--accent-secondary": "#5f7a5f",
        "--accent-tertiary": "#6a8a6a",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.1)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.15)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
    oceanBreeze: {
      name: "Ocean Breeze",
      properties: {
        "--bg-primary": "#1a2a38",
        "--bg-secondary": "#253545",
        "--bg-tertiary": "#2f4055",
        "--bg-hover": "#3a4c66",
        "--bg-active": "#4a5c76",
        "--text-primary": "#e0eaf0",
        "--text-secondary": "#b0c0d0",
        "--text-accent": "#68c1e8",
        "--text-warning": "#f0ad4e",
        "--border-color": "#456277",
        "--accent-primary": "#68c1e8",
        "--accent-secondary": "#4fa3ca",
        "--accent-tertiary": "#3a89b0",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.3)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.4)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  
    sunsetGlow: {
      name: "Sunset Glow",
      properties: {
        "--bg-primary": "#272133",
        "--bg-secondary": "#332940",
        "--bg-tertiary": "#3e304d",
        "--bg-hover": "#4b3a5d",
        "--bg-active": "#57446d",
        "--text-primary": "#f5e6ff",
        "--text-secondary": "#d1b6e1",
        "--text-accent": "#ff9e7a",
        "--text-warning": "#ffcc66",
        "--border-color": "#5d4970",
        "--accent-primary": "#ff9e7a",
        "--accent-secondary": "#e68a6a",
        "--accent-tertiary": "#cc775a",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.35)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.45)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  
    mintChocolate: {
      name: "Mint Chocolate",
      properties: {
        "--bg-primary": "#1e2721",
        "--bg-secondary": "#26322a",
        "--bg-tertiary": "#2e3d33",
        "--bg-hover": "#38493e",
        "--bg-active": "#425548",
        "--text-primary": "#e0f0e8",
        "--text-secondary": "#b0c5b8",
        "--text-accent": "#7fe0c3",
        "--text-warning": "#d9b382",
        "--border-color": "#3d4940",
        "--accent-primary": "#7fe0c3",
        "--accent-secondary": "#58c4a3",
        "--accent-tertiary": "#3ba483",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.3)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.4)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  
    lavenderFields: {
      name: "Lavender Fields",
      properties: {
        "--bg-primary": "#2b2440",
        "--bg-secondary": "#352e4e",
        "--bg-tertiary": "#3f385c",
        "--bg-hover": "#4a426a",
        "--bg-active": "#554c78",
        "--text-primary": "#ece8ff",
        "--text-secondary": "#c7c0e3",
        "--text-accent": "#b28dff",
        "--text-warning": "#ffad9c",
        "--border-color": "#4d4566",
        "--accent-primary": "#b28dff",
        "--accent-secondary": "#9a77e0",
        "--accent-tertiary": "#8360c6",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.35)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.45)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  
    nordAurora: {
      name: "Kirkeness",
      properties: {
        "--bg-primary": "#2e3440",
        "--bg-secondary": "#3b4252",
        "--bg-tertiary": "#434c5e",
        "--bg-hover": "#4c566a",
        "--bg-active": "#5e6779",
        "--text-primary": "#eceff4",
        "--text-secondary": "#d8dee9",
        "--text-accent": "#88c0d0",
        "--text-warning": "#ebcb8b",
        "--border-color": "#4c566a",
        "--accent-primary": "#88c0d0",
        "--accent-secondary": "#81a1c1",
        "--accent-tertiary": "#5e81ac",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.25)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.35)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    }
  };
  
  

// Get theme names as an array
const themeNames = Object.keys(themes);

// Current theme state
let currentTheme = 'forest';

// Apply theme to document
function applyTheme(themeName) {
    const theme = themes[themeName];
    if (!theme) return;

    currentTheme = themeName;
    
    // Apply each CSS variable
    const root = document.documentElement;
    Object.entries(theme.properties).forEach(([property, value]) => {
        root.style.setProperty(property, value);
    });
    
    // Save to localStorage
    localStorage.setItem('theme', themeName);
}

// Rotate to next theme
function nextTheme() {
    const currentIndex = themeNames.indexOf(currentTheme);
    const nextIndex = (currentIndex + 1) % themeNames.length;
    applyTheme(themeNames[nextIndex]);
}

// Get current theme name
function getCurrentTheme() {
    return currentTheme;
}

// Get all theme names
function getThemes() {
    return themeNames;
}

// Initialize theme from localStorage or default
function initTheme() {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme && themes[savedTheme]) {
        applyTheme(savedTheme);
    } else {
        applyTheme(currentTheme);
    }
}

// Export the service functions
export default {
    initTheme,
    applyTheme,
    nextTheme,
    getCurrentTheme,
    getThemes
};