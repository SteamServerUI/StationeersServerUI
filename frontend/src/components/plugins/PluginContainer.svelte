<script>
  import { getCurrentBackendUrl } from '../../services/api-v7';
  import { onMount } from 'svelte';
  
  let { pluginPath = '' } = $props();
  
  let iframeRef;
  let loading = $state(true);
  let iframeSrc = $state('');

  async function constructIframeSrc() {
    try {
      const backendUrl = getCurrentBackendUrl();
      iframeSrc = `${backendUrl}${pluginPath}`;
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

  $effect(() => {
    if (pluginPath) {
      constructIframeSrc();
    }
  }); 
</script>

<div class="iframe-container" class:hidden={loading}>
  <iframe
  bind:this={iframeRef}
  src={iframeSrc}
  title="Plugin Container"
  onload={handleIframeLoad}
  onerror={handleIframeError}
  ></iframe>
</div>

{#if loading}
  <div class="loading">
    <p>Loading plugin...</p>
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