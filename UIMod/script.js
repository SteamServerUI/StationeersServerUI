// /static/script.js
document.addEventListener('DOMContentLoaded', () => {
    typeText(document.querySelector('h1'), 30);  // Type out the h1 on load
    setDefaultConsoleMessage();
    setupTabs();
    fetchOutput();
    fetchDetectionEvents();
    fetchBackups();
});

// Global references to EventSource objects
let outputEventSource = null;
let detectionEventSource = null;

function setupTabs() {
    showTab('console-tab');
}

function showTab(tabId) {
    document.querySelectorAll('.tab-content').forEach(tab => tab.classList.remove('active'));
    document.querySelectorAll('.tab-button').forEach(button => button.classList.remove('active'));
    document.getElementById(tabId).classList.add('active');
    document.querySelector(`.tab-button[onclick*="showTab('${tabId}')"]`).classList.add('active');
}

function typeText(element, speed) {
    const fullText = element.textContent; // Use textContent instead of innerHTML
    console.log(fullText);
    element.textContent = ''; // Clear using textContent
    let i = 0;
    
    function typeChar() {
        if (i < fullText.length) {
            element.textContent += fullText.charAt(i);
            i++;
            setTimeout(typeChar, speed);
            console.log(i, element.textContent)


        }
    }
    typeChar();
}

function typeTextWithCallback(element, text, speed, callback) {
    element.innerHTML = ''; // Clear the element
    let i = 0;
    
    function typeChar() {
        if (i < text.length) {
            element.innerHTML += text.charAt(i);
            i++;
            setTimeout(typeChar, speed);
        } else if (callback) {
            setTimeout(callback, 50);
        }
    }
    typeChar();
}

function startServer() {
    fetch('/start')
        .then(response => response.text())
        .then(document.getElementById('status').hidden = false)
        .then(data => typeTextWithCallback(document.getElementById('status'), data, 20))
        .then(setTimeout(() => document.getElementById('status').hidden = true, 10000));
}

function stopServer() {
    fetch('/stop')
        .then(response => response.text())
        .then(document.getElementById('status').hidden = false)
        .then(data => typeTextWithCallback(document.getElementById('status'), data, 20))
        .then(setTimeout(() => document.getElementById('status').hidden = true, 10000));
}

// Close SSE connections before navigation
function closeEventSources() {
    if (outputEventSource) {
        outputEventSource.close();
        outputEventSource = null;
        console.log("Output stream closed");
    }
    if (detectionEventSource) {
        detectionEventSource.close();
        detectionEventSource = null;
        console.log("Detection events stream closed");
    }
}

// Wrap navigation in a function to handle cleanup
function navigateTo(url) {
    closeEventSources();
    window.location.href = url;
}

function fetchOutput() {
    outputEventSource = new EventSource('/console');
    outputEventSource.onmessage = function(event) {
        const consoleElement = document.getElementById('console');
        const message = document.createElement('div');
        message.textContent = event.data;
        consoleElement.appendChild(message);
        consoleElement.scrollTop = consoleElement.scrollHeight;
    };
    outputEventSource.onerror = function() {
        console.error("Output stream disconnected");
        outputEventSource.close();
        outputEventSource = null;
        // Only reconnect if still on the main page
        if (window.location.pathname === '/') {
            setTimeout(fetchOutput, 2000);
        }
    };
}

function fetchDetectionEvents() {
    const maxMessages = 500;
    function connect() {
        detectionEventSource = new EventSource('/events');
        detectionEventSource.onmessage = function(event) {
            const detectionConsole = document.getElementById('detection-console');
            const message = document.createElement('div');
            message.className = getEventClassName(event.data);
            message.classList.add('detection-event');
            
            const timestamp = document.createElement('span');
            timestamp.className = 'event-timestamp';
            timestamp.textContent = new Date().toLocaleTimeString() + ': ';
            
            const content = document.createElement('span');
            content.textContent = event.data;
            
            message.appendChild(timestamp);
            message.appendChild(content);
            detectionConsole.appendChild(message);
            
            while (detectionConsole.children.length > maxMessages) {
                detectionConsole.removeChild(detectionConsole.firstChild);
            }
            detectionConsole.scrollTop = detectionConsole.scrollHeight;
            
            if (!document.getElementById('detection-tab').classList.contains('active')) {
                const tabButton = document.querySelector('.tab-button[onclick*="detection-tab"]');
                tabButton.classList.add('notification');
                setTimeout(() => tabButton.classList.remove('notification'), 3000);
            }
        };
        detectionEventSource.onopen = function() {
            console.log("Detection events stream connected");
        };
        detectionEventSource.onerror = function() {
            console.error("Detection events stream disconnected");
            detectionEventSource.close();
            detectionEventSource = null;
            // Only reconnect if still on the main page
            if (window.location.pathname === '/') {
                setTimeout(connect, 2000);
            }
        };
    }
    connect();
}

function getEventClassName(eventText) {
    if (eventText.includes('Server is ready')) return 'event-server-ready';
    if (eventText.includes('Server is starting')) return 'event-server-starting';
    if (eventText.includes('Server error')) return 'event-server-error';
    if (eventText.includes('Player') && eventText.includes('connecting')) return 'event-player-connecting';
    if (eventText.includes('Player') && eventText.includes('ready')) return 'event-player-ready';
    if (eventText.includes('Player') && eventText.includes('disconnected')) return 'event-player-disconnect';
    if (eventText.includes('World Saved')) return 'event-world-saved';
    if (eventText.includes('Exception')) return 'event-exception';
    return '';
}

