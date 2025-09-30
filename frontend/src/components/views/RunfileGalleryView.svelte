<script lang="js">
  import { apiFetch } from '../../services/api';
  import RunfileCard from '../Gallery/RunfileCard.svelte';

  // State
  let runfiles = $state([]);
  let loading = $state(true);
  let error = $state(null);

  // Fetch runfiles on mount
  async function fetchRunfiles(forceUpdate = false) {
    loading = true;
    error = null;
    
    try {
      // Get the response from apiFetch
      const response = await apiFetch(`/api/v2/gallery${forceUpdate ? '?forceUpdate=true' : ''}`);
      
      let data;
      if (response instanceof Response) {
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        data = await response.json();
      } else {
        data = response;
      }
      
      // Handle different response formats
      if (data && data.data && Array.isArray(data.data)) {
        // Format: { data: [...runfiles] }
        runfiles = data.data;
      } else if (Array.isArray(data)) {
        // Format: [...runfiles]
        runfiles = data;
      } else {
        console.warn('Unexpected response format:', data);
        runfiles = [];
      }
      
      // Add missing identifiers if needed
      runfiles = runfiles.map(runfile => {
        if (!runfile.identifier) {
          runfile.identifier = runfile.filename || runfile.name;
        }
        return runfile;
      });
      
    } catch (err) {
      error = err.message || 'Failed to fetch runfiles';
      console.error('Error fetching runfiles:', err);
      runfiles = []; // Ensure runfiles is always an array
    } finally {
      loading = false;
    }
  }

  // Refresh gallery
  function refreshGallery() {
    fetchRunfiles(true);
  }

  // Fetch runfiles on mount
  $effect(() => {
    fetchRunfiles();
  });
</script>

