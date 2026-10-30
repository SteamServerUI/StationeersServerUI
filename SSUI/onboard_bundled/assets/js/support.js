(function () {
    const panel = document.getElementById('support-config-tab');
    if (!panel) return;
    const overall = document.getElementById('support-overall');
    const findings = document.getElementById('support-findings');
    const recordButton = document.getElementById('support-record');
    const recordingButton = document.getElementById('support-recording-package');
    const refreshButton = document.getElementById('support-refresh');
    const canControl = window.SSUIAccess.can('server.control');
    const text = panel.dataset;
    let recording = null;
    let pollTimer;

    function escapeHTML(value) {
        return String(value ?? '').replace(/[&<>"']/g, character => ({
            '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#039;'
        })[character]);
    }

    async function request(url, options) {
        const response = await fetch(url, options);
        const result = await response.json().catch(() => ({}));
        if (!response.ok) throw new Error(result.error || 'The request failed.');
        return result;
    }

    async function loadDiagnostics() {
        refreshButton.disabled = true;
        overall.textContent = text.running || 'Running diagnostics...';
        delete overall.dataset.overall;
        try {
            const report = await request('/api/v3/support/diagnostics');
            document.getElementById('support-version').textContent = `SSUI ${report.system.version} (${report.system.branch})`;
            overall.textContent = report.overall;
            overall.dataset.overall = report.overall;
            findings.innerHTML = report.findings.length ? report.findings.map(finding => `
                <div class="support-finding is-${escapeHTML(finding.status)}">
                    <strong>${escapeHTML(finding.code)}</strong>
                    <p>${escapeHTML(finding.summary)}</p>
                    ${finding.fix ? `<small>${escapeHTML(finding.fix)}</small>` : ''}
                </div>`).join('') : `<div class="support-finding is-ok"><p>${escapeHTML(text.allGood || 'Nothing needs attention.')}</p></div>`;
        } catch (error) {
            overall.textContent = error.message;
            overall.dataset.overall = 'failed';
            findings.innerHTML = '';
        } finally {
            refreshButton.disabled = false;
        }
    }

    function renderRecording() {
        const active = Boolean(recording?.active);
        const size = `${((recording?.bytes || 0) / 1048576).toFixed(1)} MB`;
        const detail = `${size}${recording?.full ? `, ${text.recordingFull || 'full'}` : ''}`;
        recordButton.textContent = active
            ? (text.stopRecording || 'Stop recording ({size})').replace('{size}', detail)
            : text.startRecording || 'Start recording';
        recordingButton.hidden = active || !recording?.stoppedAt;
        window.clearTimeout(pollTimer);
        if (active && isOpen()) pollTimer = window.setTimeout(loadRecording, 5000);
    }

    async function loadRecording() {
        try {
            recording = await request('/api/v3/support/recording');
            renderRecording();
        } catch (error) {
            window.SSUIAccess.notify(error.message);
        }
    }

    async function toggleRecording() {
        recordButton.disabled = true;
        try {
            const action = recording?.active ? 'stop' : 'start';
            recording = await request(`/api/v3/support/recording/${action}`, { method: 'POST' });
            renderRecording();
            if (action === 'stop') loadDiagnostics();
        } catch (error) {
            window.SSUIAccess.notify(error.message);
        } finally {
            recordButton.disabled = false;
        }
    }

    async function download(button, url) {
        button.disabled = true;
        try {
            const response = await fetch(url);
            if (!response.ok) {
                const result = await response.json().catch(() => ({}));
                throw new Error(result.error || 'The download failed.');
            }
            const blob = await response.blob();
            const name = /filename="([^"]+)"/.exec(response.headers.get('Content-Disposition') || '')?.[1] || 'ssui-support.zip';
            const link = Object.assign(document.createElement('a'), { href: URL.createObjectURL(blob), download: name });
            link.click();
            URL.revokeObjectURL(link.href);
        } catch (error) {
            window.SSUIAccess.notify(error.message);
        } finally {
            button.disabled = false;
        }
    }

    const link = document.getElementById('support-link');
    const linkText = link.dataset;
    let linkTimer;

    function renderLink(status) {
        const active = Boolean(status?.active);
        document.getElementById('support-link-idle').hidden = active;
        document.getElementById('support-link-active').hidden = !active;
        window.clearTimeout(linkTimer);
        if (!active) return;

        const expires = status.expiresAt ? new Date(status.expiresAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : '';
        const state = status.connected ? (status.claimed ? linkText.claimed : linkText.waiting) : status.error ? linkText.offline : linkText.waiting;
        document.getElementById('support-link-state').textContent = `${state} · ${(linkText.expires || 'Ends at {time}').replace('{time}', expires)}`;
        document.getElementById('support-link-code-box').hidden = !status.claimCode;
        document.getElementById('support-link-code').textContent = status.claimCode || '';
        const activity = status.activity || [];
        document.getElementById('support-link-activity').innerHTML = activity.length
            ? activity.slice().reverse().map(entry => `
                <div class="support-activity-item${entry.ok ? ' is-ok' : ' is-fail'}">
                    <span>${escapeHTML(entry.command)}</span>
                    <time>${escapeHTML(new Date(entry.at).toLocaleTimeString())}</time>
                </div>`).join('')
            : `<p class="support-empty">${escapeHTML(linkText.noActivity || 'No requests yet.')}</p>`;
        if (isOpen()) linkTimer = window.setTimeout(loadLink, status.connected || status.error ? 5000 : 1000);
    }

    async function loadLink() {
        try {
            renderLink(await request('/api/v3/support/connection'));
        } catch (error) {
            window.SSUIAccess.notify(error.message);
        }
    }

    async function changeLink(button, options) {
        button.disabled = true;
        try {
            renderLink(await request('/api/v3/support/connection', options));
        } catch (error) {
            window.SSUIAccess.notify(error.message);
        } finally {
            button.disabled = false;
        }
    }

    document.getElementById('support-link-open').addEventListener('click', event => changeLink(event.currentTarget, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ minutes: Number(document.getElementById('support-link-minutes').value) })
    }));
    document.getElementById('support-link-close').addEventListener('click', event =>
        changeLink(event.currentTarget, { method: 'DELETE' }));
    document.getElementById('support-link-copy').addEventListener('click', async () => {
        try {
            await navigator.clipboard.writeText(document.getElementById('support-link-code').textContent);
        } catch (_) {
            window.getSelection().selectAllChildren(document.getElementById('support-link-code'));
        }
    });

    link.hidden = !canControl;
    document.getElementById('support-actions').hidden = !canControl;
    document.getElementById('support-package').addEventListener('click', event =>
        download(event.currentTarget, '/api/v3/support/package'));
    recordingButton.addEventListener('click', () => download(recordingButton, '/api/v3/support/recording/package'));
    recordButton.addEventListener('click', toggleRecording);
    refreshButton.addEventListener('click', loadDiagnostics);

    function isOpen() {
        return panel.classList.contains('active');
    }

    function open() {
        window.clearTimeout(linkTimer);
        window.clearTimeout(pollTimer);
        if (!isOpen()) return;
        loadDiagnostics();
        if (canControl) {
            loadRecording();
            loadLink();
        }
    }

    document.addEventListener('config-tab', open);
    open();
})();
