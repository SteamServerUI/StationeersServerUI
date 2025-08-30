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
    
    fetch(url)
        .then(response => response.json())
        .then(data => {
            const backupList = document.getElementById('backupList');
            backupList.innerHTML = '';
            
            if (!data || data.length === 0) {
                backupList.innerHTML = '<li class="no-backups">No valid backup files found.</li>';
                return;
            }
            
            let animationCount = 0;
            data.forEach((backup) => {
                const li = document.createElement('li');
                li.className = 'backup-item';
                
                const backupType = getBackupType(backup);
                const fileName = backup.BinFile.split('/').pop();
                const formattedDate = new Date(backup.ModTime).toLocaleString();
                
                li.innerHTML = `
                    <div class="backup-info">
                        <div class="backup-header">
                            <span class="backup-name">${fileName}</span>
                            <span class="backup-type ${backupType.toLowerCase()}">${backupType}</span>
                        </div>
                        <div class="backup-date">${formattedDate}</div>
                    </div>
                    <button class="restore-btn" onclick="restoreBackup(${backup.Index})">Restore</button>
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
            document.getElementById('backupList').innerHTML = '<li class="error">Failed to load backups</li>';
        });
}

function getBackupType(backup) {
    if (backup.BinFile && backup.XMLFile && backup.MetaFile) {
        return 'Legacy';
    } else if (backup.BinFile && !backup.XMLFile && !backup.MetaFile) {
        return 'Dotsave';
    }
    return 'Unknown';
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
        })
        .catch(err => console.error(`Failed to restore backup ${index}:`, err));
}

function pollRecurringTasks() {
    window.gamserverstate = false;

    // Poll server status every 3.5 seconds
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