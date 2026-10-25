const ROUTINE_ACTIONS = {
    eventlog: { group: 'messages', field: 'text' },
    discord_hub: { group: 'messages', field: 'text' },
    announce: { group: 'ingame', field: 'text' },
    command: { group: 'ingame', field: 'text', placeholder: 'kick {1}' },
    save: { group: 'ingame' },
    wait: { group: 'flow', field: 'seconds' },
    webhook: { group: 'flow', field: 'url', placeholder: 'https://example.com/hook' },
    restart: { group: 'server' },
    stop: { group: 'server' },
    start: { group: 'server' },
    restore: { group: 'server' }
};

let text = {};

function actionKey(type) {
    return type.replace(/_(.)/g, (_, letter) => letter.toUpperCase());
}

function actionLabel(type) {
    const key = actionKey(type);
    return text['action' + key[0].toUpperCase() + key.slice(1)] || type;
}

function actionPlaceholder(type) {
    const key = actionKey(type);
    return text['placeholder' + key[0].toUpperCase() + key.slice(1)] || ROUTINE_ACTIONS[type].placeholder || '';
}

function groupLabel(group) {
    return text['group' + group[0].toUpperCase() + group.slice(1)] || group;
}

let triggerMode = 'keyword';

function setTriggerMode(mode) {
    triggerMode = mode;
    let modeButton;
    document.querySelectorAll('.routine-mode').forEach(button => {
        const active = button.dataset.mode === mode;
        button.classList.toggle('active', active);
        if (active) modeButton = button;
    });
    const pattern = document.getElementById('pattern');
    pattern.hidden = mode === 'event';
    pattern.placeholder = modeButton?.dataset.placeholder || '';
    pattern.classList.toggle('routine-mono', mode === 'regex');
    document.getElementById('trigger-event').hidden = mode !== 'event';
    document.getElementById('placeholder-hint').textContent = modeButton?.dataset.hint || '';
}

function actionOptions(selected) {
    const groups = {};
    Object.entries(ROUTINE_ACTIONS).forEach(([type, action]) => {
        (groups[action.group] ||= []).push(`<option value="${type}"${type === selected ? ' selected' : ''}>${escapeHtml(actionLabel(type))}</option>`);
    });
    return Object.entries(groups).map(([group, options]) => `<optgroup label="${escapeHtml(groupLabel(group))}">${options.join('')}</optgroup>`).join('');
}

function actionFields(type) {
    const action = ROUTINE_ACTIONS[type];
    if (action.field === 'seconds') {
        return `<span class="routine-seconds"><input type="number" class="routine-input action-value" min="1" max="3600" value="30"> ${escapeHtml(text.seconds || 'seconds')}</span>`;
    }
    if (action.field === 'url') {
        return `<input type="url" class="routine-input action-value" placeholder="${escapeHtml(actionPlaceholder(type))}">
            <input type="text" class="routine-input action-text" placeholder="${escapeHtml(text.webhookMessage || 'Message to send (optional)')}">`;
    }
    if (action.field === 'text') {
        return `<input type="text" class="routine-input action-value" maxlength="1000" placeholder="${escapeHtml(actionPlaceholder(type))}">`;
    }
    return `<span class="routine-no-input">${escapeHtml(text.noSettings || 'No settings needed')}</span>`;
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
            <button type="button" class="routine-icon-button" data-move="-1" title="${escapeHtml(text.moveUp || 'Move up')}">↑</button>
            <button type="button" class="routine-icon-button" data-move="1" title="${escapeHtml(text.moveDown || 'Move down')}">↓</button>
            <button type="button" class="routine-icon-button danger" data-remove title="${escapeHtml(text.remove || 'Remove')}">✕</button>
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
        return { badge: 'event', text: option ? option.textContent : detection.pattern, mono: false };
    }
    return { badge: detection.type === 'regex' ? 'regex' : 'keyword', text: detection.pattern, mono: true };
}

function actionSummary(action) {
    const meta = ROUTINE_ACTIONS[action.type] || {};
    const label = actionLabel(action.type);
    let summary = label;
    if (action.type === 'wait') summary = (text.waitSummary || 'Wait {s}s').replace('{s}', action.value);
    else if (action.type === 'webhook') { try { summary = new URL(action.value).host; } catch { summary = action.value; } }
    else if (meta.field === 'text') summary = action.value;
    return `<span class="routine-chip${meta.group === 'server' ? ' routine-chip-server' : ''}" title="${escapeHtml(label)}"><span>${escapeHtml(summary)}</span></span>`;
}

function renderRoutine(detection) {
    const trigger = triggerSummary(detection);
    const actions = (detection.actions || []).map(actionSummary).join('<span class="routine-arrow">→</span>');
    const cooldown = detection.cooldownSeconds ? `<span class="routine-meta">${escapeHtml((text.cooldownMeta || '{s}s cooldown').replace('{s}', detection.cooldownSeconds))}</span>` : '';
    return `
        <article class="routine-card">
            <div class="routine-card-header">
                <h4${detection.name ? '' : ' class="routine-untitled"'}>${escapeHtml(detection.name || text.untitled || 'Untitled routine')}</h4>
                <button type="button" class="slp-button slp-button-small danger routine-delete" data-id="${escapeHtml(detection.id)}">${escapeHtml(text.delete || 'Delete')}</button>
            </div>
            <div class="routine-trigger">
                <span class="type-badge type-${trigger.badge}">${escapeHtml(text['badge' + trigger.badge[0].toUpperCase() + trigger.badge.slice(1)] || trigger.badge)}</span>
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
                : `<div class="empty-list">${escapeHtml(text.noRoutines || 'No routines yet. Build one below.')}</div>`;
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
    text = document.getElementById('detectionmanager-config-tab')?.dataset || {};
    document.querySelectorAll('.routine-mode').forEach(button => button.addEventListener('click', () => setTriggerMode(button.dataset.mode)));
    document.getElementById('add-action').addEventListener('click', () => addActionRow());
    document.querySelector('.add-button').addEventListener('click', submitDetection);
    resetBuilder();
    loadDetections();
});
