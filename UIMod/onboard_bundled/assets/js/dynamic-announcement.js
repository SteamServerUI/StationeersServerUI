// dynamic-announcement.js
(function () {
    // Configuration - change only these values if needed
    const ANNOUNCEMENT_ID = 'dynamic-announcement';
    const JSON_URL = 'https://steamserverui.github.io/StationeersServerUI/dynamic-announcement.json';
    const FETCH_TIMEOUT = 8000; // ms

    // Find the announcement container
    const container = document.getElementById(ANNOUNCEMENT_ID);
    if (!container) {
        console.warn(`[Dynamic Announcement] Element #${ANNOUNCEMENT_ID} not found on page.`);
        return;
    }

    // Hide it initially (in case CSS shows it by default)
    container.style.display = 'none';

    // Helper: simple timeout for fetch
    const fetchWithTimeout = (url, options = {}, timeout = FETCH_TIMEOUT) => {
        return Promise.race([
            fetch(url, options),
            new Promise((_, reject) =>
                setTimeout(() => reject(new Error('Fetch timeout')), timeout)
            )
        ]);
    };

    // Main logic
    fetchWithTimeout(JSON_URL, { method: 'GET', cache: 'no-cache' })
        .then(response => {
            if (!response.ok) {
                if (response.status === 404) {
                    console.info('[Dynamic Announcement] No announcement (404) - staying hidden.');
                    return null;
                }
                throw new Error(`HTTP ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            if (!data) return; // 404 or empty

            // Validate required fields
            if (!data.headline || !data.bodyHtml) {
                console.warn('[Dynamic Announcement] JSON is missing required fields.');
                return;
            }

            // Optional date range check
            const now = Date.now();
            const start = data.validFrom ? new Date(data.validFrom).getTime() : null;
            const end = data.validUntil ? new Date(data.validUntil).getTime() : null;

            if ((start !== null && now < start) || (end !== null && now > end)) {
                console.info('[Dynamic Announcement] Current date is outside the valid range.');
                return;
            }

            // Fill the template
            const headerElement = container.querySelector('h3');
            if (headerElement) {
                headerElement.innerHTML = `
                    <span class="notice-icon">📢</span>
                    ${escapeHtml(data.headline)}
                `;
            }

            const contentDiv = container.querySelector('.collapsible-content');

            // Short description (optional)
            let shortHtml = '';
            if (data.shortDescription) {
                shortHtml = `<p>${escapeHtml(data.shortDescription)}</p>`;
            }

            // Warning (optional)
            let warningHtml = '';
            if (data.warningHtml) {
                warningHtml = `<p class="status-bad">${data.warningHtml}</p>`;
            }

            // Signature (optional)
            let signatureHtml = '';
            if (data.author || data.authorRole) {
                const author = data.author ? escapeHtml(data.author) : '';
                const role = data.authorRole ? escapeHtml(data.authorRole) : '';
                signatureHtml = `<p><i>${author}${author && role ? ' - ' : ''}${role}</i></p>`;
            }

            contentDiv.innerHTML = `
                ${shortHtml}
                <p>${data.bodyHtml}</p>
                ${warningHtml}
                <p></p>
                ${signatureHtml}
            `;

            // Show the announcement
            container.style.display = ''; // revert to CSS default (usually block)
            console.info('[Dynamic Announcement] Announcement loaded and displayed.');
        })
        .catch(err => {
            // On any error (network, timeout, JSON parse, etc.) just keep it hidden
            console.info('[Dynamic Announcement] Failed to load announcement:', err.message);
        });

    // Simple HTML escape utility (prevents XSS if you trust the JSON source)
    function escapeHtml(text) {
        if (typeof text !== 'string') return text;
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
})();