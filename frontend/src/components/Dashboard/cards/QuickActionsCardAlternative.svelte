<script>
  import { onMount } from 'svelte';
  import { apiFetch } from '../../../services/api-v7';
  import ToggleServer from '../../resuables/ToggleServer.svelte';
  import Loaders from '../../resuables/Loaders.svelte';

  let status = $state('loading');
  let errorMessage = $state('');
  let metadata = $state(null);

  onMount(async () => {
    try {
      const response = await apiFetch('/api/v3/runfile/meta', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          fields: ['name', 'version', 'image', 'logo']
        })
      });

      const result = await response.json();

      if (result.data.status === 'failed') {
        if (result.data.message === 'runfile not loaded') {
          status = 'no-runfile';
        } else {
          status = 'error';
          errorMessage = 'Error loading runfile data';
        }
      } else if (result.data.status === 'success') {
        status = 'success';
        metadata = result.data.values;
      }
    } catch (err) {
      status = 'error';
      errorMessage = 'Failed to fetch runfile data';
    }
  });
</script>

<div class="card quick-actions">
  {#if metadata?.image}
    <div class="background-image" style="background-image: url('{metadata.image}')"></div>
    <div class="overlay"></div>
  {/if}

  <div class="card-content">
    {#if status === 'loading'}
      <div class="state-content">
        <div class="spinner"></div>
        <p class="state-text">Loading...</p>
      </div>
    {:else if status === 'no-runfile'}
      <div class="state-content">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
        </svg>
        <h3 class="state-title">No Runfile Selected</h3>
        <p class="state-text">Select a runfile from the Gallery to get started</p>
      </div>
    {:else if status === 'error'}
      <div class="state-content">
        <svg class="icon error-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <line x1="12" y1="8" x2="12" y2="12"/>
          <line x1="12" y1="16" x2="12.01" y2="16"/>
        </svg>
        <p class="state-text">{errorMessage}</p>
      </div>
    {:else if status === 'success' && metadata}
      <div class="card-header">
        <div class="header-content">
          {#if metadata.logo}
            <img src={metadata.logo} alt="{metadata.name} logo" class="logo" />
          {/if}
          <div class="title-section">
            <h3>{metadata.name || 'Unnamed Runfile'}</h3>
            {#if metadata.version}
              <span class="version">v{metadata.version}</span>
            {/if}
          </div>
        </div>
        <div class="card-icon">⚡</div>
      </div>
      
      <ToggleServer />
      
      <div class="action-buttons">
        <Loaders />
      </div>
    {/if}
  </div>
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
    position: relative;
    overflow: hidden;
    min-height: 350px;
  }

  .card:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-medium);
  }

  .background-image {
    position: absolute;
    top: -20px;
    left: -20px;
    right: -20px;
    bottom: -20px;
    background-size: cover;
    background-position: center;
    opacity: 0.2;
    z-index: 0;
  }

  .overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, 
      var(--bg-secondary) 0%, 
      transparent 50%, 
      var(--bg-secondary) 100%);
    opacity: 0.9;
    z-index: 1;
  }

  .card-content {
    position: relative;
    z-index: 2;
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    border-bottom: 1px solid var(--border-color);
    padding-bottom: 0.75rem;
  }

  .header-content {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex: 1;
  }

  .logo {
    width: 48px;
    height: 48px;
    object-fit: cover;
    border-radius: 6px;
    border: 2px solid var(--border-color);
    box-shadow: var(--shadow-light);
  }

  .title-section {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-wrap: wrap;
  }

  .card h3 {
    margin: 0;
    color: var(--text-accent);
    font-size: 1.1rem;
    font-weight: 600;
  }

  .version {
    display: inline-block;
    padding: 0.15rem 0.5rem;
    background-color: var(--accent-primary);
    color: var(--bg-primary);
    border-radius: 10px;
    font-size: 0.75rem;
    font-weight: 600;
    box-shadow: var(--shadow-light);
  }

  .card-icon {
    font-size: 1.25rem;
    opacity: 0.8;
  }

  .action-buttons {
    display: grid;
    grid-template-columns: 1fr;
    gap: 0.5rem;
    margin-top: 1rem;
  }

  .state-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    gap: 1rem;
    color: var(--text-secondary);
  }

  .spinner {
    width: 48px;
    height: 48px;
    border: 4px solid var(--bg-tertiary);
    border-top-color: var(--accent-primary);
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .icon {
    width: 56px;
    height: 56px;
    color: var(--text-secondary);
  }

  .error-icon {
    color: var(--text-warning);
  }

  .state-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0;
  }

  .state-text {
    font-size: 0.95rem;
    color: var(--text-secondary);
    text-align: center;
    margin: 0;
    max-width: 280px;
  }
</style>