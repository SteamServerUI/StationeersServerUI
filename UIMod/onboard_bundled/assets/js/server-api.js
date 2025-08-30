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

function fetchPlayers() {
    const playersDiv = document.getElementById('players');
    const playerList = document.getElementById('playerList');
    
    const playerImages = [
        "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/2089290/ss_5ccc7eafd0d54f887595b688d008debf7dd3c398.600x338.jpg?t=1658208343",
        "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/2089290/ss_178b7c3190794ee5bcf93a63536c4e4f5fae747d.600x338.jpg?t=1658208343",
        "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/2089290/ss_25dc8c613e7507d1215f27884866b37279c66f99.600x338.jpg?t=1658208343",
        "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/2089290/ss_d8a6ff6c070a6f6a51730797640b6f698b226b51.600x338.jpg?t=1658208343",
        "https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/1038500/ss_adc371d38cd1dcdd268bd7907ff1473615779ad8.600x338.jpg?t=1693478415"
    ];

    fetch('/api/v2/server/status/connectedplayers')
        .then(response => response.json())
        .then(data => {
            playerList.innerHTML = '';
            
            if (!data || data.length === 0) {
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
                // Use player avatar or random fallback
                avatar.src = playerImages[Math.floor(Math.random() * playerImages.length)];
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