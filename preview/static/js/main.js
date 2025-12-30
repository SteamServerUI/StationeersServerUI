function resourceSaver(pause) {
    // Get space background once outside the loop
    const spaceBackground = document.getElementById('space-background');
    
    // Handle animation states for all elements
    document.querySelectorAll('*').forEach(element => {
      element.style.animationPlayState = pause ? 'paused' : 'running';
    });
    
    // Fade the space background in/out instead of abrupt display change
    if (pause) {
      // Fade out
      spaceBackground.style.transition = 'opacity 0.5s ease';
      spaceBackground.style.opacity = '0';
      // Only hide it after the fade completes
      setTimeout(() => {
        if (document.hasFocus() === false) { // Double-check we're still unfocused
          spaceBackground.style.display = 'none';
        }
      }, 500);
    } else {
      // Make it visible first, then fade in
      spaceBackground.style.display = 'block';
      // Use setTimeout to ensure the display change is processed before starting the fade
      setTimeout(() => {
        spaceBackground.style.transition = 'opacity 0.5s ease';
        spaceBackground.style.opacity = '1';
      }, 10);
    }
}

function toggleGPUSaver() {
    window.GPUSaverEnabled = !window.GPUSaverEnabled;
    localStorage.setItem('GPUSaverEnabled', window.GPUSaverEnabled);
}

// Event listeners for window focus and blur
window.addEventListener('focus', () => {
    if (window.GPUSaverEnabled) {
        resourceSaver(false); // Resume animations when page is in focus
    }
});

window.addEventListener('blur', () => {
    if (window.GPUSaverEnabled) {
        resourceSaver(true); // Pause animations when page loses focus
    }
});