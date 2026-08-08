<script lang="js">
  import { apiFetch } from '../../../services/api-v7';
  import ErrorPopup from './ErrorPopup.svelte';

  // Props using $props in runes mode
  const { plugin } = $props();

  // State
  let isFlipped = $state(false);
  let isLoading = $state(false);
  let error = $state(null);
  let showErrorPopup = $state(false);

  $effect(() => {
    if (plugin && !plugin.identifier) {
      plugin.identifier = plugin.name || plugin.name;
    }
  });

  // Handle card flip
  function toggleFlip() {
    isFlipped = !isFlipped;
  }

  // Apply and download plugin in a single step
  async function applyAndDownload(event) {
    // Stop event propagation to prevent card flip
    event.stopPropagation();
    
    isLoading = true;

    try {
      // First select the plugin
      const response = await apiFetch('/api/v2/plugingallery/select', {
        method: 'POST',
        body: JSON.stringify({ name: plugin.name, redownload: false })
      });

      // Check for HTTP 409 (Conflict) indicating plugin already exists
      if (response.status === 409) {
        const confirmRedownload = window.confirm(`Plugin ${plugin.name} already exists. Would you like to re-download it?`);
        if (confirmRedownload) {
          // Retry with redownload: true
          await apiFetch('/api/v2/plugingallery/select', {
            method: 'POST',
            body: JSON.stringify({ name: plugin.name, redownload: true })
          });
        } else {
          // User cancelled redownload
          throw new Error('Plugin download cancelled');
        }
      }

      showSuccessPopup('Plugin downloaded and applied successfully!');
    } catch (err) {
      error = err.message || 'Failed to download and apply plugin';
      showErrorPopup = true;
    } finally {
      isLoading = false;
    }
  }

  // Success popup helper
  function showSuccessPopup(message) {
    error = message;
    showErrorPopup = true;
  }

  // Auto-dismiss error after 5 seconds
  $effect(() => {
    if (showErrorPopup) {
      const timer = setTimeout(() => {
        showErrorPopup = false;
      }, 5000);
      
      return () => clearTimeout(timer);
    }
  });

</script>

<div class="card-container">
  <div class="card" class:flipped={isFlipped} onclick={toggleFlip} onkeydown={(e) => e.key === 'Enter'} tabindex="0" role="button" aria-label={`Runfile card for ${plugin.name}`}>
    <!-- Front Side -->
    <div class="card-front" style="background-image: url('{plugin.background_url || ''}')">
    </div>
    
    <!-- Back Side - Now with side-by-side layout -->
    <div class="card-back">
      <h3 class="plugin-name">{plugin.name}</h3>
      
      <div class="card-back-content">
        <div class="metadata">
          <p><strong>Version:</strong> {plugin.version || 'N/A'}</p>
          <p><strong>Min Version:</strong> {plugin.min_version || 'N/A'}</p>
          <p><strong>Supported OS:</strong> {plugin.supported_os || 'N/A'}</p>
          {#if plugin.filename}
            <p><strong>Filename:</strong> {plugin.filename}</p>
          {/if}
        </div>
        
        <div class="actions-container">
          <div class="actions">
            {#if isLoading}
                <div class="spinner"></div>
            {:else}
              <button 
                class="download-button" 
                onclick={applyAndDownload}
                aria-label={`Download and apply ${plugin.name}`}
              >
                Download & Apply
              </button>
            {/if}
          </div>
          
          {#if showErrorPopup}
            <div class="popup-container">
              <ErrorPopup 
                message={error} 
                isSuccess={error && error.includes('successfully')}/>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .card-container {
    width: 460px;
    height: 215px;
    perspective: 1000px;
    margin: 1rem;
  }
  
  .card {
    position: relative;
    width: 100%;
    height: 100%;
    transition: transform var(--transition-speed, 0.6s);
    transform-style: preserve-3d;
    will-change: transform;
    cursor: pointer;
    box-shadow: var(--shadow-light, 0 4px 8px rgba(0,0,0,0.1));
    border-radius: 8px;
  }
  
  .card:hover {
    box-shadow: var(--shadow-medium, 0 8px 16px rgba(0,0,0,0.15));
  }
  
  .card.flipped {
    transform: rotateY(180deg);
  }
  
  .card-front, .card-back {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    -webkit-backface-visibility: hidden;
    backface-visibility: hidden;
    display: flex;
    flex-direction: column;
    align-items: center;
    background-color: var(--bg-secondary, #ffffff);
    border-radius: 8px;
    padding: 1rem;
    overflow: hidden;
  }
  
  .card-front {
    background-size: cover;
    background-position: center;
    background-repeat: no-repeat;
    position: relative;
    z-index: 1;
  }
  
  .card-front::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.3);
    border-radius: 8px;
    z-index: 0;
  }
  
  .plugin-name {
    color: var(--text-primary, #ffffff);
    text-align: center;
    font-size: 1.2rem;
    margin: 0 0 0.5rem 0;
    position: relative;
    z-index: 1;
  }
  
  .card-back {
    transform: rotateY(180deg);
    color: var(--text-primary, #333333);
    z-index: 0;
  }
  
  .card-back .plugin-name {
    color: var(--text-primary, #333333);
    margin-top: 0.5rem;
  }

  .card-back-content {
    display: flex;
    width: 100%;
    flex-direction: row;
    justify-content: space-between;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .metadata {
    flex: 1;
    margin: 0;
  }
  
  .metadata p {
    margin: 0.25rem 0;
    color: var(--text-secondary, #666666);
    font-size: 0.9rem;
  }
  
  .actions-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 100px;
  }
  
  .actions {
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 44px;
    position: relative;
  }
  
  .popup-container {
    margin-top: 0.5rem;
    width: 100%;
  }
  
  .download-button {
    background-color: var(--accent-primary, #2C6E49);
    color: white;
    border: none;
    border-radius: 4px;
    padding: 0.5rem 1rem;
    cursor: pointer;
    font-weight: bold;
    transition: all 0.2s ease;
    width: 100%;
  }
  
  .download-button:hover {
    transform: scale(1.05);
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
  
  @media (max-width: 600px) {
    .card-container {
      width: 100%;
      max-width: 250px;
      margin: 0.5rem auto;
    }
    
    /* Make layout vertical on small screens */
    .card-back-content {
      flex-direction: column;
    }
  }
</style>