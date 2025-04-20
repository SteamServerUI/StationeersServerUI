// /static/main.js
document.addEventListener('DOMContentLoaded', () => {
    window.GPUSaverEnabled = localStorage.getItem('GPUSaverEnabled') !== 'false';
    typeh1(document.querySelector('h1'), 30);
    if (window.location.pathname == '/') {
        setupTabs();
        fetchDetectionEvents();
        fetchBackups();
        handleConsole();
        pollServerStatus();
        // Create planets with size, orbit radius, speed, and color
        const planetContainer = document.getElementById('planet-container');
        createPlanet(planetContainer, 80, 650, 34, 'rgba(200, 100, 50, 0.7)');
        createPlanet(planetContainer, 50, 1000, 46, 'rgba(100, 200, 150, 0.5)');
        createPlanet(planetContainer, 30, 1250, 63, 'rgba(50, 150, 250, 0.6)');
        createPlanet(planetContainer, 70, 400, 28, 'rgba(200, 150, 200, 0.7)'); 
        console.warn("If you see errors for sscm.js or sscm.css, you may want to enable SSCM.");
    }
});

// Global references to EventSource objects
let outputEventSource = null;
let detectionEventSource = null;

function closeEventSources() {
    [outputEventSource, detectionEventSource].forEach(source => {
        if (source) {
            source.close();
            console.log(`${source === outputEventSource ? 'Output' : 'Detection events'} stream closed`);
        }
    });
    outputEventSource = detectionEventSource = null;
}

function typeh1(element, speed) {
  // Check if typing is already in progress
  if (element.dataset.isTyping === 'true') {
      // Optionally, clear the previous timeout (requires storing it)
      clearTimeout(element.dataset.timeoutId);
  }

  const fullText = element.textContent;
  element.textContent = '';
  element.dataset.isTyping = 'true'; // Mark as typing
  let i = 0;
  
  const typeChar = () => {
      if (i < fullText.length) {
          element.textContent += fullText.charAt(i++);
          const timeoutId = setTimeout(typeChar, speed);
          element.dataset.timeoutId = timeoutId; // Store timeout ID
      } else {
          element.dataset.isTyping = 'false'; // Done typing
          delete element.dataset.timeoutId;
      }
  };
  typeChar();
}

function navigateTo(url) {
    closeEventSources();
    window.location.href = url;
}

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