function fetchBackups() {
    fetch('/backups')
        .then(response => response.text())
        .then(data => {
            const backupList = document.getElementById('backupList');
            backupList.innerHTML = '';
            if (data.trim() === "No valid backup files found.") {
                backupList.innerHTML = data;
            } else {
                const backups = data.split('\n');
                backups.forEach(backup => {
                    if (backup.trim()) {
                        const listItem = document.createElement('li');
                        listItem.classList.add('backup-item');
                        listItem.innerHTML = backup + ' <button onclick="restoreBackup(' + extractIndex(backup) + ')">Restore</button>';
                        backupList.appendChild(listItem);
                    }
                });
            }
        });
}

function extractIndex(backupText) {
    const match = backupText.match(/Index: (\d+)/);
    return match ? match[1] : null;
}

function restoreBackup(index) {
    fetch(`/restore?index=${index}`)
        .then(response => response.text())
        .then(document.getElementById('status').hidden = false)
        .then(data => typeTextWithCallback(document.getElementById('status'), data, 20))
        .then(setTimeout(() => document.getElementById('status').hidden = true, 30000));
}

function setDefaultConsoleMessage() {
    const consoleElement = document.getElementById('console');
    consoleElement.innerHTML = '';
    const bootTitle = "Interface booting...";
    const bootProgressStages = [
        "[                       ] 0%",
        "[####                   ] 20%",
        "[#####                  ] 30%",
        "[########               ] 40%",
        "[#############          ] 60%",
        "[##################     ] 80%",
        "[#######################] 100%"
    ];
    const bootCompleteMessage = "Interface boot complete.";
    const bugChance = Math.random();
    const bugMessage = "ERROR: Nuclear parts in airflow detected! Initiating repair sequence...";
    
    // Random fun messages that appear during boot
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
    
    // Add note that user can interact while booting
    const noteElement = document.createElement('div');
    noteElement.innerHTML = "<em>Note: You can use all controls while boot animation plays</em>";
    noteElement.style.opacity = "0.7";
    noteElement.style.fontSize = "0.9em";
    noteElement.style.marginTop = "5px";
    consoleElement.appendChild(noteElement);
    
    typeTextWithCallback(consoleElement, bootTitle, 50, () => {
        setTimeout(() => {
            let index = 0;
            const progressElement = document.createElement('div');
            consoleElement.appendChild(progressElement);
            
            // Add a random fun message
            const funMessageElement = document.createElement('div');
            funMessageElement.style.color = '#0af';
            funMessageElement.style.fontStyle = 'italic';
            funMessageElement.textContent = funMessages[Math.floor(Math.random() * funMessages.length)];
            consoleElement.appendChild(funMessageElement);
            
            const bootInterval = setInterval(() => {
                if (index < bootProgressStages.length) {
                    progressElement.textContent = bootProgressStages[index];
                    
                    // Add another random message halfway through
                    if (index === 2) {
                        const anotherFunMessage = document.createElement('div');
                        anotherFunMessage.style.color = '#0af';
                        anotherFunMessage.style.fontStyle = 'italic';
                        
                        // Get a different message than the first one
                        let messageIndex;
                        do {
                            messageIndex = Math.floor(Math.random() * funMessages.length);
                        } while (funMessageElement.textContent === funMessages[messageIndex]);
                        
                        anotherFunMessage.textContent = funMessages[messageIndex];
                        consoleElement.appendChild(anotherFunMessage);
                    }
                    
                    consoleElement.scrollTop = consoleElement.scrollHeight;
                    index++;
                } else {
                    clearInterval(bootInterval);
                    completeBootAnimation();
                }
            }, 100); // Faster interval
        }, 50); // Shorter initial delay
    });
    
    function completeBootAnimation() {
        if (bugChance < 0.05) {
            const bugElement = document.createElement('div');
            bugElement.style.color = 'red';
            bugElement.textContent = bugMessage;
            consoleElement.appendChild(bugElement);
            consoleElement.scrollTop = consoleElement.scrollHeight;
            setTimeout(() => {
                const repairMessage = "Repair complete. Continuing boot sequence...";
                const repairElement = document.createElement('div');
                repairElement.style.color = 'green';
                repairElement.textContent = repairMessage;
                consoleElement.appendChild(repairElement);
                consoleElement.scrollTop = consoleElement.scrollHeight;
            }, 2000);
        }
        // Clear any running intervals
        for (let i = 1; i < 9999; i++) window.clearInterval(i);
        
        // Complete the animation
        const completionElement = document.createElement('div');
        completionElement.innerHTML = bootCompleteMessage.replace(/\n/g, '<br>');
        completionElement.style.color = '#0f0';
        completionElement.style.fontWeight = 'bold';
        consoleElement.appendChild(completionElement);
        
        // Add a small celebratory message
        const celebrationElement = document.createElement('div');
        celebrationElement.innerHTML = "🎮 Happy gaming! 🎮";
        celebrationElement.style.textAlign = 'center';
        celebrationElement.style.marginTop = '10px';
        consoleElement.appendChild(celebrationElement);
        
        consoleElement.scrollTop = consoleElement.scrollHeight;
    }
}

const additionalCSS = `
.skip-animation-btn {
    position: absolute;
    top: 10px;
    right: 10px;
    background: rgba(0, 170, 255, 0.3);
    border: 1px solid rgba(0, 170, 255, 0.6);
    color: rgba(0, 255, 255, 0.8);
    border-radius: 4px;
    padding: 5px 10px;
    cursor: pointer;
    transition: all 0.2s;
}
.skip-animation-btn:hover {
    background: rgba(0, 170, 255, 0.5);
    color: rgba(0, 255, 255, 1);
}
`;

const extraStyle = document.createElement('style');
extraStyle.textContent = additionalCSS;
document.head.appendChild(extraStyle);