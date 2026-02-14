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
    const limit = document.getElementById('backupLimit').value;
    const url = limit ? `/api/v2/backups?limit=${limit}` : '/api/v2/backups';
    
    return fetch(url)
        .then(response => {
            const contentType = response.headers.get('Content-Type');
            if (contentType && contentType.includes('application/json')) {
                return response.json().then(data => ({ status: response.ok, data }));
            } else {
                return response.text().then(text => ({ status: response.ok, text }));
            }
        })
        .then(result => {
            const backupList = document.getElementById('backupList');
            backupList.innerHTML = '';
            
            if (!result.status || result.text) {
                backupList.innerHTML = `<li class="backuperror">${result.text || 'Failed to load backups'}</li>`;
                return;
            }
            
            const data = result.data;
            if (!data || data.length === 0) {
                backupList.innerHTML = '<li class="no-backups">No valid backup files found.</li>';
                return;
            }
            
            let animationCount = 0;
            data.forEach((backup) => {
                const li = document.createElement('li');
                li.className = 'backup-item';
                
                const backupType = "Dotsave"
                const fileName = "Backup Index: " + backup.Index;
                const formattedDate = "Created: " + new Date(backup.SaveTime).toLocaleString();
                const isDotsave = backupType === "Dotsave";
                
                li.innerHTML = `
                    <div class="backup-info">
                        <div class="backup-header">
                            <span class="backup-name">${fileName}</span>
                            <span class="backup-type ${backupType.toLowerCase()}">${backupType}</span>
                        </div>
                        <div class="backup-date">${formattedDate}</div>
                    </div>
                    <div class="backup-actions">
                        ${isDotsave ? `<button class="download-btn" onclick="downloadBackup(${backup.Index})">Download</button>` : ''}
                        <button class="restore-btn" onclick="restoreBackup(${backup.Index})">Restore</button>
                    </div>
                `;
                
                backupList.appendChild(li);
                
                if (animationCount < 20) {
                    setTimeout(() => {
                        li.classList.add('animate-in');
                    }, animationCount * 50);
                    animationCount++;
                }
            });
        })
        .catch(err => {
            console.error("Failed to fetch backups:", err);
            document.getElementById('backupList').innerHTML = '<li class="backuperror">Failed to load backups</li>';
        });
}

function getBackupType(backup) {
    return "Dotsave";
}

function fetchPlayers() {
    const playersDiv = document.getElementById('players');
    const playerList = document.getElementById('playerList');
    
    const playerImages = [
        "/static/playerimages/anna.webp",
        "/static/playerimages/dan.webp",
        "/static/playerimages/darragh.webp",
        "/static/playerimages/david.webp",
        "/static/playerimages/dean.webp",
        "/static/playerimages/garrison.webp",
        "/static/playerimages/ivette.webp",
        "/static/playerimages/john.webp",
        "/static/playerimages/julia.webp",
        "/static/playerimages/ove.webp",
        "/static/playerimages/pierre.webp",
        "/static/playerimages/rolf.webp",
        "/static/playerimages/ronald.webp",
    ];

    return fetch('/api/v2/server/status/connectedplayers')
        .then(response => response.json())
        .then(data => {
            playerList.innerHTML = '';
            
            if (!Array.isArray(data) || data.length === 0) {
                playersDiv.style.display = 'none';
                return;
            }

            playersDiv.style.display = 'block';
            let animationCount = 0;
            data.forEach(playerObj => {
                const player = Object.values(playerObj)[0];
                const li = document.createElement('li');
                li.className = 'player-item';
                
                // Create player item content
                const playerContent = document.createElement('div');
                playerContent.className = 'player-content';
                
                // Avatar
                const avatar = document.createElement('img');
                let persistedImage = sessionStorage.getItem(`playerImage_${player.steamID}`);
                if (!persistedImage) {
                    // Assign rnd image and persist it until page reload
                    persistedImage = playerImages[Math.floor(Math.random() * playerImages.length)];
                    sessionStorage.setItem(`playerImage_${player.steamID}`, persistedImage);
                }
                avatar.src = persistedImage;
                avatar.alt = `${player.username}'s avatar`;
                avatar.className = 'player-avatar';
                avatar.title = player.steamID;
                avatar.addEventListener('click', () => {
                    window.open(`https://steamcommunity.com/profiles/${player.steamID}`, '_blank');
                });
                
                const name = document.createElement('span');
                name.textContent = player.username;
                name.className = 'player-name';
                
                playerContent.appendChild(avatar);
                playerContent.appendChild(name);
                li.appendChild(playerContent);
                playerList.appendChild(li);
                
                // Animation
                if (animationCount < 20) {
                    setTimeout(() => {
                        li.classList.add('animate-in');
                    }, animationCount * 100);
                    animationCount++;
                }
            });
        })
        .catch(err => {
            console.error("Failed to fetch players:", err);
            playersDiv.style.display = 'none';
            playerList.textContent = 'Error loading players.';
        });
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
            showPopup('info', data);
        })
        .catch(err => console.error(`Failed to restore backup ${index}:`, err));
}

