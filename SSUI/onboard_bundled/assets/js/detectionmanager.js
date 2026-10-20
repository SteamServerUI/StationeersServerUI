const ROUTINE_ACTIONS = {
    eventlog: { group: 'Messages', label: 'Post to event log', field: 'text', placeholder: '{1} tried to spawn {2}' },
    discord_hub: { group: 'Messages', label: 'Post to Discord hub', field: 'text', placeholder: 'Heads up: {1} did something' },
    announce: { group: 'In game', label: 'Announce to players', field: 'text', placeholder: 'Nice try, {1}' },
    command: { group: 'In game', label: 'Run console command', field: 'text', placeholder: 'kick {1}' },
    save: { group: 'In game', label: 'Save world' },
    wait: { group: 'Flow', label: 'Wait', field: 'seconds' },
    webhook: { group: 'Flow', label: 'Call webhook', field: 'url', placeholder: 'https://example.com/hook' },
    restart: { group: 'Server', label: 'Restart server' },
    stop: { group: 'Server', label: 'Stop server' },
    start: { group: 'Server', label: 'Start server' },
    restore: { group: 'Server', label: 'Restore 2nd-newest backup' }
};

const TRIGGER_HINTS = {
    keyword: 'Matches when a log line contains this text (case-sensitive). Use {0} in actions for the whole line.',
    regex: 'Use {1}, {2}, … in actions for capture groups and {0} for the whole match. regex101.com is handy for testing.',
    event: 'For player events, use {1} for the player name and {2} for the SteamID.'
};

let triggerMode = 'keyword';

function setTriggerMode(mode) {
    triggerMode = mode;
    document.querySelectorAll('.routine-mode').forEach(button => button.classList.toggle('active', button.dataset.mode === mode));
    const pattern = document.getElementById('pattern');
    pattern.hidden = mode === 'event';
    pattern.placeholder = mode === 'regex' ? "Player '(.+?)' tried to thing-spawn '(.+?)'" : 'Text to look for in the log';
    pattern.classList.toggle('routine-mono', mode === 'regex');
    document.getElementById('trigger-event').hidden = mode !== 'event';
    document.getElementById('placeholder-hint').textContent = TRIGGER_HINTS[mode];
}

function actionOptions(selected) {
    const groups = {};
    Object.entries(ROUTINE_ACTIONS).forEach(([type, action]) => {
        (groups[action.group] ||= []).push(`<option value="${type}"${type === selected ? ' selected' : ''}>${action.label}</option>`);
    });
    return Object.entries(groups).map(([group, options]) => `<optgroup label="${group}">${options.join('')}</optgroup>`).join('');
}

function actionFields(type) {
    const action = ROUTINE_ACTIONS[type];
    if (action.field === 'seconds') {
        return '<span class="routine-seconds"><input type="number" class="routine-input action-value" min="1" max="3600" value="30"> seconds</span>';
    }
    if (action.field === 'url') {
        return `<input type="url" class="routine-input action-value" placeholder="${action.placeholder}">
            <input type="text" class="routine-input action-text" placeholder="Message to send (optional)">`;
    }
    if (action.field === 'text') {
        return `<input type="text" class="routine-input action-value" maxlength="1000" placeholder="${escapeHtml(action.placeholder)}">`;
    }
    return '<span class="routine-no-input">No settings needed</span>';
}

function addActionRow(type = 'eventlog') {
    const row = document.createElement('li');
    row.className = 'routine-action';
    row.innerHTML = `
        <span class="routine-action-number"></span>
        <div class="routine-action-body">
            <select class="routine-input action-type">${actionOptions(type)}</select>
            <div class="routine-action-fields">${actionFields(type)}</div>
        </div>
        <div class="routine-action-controls">
            <button type="button" class="routine-icon-button" data-move="-1" title="Move up">↑</button>
            <button type="button" class="routine-icon-button" data-move="1" title="Move down">↓</button>
            <button type="button" class="routine-icon-button danger" data-remove title="Remove">✕</button>
        </div>`;
    row.querySelector('.action-type').addEventListener('change', event => {
        row.querySelector('.routine-action-fields').innerHTML = actionFields(event.target.value);
    });
    row.querySelector('[data-remove]').addEventListener('click', () => {
        row.remove();
        renumberActions();
    });
    row.querySelectorAll('[data-move]').forEach(button => button.addEventListener('click', () => {
        const sibling = button.dataset.move === '-1' ? row.previousElementSibling : row.nextElementSibling;
        if (sibling) {
            button.dataset.move === '-1' ? sibling.before(row) : sibling.after(row);
            renumberActions();
        }
    }));
    document.getElementById('routine-actions').appendChild(row);
    renumberActions();
}

function renumberActions() {
    const rows = document.querySelectorAll('#routine-actions .routine-action');
    rows.forEach((row, index) => { row.querySelector('.routine-action-number').textContent = index + 1; });
    document.getElementById('add-action').disabled = rows.length >= 10;
}

