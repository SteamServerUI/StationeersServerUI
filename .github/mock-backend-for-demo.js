// /static/script.js - MOCKED VERSION FOR GITHUB PREVIEW
// This is in the .github/workflows/ folder because I really couldn't figure out how where else to put it, realistically.
document.addEventListener('DOMContentLoaded', () => {
    typeText(document.querySelector('h1'), 30);
    setupTabs();
    mockDetectionEvents();
    mockBackups();
    mockConsole();
    // Create planets with size, orbit radius, speed, and color
    const planetContainer = document.getElementById('planet-container');
    createPlanet(planetContainer, 80, 650, 34, 'rgba(200, 100, 50, 0.7)');
    createPlanet(planetContainer, 50, 1000, 46, 'rgba(100, 200, 150, 0.5)');
    createPlanet(planetContainer, 30, 1250, 63, 'rgba(50, 150, 250, 0.6)');
    createPlanet(planetContainer, 70, 400, 28, 'rgba(200, 150, 200, 0.7)'); 
});

// Mock data for preview
const MOCK = {
    detectionEvents: [
        "🎮 [Gameserver] 🕑 Server is starting up...",
        "🎮 [Gameserver] ✅ Server process has started!",
        "🎮 [Gameserver] ⚙️ Setting StartLocalHost changed from False to True",
        "🎮 [Gameserver] 🔔 Server is ready to connect!",
        "🎮 [Gameserver] 💾 World Saved: BackupIndex: 9 UTC Time: 2025-03-30T12:40:08Z",
        "🎮 [Gameserver] 📡 Player BobTheBuilder connecting from 192.168.1.100",
        "🎮 [Gameserver] 📡 Player BobTheBuilder ready",
        "🎮 [Gameserver] 📡 Player SpaceCowboy connecting from 192.168.1.101",
        "🎮 [Gameserver] 📡 Player SpaceCowboy ready",
        "🎮 [Gameserver] 💀 Player BobTheBuilder disconnected",
        "🎮 [Gameserver] ❌ Exception in thread 'main': unity.Exception: random.unity.exeption  caught and handled",
        "🎮 [Gameserver] 🚨 Server is stopping...",
        "🎮 [Gameserver] 🚨 Server process has stopped!",
        "🎮 [Gameserver] 🕑 Server is starting up...",
        "🎮 [Gameserver] ✅ Server process has started!",
    ],
    backups: [
        "BackupIndex: 8, Created: 29.03.2025 18:59:47",
        "BackupIndex: 7, Created: 28.03.2025 12:30:45",
        "BackupIndex: 6, Created: 27.03.2025 08:15:22",
        "BackupIndex: 5, Created: 26.03.2025 19:20:00",
        "BackupIndex: 4, Created: 25.03.2025 16:45:11",
        "BackupIndex: 3, Created: 24.03.2025 14:30:00",
        "BackupIndex: 2, Created: 23.03.2025 10:15:22",
        "BackupIndex: 1, Created: 22.03.2025 08:45:11"
    ],
    consoleMessages: [
        "Preview mode active, simulating after-startconsole output",
        "***Stationeers - 0.2.5499.24517***",
        "loaded 48 systems successfully",
        "game manager initialized",
        "World Loaded in 0:0",
        "RocketNet Succesfully hosted with Address: 0.0.0.0 Port: 27016",
        "14:40:06: StartSession. config:",
        "gameName: Preview Server",
        "mapName: Preview Server",
        "No clients connected. Auto pause timer started (10000ms)",
        "Ready"
    ]
};

// Utility function for typing text
function typeText(element, speed) {
    const fullText = element.textContent;
    element.textContent = '';
    let i = 0;
    
    const typeChar = () => {
        if (i < fullText.length) {
            element.textContent += fullText.charAt(i++);
            setTimeout(typeChar, speed);
        }
    };
    typeChar();
}