<div class="gallery-page">
  <div class="header">
    <div class="header-content">
    </div>
    <button 
      class="refresh-button" 
      onclick={refreshGallery} 
      disabled={loading}
      aria-label="Refresh runfile gallery"
    >
      <svg class="refresh-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/>
        <path d="M21 3v5h-5"/>
        <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/>
        <path d="M3 21v-5h5"/>
      </svg>
      {loading ? 'Loading...' : 'Refresh'}
    </button>
  </div>
  
  {#if error}
    <div class="error-banner">
      <svg class="error-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"/>
        <line x1="15" y1="9" x2="9" y2="15"/>
        <line x1="9" y1="9" x2="15" y2="15"/>
      </svg>
      <p>{error}</p>
    </div>
  {/if}
  
  {#if loading}
    <div class="loading">
      <div class="spinner"></div>
      <p>Loading runfiles...</p>
    </div>
  {:else if !runfiles || runfiles.length === 0}
    <div class="empty-state">
      <svg class="empty-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
        <circle cx="9" cy="9" r="2"/>
        <path d="M21 15l-3.086-3.086a2 2 0 0 0-2.828 0L6 21"/>
      </svg>
      <h3>No runfiles available</h3>
      <p>Try refreshing the gallery to see available runfiles. Is this backend outdated?</p>
    </div>
  {:else}
    <div class="gallery-stats">
      <span class="stats-text">{runfiles.length} runfile{runfiles.length !== 1 ? 's' : ''} available</span>
    </div>
    <div class="gallery-grid">
      {#each runfiles as runfile (runfile.identifier)}
        <RunfileCard {runfile} />
      {/each}
    </div>
  {/if}
</div>

<style>
  .gallery-page {
    padding: 2rem;
    background: linear-gradient(135deg, var(--bg-primary, #f0f0f0) 0%, var(--bg-secondary, #ffffff) 100%);
    min-height: 100vh;
  }
  
  .header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 3rem;
    padding-bottom: 1rem;
    border-bottom: 2px solid var(--border-color, #e0e0e0);
  }

  .header-content {
    flex: 1;
  }

  .gallery-title {
    font-size: 2.5rem;
    font-weight: 700;
    color: var(--text-primary, #333);
    margin: 0 0 0.5rem 0;
    background: linear-gradient(45deg, var(--accent-primary, #2C6E49), var(--accent-secondary, #4A90E2));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .gallery-subtitle {
    font-size: 1.1rem;
    color: var(--text-secondary, #666);
    margin: 0;
    font-weight: 400;
  }

  .refresh-button {
    background: linear-gradient(135deg, var(--accent-primary, #2C6E49), var(--accent-secondary, #4A90E2));
    color: white;
    border: none;
    border-radius: 12px;
    padding: 0.75rem 1.5rem;
    cursor: pointer;
    font-weight: 600;
    font-size: 1rem;
    transition: all var(--transition-speed, 0.3s) ease;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    box-shadow: var(--shadow-light, 0 4px 8px rgba(0,0,0,0.1));
  }

  .refresh-icon {
    width: 18px;
    height: 18px;
    transition: transform var(--transition-speed, 0.3s) ease;
  }
  
  .refresh-button:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: var(--shadow-medium, 0 8px 16px rgba(0,0,0,0.15));
  }

  .refresh-button:hover:not(:disabled) .refresh-icon {
    transform: rotate(180deg);
  }
  
  .refresh-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
  }
  
  .error-banner {
    background: linear-gradient(135deg, var(--bg-tertiary, #ffeded), #fff0f0);
    color: var(--text-warning, #FF6B6B);
    padding: 1.5rem;
    border-radius: 12px;
    margin-bottom: 2rem;
    border: 1px solid rgba(255, 107, 107, 0.2);
    display: flex;
    align-items: center;
    gap: 1rem;
    box-shadow: var(--shadow-light, 0 4px 8px rgba(0,0,0,0.1));
  }

  .error-icon {
    width: 24px;
    height: 24px;
    flex-shrink: 0;
  }
  
  .error-banner p {
    margin: 0;
    font-weight: 500;
  }
  
  .loading, .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 4rem 2rem;
    text-align: center;
    background: var(--bg-secondary, #ffffff);
    border-radius: 16px;
    box-shadow: var(--shadow-light, 0 4px 8px rgba(0,0,0,0.1));
  }

  .empty-state {
    border: 2px dashed var(--border-color, #e0e0e0);
  }

  .empty-icon {
    width: 64px;
    height: 64px;
    color: var(--text-secondary, #666);
    margin-bottom: 1rem;
  }

  .empty-state h3 {
    color: var(--text-primary, #333);
    margin: 0 0 0.5rem 0;
    font-size: 1.5rem;
  }

  .empty-state p {
    color: var(--text-secondary, #666);
    margin: 0;
    font-size: 1rem;
  }
  
  .spinner {
    width: 48px;
    height: 48px;
    border: 4px solid var(--bg-tertiary, #eeeeee);
    border-top: 4px solid var(--accent-primary, #2C6E49);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 1.5rem;
  }

  .loading p {
    color: var(--text-secondary, #666);
    font-size: 1.1rem;
    margin: 0;
  }

  .gallery-stats {
    margin-bottom: 2rem;
    display: flex;
    justify-content: center;
  }

  .stats-text {
    background: var(--bg-secondary, #ffffff);
    color: var(--text-secondary, #666);
    padding: 0.5rem 1rem;
    border-radius: 20px;
    font-size: 0.9rem;
    font-weight: 500;
    border: 1px solid var(--border-color, #e0e0e0);
    box-shadow: var(--shadow-light, 0 2px 4px rgba(0,0,0,0.05));
  }
  
  .gallery-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(460px, 1fr));
    gap: 2rem;
    justify-items: center;
    padding: 1rem 0;
  }
  
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
  
  @media (max-width: 1200px) {
    .gallery-grid {
      grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
      gap: 1.5rem;
    }
  }
  
  @media (max-width: 768px) {
    .gallery-page {
      padding: 1rem;
    }

    .header {
      flex-direction: column;
      gap: 1rem;
      align-items: flex-start;
    }

    .gallery-title {
      font-size: 2rem;
    }

    .gallery-subtitle {
      font-size: 1rem;
    }

    .refresh-button {
      align-self: stretch;
      justify-content: center;
    }
    
    .gallery-grid {
      grid-template-columns: 1fr;
      gap: 1rem;
    }
  }
  
  @media (max-width: 520px) {
    .gallery-grid {
      grid-template-columns: 1fr;
    }
  }
</style>