function collectActions() {
    return [...document.querySelectorAll('#routine-actions .routine-action')].map(row => {
        const action = { type: row.querySelector('.action-type').value };
        const value = row.querySelector('.action-value');
        const text = row.querySelector('.action-text');
        if (value) action.value = value.value.trim();
        if (text && text.value.trim()) action.text = text.value.trim();
        return action;
    });
}

function resetBuilder() {
    document.getElementById('routine-name').value = '';
    document.getElementById('pattern').value = '';
    document.getElementById('routine-cooldown').value = 0;
    document.getElementById('routine-actions').innerHTML = '';
    addActionRow();
    setTriggerMode('keyword');
}

function triggerSummary(detection) {
    if (detection.type === 'event') {
        const option = document.querySelector(`#trigger-event option[value="${CSS.escape(detection.pattern)}"]`);
        return { badge: 'EVENT', text: option ? option.textContent : detection.pattern, mono: false };
    }
    return { badge: detection.type === 'regex' ? 'REGEX' : 'KEYWORD', text: detection.pattern, mono: true };
}

function actionSummary(action) {
    const meta = ROUTINE_ACTIONS[action.type] || { label: action.type };
    let text = meta.label;
    if (action.type === 'wait') text = `Wait ${action.value}s`;
    else if (action.type === 'webhook') { try { text = new URL(action.value).host; } catch { text = action.value; } }
    else if (meta.field === 'text') text = action.value;
    return `<span class="routine-chip${meta.group === 'Server' ? ' routine-chip-server' : ''}" title="${escapeHtml(meta.label)}"><span>${escapeHtml(text)}</span></span>`;
}

function renderRoutine(detection) {
    const trigger = triggerSummary(detection);
    const actions = (detection.actions || []).map(actionSummary).join('<span class="routine-arrow">→</span>');
    const cooldown = detection.cooldownSeconds ? `<span class="routine-meta">${detection.cooldownSeconds}s cooldown</span>` : '';
    return `
        <article class="routine-card">
            <div class="routine-card-header">
                <h4${detection.name ? '' : ' class="routine-untitled"'}>${escapeHtml(detection.name || 'Untitled routine')}</h4>
                <button type="button" class="slp-button slp-button-small danger routine-delete" data-id="${escapeHtml(detection.id)}">Delete</button>
            </div>
            <div class="routine-trigger">
                <span class="type-badge type-${trigger.badge.toLowerCase()}">${trigger.badge}</span>
                <span class="${trigger.mono ? 'routine-mono' : ''}">${escapeHtml(trigger.text)}</span>
            </div>
            <div class="routine-chain">${actions}</div>
            ${cooldown}
        </article>`;
}

function loadDetections() {
    const loader = document.getElementById('list-loader');
    const items = document.getElementById('detection-items');
    loader.style.display = 'block';

    fetch('/api/v3/detections')
        .then(response => {
            if (!response.ok) throw new Error('Could not load routines');
            return response.json();
        })
        .then(result => {
            loader.style.display = 'none';
            const detections = result.detections || [];
            items.innerHTML = detections.length
                ? detections.map(renderRoutine).join('')
                : '<div class="empty-list">No routines yet. Build one below.</div>';
            items.querySelectorAll('.routine-delete').forEach(button => button.addEventListener('click', () => deleteDetection(button.dataset.id)));
        })
        .catch(error => {
            loader.style.display = 'none';
            showNotification('Error: ' + error.message, 'error');
        });
}

function submitDetection() {
    const data = {
        name: document.getElementById('routine-name').value.trim(),
        type: triggerMode,
        pattern: triggerMode === 'event' ? document.getElementById('trigger-event').value : document.getElementById('pattern').value.trim(),
        actions: collectActions(),
        cooldownSeconds: parseInt(document.getElementById('routine-cooldown').value, 10) || 0
    };

    fetch('/api/v3/detections', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    })
        .then(response => response.json().catch(() => ({})).then(result => {
            if (!response.ok) throw new Error(result.error || 'Could not save the routine');
        }))
        .then(() => {
            resetBuilder();
            showNotification('Routine saved', 'success');
            loadDetections();
        })
        .catch(error => showNotification('Error: ' + error.message, 'error'));
}

function deleteDetection(id) {
    fetch(`/api/v3/detections/${encodeURIComponent(id)}`, { method: 'DELETE' })
        .then(response => {
            if (!response.ok) throw new Error('Could not delete the routine');
            showNotification('Routine deleted', 'success');
            loadDetections();
        })
        .catch(error => showNotification('Error: ' + error.message, 'error'));
}

function showNotification(message, type) {
    const notification = document.getElementById('notification');
    notification.textContent = message;
    notification.className = `notification notification-${type}`;
    notification.style.display = 'block';

    setTimeout(() => {
        notification.style.display = 'none';
    }, 5000);
}

function escapeHtml(unsafe) {
    return String(unsafe ?? '')
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.routine-mode').forEach(button => button.addEventListener('click', () => setTriggerMode(button.dataset.mode)));
    document.getElementById('add-action').addEventListener('click', () => addActionRow());
    document.querySelector('.add-button').addEventListener('click', submitDetection);
    resetBuilder();
    loadDetections();
});
