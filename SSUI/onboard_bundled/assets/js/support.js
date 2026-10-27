(function () {
    const panel = document.getElementById('support-panel');
    if (!panel) return;
    const overall = document.getElementById('support-overall');
    const findings = document.getElementById('support-findings');
    const recordButton = document.getElementById('support-record');
    const recordingButton = document.getElementById('support-recording-package');
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
        overall.textContent = text.running || 'Running diagnostics...';
        try {
            const report = await request('/api/v3/support/diagnostics');
            overall.textContent = `SSUI ${report.system.version} (${report.system.branch}) · ${report.overall}`;
            overall.dataset.overall = report.overall;
            findings.innerHTML = report.findings.length ? report.findings.map(finding => `
                <div class="access-list-item support-finding is-${escapeHTML(finding.status)}">
                    <div><strong>${escapeHTML(finding.code)}</strong> ${escapeHTML(finding.summary)}
                        ${finding.fix ? `<small>${escapeHTML(finding.fix)}</small>` : ''}</div>
                </div>`).join('') : `<p class="access-empty">${escapeHTML(text.allGood || 'Nothing needs attention.')}</p>`;
        } catch (error) {
            overall.textContent = error.message;
            findings.innerHTML = '';
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
        if (active && panel.open) pollTimer = window.setTimeout(loadRecording, 5000);
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
        const state = !status.connected ? linkText.offline : status.claimed ? linkText.claimed : linkText.waiting;
        document.getElementById('support-link-state').textContent = `${state} · ${(linkText.expires || 'Ends at {time}').replace('{time}', expires)}`;
        document.getElementById('support-link-code-box').hidden = !status.claimCode;
        document.getElementById('support-link-code').textContent = status.claimCode || '';
        const activity = status.activity || [];
        document.getElementById('support-link-activity').innerHTML = activity.length
            ? activity.slice().reverse().map(entry => `
                <div class="access-list-item support-finding${entry.ok ? ' is-ok' : ' is-fail'}">
                    <div><strong>${escapeHTML(entry.command)}</strong>
                        <small>${escapeHTML(new Date(entry.at).toLocaleTimeString())}</small></div>
                </div>`).join('')
            : `<p class="access-empty">${escapeHTML(linkText.noActivity || 'No requests yet.')}</p>`;
        if (panel.open) linkTimer = window.setTimeout(loadLink, 5000);
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
    panel.addEventListener('toggle', () => {
        if (!panel.open) {
            window.clearTimeout(linkTimer);
            return window.clearTimeout(pollTimer);
        }
        loadDiagnostics();
        if (canControl) {
            loadRecording();
            loadLink();
        }
    });
})();
