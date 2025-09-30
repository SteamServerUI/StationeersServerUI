<script>
    import { apiFetch } from "../../../services/api";
  
    let systemInfo = $state(null);
    let error = $state(null);
    let debugInfo = $state(null);
    let showDebug = $state(false);
  
    async function fetchSystemInfo() {
        try {
            const response = await apiFetch('/api/v2/osstats', { raw: true });
            const data = await response.json();
            debugInfo = JSON.stringify(
                {
                    status: response.status,
                    headers: Object.fromEntries(response.headers.entries()),
                    body: data
                },
                null,
                2
            );
            if (data && typeof data === 'object' && Object.keys(data).length > 0) {
                systemInfo = { ...data };
                error = null;
            } else {
                throw new Error('Empty or invalid response data');
            }
        } catch (err) {
            error = `Failed to fetch system information: ${err.message}`;
            systemInfo = null;
            if (!debugInfo) {
                debugInfo = `Error: ${err.message}\nStack: ${err.stack}`;
            }
            console.error('Fetch error:', err);
        }
    }
  
    $effect(() => {
        fetchSystemInfo();
        const interval = setInterval(fetchSystemInfo, 30000);
        return () => clearInterval(interval);
    });
  </script>
  
  <div class="card info-card">
    <div class="card-header">
        <h3>System Information</h3>
        <div class="card-icon">ℹ️</div>
    </div>
    {#if error}
        <div class="error">{error}</div>
        <button class="retry-button" onclick={fetchSystemInfo}>Retry</button>
    {:else if systemInfo}
        <div class="resource-metrics">
            <div class="metric">
                <div class="metric-value">{typeof systemInfo.cpuUsage === 'number' ? systemInfo.cpuUsage.toFixed(2) : 'N/A'}%</div>
                <div class="metric-label">System CPU Usage</div>
                <div class="metric-bar" style="--value: {typeof systemInfo.cpuUsage === 'number' ? systemInfo.cpuUsage : 0}%; --color: var(--accent-primary);"></div>
            </div>
            <div class="metric">
                <div class="metric-value">{typeof systemInfo.memoryUsage === 'number' ? (systemInfo.memoryUsage/1024).toFixed(2) : 'N/A'} GB</div>
                <div class="metric-label">Memory Consumption</div>
                <div class="metric-bar" style="--value: {typeof systemInfo.memoryUsage === 'number' ? (systemInfo.memoryUsage/1024)*10 : 0}%; --color: var(--accent-secondary);"></div>
            </div>
            <div class="metric">
                <div class="metric-value">{typeof systemInfo.diskUsage === 'number' ? (100-systemInfo.diskUsage).toFixed(2) : 'N/A'}%</div>
                <div class="metric-label">Disk Space Unused</div>
                <div class="metric-bar" style="--value: {typeof systemInfo.diskUsage === 'number' ? (100-systemInfo.diskUsage) : 0}%; --color: var(--accent-tertiary);"></div>
            </div>
        </div>
        <div class="info-grid">
            <div class="info-item">
                <div class="info-label">OS Version</div>
                <div class="info-value">{systemInfo.osName || 'Temple OS v3.14.15'} {systemInfo.osVersion || ''}</div>
            </div>
            <div class="info-item">
                <div class="info-label">Kernel</div>
                <div class="info-value">{systemInfo.kernel || 'Cannot detect'}</div>
            </div>
            <div class="info-item">
                <div class="info-label">Uptime</div>
                <div class="info-value">{systemInfo.uptime || 'Cannot detect'}</div>
            </div>
            <div class="info-item">
                <div class="info-label">Backend IP Address</div>
                <div class="info-value">{systemInfo.backendIpAddress || 'Cannot detect'}</div>
            </div>
        </div>
        <div class="info-grid">
            <div class="info-item">
                <div class="info-label">Last Refresh</div>
                <div class="info-value">{systemInfo.lastRefreshTime ? new Date(systemInfo.lastRefreshTime).toLocaleString() : 'Cannot detect'}</div>
            </div>
        </div>
    {:else}
        <div class="loading">Loading system information...</div>
    {/if}
  </div>
  
  <style>
    .card {
        background-color: var(--bg-secondary);
        border-radius: 8px;
        padding: 1.25rem;
        box-shadow: var(--shadow-light);
        transition: transform 0.2s ease, box-shadow 0.2s ease;
        display: flex;
        flex-direction: column;
        border: 1px solid var(--border-color);
    }
  
    .card:hover {
        transform: translateY(-2px);
        box-shadow: var(--shadow-medium);
    }
  
    .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1rem;
        border-bottom: 1px solid var(--border-color);
        padding-bottom: 0.75rem;
    }
  
    .card h3 {
        margin: 0;
        color: var(--text-accent);
        font-size: 1.1rem;
        font-weight: 600;
    }
  
    .card-icon {
        font-size: 1.25rem;
        opacity: 0.8;
    }
  
    .info-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 0.75rem;
    }
  
    .info-item {
        margin-bottom: 0.5rem;
    }
  
    .info-label {
        font-size: 0.75rem;
        color: var(--text-secondary);
    }
  
    .info-value {
        font-size: 0.9rem;
        font-weight: 500;
    }
  
    .error {
        color: var(--error-color, #ff4d4f);
        font-size: 0.9rem;
        text-align: center;
        padding: 1rem;
    }
  
    .loading {
        font-size: 0.9rem;
        text-align: center;
        padding: 1rem;
        color: var(--text-secondary);
    }
  
    .resource-metrics {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 1rem;
        margin-bottom: 1.5rem;
    }
  
    .metric {
        text-align: center;
        align-content: end;
    }
  
    .metric-value {
        font-size: 1.2rem;
        font-weight: 600;
        margin-bottom: 0.25rem;
    }
  
    .metric-label {
        font-size: 0.8rem;
        color: var(--text-secondary);
        margin-bottom: 0.5rem;
    }
  
    .metric-bar {
        height: 4px;
        background-color: var(--bg-hover);
        border-radius: 2px;
        overflow: hidden;
    }
  
    .metric-bar::after {
        content: '';
        display: block;
        height: 100%;
        width: var(--value);
        background-color: var(--color);
        border-radius: 2px;
    }
  
    @media (max-width: 1024px) {
        .info-grid {
            grid-template-columns: 1fr;
        }
        .resource-metrics {
            grid-template-columns: 1fr;
        }
    }
  </style>