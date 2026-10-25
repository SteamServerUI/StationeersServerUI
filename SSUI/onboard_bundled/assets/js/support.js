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

    document.getElementById('support-actions').hidden = !canControl;
    document.getElementById('support-package').addEventListener('click', event =>
        download(event.currentTarget, '/api/v3/support/package'));
    recordingButton.addEventListener('click', () => download(recordingButton, '/api/v3/support/recording/package'));
    recordButton.addEventListener('click', toggleRecording);
    panel.addEventListener('toggle', () => {
        if (!panel.open) return window.clearTimeout(pollTimer);
        loadDiagnostics();
        if (canControl) loadRecording();
    });
})();
