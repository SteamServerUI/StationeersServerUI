// /static/server-api.js

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

function triggerSteamCMD() {
    const status = document.getElementById('status');
    status.hidden = false;
    typeTextWithCallback(status, 'Running SteamCMD, please wait... ', 20, () => {
        fetch('/api/v2/steamcmd/run')
            .then(response => response.json())
            .then(data => {
                showPopup("info", data.message);
            })
            .catch(err => {
                typeTextWithCallback(status, 'Error: Failed to trigger SteamCMD', 20, () => {
                    setTimeout(() => status.hidden = true, 10000);
                });
                console.error(`Failed to trigger SteamCMD:`, err);
            });
    });
}

function fetchBackups() {
    fetch('/api/v2/backups?mode=classic')
        .then(response => response.text())
        .then(data => {
            const backupList = document.getElementById('backupList');
            backupList.innerHTML = '';
            
            if (data.trim() === "No valid backup files found.") {
                backupList.textContent = data;
            } else {
                let animationCount = 0; // Track number of animated items
                data.split('\n').filter(Boolean).forEach((backup) => {
                    const li = document.createElement('li');
                    li.className = 'backup-item';
                    li.innerHTML = `${backup} <button onclick="restoreBackup(${extractIndex(backup)})">Restore</button>`;
                    backupList.appendChild(li);
                    if (animationCount < 20) {
                        setTimeout(() => {
                            li.classList.add('animate-in');
                        }, animationCount * 100);
                        animationCount++;
                    }
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

function pollServerStatus() {
    window.gamserverstate = false;
    const statusInterval = setInterval(() => {
        fetch('/api/v2/server/status')
            .then(response => response.json())
            .then(data => {
                updateStatusIndicator(data.isRunning);
                if (data.uuid) {
                    localStorage.setItem('gameserverrunID', data.uuid);
                }
            })
            .catch(err => {
                console.error("Failed to fetch server status:", err);
                updateStatusIndicator(false, true); // Set error state
            });
    }, 3500); // Poll every 3.5 seconds (adjusted from 1000 to reduce server load checking the status each time)

    // Store the interval ID so we can clear it if needed
    window.statusPollingInterval = statusInterval;
}

function updateStatusIndicator(isRunning, isError = false) {
    const indicator = document.getElementById('status-indicator');
    
    if (isError) {
        indicator.className = 'status-indicator error';
        indicator.title = 'Error fetching server status';
        window.gamserverstate = false;
        return;
    }
    
    if (isRunning) {
        indicator.className = 'status-indicator online';
        indicator.title = 'Server is running';
        window.gamserverstate = true;
    } else {
        indicator.className = 'status-indicator offline';
        indicator.title = 'Server is offline';
        window.gamserverstate = false;
    }
}