<script>
  import { getCurrentBackendUrl } from '../../services/api';
  import { onMount } from 'svelte';
  
  let iframeRef;
  let loading = true;
  let iframeSrc = '';

  async function constructIframeSrc() {
    try {

      const backendUrl = getCurrentBackendUrl();
      // Construct the full URL with the folder query parameter
      iframeSrc = `${backendUrl}/detectionmanager`;
    } catch (error) {
      console.error('Error constructing iframe src:', error);
    }
  }

  function handleIframeLoad() {
    loading = false;
  }

  function handleIframeError() {
    loading = false;
    console.error('Iframe failed to load');
  }

  // Run on component mount
  onMount(() => {
    constructIframeSrc();
  }); 
</script>


<div class="iframe-container" class:hidden={loading}>
  <iframe
  bind:this={iframeRef}
  src={iframeSrc}
  title="Detection Manager"
  on:load={handleIframeLoad}
  on:error={handleIframeError}
  ></iframe>
</div>

{#if loading}
  <div class="loading">
    <p>Loading Detection Manager...</p>
  </div>
{/if}

<style>
  .loading {
    padding: 2rem;
    text-align: center;
  }
  
  .iframe-container {
    width: 100%;
    height: 100%;
    border: none;
    border-radius: 8px;
    overflow: hidden;
  }
  
  .iframe-container.hidden {
    display: none;
  }
  
  iframe {
    width: 100%;
    height: 100%;
    border: none;
    background: none;
  }
</style>