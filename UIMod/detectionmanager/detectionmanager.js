// Show active tab
function showTab(tabId) {
    document.querySelectorAll('.tab-content').forEach(tab => tab.classList.remove('active'));
    document.querySelectorAll('.tab-button').forEach(button => button.classList.remove('active'));

    document.getElementById(tabId).classList.add('active');
    document.querySelector(`.tab-button[data-tab="${tabId}"]`).classList.add('active');

    if (tabId === 'detection-list-tab') {
        loadDetections();
    }
}

// Toggle detection type
function setupDetectionTypeToggle() {
    const toggle = document.getElementById('detection-type-toggle');
    const typeLabel = document.getElementById('detection-type-label');
    const typeInput = document.getElementById('type');
    const patternInfo = document.getElementById('pattern-info');
    const messageInfo = document.getElementById('message-info');

    toggle.addEventListener('change', function() {
        if (this.checked) {
            typeLabel.textContent = 'Regex';
            typeInput.value = 'regex';
            patternInfo.textContent = 'Regular expression pattern (e.g., "Player (.+) has reached level (\\d+)")';
            messageInfo.textContent = 'Message to display when pattern is detected. Use {1}, {2}, etc. for captured groups';
        } else {
            typeLabel.textContent = 'Keyword';
            typeInput.value = 'keyword';
            patternInfo.textContent = 'Text to match exactly (case-sensitive)';
            messageInfo.textContent = 'Message to display when pattern is detected';
        }
    });
}

// Load detections
function loadDetections() {
    const loader = document.getElementById('list-loader');
    const detectionItems = document.getElementById('detection-items');

    loader.style.display = 'block';

    fetch('/api/v2/custom-detections')
        .then(response => {
            if (!response.ok) throw new Error('Failed to load detections');
            return response.json();
        })
        .then(detections => {
            loader.style.display = 'none';

            if (detections.length === 0) {
                detectionItems.innerHTML = '<div class="empty-list">No custom detections found. Add one to get started.</div>';
                return;
            }

            detectionItems.innerHTML = '';
            detections.forEach(detection => {
                const item = document.createElement('div');
                item.className = 'detection-item';
                item.innerHTML = `
                    <div><span class="type-badge type-${detection.type}">${detection.type}</span></div>
                    <div class="detection-item-pattern" title="${escapeHtml(detection.pattern)}">${escapeHtml(detection.pattern)}</div>
                    <div class="detection-item-message" title="${escapeHtml(detection.message)}">${escapeHtml(detection.message)}</div>
                    <div class="action-buttons">
                        <button class="delete-button" onclick="deleteDetection('${detection.id}')">Delete</button>
                    </div>
                `;
                detectionItems.appendChild(item);
            });
        })
        .catch(error => {
            loader.style.display = 'none';
            showNotification('Error: ' + error.message, 'error');
            console.error('Error loading detections:', error);
        });
}

// Submit detection
function submitDetection() {
    const form = document.getElementById('detection-form');
    const type = document.getElementById('type').value;
    const pattern = document.getElementById('pattern').value.trim();
    const message = document.getElementById('message').value.trim();

    if (!pattern || !message) {
        showNotification('Please fill in all fields', 'error');
        return;
    }

    const data = {
        type: type,
        pattern: pattern,
        eventType: 'CUSTOM_DETECTION',
        message: message
    };

    fetch('/api/v2/custom-detections', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            return response.text().then(text => { throw new Error(text || 'Failed to add detection'); });
        }
        return response.json();
    })
    .then(() => {
        form.reset();
        document.getElementById('detection-type-toggle').checked = false;
        document.getElementById('detection-type-label').textContent = 'Keyword';
        document.getElementById('type').value = 'keyword';
        showNotification('Detection added successfully', 'success');
        showTab('detection-list-tab');
    })
    .catch(error => {
        showNotification('Error: ' + error.message, 'error');
        console.error('Error adding detection:', error);
    });
}

// Delete detection
function deleteDetection(id) {

    fetch(`/api/v2/custom-detections/delete/?id=${id}`, { method: 'DELETE' })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text || 'Failed to delete detection'); });
            }
            showNotification('Detection deleted successfully', 'success');
            loadDetections();
        })
        .catch(error => {
            showNotification('Error: ' + error.message, 'error');
            console.error('Error deleting detection:', error);
        });
}

// Show notification
function showNotification(message, type) {
    const notification = document.getElementById('notification');
    notification.textContent = message;
    notification.className = `notification notification-${type}`;
    notification.style.display = 'block';

    setTimeout(() => {
        notification.style.display = 'none';
    }, 5000);
}

// Helper function to escape HTML
function escapeHtml(unsafe) {
    return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

// Event listeners
document.addEventListener('DOMContentLoaded', () => {
    loadDetections();
    setupDetectionTypeToggle();
    document.querySelectorAll('.tab-button').forEach(button => {
        button.addEventListener('click', () => showTab(button.getAttribute('data-tab')));
    });
    document.querySelector('.add-button').addEventListener('click', submitDetection);
});


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

// Event listeners for window focus and blur
window.addEventListener('focus', () => {
  resourceSaver(false); // Resume animations when page is in focus
});

window.addEventListener('blur', () => {
  resourceSaver(true); // Pause animations when page loses focus
});