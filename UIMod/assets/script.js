// /static/script.js
document.addEventListener('DOMContentLoaded', () => {
    typeText(document.querySelector('h1'), 30);
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

});

// Global references to EventSource objects
let outputEventSource = null;
let detectionEventSource = null;

// Utility function for typing text
function typeText(element, speed) {
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

// Utility function for typing text with a callback
function typeTextWithCallback(element, text, speed, callback) {
    if (element.dataset.isTyping === 'true') {
        clearTimeout(element.dataset.timeoutId);
    }

    element.textContent = '';
    element.dataset.isTyping = 'true';
    let i = 0;
    
    const typeChar = () => {
        if (i < text.length) {
            element.textContent += text.charAt(i++);
            const timeoutId = setTimeout(typeChar, speed);
            element.dataset.timeoutId = timeoutId;
        } else {
            element.dataset.isTyping = 'false';
            delete element.dataset.timeoutId;
            if (callback) setTimeout(callback, 50);
        }
    };
    typeChar();
}

// Tab management
function setupTabs() {
    showTab('console-tab');
}

function showTab(tabId) {
    document.querySelectorAll('.tab-content').forEach(tab => tab.classList.remove('active'));
    document.querySelectorAll('.tab-button').forEach(btn => btn.classList.remove('active'));
    const tab = document.getElementById(tabId);
    tab.classList.add('active');
    document.querySelector(`.tab-button[onclick*="showTab('${tabId}')"]`).classList.add('active');
}

// Server control functions
function startServer() {
    toggleServer('/start');
}

function stopServer() {
    toggleServer('/stop');
}

function toggleServer(endpoint) {
    const status = document.getElementById('status');
    fetch(endpoint)
        .then(response => response.text())
        .then(data => {
            status.hidden = false;
            typeTextWithCallback(status, data, 20, () => {
                setTimeout(() => status.hidden = true, 10000);
            });
        })
        .catch(err => console.error(`Failed to ${endpoint}:`, err));
}

// EventSource management
function closeEventSources() {
    [outputEventSource, detectionEventSource].forEach(source => {
        if (source) {
            source.close();
            console.log(`${source === outputEventSource ? 'Output' : 'Detection events'} stream closed`);
        }
    });
    outputEventSource = detectionEventSource = null;
}

function navigateTo(url) {
    closeEventSources();
    window.location.href = url;
}

// Detection events streaming
function fetchDetectionEvents() {
    const maxMessages = 500;
    const detectionConsole = document.getElementById('detection-console');
    
    const connect = () => {
        detectionEventSource = new EventSource('/events');
        
        detectionEventSource.onmessage = event => {
            const message = document.createElement('div');
            message.className = `detection-event ${getEventClassName(event.data)}`;
            
            const timestamp = document.createElement('span');
            timestamp.className = 'event-timestamp';
            timestamp.textContent = `${new Date().toLocaleTimeString()}: `;
            
            const content = document.createElement('span');
            content.textContent = event.data;
            
            message.append(timestamp, content);
            detectionConsole.appendChild(message);
            
            while (detectionConsole.childElementCount > maxMessages) {
                detectionConsole.firstChild.remove();
            }
            detectionConsole.scrollTop = detectionConsole.scrollHeight;
            
            const detectionTab = document.getElementById('detection-tab');
            if (!detectionTab.classList.contains('active')) {
                const tabButton = document.querySelector('.tab-button[onclick*="detection-tab"]');
                tabButton.classList.add('notification');
                setTimeout(() => tabButton.classList.remove('notification'), 3000);
            }
        };
        
        detectionEventSource.onopen = () => console.log("Detection events stream connected");
        
        detectionEventSource.onerror = () => {
            console.error("Detection events stream disconnected");
            detectionEventSource.close();
            detectionEventSource = null;
            if (window.location.pathname === '/') {
                setTimeout(connect, 2000);
            }
        };
    };
    connect();
}

function getEventClassName(eventText) {
    const checks = [
        ['Server is ready', 'event-server-ready'],
        ['Server is starting', 'event-server-starting'],
        ['Server error', 'event-server-error'],
        ['Player', 'connecting', 'event-player-connecting'],
        ['Player', 'ready', 'event-player-ready'],
        ['Player', 'disconnected', 'event-player-disconnect'],
        ['World Saved', 'event-world-saved'],
        ['Exception', 'event-exception']
    ];
    
    return checks.find(([text, , condition]) => 
        condition ? eventText.includes(text) && eventText.includes(condition) : eventText.includes(text)
    )?.[1] || '';
}

// Backup management
function fetchBackups() {
    fetch('/api/v2/backups?mode=classic')
        .then(response => response.text())
        .then(data => {
            const backupList = document.getElementById('backupList');
            backupList.innerHTML = '';
            
            if (data.trim() === "No valid backup files found.") {
                backupList.textContent = data;
            } else {
                data.split('\n').filter(Boolean).forEach(backup => {
                    const li = document.createElement('li');
                    li.className = 'backup-item';
                    li.innerHTML = `${backup} <button onclick="restoreBackup(${extractIndex(backup)})">Restore</button>`;
                    backupList.appendChild(li);
                });
            }
        })
        .catch(err => console.error("Failed to fetch backups:", err));
}

function extractIndex(backupText) {
    return backupText.match(/Index: (\d+)/)?.[1] || null;
}

function restoreBackup(index) {
    const status = document.getElementById('status');
    fetch(`/api/v2/backups/restore?index=${index}`)
        .then(response => response.text())
        .then(data => {
            status.hidden = false;
            typeTextWithCallback(status, data, 20, () => {
                setTimeout(() => status.hidden = true, 30000);
            });
        })
        .catch(err => console.error(`Failed to restore backup ${index}:`, err));
}

// Console initialization with SSE stream setup
function handleConsole() {
    const consoleElement = document.getElementById('console');
    consoleElement.innerHTML = '';
    const bootTitle = "Interface initializing...";
    const bootCompleteMessage = "Interface ready.🎮 Happy gaming! 🎮";
    const bugChance = Math.random();
    const bugMessage = "ERROR: Nuclear parts in airflow detected! Initiating repair sequence...";
    
    const funMessages = [
        "Calibrating quantum flux capacitors...",
        "Initializing player happiness modules...",
        "Checking for monsters under the server...",
        "Brewing coffee for the CPU...",
        "Charging laser sharks...",
        "Teaching AI to say 'please' and 'thank you'...",
        "Polishing pixels to a mirror shine...",
        "Convincing electrons to flow in the right direction...",
        "Rebooting atmospheric systems for the 17th time...",
        "Attempting to locate your body after that last airlock malfunction...",
        "Converting oxygen to errors at alarming efficiency...",
        "Persuading physics engine to acknowledge gravity exists...",
        "Calculating ways your base will catastrophically depressurize...",
        "Optimizing unity garbage collection (good luck with that)...",
        "Aligning planetary rotation with server tick rate...",
        "Patching holes in space-time continuum and your habitat...",
        "Convincing solar panels that 'sun' is not just a theoretical concept...",
        "Negotiating peace treaty between logic circuits and the laws of thermodynamics...",
        "Compressing atmosphere until your CPU begs for mercy...",
        "Measuring distance between you and nearest fatal bug...",
        "Attempting to explain 'pipe networks' to confused server hamsters...",
        "Calculating probability of survival (spoiler: it's low)...",
        "Wrangling rogue Unity instances back into containment...",
        "Sacrificing RAM to the gods of stable framerates...",
        "Convincing electrons to flow in the right direction... nope, the power grid's borked.",
        "Patching hull breaches with duct tape and prayers...",
        "Recalculating O2 levels... wait, why is it all CO2 now?",
        "Spinning up the fabricator... hope it doesn’t eat the server this time.",
        "Debugging Unity physics... object launched into orbit, send help.",
        "Warming up the furnace... or just setting the base on fire, 50/50 shot.",
        "Rerouting pipes... because who needs logical fluid dynamics anyway?",
        "Loading terrain... oh look, it’s floating 3 meters above the ground again.",
        "Processing ore... into a fine paste of lag and despair.",
        "Stabilizing frame rate... lol, just kidding, welcome to 12 FPS city.",
        "Checking for updates... new bug introduced, feature still broken!",
        "Assembling solar tracker... now it’s tracking the admin instead.",
        "Balancing gas mixtures... kaboom imminent, run you fool!"
    ];

    const addMessage = (text, color, style = 'normal') => {
        const div = document.createElement('div');
        div.textContent = text;
        div.style.color = color;
        div.style.fontStyle = style;
        consoleElement.appendChild(div);
        consoleElement.scrollTop = consoleElement.scrollHeight;
    };

    // Start with initializing message
    typeTextWithCallback(consoleElement, bootTitle, 30, () => {
        // Show two funny messages while connecting
        const messageIndex1 = Math.floor(Math.random() * funMessages.length);
        addMessage(funMessages[messageIndex1], '#0af', 'italic');

        let messageIndex2;
        do {
            messageIndex2 = Math.floor(Math.random() * funMessages.length);
        } while (messageIndex2 === messageIndex1);
        addMessage(funMessages[messageIndex2], '#0af', 'italic');

        // Set up the persistent console stream
        outputEventSource = new EventSource('/console');
        
        // Persistent message handler
        outputEventSource.onmessage = event => {
            const message = document.createElement('div');
            message.textContent = event.data;
            consoleElement.appendChild(message);
            consoleElement.scrollTop = consoleElement.scrollHeight;
        };

        outputEventSource.onopen = () => {
            console.log("Console stream connected");
            finishInitialization();
        };

        outputEventSource.onerror = () => {
            console.error("Console stream disconnected");
            outputEventSource.close();
            outputEventSource = null;
            addMessage("Warning: Console stream unavailable. Retrying...", '#ff0');
            if (window.location.pathname === '/') {
                setTimeout(() => {
                    if (!outputEventSource) {
                        // Re-run setup to reconnect
                        consoleElement.innerHTML = ''; // Clear console for fresh start
                        handleConsole();
                    }
                }, 2000);
            }
        };
    });

    function finishInitialization() {
        if (bugChance < 0.05) {
            addMessage(bugMessage, 'red');
            setTimeout(() => {
                addMessage("Repair complete. Continuing initialization...", 'green');
                completeBoot();
            }, 1000);
        } else {
            completeBoot();
        }
    }

    function completeBoot() {
        setTimeout(() => {
            addMessage(bootCompleteMessage, '#0f0');
            consoleElement.scrollTop = consoleElement.scrollHeight;
        }, 500);
    }
}

function createPlanet(container, size, orbitRadius, speed, color) {
    const orbit = document.createElement('div');
    orbit.classList.add('orbit');
    orbit.style.width = `${orbitRadius * 2}px`;
    orbit.style.height = `${orbitRadius * 2}px`;
    orbit.style.position = 'absolute';
    orbit.style.left = '50%';
    orbit.style.top = '50%';
    orbit.style.transform = 'translate(-50%, -50%)';
    
    // Add random delay to start animation at different points
    const randomDelay = -(Math.random() * speed); // Negative delay to offset start
    orbit.style.animation = `orbit ${speed}s linear infinite ${randomDelay}s`;

    const planet = document.createElement('div');
    planet.classList.add('planet');
    planet.style.width = `${size}px`;
    planet.style.height = `${size}px`;
    planet.style.position = 'absolute';
    planet.style.left = '0%';
    planet.style.top = '50%';
    planet.style.backgroundColor = color;
    planet.style.borderRadius = '50%';
    planet.style.boxShadow = `0 0 20px ${color}`;
    
    orbit.appendChild(planet);
    container.appendChild(orbit);
}


function pollServerStatus() {
    const statusInterval = setInterval(() => {
        fetch('/api/v2/server/status')
            .then(response => response.json())
            .then(data => {
                updateStatusIndicator(data.isRunning);
            })
            .catch(err => {
                console.error("Failed to fetch server status:", err);
                updateStatusIndicator(false, true); // Set error state
            });
    }, 1000); // Poll every second

    // Store the interval ID so we can clear it if needed
    window.statusPollingInterval = statusInterval;
}

function updateStatusIndicator(isRunning, isError = false) {
    const indicator = document.getElementById('status-indicator');
    
    if (isError) {
        indicator.className = 'status-indicator error';
        indicator.title = 'Error fetching server status';
        return;
    }
    
    if (isRunning) {
        indicator.className = 'status-indicator online';
        indicator.title = 'Server is running';
    } else {
        indicator.className = 'status-indicator offline';
        indicator.title = 'Server is offline';
    }
}