// Utility function for typing text with a callback
function typeTextWithCallback(element, text, speed, callback) {
    element.textContent = '';
    let i = 0;
    
    const typeChar = () => {
        if (i < text.length) {
            element.textContent += text.charAt(i++);
            setTimeout(typeChar, speed);
        } else if (callback) {
            setTimeout(callback, 50);
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

// Mocked server control functions
function startServer() {
    mockToggleServer('start');
}

function stopServer() {
    mockToggleServer('stop');
}

function mockToggleServer(action) {
    const status = document.getElementById('status');
    const messages = {
        start: "Server started. Tthis is a preview - no actual server is running. But try to set this tool up yourself - it's easy! 😉",
        stop: "Server stopped. This is a preview - no actual server was running. Have you tried to set this tool up yourself yet? Come on, you even hit this button! 😉"
    };
    
    status.hidden = false;
    typeTextWithCallback(status, messages[action], 20, () => {
        setTimeout(() => status.hidden = true, 20000);
    });
}

function navigateTo(url) {
    window.location.href = url;
}

// Mock detection events
function mockDetectionEvents() {
    const detectionConsole = document.getElementById('detection-console');
    
    // Initial population
    addInitialEvents();
    
    // Then add new events periodically
    setInterval(() => {
        const randomEvent = MOCK.detectionEvents[Math.floor(Math.random() * MOCK.detectionEvents.length)];
        addDetectionEvent(randomEvent);
    }, 5000);
    
    function addInitialEvents() {
        // Add a subset of events initially
        const initialEvents = [
            "Preview mode active, simulating console output",
        ];
        
        initialEvents.forEach(eventText => {
            addDetectionEvent(eventText);
        });
    }
    
    function addDetectionEvent(eventText) {
        const message = document.createElement('div');
        message.className = `detection-event ${getEventClassName(eventText)}`;
        
        const timestamp = document.createElement('span');
        timestamp.className = 'event-timestamp';
        timestamp.textContent = `${new Date().toLocaleTimeString()}: `;
        
        const content = document.createElement('span');
        content.textContent = eventText;
        
        message.append(timestamp, content);
        detectionConsole.appendChild(message);
        
        const maxMessages = 500;
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
    }
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
    
    return checks.find(([text, condition]) => 
        condition ? eventText.includes(text) && eventText.includes(condition) : eventText.includes(text)
    )?.[1] || '';
}

// Mock backups functionality
function mockBackups() {
    const backupList = document.getElementById('backupList');
    backupList.innerHTML = '';
    
    MOCK.backups.forEach(backup => {
        const li = document.createElement('li');
        li.className = 'backup-item';
        li.innerHTML = `${backup} <button onclick="restoreBackup(${extractIndex(backup)})">Restore</button>`;
        backupList.appendChild(li);
    });
}

function extractIndex(backupText) {
    return backupText.match(/Index: (\d+)/)?.[1] || null;
}

function restoreBackup(index) {
    const status = document.getElementById('status');
    const message = `Restored backup with index ${index}. (This is a preview - no actual restoration occurred)`;
    
    status.hidden = false;
    typeTextWithCallback(status, message, 20, () => {
        setTimeout(() => status.hidden = true, 30000);
    });
}

// Mock console with simulated output
function mockConsole() {
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
        "Spinning up the fabricator... hope it doesn't eat the server this time.",
        "Debugging Unity physics... object launched into orbit, send help.",
        "Warming up the furnace... or just setting the base on fire, 50/50 shot.",
        "Rerouting pipes... because who needs logical fluid dynamics anyway?",
        "Loading terrain... oh look, it's floating 3 meters above the ground again.",
        "Processing ore... into a fine paste of lag and despair.",
        "Stabilizing frame rate... lol, just kidding, welcome to 12 FPS city.",
        "Checking for updates... new bug introduced, feature still broken!",
        "Assembling solar tracker... now it's tracking the admin instead.",
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

        // Simulate connection to server
        setTimeout(() => {
            // Add bug message occasionally
            if (bugChance < 0.05) {
                addMessage(bugMessage, 'red');
                setTimeout(() => {
                    addMessage("Repair complete. Continuing initialization...", 'green');
                    completeBootAndSimulateMessages();
                }, 1000);
            } else {
                completeBootAndSimulateMessages();
            }
        }, 1000);
    });

    function completeBootAndSimulateMessages() {
        // Show boot complete message
        addMessage(bootCompleteMessage, '#0f0');
        
        // Start simulating console messages
        let messageIndex = 0;
        
        function addNextMessage() {
            if (messageIndex < MOCK.consoleMessages.length) {
                addMessage(MOCK.consoleMessages[messageIndex], '#fff');
                messageIndex++;
                
                const delay = 2000 + Math.random() * 3000; // Random delay between 2-5 seconds
                setTimeout(addNextMessage, delay);
            } else {
                // When we've shown all messages, start again with random selection
                setTimeout(() => {
                    const randomMessage = MOCK.consoleMessages[Math.floor(Math.random() * MOCK.consoleMessages.length)];
                    addMessage(randomMessage, '#fff');
                    setTimeout(addNextMessage, 3000 + Math.random() * 5000);
                }, 3000);
            }
        }
        
        // Start showing console messages after a delay
        setTimeout(addNextMessage, 1500);
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