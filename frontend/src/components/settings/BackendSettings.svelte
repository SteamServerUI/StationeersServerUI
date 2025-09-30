<!-- BackendSettings.svelte -->
<script>
  import { onMount, onDestroy } from 'svelte';
  import { backendConfig, setBackend, setActiveBackend, initializeApiService, apiFetch } from '../../services/api';
  
  let currentConfig = $state();
  let newBackendId = $state('');
  let newBackendUrl = $state('');
  let activeBackend = $state('');
  let backends = $state([]);
  let backendStatus = {};
  let testing = $state(false);
  let activeBackendStatus = $state({ status: 'unknown', lastChecked: null }); // New status for active backend only
  
  const unsubscribe = backendConfig.subscribe(value => {
    currentConfig = value;
    activeBackend = value.active;
    backends = Object.keys(value.backends);
    backends.forEach(id => {
      if (!(id in backendStatus)) {
        backendStatus[id] = { status: 'unknown', lastChecked: null };
      }
    });
  });
  
  onMount(() => {
    initializeApiService();
    const interval = setInterval(checkAllBackends, 10000);
    checkAllBackends();
    
    return () => {
      unsubscribe();
      clearInterval(interval);
    };
  });
  
  async function testBackend(id, silent = false) {
    if (testing) return;
    testing = true;
    
    const backendUrl = currentConfig.backends[id].url;
    try {
      const response = await apiFetch(`/api/v2/server/status`, {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' }
      });
      
      const status = {
        status: response.ok ? 'online' : 'offline',
        lastChecked: new Date()
      };
      
      backendStatus[id] = status;
      
      if (id === activeBackend) {
        activeBackendStatus = status;
        if (!silent) {
          //alert(`Connection to ${id} ${response.ok ? 'successful' : 'failed'}!`);
        }
      }
    } catch (error) {
      const status = { status: 'error', lastChecked: new Date() };
      backendStatus[id] = status;
      
      if (id === activeBackend) {
        activeBackendStatus = status;
        if (!silent) {
          //alert(`Error connecting to ${id}: ${error.message}`);
        }
      }
    } finally {
      testing = false;
      backendStatus = { ...backendStatus };
      activeBackendStatus = { ...activeBackendStatus };
    }
  }
  
  async function checkAllBackends() {
    for (const id of backends) {
      await testBackend(id, true); // Silent checks for auto-tests
    }
  }
  
  function addBackend() {
    if (!newBackendId || !newBackendUrl) return;
    setBackend(newBackendId, newBackendUrl);
    newBackendId = '';
    newBackendUrl = '';
  }
  
  function changeActiveBackend() {
    setActiveBackend(activeBackend);
    testBackend(activeBackend); // Non-silent test when manually switching
  }
  
  function removeBackend(id) {
    if (id === 'default') {
      alert('Cannot remove the default backend');
      return;
    }
    
    backendConfig.update(config => {
      if (config.active === id) {
        config.active = 'default';
        activeBackend = 'default';
        activeBackendStatus = backendStatus['default'] || { status: 'unknown', lastChecked: null };
      }
      delete config.backends[id];
      delete backendStatus[id];
      return config;
    });
  }
</script>

