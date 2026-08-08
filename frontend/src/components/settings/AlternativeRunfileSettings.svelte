<!-- RunfileSettings.svelte -->
<script>
    import { onMount } from 'svelte';
    import { apiFetch } from '../../services/api-v7';
    
    // State management
    let runfileGroups = $state([]);
    let activeRunfileGroup = $state('');
    let statusMessage = $state('');
    let isError = $state(false);
    let statusTimeout;
    let runfileArgs = $state([]);
    let isLoading = $state(true);
    let isSaving = $state(false);
    let expandedSections = $state(new Set());
    
    // Fetch runfile groups on mount
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
        
        // Sort groups alphabetically for consistent ordering
        runfileGroups = data.sort((a, b) => a.localeCompare(b));
        
        // Select the first group by default if available
        if (runfileGroups.length > 0) {
          await selectRunfileGroup(runfileGroups[0]);
          expandedSections.add(runfileGroups[0]);
        }
      } catch (e) {
        showStatus(`Error fetching runfile groups: ${e.message}`, true);
      } finally {
        isLoading = false;
      }
    }
    
    // Handle runfile group selection
    async function selectRunfileGroup(group) {
      // Toggle expansion
      if (expandedSections.has(group)) {
        expandedSections.delete(group);
        activeRunfileGroup = '';
        expandedSections = new Set(expandedSections);
      } else {
        activeRunfileGroup = group;
        isLoading = true;
        await fetchRunfileArgs(group);
        expandedSections.add(group);
        expandedSections = new Set(expandedSections);
        isLoading = false;
      }
    }
    
    // Fetch args for selected group
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
        
        showStatus(`Updated ${flag}`, false);
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
      
      // Auto-hide after 5 seconds
      statusTimeout = setTimeout(() => {
        statusMessage = '';
      }, 5000);
    }
  </script>
  
  <div class="runfile-settings">
    <div class="settings-header">
      <h2>Game Server Configuration</h2>
      <p class="settings-intro">
        Configure command line arguments for your game server instances
      </p>
    </div>
    
    {#if isLoading && runfileGroups.length === 0}
      <div class="loading-container">
        <div class="loading-spinner"></div>
        <p>Loading configuration...</p>
      </div>
    {:else if runfileGroups.length === 0}
      <div class="empty-state">
        <div class="empty-icon">⚙️</div>
        <h3>No Configuration Available</h3>
        <p>Select a runfile from the Runfile Gallery to begin configuring your server.</p>
      </div>
    {:else}
      <div class="runfile-sections">
        {#each runfileGroups as group}
          <div class="section-card">
            <button 
              class="section-header {expandedSections.has(group) ? 'expanded' : ''}" 
              onclick={() => selectRunfileGroup(group)}
            >
              <span class="section-title">{group}</span>
              <span class="expand-icon">{expandedSections.has(group) ? '−' : '+'}</span>
            </button>
            
            {#if expandedSections.has(group) && activeRunfileGroup === group}
              <div class="section-content">
                {#if isLoading && activeRunfileGroup === group}
                  <div class="loading-container small">
                    <div class="loading-spinner small"></div>
                    <p>Loading settings...</p>
                  </div>
                {:else if activeRunfileGroup === group}
                  {#if runfileArgs.length === 0}
                    <div class="empty-args">
                      <p>No configurable arguments available for this group.</p>
                    </div>
                  {:else}
                    <div class="settings-grid">
                      {#each runfileArgs as arg}
                        <div class="setting-field">
                          <div class="field-header">
                            <label for={arg.flag} class="field-label">
                              {arg.ui_label || arg.flag}
                            </label>
                            {#if arg.required}
                              <span class="required-badge">Required</span>
                            {/if}
                          </div>
                          
                          {#if arg.description}
                            <div class="field-description">{arg.description}</div>
                          {/if}
                          
                          <div class="field-input">
                            {#if arg.type === 'bool'}
                              <label class="toggle-switch">
                                <input 
                                  type="checkbox" 
                                  id={arg.flag} 
                                  checked={arg.runtime_value === 'true'} 
                                  disabled={arg.disabled} 
                                  onchange={(e) => handleInputChange(arg, e)} 
                                />
                                <span class="toggle-slider"></span>
                              </label>
                            {:else if arg.type === 'int'}
                              <input 
                                type="number" 
                                id={arg.flag} 
                                class="input-field"
                                value={arg.runtime_value || ''} 
                                min={arg.min} 
                                max={arg.max} 
                                placeholder={arg.min !== undefined ? `Min: ${arg.min}` : ''}
                                required={arg.required} 
                                disabled={arg.disabled} 
                                onchange={(e) => handleInputChange(arg, e)} 
                              />
                            {:else}
                              <input 
                                type="text" 
                                id={arg.flag} 
                                class="input-field"
                                value={arg.runtime_value || ''} 
                                placeholder="Enter value..."
                                required={arg.required} 
                                disabled={arg.disabled} 
                                onchange={(e) => handleInputChange(arg, e)} 
                              />
                            {/if}
                          </div>
                        </div>
                      {/each}
                    </div>
                  {/if}
                {/if}
              </div>
            {/if}
          </div>
        {/each}
      </div>
      
      <div class="actions-bar">
        <button class="save-button" onclick={saveRunfile} disabled={isSaving}>
          {isSaving ? 'Saving...' : 'Save All Changes'}
        </button>
      </div>
    {/if}
    
    {#if statusMessage}
      <div class="toast {isError ? 'error' : 'success'}">
        <span class="toast-icon">{isError ? '✕' : '✓'}</span>
        <span class="toast-message">{statusMessage}</span>
        <button class="toast-close" onclick={() => statusMessage = ''}>×</button>
      </div>
    {/if}
  </div>
  
  <style>
    .runfile-settings {
      max-width: 1200px;
      margin: 0 auto;
      padding: 0 1rem 2rem;
    }
    
    .settings-header {
      margin-bottom: 2rem;
    }
    
    h2 {
      margin: 0 0 0.5rem 0;
      font-size: 1.75rem;
      font-weight: 600;
      color: var(--text-primary);
    }
    
    .settings-intro {
      margin: 0;
      color: var(--text-secondary);
      font-size: 0.95rem;
      line-height: 1.5;
    }
    
    .runfile-sections {
      display: flex;
      flex-direction: column;
      gap: 1rem;
      margin-bottom: 1.5rem;
    }
    
    .section-card {
      background-color: var(--bg-secondary);
      border-radius: 8px;
      border: 1px solid var(--border-color);
      overflow: hidden;
      transition: box-shadow 0.2s ease;
    }
    
    .section-card:hover {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    }
    
    .section-header {
      width: 100%;
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1rem 1.25rem;
      background: none;
      border: none;
      cursor: pointer;
      transition: background-color 0.2s ease;
      text-align: left;
    }
    
    .section-header:hover {
      background-color: var(--bg-hover);
    }
    
    .section-header.expanded {
      border-bottom: 1px solid var(--border-color);
    }
    
    .section-title {
      font-size: 1.1rem;
      font-weight: 500;
      color: var(--text-primary);
    }
    
    .expand-icon {
      font-size: 1.5rem;
      color: var(--text-secondary);
      font-weight: 300;
      transition: transform 0.2s ease;
    }
    
    .section-header.expanded .expand-icon {
      transform: rotate(0deg);
    }
    
    .section-content {
      padding: 1.5rem;
      animation: slideDown 0.2s ease;
    }
    
    .settings-grid {
      display: grid;
      gap: 1.5rem;
    }
    
    @media (min-width: 768px) {
      .settings-grid {
        grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
      }
    }
    
    .setting-field {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }
    
    .field-header {
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }
    
    .field-label {
      font-size: 0.9rem;
      font-weight: 500;
      color: var(--text-primary);
    }
    
    .required-badge {
      font-size: 0.7rem;
      padding: 0.15rem 0.4rem;
      background-color: var(--accent-primary);
      color: white;
      border-radius: 3px;
      text-transform: uppercase;
      font-weight: 600;
    }
    
    .field-description {
      font-size: 0.85rem;
      color: var(--text-secondary);
      line-height: 1.4;
    }
    
    .field-input {
      margin-top: 0.25rem;
    }
    
    .input-field {
      width: 100%;
      padding: 0.625rem 0.75rem;
      background-color: var(--bg-tertiary);
      color: var(--text-primary);
      border: 1px solid var(--border-color);
      border-radius: 6px;
      font-size: 0.95rem;
      transition: all 0.2s ease;
    }
    
    .input-field:focus {
      border-color: var(--accent-primary);
      outline: none;
      box-shadow: 0 0 0 3px rgba(106, 153, 85, 0.1);
    }
    
    .input-field:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
    
    /* Toggle Switch Styling */
    .toggle-switch {
      position: relative;
      display: inline-block;
      width: 48px;
      height: 26px;
    }
    
    .toggle-switch input {
      opacity: 0;
      width: 0;
      height: 0;
    }
    
    .toggle-slider {
      position: absolute;
      cursor: pointer;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background-color: var(--bg-hover);
      border: 1px solid var(--border-color);
      transition: 0.3s;
      border-radius: 26px;
    }
    
    .toggle-slider:before {
      position: absolute;
      content: "";
      height: 18px;
      width: 18px;
      left: 3px;
      bottom: 3px;
      background-color: white;
      transition: 0.3s;
      border-radius: 50%;
    }
    
    input:checked + .toggle-slider {
      background-color: var(--accent-primary);
      border-color: var(--accent-primary);
    }
    
    input:checked + .toggle-slider:before {
      transform: translateX(22px);
    }
    
    input:disabled + .toggle-slider {
      opacity: 0.5;
      cursor: not-allowed;
    }
    
    .actions-bar {
      display: flex;
      justify-content: flex-end;
      padding: 1.5rem 0;
      border-top: 1px solid var(--border-color);
    }
    
    .save-button {
      background-color: var(--accent-primary);
      color: white;
      border: none;
      border-radius: 6px;
      padding: 0.75rem 2rem;
      font-weight: 500;
      font-size: 0.95rem;
      cursor: pointer;
      transition: all 0.2s ease;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }
    
    .save-button:hover:not(:disabled) {
      background-color: var(--accent-secondary);
      transform: translateY(-1px);
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
    }
    
    .save-button:disabled {
      background-color: var(--bg-hover);
      color: var(--text-secondary);
      cursor: not-allowed;
      transform: none;
      box-shadow: none;
    }
    
    .loading-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 3rem 1rem;
      color: var(--text-secondary);
    }
    
    .loading-container.small {
      padding: 2rem 1rem;
    }
    
    .loading-spinner {
      width: 40px;
      height: 40px;
      border: 3px solid var(--border-color);
      border-radius: 50%;
      border-top-color: var(--accent-primary);
      animation: spin 0.8s linear infinite;
      margin-bottom: 1rem;
    }
    
    .loading-spinner.small {
      width: 28px;
      height: 28px;
      border-width: 2px;
    }
    
    .empty-state {
      text-align: center;
      padding: 4rem 2rem;
      color: var(--text-secondary);
      background-color: var(--bg-secondary);
      border-radius: 8px;
      border: 2px dashed var(--border-color);
    }
    
    .empty-icon {
      font-size: 3rem;
      margin-bottom: 1rem;
      opacity: 0.5;
    }
    
    .empty-state h3 {
      margin: 0 0 0.5rem 0;
      color: var(--text-primary);
      font-weight: 500;
    }
    
    .empty-state p {
      margin: 0;
    }
    
    .empty-args {
      text-align: center;
      padding: 2rem;
      color: var(--text-secondary);
      font-size: 0.9rem;
    }
    
    /* Toast Notification */
    .toast {
      position: fixed;
      bottom: 2rem;
      right: 2rem;
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 1rem 1.25rem;
      border-radius: 8px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      animation: slideUp 0.3s ease;
      z-index: 1000;
      max-width: 400px;
    }
    
    .toast.success {
      background-color: var(--accent-primary);
      color: white;
    }
    
    .toast.error {
      background-color: #d32f2f;
      color: white;
    }
    
    .toast-icon {
      font-size: 1.25rem;
      font-weight: bold;
    }
    
    .toast-message {
      flex: 1;
      font-size: 0.9rem;
    }
    
    .toast-close {
      background: none;
      border: none;
      color: white;
      font-size: 1.5rem;
      cursor: pointer;
      opacity: 0.8;
      padding: 0;
      width: 24px;
      height: 24px;
      display: flex;
      align-items: center;
      justify-content: center;
      transition: opacity 0.2s;
    }
    
    .toast-close:hover {
      opacity: 1;
    }
    
    @keyframes slideDown {
      from {
        opacity: 0;
        transform: translateY(-10px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    }
    
    @keyframes slideUp {
      from {
        opacity: 0;
        transform: translateY(20px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    }
    
    @keyframes spin {
      to { transform: rotate(360deg); }
    }
    
    @media (max-width: 768px) {
      h2 {
        font-size: 1.5rem;
      }
      
      .section-header {
        padding: 0.875rem 1rem;
      }
      
      .section-content {
        padding: 1rem;
      }
      
      .toast {
        bottom: 1rem;
        right: 1rem;
        left: 1rem;
        max-width: none;
      }
    }
  </style>