function downloadBackup(index) {
    const status = document.getElementById('status');
    status.hidden = false;
    typeTextWithCallback(status, 'Preparing download...', 20, () => {});
    
    fetch('/api/v2/backups/download', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ index: index })
    })
    .then(response => {
        if (!response.ok) {
            return response.json().then(err => { throw new Error(err.error || 'Download failed'); });
        }
        const disposition = response.headers.get('Content-Disposition');
        let filename = `backup_${index}.save`;
        if (disposition) {
            const match = disposition.match(/filename="(.+)"/);
            if (match) filename = match[1];
        }
        return response.blob().then(blob => ({ blob, filename }));
    })
    .then(({ blob, filename }) => {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        a.remove();
        status.hidden = true;
    })
    .catch(err => {
        console.error(`Failed to download backup ${index}:`, err);
        showPopup('error', 'Download failed: ' + err.message);
        status.hidden = true;
    });
}

function pollRecurringTasks() {
    window.gamserverstate = false;

    // Poll server status every 3.5 seconds
    const statusInterval = setInterval(() => {
        fetch('/api/v2/server/status')
            .then(response => response.json())
            .then(data => {
                updateStatusIndicator(data.isRunning, false, data.uptime);
                if (data.uuid) {
                    localStorage.setItem('gameserverrunID', data.uuid);
                }
            })
            .catch(err => {
                console.error("Failed to fetch server status:", err);
                updateStatusIndicator(false, true); // Set error state
            });
    }, 3500);

    // Poll connectred players every 10 seconds
    const playersInterval = setInterval(() => {
        fetchPlayers()
            .catch(err => {
                console.error("Failed to fetch connectedplayers:", err);
            });
    }, 10000);

    // Poll backups every 30 seconds
    const backupsInterval = setInterval(() => {
        fetchBackups()
            .catch(err => {
                console.error("Failed to fetch backups:", err);
            });
    }, 30000);
}

function updateStatusIndicator(isRunning, isError = false, uptime = '') {
    const indicator = document.getElementById('status-indicator');
    const uptimeDisplay = document.getElementById('uptime-display');
    
    if (isError) {
        indicator.className = 'status-indicator error';
        indicator.title = 'Error fetching server status';
        window.gamserverstate = false;
        if (uptimeDisplay) uptimeDisplay.style.display = 'none';
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

    // Show uptime only when server is running and uptime is not "0s"
    if (uptimeDisplay) {
        if (isRunning && uptime && uptime !== '0s') {
            uptimeDisplay.textContent = uptime;
            uptimeDisplay.style.display = 'inline-block';
        } else {
            uptimeDisplay.style.display = 'none';
        }
    }
}