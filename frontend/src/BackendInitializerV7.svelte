<script>
  import { onMount } from 'svelte';
  import { initializeApiService, apiFetchTimeout } from './services/api-v7';
  import InitializingView from './components/resuables/InitializingView.svelte';

  let { onStatusChange = () => {}, children } = $props();
  let ready = $state(false);
  let status = $state('checking');
  let errorMessage = $state(null);

  onMount(async () => {
    initializeApiService();
    try {
      const response = await apiFetchTimeout('/api/v2/auth/setup/status', {}, 5000);
      if (!response.ok) throw new Error(`Backend returned ${response.status}`);
      status = 'online';
    } catch (error) {
      status = 'error';
      errorMessage = error.name === 'AbortError' ? 'Connection timed out.' : error.message;
    }
    onStatusChange({ status, error: errorMessage });
    ready = true;
  });
</script>

{#if ready}
  {@render children?.()}
{:else}
  <InitializingView serverStatus={status} {errorMessage} />
{/if}
