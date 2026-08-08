<script>
  import { onMount } from 'svelte';
  import { authState, syncAuthState } from './services/api-v7';
  import LoginV7 from './LoginV7.svelte';
  import InitializingView from './components/resuables/InitializingView.svelte';

  let { serverStatus = 'online', serverError = null, children } = $props();
  let checked = $state(false);

  onMount(async () => {
    if (serverStatus === 'online') await syncAuthState();
    checked = true;
  });
</script>

{#if !checked || $authState.isAuthenticating}
  <InitializingView serverStatus="checking" />
{:else if serverStatus === 'error'}
  <InitializingView serverStatus="error" errorMessage={serverError?.toString() || 'Could not connect to the backend.'} />
{:else if !$authState.isAuthenticated}
  <LoginV7 />
{:else}
  {@render children?.()}
{/if}
