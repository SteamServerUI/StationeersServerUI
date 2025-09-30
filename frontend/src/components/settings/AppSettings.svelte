<!-- AppSettings.svelte -->
<script>
    import { onMount } from 'svelte';
    import { apiFetch } from '../../services/api';
    
    
  /**
   * @typedef {Object} Props
   * @property {any} activeSidebarTab - Props
   */

  /** @type {Props} */
  let { activeSidebarTab } = $props();
    
    // State management
    let settingsData = $state([]);
    let settingsGroups = $state([]);
    let activeSettingsGroup = $state('');
    let statusMessage = $state('');
    let isError = $state(false);
    let statusTimeout;
    
    // apiFetch settings data on component mount
    onMount(async () => {
      await fetchSettings();
    });
    
    async function fetchSettings() {
      try {
        const response = await apiFetch('/api/v2/settings');
        const { data, error } = await response.json();
        
        if (error) {
          showStatus(`Failed to load settings: ${error}`, true);
          return;
        }
        
        settingsData = data;
        // Extract unique groups
        settingsGroups = [...new Set(settingsData.map(s => s.group))];
        
        if (settingsGroups.length > 0) {
          activeSettingsGroup = settingsGroups[0];
        }
      } catch (e) {
        showStatus(`Error fetching settings: ${e.message}`, true);
      }
    }
    
    // Handle settings group selection
    function selectSettingsGroup(group) {
      activeSettingsGroup = group;
    }
    
    // Update a setting
    async function updateSetting(name, value) {
      try {
        const response = await apiFetch('/api/v2/settings/save', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ [name]: value })
        });
        
        const { status, message } = await response.json();
        if (status === 'error') {
          showStatus(`Failed to update ${name}: ${message}`, true);
          return;
        }
        
        showStatus(`Updated ${name} successfully`, false);
      } catch (e) {
        showStatus(`Error updating ${name}: ${e.message}`, true);
      }
    }
    
    // Handle various input types
    function handleInputChange(setting, event) {
      const target = event.target;
      let value;
      
      if (setting.type === 'bool') {
        value = target.checked;
      } else if (setting.type === 'int') {
        value = target.value ? parseInt(target.value) : null;
        if (setting.required && !value) {
          showStatus(`Value for ${setting.name} is required`, true);
          return;
        }
      } else if (setting.type === 'array') {
        value = target.value ? target.value.split(',').map(s => s.trim()) : [];
      } else if (setting.type === 'map') {
        try {
          value = JSON.parse(target.value);
        } catch (e) {
          showStatus(`Invalid JSON for ${setting.name}: ${e.message}`, true);
          return;
        }
      } else {
        value = target.value;
      }
      
      updateSetting(setting.name, value);
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
  
  {#if activeSidebarTab === 'General' || activeSidebarTab === 'SSUI Settings'}
  <div class="settings-container">
    <h2>SSUI Settings</h2>
    
    <p class="settings-intro">
      Configure general application settings. Changes will be applied immediately.
    </p>
    
    {#if settingsGroups.length > 0}
      <div class="settings-group-nav">
        {#each settingsGroups as group}
          <button 
            class="section-nav-button {activeSettingsGroup === group ? 'active' : ''}" 
            onclick={() => selectSettingsGroup(group)}>
            {group}
          </button>
        {/each}
      </div>
      
      {#each settingsGroups as group}
        {#if activeSettingsGroup === group}
          <div class="settings-group">
            <h3>{group}</h3>
            <div class="settings-grid">
              {#each settingsData.filter(s => s.group === group) as setting}
                <div class="setting-card">
                  <div class="setting-info">
                    <h4>{setting.name}</h4>
                    {#if setting.description}
                      <p class="description">{setting.description}</p>
                    {/if}
                  </div>
                  
                  <div class="setting-control">
                    {#if setting.type === 'bool'}
                      <label class="switch">
                        <input 
                          type="checkbox" 
                          checked={setting.value === true} 
                          onchange={(e) => handleInputChange(setting, e)} 
                        />
                        <span class="slider"></span>
                      </label>
                    {:else if setting.type === 'int'}
                      <input 
                        type="number" 
                        value={setting.value ?? ''} 
                        min={setting.min} 
                        max={setting.max} 
                        required={setting.required} 
                        onchange={(e) => handleInputChange(setting, e)} 
                        class="number-input"
                      />
                    {:else if setting.type === 'array'}
                      <textarea 
                        value={setting.value?.join(', ') || ''} 
                        onchange={(e) => handleInputChange(setting, e)} 
                        class="text-area"
                        placeholder="Enter comma-separated values"
                      ></textarea>
                    {:else if setting.type === 'map'}
                      <textarea 
                        value={JSON.stringify(setting.value, null, 2) || '{}'} 
                        onchange={(e) => handleInputChange(setting, e)} 
                        class="code-area"
                        placeholder="Enter valid JSON"
                      ></textarea>
                    {:else}
                      <input 
                        type="text" 
                        value={setting.value || ''} 
                        required={setting.required} 
                        onchange={(e) => handleInputChange(setting, e)} 
                        class="text-input"
                      />
                    {/if}
                  </div>
                </div>
              {/each}
            </div>
          </div>
        {/if}
      {/each}
      
      {#if statusMessage}
        <div class="status-message" class:error={isError}>
          <span class="status-icon">{isError ? '⚠️' : '✓'}</span>
          <span>{statusMessage}</span>
          <button class="close-status" onclick={() => statusMessage = ''}>×</button>
        </div>
      {/if}
    {:else}
      <div class="loading-container">
        <div class="loading-spinner"></div>
        <p>Loading settings...</p>
      </div>
    {/if}
  </div>
{/if}
  
  <style>
    h2 {
      margin-top: 0;
      margin-bottom: 1.5rem;
      font-size: 1.5rem;
      font-weight: 500;
      color: var(--text-primary);
    }
    
    .settings-container {
      background-color: var(--bg-tertiary);
      border-radius: 8px;
      padding: 1.5rem;
      box-shadow: var(--shadow-light);
    }
    
    .settings-intro {
      margin-bottom: 2rem;
      color: var(--text-secondary);
      line-height: 1.5;
    }
    
    .settings-group {
      margin-bottom: 2rem;
      animation: fadeIn 0.3s ease;
    }
    
    .settings-group h3 {
      font-size: 1.1rem;
      font-weight: 500;
      margin-bottom: 1.5rem;
      color: var(--text-accent);
      border-bottom: 1px solid var(--border-color);
      padding-bottom: 0.5rem;
    }
    
    .settings-group-nav {
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;
      margin-bottom: 1.5rem;
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
    
    .settings-grid {
      display: grid;
      gap: 1.25rem;
      grid-template-columns: 1fr;
    }
    
    @media (min-width: 1024px) {
      .settings-grid {
        grid-template-columns: repeat(2, 1fr);
      }
    }
    
    .setting-card {
      background-color: var(--bg-secondary);
      border: 1px solid var(--border-color);
      border-radius: 8px;
      padding: 1.25rem;
      transition: all var(--transition-speed) ease;
    }
    
    .setting-card:hover {
      border-color: var(--accent-primary);
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
    }
    
    .setting-info {
      margin-bottom: 1rem;
    }
    
    .setting-info h4 {
      margin: 0 0 0.5rem 0;
      font-size: 1rem;
      font-weight: 600;
      color: var(--text-primary);
    }
    
    .description {
      margin: 0;
      font-size: 0.875rem;
      color: var(--text-secondary);
      line-height: 1.4;
    }
    
    .setting-control {
      width: 100%;
    }
    
    /* Switch Toggle for Boolean */
    .switch {
      position: relative;
      display: inline-block;
      width: 50px;
      height: 24px;
    }
    
    .switch input {
      opacity: 0;
      width: 0;
      height: 0;
    }
    
    .slider {
      position: absolute;
      cursor: pointer;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background-color: #ccc;
      transition: 0.3s;
      border-radius: 24px;
    }
    
    .slider:before {
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
    
    input:checked + .slider {
      background-color: var(--accent-primary);
    }
    
    input:checked + .slider:before {
      transform: translateX(26px);
    }
    
    /* Input Styles */
    .text-input,
    .number-input {
      width: 100%;
      padding: 0.75rem;
      border: 1px solid var(--border-color);
      border-radius: 6px;
      background-color: var(--bg-tertiary);
      color: var(--text-primary);
      font-size: 0.9rem;
      transition: all var(--transition-speed) ease;
    }
    
    .text-area,
    .code-area {
      width: 100%;
      min-height: 80px;
      padding: 0.75rem;
      border: 1px solid var(--border-color);
      border-radius: 6px;
      background-color: var(--bg-tertiary);
      color: var(--text-primary);
      font-size: 0.9rem;
      line-height: 1.4;
      resize: vertical;
      transition: all var(--transition-speed) ease;
    }
    
    .code-area {
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
      font-size: 0.85rem;
      min-height: 100px;
    }
    
    .text-input:focus,
    .number-input:focus,
    .text-area:focus,
    .code-area:focus {
      border-color: var(--accent-primary);
      outline: none;
      box-shadow: 0 0 0 3px rgba(106, 153, 85, 0.1);
    }
    
    .status-message {
      margin-top: 1.5rem;
      padding: 1rem;
      display: flex;
      align-items: center;
      border-radius: 6px;
      background-color: rgba(106, 153, 85, 0.1);
      color: var(--accent-primary);
      animation: slideIn 0.3s ease;
      position: relative;
    }
    
    .status-message.error {
      background-color: rgba(206, 145, 120, 0.1);
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
      border: 3px solid rgba(106, 153, 85, 0.3);
      border-radius: 50%;
      border-top-color: var(--accent-primary);
      animation: spin 1s ease-in-out infinite;
      margin-bottom: 1rem;
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
      .settings-container {
        padding: 1rem;
      }
      
      .setting-card {
        padding: 1rem;
      }
    }
  </style>