<div class="backend-settings">
  <div class="section-header">
    <h2>Backend Connections</h2>
  </div>
  
  <div class="current-backends">
    <div class="backend-controls">
      <h3>Configured Backends</h3>
      <div class="active-status">
        {#if activeBackendStatus.status === 'online'}
          🟢 Online
        {:else if activeBackendStatus.status === 'offline'}
          🔴 Offline
        {:else if activeBackendStatus.status === 'error'}
          ⚠️ Error
        {:else}
          ⚪ Unknown
        {/if}
        {#if activeBackendStatus.lastChecked}
          <span class="status-timestamp">
            (Last checked: {new Date(activeBackendStatus.lastChecked).toLocaleTimeString()})
          </span>
        {/if}
      </div>
      <button 
        class="action-button test-all" 
        onclick={() => testBackend(activeBackend)}
        disabled={testing}
      >
        {testing ? 'Testing...' : 'Test Active Backend'}
      </button>
    </div>
    
    <div class="backend-list">
      {#each backends as backendId}
        <div class="backend-item">
          <div class="backend-info">
            <h3>{backendId}</h3>
            <p>{currentConfig.backends[backendId].url}</p>
          </div>
          <div class="backend-status">
            <label class="status-toggle">
              <input 
                type="radio" 
                name="activeBackend" 
                value={backendId} 
                bind:group={activeBackend} 
                onchange={changeActiveBackend}
              />
              <span class="status-label">{backendId === activeBackend ? 'Active' : 'Set Active'}</span>
            </label>
          </div>
          <div class="backend-actions">
            {#if backendId !== 'default'}
              <button class="action-button danger" onclick={() => removeBackend(backendId)}>Remove</button>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  </div>
  
  <div class="add-backend">
    <div class="section-header">
      <h3>Add New Backend</h3>
    </div>
    
    <div class="form-container">
      <div class="form-group">
        <label for="backend-id">Backend ID</label>
        <input 
          type="text" 
          id="backend-id" 
          bind:value={newBackendId} 
          placeholder="e.g., production"
        />
      </div>
      
      <div class="form-group">
        <label for="backend-url">Backend URL</label>
        <input 
          type="text" 
          id="backend-url" 
          bind:value={newBackendUrl} 
          placeholder="e.g., https://api.example.com"
        />
      </div>
      
      <button class="primary-button" onclick={addBackend}>+ Add Backend</button>
    </div>
  </div>
</div>
<style>
  .backend-settings {
    display: flex;
    flex-direction: column;
    gap: 2rem;
    padding: 1.5rem;
    max-width: 800px;
    margin: 0 auto;
  }
  
  .section-header {
    padding-bottom: 0.75rem;
    border-bottom: 2px solid var(--border-color);
  }
  
  .section-header h2, .section-header h3 {
    margin: 0;
    color: var(--text-primary);
    font-weight: 600;
  }
  
  .backend-controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    gap: 1rem;
    flex-wrap: wrap;
  }
  
  .backend-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .backend-item {
    display: grid;
    grid-template-columns: 1fr auto auto;
    align-items: center;
    gap: 1.5rem;
    background-color: var(--bg-secondary);
    border-radius: 8px;
    padding: 1.25rem;
    box-shadow: var(--shadow-light);
    transition: transform var(--transition-speed) ease;
  }
  
  .backend-item:hover {
    transform: translateY(-2px);
  }
  
  .backend-info {
    flex: 1;
  }
  
  .backend-info h3 {
    margin: 0 0 0.5rem 0;
    font-size: 1.2rem;
    font-weight: 500;
    color: var(--text-primary);
  }
  
  .backend-info p {
    margin: 0;
    color: var(--text-secondary);
    font-size: 0.95rem;
  }
  
  .status-timestamp {
    font-size: 0.85rem;
    color: var(--text-secondary);
  }
  
  .backend-status {
    display: flex;
    align-items: center;
  }
  
  .status-toggle {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 4px;
    transition: background-color var(--transition-speed) ease;
  }
  
  .status-toggle:hover {
    background-color: var(--bg-hover);
  }
  
  .status-label {
    color: var(--text-primary);
    font-size: 0.95rem;
    font-weight: 500;
  }
  
  .backend-actions {
    display: flex;
    gap: 0.75rem;
  }

  .active-status {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1rem;
    color: var(--text-primary);
  }
  
  .action-button {
    background-color: var(--bg-tertiary);
    color: var(--text-primary);
    border: 1px solid var(--border-color);
    padding: 0.5rem 1rem;
    border-radius: 6px;
    font-size: 0.95rem;
    cursor: pointer;
    transition: all var(--transition-speed) ease;
  }
  
  .action-button:hover {
    background-color: var(--bg-hover);
    border-color: var(--accent-primary);
  }
  
  .action-button.test-all {
    background-color: var(--accent-primary);
    color: white;
    border: none;
  }
  
  .action-button.test-all:hover {
    background-color: var(--accent-secondary);
  }
  
  .action-button.test-all:disabled {
    background-color: var(--bg-hover);
    color: var(--text-secondary);
    cursor: not-allowed;
  }
  
  .action-button.danger {
    color: var(--text-warning);
    border-color: var(--text-warning);
  }
  
  .action-button.danger:hover {
    background-color: rgba(206, 145, 120, 0.1); /* Using text-warning with opacity */
    color: var(--text-warning);
  }
  
  .form-container {
    display: grid;
    grid-template-columns: 1fr 1fr auto;
    gap: 1.5rem;
    background-color: var(--bg-secondary);
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: var(--shadow-light);
  }
  
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .form-group label {
    color: var(--text-primary);
    font-size: 0.95rem;
    font-weight: 500;
  }
  
  .form-group input {
    padding: 0.75rem;
    background-color: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    color: var(--text-primary);
    border-radius: 6px;
    transition: border-color var(--transition-speed) ease;
  }
  
  .form-group input:focus {
    border-color: var(--accent-primary);
    outline: none;
    box-shadow: 0 0 0 2px rgba(106, 153, 85, 0.2); /* Using accent-primary with opacity */
  }
  
  .primary-button {
    background-color: var(--accent-primary);
    color: white;
    border: none;
    padding: 0.75rem 1.5rem;
    font-weight: 500;
    border-radius: 6px;
    cursor: pointer;
    align-self: center;
    transition: background-color var(--transition-speed) ease;
  }
  
  .primary-button:hover {
    background-color: var(--accent-secondary);
  }
  
  .primary-button:disabled {
    background-color: var(--bg-hover);
    color: var(--text-secondary);
    cursor: not-allowed;
  }
  
  /* Animation for status changes */
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  
  @media (max-width: 768px) {
    .backend-item {
      grid-template-columns: 1fr;
      grid-template-rows: auto auto auto;
      gap: 1rem;
    }
    
    .backend-status {
      grid-column: 1;
      grid-row: 2;
    }
    
    .backend-actions {
      grid-column: 1 / -1;
      grid-row: 3;
      justify-content: flex-end;
    }
    
    .form-container {
      grid-template-columns: 1fr;
    }
    
    .primary-button {
      width: 100%;
    }
  }
</style>