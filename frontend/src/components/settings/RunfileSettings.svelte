<!-- RunfileSettings.svelte -->
<script>
    import { onMount } from 'svelte';
    import { apiFetch } from '../../services/api';
    
    // State management
    let runfileGroups = $state([]);
    let activeRunfileGroup = $state('');
    let statusMessage = $state('');
    let isError = $state(false);
    let statusTimeout;
    let runfileArgs = $state([]);
    let isLoading = $state(true);
    let isSaving = $state(false);
    
    // apiFetch runfile groups on mount
    onMount(async () => {
      await fetchRunfileGroups();
    });
    
    async function fetchRunfileGroups() {
      isLoading = true;
      try {
        const response = await apiFetch('/api/v2/runfile/groups');
        const { data, error } = await response.json();
        
        if (error) {
          showStatus(`Failed to load groups: ${error}`, true);
          return;
        }
        
        runfileGroups = data;
        
        // Select the first group by default if available
        if (runfileGroups.length > 0) {
          await selectRunfileGroup(runfileGroups[0]);
        }
      } catch (e) {
        showStatus(`Error fetching runfile groups: ${e.message}`, true);
      } finally {
        isLoading = false;
      }
    }
    
    // Handle runfile group selection
    async function selectRunfileGroup(group) {
      isLoading = true;
      activeRunfileGroup = group;
      await fetchRunfileArgs(group);
      isLoading = false;
    }
    
    // apiFetch args for selected group
    async function fetchRunfileArgs(group) {
      try {
        const response = await apiFetch(`/api/v2/runfile/args?group=${encodeURIComponent(group)}`);
        const { data, error } = await response.json();
        
        if (error) {
          showStatus(`Failed to load args: ${error}`, true);
          return;
        }
        
        runfileArgs = data;
      } catch (e) {
        showStatus(`Error fetching runfile args: ${e.message}`, true);
      }
    }
    
    // Update a runfile arg
    async function updateRunfileArg(flag, value) {
      try {
        const response = await apiFetch('/api/v2/runfile/args/update', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ flag, value })
        });
        
        const { error } = await response.json();
        if (error) {
          showStatus(`Failed to update ${flag}: ${error}`, true);
          return;
        }
        
        showStatus(`Updated ${flag} successfully`, false);
      } catch (e) {
        showStatus(`Error updating ${flag}: ${e.message}`, true);
      }
    }
    
    // Handle input changes
    function handleInputChange(arg, event) {
      const target = event.target;
      let value;
      
      if (arg.type === 'bool') {
        value = target.checked.toString();
      } else {
        value = target.value;
      }
      
      updateRunfileArg(arg.flag, value);
    }
    
    // Save the entire runfile
    async function saveRunfile() {
      isSaving = true;
      try {
        const response = await apiFetch('/api/v2/runfile/save', {
          method: 'POST'
        });
        
        const { error } = await response.json();
        if (error) {
          showStatus(`Failed to save runfile: ${error}`, true);
          return;
        }
        
        showStatus('Runfile saved successfully', false);
      } catch (e) {
        showStatus(`Error saving runfile: ${e.message}`, true);
      } finally {
        isSaving = false;
      }
    }
    
    // Show status message
    function showStatus(message, error) {
      statusMessage = message;
      isError = error;
      
      // Clear any existing timeout
      if (statusTimeout) clearTimeout(statusTimeout);
      
      // Auto-hide after 30 seconds
      statusTimeout = setTimeout(() => {
        statusMessage = '';
      }, 30000);
    }
  </script>
  
  <div class="runfile-settings">
    <h2>Game Settings</h2>
    
    <p class="settings-intro">
      Configure command line arguments that will be used when launching the game server.
    </p>
    
    {#if isLoading && runfileGroups.length === 0}
      <div class="loading-container">
        <div class="loading-spinner"></div>
        <p>Loading runfile configuration...</p>
      </div>
    {:else if runfileGroups.length === 0}
      <div class="empty-state">
        <h3>No Runfile Groups Available</h3>
        <p>Select a runfile from the Runfile Gallery to get started!</p>
      </div>
    {:else}
      <div class="runfile-container">
        <div class="runfile-group-nav">
          {#each runfileGroups as group}
            <button 
              class="section-nav-button {activeRunfileGroup === group ? 'active' : ''}" 
              onclick={() => selectRunfileGroup(group)}>
              {group}
            </button>
          {/each}
        </div>
        
        {#if isLoading && activeRunfileGroup}
          <div class="loading-container">
            <div class="loading-spinner"></div>
            <p>Loading {activeRunfileGroup} settings...</p>
          </div>
        {:else if activeRunfileGroup}
          <div class="runfile-section">
            <h3>{activeRunfileGroup} GameServer Settings</h3>
            
            {#if runfileArgs.length === 0}
              <div class="empty-state">
                <p>No configurable arguments found for this group.</p>
              </div>
            {:else}
              <div class="channel-grid">
                {#each runfileArgs as arg}
                  <div class="setting-item">
                    <label>
                      <span>{arg.ui_label || arg.flag}</span>
                      
                      {#if arg.type === 'bool'}
                        <input 
                          type="checkbox" 
                          id={arg.flag} 
                          checked={arg.runtime_value === 'true'} 
                          disabled={arg.disabled} 
                          onchange={(e) => handleInputChange(arg, e)} 
                        />
                      {:else if arg.type === 'int'}
                        <input 
                          type="number" 
                          id={arg.flag} 
                          value={arg.runtime_value || ''} 
                          min={arg.min} 
                          max={arg.max} 
                          required={arg.required} 
                          disabled={arg.disabled} 
                          onchange={(e) => handleInputChange(arg, e)} 
                        />
                      {:else}
                        <input 
                          type="text" 
                          id={arg.flag} 
                          value={arg.runtime_value || ''} 
                          required={arg.required} 
                          disabled={arg.disabled} 
                          onchange={(e) => handleInputChange(arg, e)} 
                        />
                      {/if}
                    </label>
                    <div class="input-info">{arg.description || 'No description available'}</div>
                  </div>
                {/each}
              </div>
              
              <div class="form-actions">
                <button class="save-button" onclick={saveRunfile} disabled={isSaving}>
                  {isSaving ? 'Saving...' : 'Save All Changes'}
                </button>
              </div>
            {/if}
          </div>
        {/if}
        
        {#if statusMessage}
          <div class="status-message" class:error={isError}>
            <span class="status-icon">{isError ? '⚠️' : '✓'}</span>
            <span>{statusMessage}</span>
            <button class="close-status" onclick={() => statusMessage = ''}>×</button>
          </div>
        {/if}
      </div>
    {/if}
  </div>
  
  <style>
    .runfile-settings {
      padding-bottom: 2rem;
    }
    
    h2 {
      margin-top: 0;
      margin-bottom: 1rem;
      font-size: 1.5rem;
      font-weight: 500;
      color: var(--text-primary);
    }
    
    .settings-intro {
      margin-bottom: 2rem;
      color: var(--text-secondary);
      line-height: 1.5;
    }
    
    .runfile-container {
      background-color: var(--bg-tertiary);
      border-radius: 8px;
      padding: 1.5rem;
      box-shadow: var(--shadow-light);
    }
    
    .runfile-group-nav {
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;
      margin-bottom: 2rem;
    }
    
    .section-nav-button {
      padding: 0.5rem 1rem;
      background-color: var(--bg-secondary);
      border: 1px solid var(--border-color);
      border-radius: 4px;
      cursor: pointer;
      transition: all var(--transition-speed) ease;
      font-weight: 500;
      color: var(--text-primary);
    }
    
    .section-nav-button:hover {
      background-color: var(--bg-hover);
    }
    
    .section-nav-button.active {
      background-color: var(--accent-primary);
      color: white;
      border-color: var(--accent-primary);
    }
    
    .runfile-section {
      animation: fadeIn 0.3s ease;
    }
    
    .runfile-section h3 {
      font-size: 1.1rem;
      font-weight: 500;
      margin-bottom: 1.5rem;
      color: var(--text-accent);
      border-bottom: 1px solid var(--border-color);
      padding-bottom: 0.5rem;
    }
    
    .channel-grid {
      display: grid;
      grid-template-columns: 1fr;
      gap: 1.5rem;
    }
    
    @media (min-width: 768px) {
      .channel-grid {
        grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
      }
    }
    
    .setting-item {
      margin-bottom: 0.5rem;
    }
    
    .setting-item label {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 0.5rem;
    }
    
    .setting-item input[type="checkbox"] {
      width: 18px;
      height: 18px;
      accent-color: var(--accent-primary);
    }
    
    .setting-item input[type="text"],
    .setting-item input[type="number"] {
      background-color: var(--bg-tertiary);
      color: var(--text-primary);
      border: 1px solid var(--border-color);
      padding: 0.5rem;
      border-radius: 4px;
      width: 200px;
      transition: border-color var(--transition-speed) ease;
    }
    
    .setting-item input:focus {
      border-color: var(--accent-primary);
      outline: none;
      box-shadow: 0 0 0 2px rgba(106, 153, 85, 0.2); /* Using accent-primary with opacity */
    }
    
    .input-info {
      font-size: 0.85rem;
      color: var(--text-secondary);
      margin-top: 0.25rem;
    }
    
    .form-actions {
      margin-top: 2rem;
      display: flex;
      justify-content: flex-end;
    }
    
    .save-button {
      background-color: var(--accent-primary);
      color: white;
      border: none;
      border-radius: 4px;
      padding: 0.75rem 1.5rem;
      font-weight: 500;
      cursor: pointer;
      transition: background-color var(--transition-speed) ease;
    }
    
    .save-button:hover {
      background-color: var(--accent-secondary);
    }
    
    .save-button:disabled {
      background-color: var(--bg-hover);
      color: var(--text-secondary);
      cursor: not-allowed;
    }
    
    .loading-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 3rem 0;
      color: var(--text-secondary);
    }
    
    .loading-spinner {
      width: 40px;
      height: 40px;
      border: 3px solid rgba(106, 153, 85, 0.3); /* Using accent-primary with opacity */
      border-radius: 50%;
      border-top-color: var(--accent-primary);
      animation: spin 1s ease-in-out infinite;
      margin-bottom: 1rem;
    }
    
    .empty-state {
      text-align: center;
      padding: 2rem;
      color: var(--text-secondary);
      background-color: var(--bg-secondary);
      border-radius: 8px;
      border: 1px dashed var(--border-color);
    }
    
    .empty-state h3 {
      margin-top: 0;
      color: var(--text-primary);
    }
    
    .status-message {
      margin-top: 1.5rem;
      padding: 1rem;
      display: flex;
      align-items: center;
      border-radius: 6px;
      background-color: rgba(106, 153, 85, 0.1); /* Using accent-primary with opacity */
      color: var(--accent-primary);
      animation: slideIn 0.3s ease;
      position: relative;
    }
    
    .status-message.error {
      background-color: rgba(206, 145, 120, 0.1); /* Using text-warning with opacity */
      color: var(--text-warning);
    }
    
    .status-icon {
      margin-right: 0.75rem;
      font-size: 1.2rem;
    }
    
    .close-status {
      position: absolute;
      right: 1rem;
      background: none;
      border: none;
      cursor: pointer;
      font-size: 1.2rem;
      color: currentColor;
      opacity: 0.6;
    }
    
    .close-status:hover {
      opacity: 1;
    }
    
    @keyframes fadeIn {
      from { opacity: 0; }
      to { opacity: 1; }
    }
    
    @keyframes slideIn {
      from { 
        transform: translateY(-10px);
        opacity: 0;
      }
      to { 
        transform: translateY(0);
        opacity: 1;
      }
    }
    
    @keyframes spin {
      to { transform: rotate(360deg); }
    }
    
    @media (max-width: 768px) {
      .setting-item label {
        flex-direction: column;
        align-items: flex-start;
        gap: 0.5rem;
      }
      
      .setting-item input[type="text"],
      .setting-item input[type="number"] {
        width: 100%;
      }
    }
  </style>