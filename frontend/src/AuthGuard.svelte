<script>
  import { onMount, onDestroy } from 'svelte';
  import { authState, syncAuthState } from './services/api';
  import Login from './Login.svelte';
  import InitializingView from './components/resuables/InitializingView.svelte';
  
  /**
   * @typedef {Object} Props
   * @property {boolean} [checkAuth] - Props - Set to false to disable auth checking
   * @property {string} [serverStatus] - 'online', 'offline', 'error'
   * @property {any} [serverError]
   * @property {import('svelte').Snippet} [children]
   */

  /** @type {Props} */
  let {
    checkAuth = true,
    serverStatus = 'online',
    serverError = null,
    children
  } = $props();
  
  // Local state
  let isChecking = $state(true);
  let unsubscribe;
  
  onMount(async () => {
    if (checkAuth && serverStatus === 'online') {
      // Subscribe to auth state changes
      unsubscribe = authState.subscribe(state => {
        isChecking = state.isAuthenticating;
      });
      
      // Sync auth state on mount
      try {
        await syncAuthState();
      } catch (err) {
        console.error('Auth check failed:', err);
        // Even on error, we're done checking
        isChecking = false;
      }
    } else {
      isChecking = false;
    }
  });
  
  onDestroy(() => {
    if (unsubscribe) unsubscribe();
  });
</script>

{#if isChecking}
  <InitializingView serverStatus="checking" />
{:else if serverStatus === 'error'}
  <InitializingView 
    serverStatus="error" 
    errorMessage={serverError ? serverError.toString() : "There was an error connecting to the server."} 
  />
{:else if !$authState.isAuthenticated && checkAuth}
  <Login />
{:else}
  {@render children?.()}
{/if}