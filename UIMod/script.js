// /static/script.js
document.addEventListener('DOMContentLoaded', () => {
    typeText(document.querySelector('h1'), 60);  // Type out the h1 on load
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
    const fullText = element.innerHTML;
    element.innerHTML = '';
    let i = 0;
    function typeChar() {
        if (i < fullText.length) {
            element.innerHTML += fullText.charAt(i);
            i++;
            setTimeout(typeChar, speed);
        }
    }
    typeChar();
}

function typeTextWithCallback(element, text, speed, callback) {
    element.innerHTML = '';
    let i = 0;
    function typeChar() {
        if (i < text.length) {
            element.innerHTML += text.charAt(i);
            i++;
            setTimeout(typeChar, speed);
        } else if (callback) {
            callback();
        }
    }
    typeChar();
}

function startServer() {
    fetch('/start')
        .then(response => response.text())
        .then(data => typeTextWithCallback(document.getElementById('status'), data, 20));
}

function stopServer() {
    fetch('/stop')
        .then(response => response.text())
        .then(data => typeTextWithCallback(document.getElementById('status'), data, 20));
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
    outputEventSource = new EventSource('/output');
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
        .then(data => typeTextWithCallback(document.getElementById('status'), data, 20));
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
    const bootCompleteMessage = "System boot complete.\nINFO: Press 'Start' to launch the game server.";
    const bugChance = Math.random();
    const bugMessage = "ERROR: Nuclear parts in airflow detected! Initiating repair sequence...";

    typeTextWithCallback(consoleElement, bootTitle, 50, () => {
        setTimeout(() => {
            let index = 0;
            const progressElement = document.createElement('div');
            consoleElement.appendChild(progressElement);
            const bootInterval = setInterval(() => {
                if (index < bootProgressStages.length) {
                    progressElement.textContent = bootProgressStages[index];
                    consoleElement.scrollTop = consoleElement.scrollHeight;
                    index++;
                } else {
                    clearInterval(bootInterval);
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
                            setTimeout(() => {
                                const completionElement = document.createElement('div');
                                completionElement.innerHTML = bootCompleteMessage.replace(/\n/g, '<br>');
                                consoleElement.appendChild(completionElement);
                                consoleElement.scrollTop = consoleElement.scrollHeight;
                            }, 1000);
                        }, 2000);
                    } else {
                        setTimeout(() => {
                            const completionElement = document.createElement('div');
                            completionElement.innerHTML = bootCompleteMessage.replace(/\n/g, '<br>');
                            consoleElement.appendChild(completionElement);
                            consoleElement.scrollTop = consoleElement.scrollHeight;
                        }, 200);
                    }
                }
            }, 150);
        }, 100);
    });
}

const cssAnimation = `
@keyframes notification-glow {
    0% { box-shadow: 0 0 5px rgba(0, 255, 171, 0.5); }
    50% { box-shadow: 0 0 20px rgba(0, 255, 171, 1); }
    100% { box-shadow: 0 0 5px rgba(0, 255, 171, 0.5); }
}
.notification {
    animation: notification-glow 1s infinite;
}`;
const style = document.createElement('style');
style.textContent = cssAnimation;
document.head.appendChild(style);