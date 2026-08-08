<script lang="js">
  import { apiFetch } from '../../../services/api-v7';
  import SteamCMDWait from './SteamCMDWait.svelte';

  const { runfile, onClose } = $props();

  let isLoading = $state(false);
  let isSteamCMDRunning = $state(false);
  let applyingSettings = $state(false);
  let applyingPlugins = $state(false);
  let reloadingBackend = $state(false);
  let results = $state(null);

  // Filter settings and plugins by OS
  const filteredSettings = $derived(
    (runfile.recommended_settings || []).filter(s => 
      !s.supported_os || s.supported_os === 'all' || s.supported_os === runfile.current_os
    )
  );

  const filteredPlugins = $derived(
    (runfile.recommended_plugins || []).filter(p => 
      !p.supported_os || p.supported_os === 'all' || p.supported_os === runfile.current_os
    )
  );

  const hasRecommendations = $derived(filteredSettings.length > 0 || filteredPlugins.length > 0);

  // Handle ESC key
  function handleKeydown(event) {
    if (event.key === 'Escape') {
      onClose();
    }
  }

  // Handle backdrop click
  function handleBackdropClick(event) {
    if (event.target === event.currentTarget) {
      onClose();
    }
  }

  // Apply runfile only
  async function applyRunfile(withRecommendations = false) {
    isLoading = true;
    isSteamCMDRunning = true;
    results = null;

    const operationResults = {
      runfile: null,
      settings: [],
      plugins: [],
      backendReload: null
    };

    try {
      // Step 1: Apply the runfile
      const response = await apiFetch('/api/v2/gallery/select', {
        method: 'POST',
        body: JSON.stringify({ identifier: runfile.name, redownload: false })
      });

      if (response.status === 409) {
        const confirmRedownload = window.confirm(
          `Runfile ${runfile.name} already exists. Do you want to re-download it? This will OVERWRITE custom settings.`
        );
        if (confirmRedownload) {
          await apiFetch('/api/v2/gallery/select', {
            method: 'POST',
            body: JSON.stringify({ identifier: runfile.name, redownload: true })
          });
          operationResults.runfile = { success: true, message: 'Runfile re-downloaded successfully' };
        } else {
          throw new Error('Runfile download cancelled');
        }
      } else {
        operationResults.runfile = { success: true, message: 'Runfile applied successfully' };
      }

      isSteamCMDRunning = false;

      // Step 2: Apply recommended settings if requested
      if (withRecommendations && filteredSettings.length > 0) {
        applyingSettings = true;
        
        for (const setting of filteredSettings) {
          try {
            const settingBody = {};
            if (setting.boolValue !== undefined) {
              settingBody[setting.name] = setting.boolValue;
            } else if (setting.stringValue !== undefined) {
              settingBody[setting.name] = setting.stringValue;
            } else if (setting.intValue !== undefined) {
              settingBody[setting.name] = setting.intValue;
            }

            const settingResponse = await apiFetch('/api/v2/settings/save', {
              method: 'POST',
              body: JSON.stringify(settingBody)
            });

            let responseData;
            if (settingResponse instanceof Response) {
              const text = await settingResponse.text();
              try {
                responseData = JSON.parse(text);
              } catch {
                responseData = { message: text };
              }
            } else {
              responseData = settingResponse;
            }

            if (responseData.status === 'success' || responseData.message?.includes('successfully')) {
              operationResults.settings.push({
                name: setting.name,
                success: true,
                message: 'Applied successfully'
              });
            } else {
              operationResults.settings.push({
                name: setting.name,
                success: false,
                message: responseData.message || 'Failed to apply'
              });
            }
          } catch (err) {
            operationResults.settings.push({
              name: setting.name,
              success: false,
              message: err.message || 'Failed to apply'
            });
          }
        }
        
        applyingSettings = false;
      }

      // Step 3: Apply recommended plugins if requested
      if (withRecommendations && filteredPlugins.length > 0) {
        applyingPlugins = true;
        
        for (const plugin of filteredPlugins) {
          try {
            const pluginResponse = await apiFetch('/api/v2/plugingallery/select', {
              method: 'POST',
              body: JSON.stringify({ name: plugin.name, redownload: true })
            });

            let responseData;
            if (pluginResponse instanceof Response) {
              const text = await pluginResponse.text();
              try {
                responseData = JSON.parse(text);
              } catch {
                responseData = { message: text };
              }
            } else {
              responseData = pluginResponse;
            }

            if (responseData.status === 'success' || pluginResponse.ok) {
              operationResults.plugins.push({
                name: plugin.name,
                success: true,
                message: 'Installed successfully'
              });
            } else {
              operationResults.plugins.push({
                name: plugin.name,
                success: false,
                message: responseData.message || 'Failed to install'
              });
            }
          } catch (err) {
            operationResults.plugins.push({
              name: plugin.name,
              success: false,
              message: err.message || 'Failed to install'
            });
          }
        }
        
        applyingPlugins = false;
      }


      // Commented out since this is not needed and causes the plugin to start twice on windows.
      //// Step 4: Reload backend to apply all changes
      //reloadingBackend = true;
      //try {
      //  const reloadResponse = await apiFetch('/api/v2/loader/reloadbackend', {
      //    method: 'GET'
      //  });
      //
      //  let reloadData;
      //  if (reloadResponse instanceof Response) {
      //    reloadData = await reloadResponse.json();
      //  } else {
      //    reloadData = reloadResponse;
      //  }
      //
      //  if (reloadData.status === 'OK' || reloadResponse.ok) {
      //    operationResults.backendReload = { success: true, message: 'Backend reloaded successfully' };
      //  } else {
      //    operationResults.backendReload = { success: false, message: 'Failed to reload backend' };
      //  }
      //} catch (err) {
      //  operationResults.backendReload = { success: false, message: err.message || 'Failed to reload backend' };
      //} finally {
      //  reloadingBackend = false;
      //}

      results = operationResults;
    } catch (err) {
      results = {
        runfile: { success: false, message: err.message || 'Failed to apply runfile' },
        settings: [],
        plugins: []
      };
    } finally {
      isLoading = false;
      isSteamCMDRunning = false;
      applyingSettings = false;
      applyingPlugins = false;
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="modal-backdrop" onclick={handleBackdropClick}>
  <div class="modal-content">
    <!-- Header with background -->
    <div class="modal-header" style="background-image: url('{runfile.background_url || ''}')">
      <h2>{runfile.name}</h2>
    </div>

    <div class="modal-body">
      <!-- Metadata Section -->
      <div class="section">
        <h3>Information</h3>
        <div class="metadata-grid">
          <div class="metadata-item">
            <span class="label">Runfile Version:</span>
            <span class="value">{runfile.version || 'N/A'}</span>
          </div>
          <div class="metadata-item">
            <span class="label">Supported OS:</span>
            <span class="value">{runfile.supported_os || 'N/A'}</span>
          </div>
          {#if runfile.filename}
            <div class="metadata-item">
              <span class="label">Filename:</span>
              <span class="value">{runfile.filename}</span>
            </div>
          {/if}
        </div>
      </div>

      <!-- Recommended Settings Section -->
      {#if filteredSettings.length > 0}
        <div class="section">
          <h3>Recommended Settings</h3>
          <div class="recommendations-list">
            {#each filteredSettings as setting}
              <div class="recommendation-item">
                <span class="rec-name">{setting.name}</span>
                <span class="rec-value">
                  {#if setting.boolValue !== undefined}
                    {setting.boolValue ? '✓ Enable' : '✗ Disable'}
                  {:else if setting.stringValue !== undefined}
                    {setting.stringValue}
                  {:else if setting.intValue !== undefined}
                    {setting.intValue}
                  {/if}
                </span>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Recommended Plugins Section -->
      {#if filteredPlugins.length > 0}
        <div class="section">
          <h3>Recommended Plugins</h3>
          <div class="recommendations-list">
            {#each filteredPlugins as plugin}
              <div class="recommendation-item">
                <span class="rec-name">{plugin.name}</span>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Results Section -->
      {#if results}
        <div class="section results-section">
          <h3>Application Results</h3>
          
          {#if results.runfile}
            <div class="result-group">
              <h4>Runfile</h4>
              <div class="result-item" class:success={results.runfile.success} class:error={!results.runfile.success}>
                <span class="result-icon">{results.runfile.success ? '✓' : '✗'}</span>
                <span>{results.runfile.message}</span>
              </div>
            </div>
          {/if}

          {#if results.settings.length > 0}
            <div class="result-group">
              <h4>Settings ({results.settings.filter(s => s.success).length}/{results.settings.length} applied)</h4>
              {#each results.settings as result}
                <div class="result-item" class:success={result.success} class:error={!result.success}>
                  <span class="result-icon">{result.success ? '✓' : '✗'}</span>
                  <span class="result-name">{result.name}:</span>
                  <span>{result.message}</span>
                </div>
              {/each}
            </div>
          {/if}

          {#if results.plugins.length > 0}
            <div class="result-group">
              <h4>Plugins ({results.plugins.filter(p => p.success).length}/{results.plugins.length} installed)</h4>
              {#each results.plugins as result}
                <div class="result-item" class:success={result.success} class:error={!result.success}>
                  <span class="result-icon">{result.success ? '✓' : '✗'}</span>
                  <span class="result-name">{result.name}:</span>
                  <span>{result.message}</span>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      {/if}

      <!-- Progress Indicators -->
      {#if isLoading || applyingSettings || applyingPlugins}
        <div class="progress-section">
          {#if isSteamCMDRunning}
            <SteamCMDWait />
          {:else if applyingSettings}
            <div class="progress-item">
              <div class="spinner"></div>
              <span>Applying settings...</span>
            </div>
          {:else if applyingPlugins}
            <div class="progress-item">
              <div class="spinner"></div>
              <span>Installing plugins...</span>
            </div>
          {/if}
        </div>
      {/if}
    </div>

    <!-- Action Buttons -->
    <div class="modal-footer">
      <button 
        class="action-button secondary" 
        onclick={() => applyRunfile(false)}
        disabled={isLoading}
      >
        Download & Apply
      </button>
      
      {#if hasRecommendations}
        <button 
          class="action-button primary" 
          onclick={() => applyRunfile(true)}
          disabled={isLoading}
        >
          Download & Apply with Recommendations
        </button>
      {/if}
    </div>
  </div>
</div>

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.7);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
    padding: 2rem;
    animation: fadeIn 0.2s ease-out;
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  .modal-content {
    background: var(--bg-secondary, #ffffff);
    border-radius: 16px;
    max-width: 800px;
    width: 100%;
    max-height: 90vh;
    display: flex;
    flex-direction: column;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    animation: slideUp 0.3s ease-out;
    overflow: hidden;
  }

  @keyframes slideUp {
    from {
      transform: translateY(20px);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 1;
    }
  }

  .modal-header {
    position: relative;
    background-size: cover;
    background-position: center;
    padding: 3rem 2rem;
    color: white;
    text-align: center;
  }

  .modal-header::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, rgba(0, 0, 0, 0.6), rgba(0, 0, 0, 0.4));
    z-index: 0;
  }

  .modal-header h2 {
    position: relative;
    z-index: 1;
    margin: 0;
    font-size: 2rem;
    text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
  }

  .modal-body {
    padding: 2rem;
    overflow-y: auto;
    flex: 1;
  }

  .section {
    margin-bottom: 2rem;
  }

  .section:last-child {
    margin-bottom: 0;
  }

  .section h3 {
    color: var(--text-primary, #333);
    margin: 0 0 1rem 0;
    font-size: 1.3rem;
    border-bottom: 2px solid var(--border-color, #e0e0e0);
    padding-bottom: 0.5rem;
  }

  .metadata-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
  }

  .metadata-item {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .label {
    font-weight: 600;
    color: var(--text-secondary, #666);
    font-size: 0.85rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .value {
    color: var(--text-primary, #333);
    font-size: 1rem;
  }

  .recommendations-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .recommendation-item {
    background: var(--bg-tertiary, #f5f5f5);
    padding: 0.75rem 1rem;
    border-radius: 8px;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .rec-name {
    font-weight: 600;
    color: var(--text-primary, #333);
  }

  .rec-value {
    color: var(--text-secondary, #666);
    font-size: 0.9rem;
  }

  .results-section {
    background: var(--bg-tertiary, #f9f9f9);
    padding: 1.5rem;
    border-radius: 12px;
    border: 2px solid var(--border-color, #e0e0e0);
  }

  .result-group {
    margin-bottom: 1.5rem;
  }

  .result-group:last-child {
    margin-bottom: 0;
  }

  .result-group h4 {
    margin: 0 0 0.75rem 0;
    color: var(--text-primary, #333);
    font-size: 1rem;
  }

  .result-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem;
    border-radius: 6px;
    margin-bottom: 0.5rem;
    font-size: 0.9rem;
  }

  .result-item.success {
    background: rgba(76, 175, 80, 0.1);
    color: #2e7d32;
  }

  .result-item.error {
    background: rgba(244, 67, 54, 0.1);
    color: #c62828;
  }

  .result-icon {
    font-weight: bold;
    font-size: 1.1rem;
  }

  .result-name {
    font-weight: 600;
  }

  .progress-section {
    margin-top: 1.5rem;
    padding: 1.5rem;
    background: var(--bg-tertiary, #f5f5f5);
    border-radius: 12px;
    display: flex;
    justify-content: center;
  }

  .progress-item {
    display: flex;
    align-items: center;
    gap: 1rem;
    color: var(--text-secondary, #666);
  }

  .spinner {
    width: 24px;
    height: 24px;
    border: 3px solid var(--bg-secondary, #eeeeee);
    border-top: 3px solid var(--accent-primary, #2C6E49);
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  .modal-footer {
    padding: 1.5rem 2rem;
    background: var(--bg-tertiary, #f5f5f5);
    border-top: 1px solid var(--border-color, #e0e0e0);
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
  }

  .action-button {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 1rem;
  }

  .action-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .action-button.secondary {
    background: var(--bg-secondary, #ffffff);
    color: var(--text-primary, #333);
    border: 2px solid var(--border-color, #e0e0e0);
  }

  .action-button.secondary:hover:not(:disabled) {
    background: var(--bg-tertiary, #f5f5f5);
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  }

  .action-button.primary {
    background: linear-gradient(135deg, var(--accent-primary, #2C6E49), var(--accent-secondary, #4A90E2));
    color: white;
  }

  .action-button.primary:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
  }

  @media (max-width: 768px) {
    .modal-backdrop {
      padding: 1rem;
    }

    .modal-content {
      max-height: 95vh;
    }

    .modal-header {
      padding: 2rem 1.5rem;
    }

    .modal-header h2 {
      font-size: 1.5rem;
    }

    .modal-body {
      padding: 1.5rem;
    }

    .metadata-grid {
      grid-template-columns: 1fr;
    }

    .modal-footer {
      flex-direction: column;
      gap: 0.75rem;
    }

    .action-button {
      width: 100%;
    }
  }
</style>