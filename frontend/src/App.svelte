<script>
  import TopNav from './components/nav/TopNav.svelte';
  import Sidebar from './components/nav/Sidebar.svelte';
  import MainContent from './components/MainContent.svelte';
  import BackendInitializer from './BackendInitializerV7.svelte';
  import AuthGuard from './AuthGuardV7.svelte';
  import { apiFetch } from './services/api-v7';
  import { pluginsList } from './services/plugins';
  import './themes/theme.css';
  import ScreenNotSupported from './components/resuables/ScreenNotSupported.svelte';

  // Track active view
  let activeView = $state('dashboard');
  let hasPlugins = $state(false);
  
  // Views available in the app
  let views = $state([
    { id: 'dashboard', name: 'Dashboard', icon: 'grid' },
    { id: 'settings', name: 'Settings', icon: 'settings' },
    { id: 'console', name: 'Console', icon: 'terminal' },
    { id: 'logs', name: 'Logs', icon: 'file-text' },
    { id: 'gallery', name: 'Gallery', icon: 'globe' },
    { id: 'backups', name: 'Backups', icon: 'archive' },
  ]);

  // Set active view function
  function setActiveView(viewId) {
    activeView = viewId;
  }

  // Server status state
  let serverStatus = $state('checking');
  let serverError = $state(null);
  
  function handleStatusChange(status) {
    serverStatus = status.status;
    serverError = status.error;
  }

  // Screen size check
  let isScreenSupported = $state(true);
  let forceShowApp = $state(false); // New state for override

  function checkScreenSize() {
    isScreenSupported = window.innerWidth >= 1024 && window.innerHeight >= 600;
  }

  // Handle the "continue anyway" action
  function handleContinueAnyway() {
    forceShowApp = true;
  }

  // Check if plugins are available
  async function checkPlugins() {
    try {
      const response = await apiFetch('/api/v3/plugins/list/apiroutes');
      const data = await response.json();
      hasPlugins = Array.isArray(data) && data.length > 0;
      
      // Store plugins data for use in PluginsView
      pluginsList.set(data || []);
      
      // Add plugins view if plugins exist
      if (hasPlugins) {
        // check if plugins view already exists
        const existingPluginsView = views.find(view => view.id === 'plugins');
        if (!existingPluginsView) {
          views = [
            ...views,
            { id: 'plugins', name: 'Plugins', icon: 'plugin' },
          ];
        }
      }
    } catch (error) {
      console.error('Error checking plugins:', error);
      hasPlugins = false;
      pluginsList.set([]);
    }
  }

  // Run on mount and on window resize
  $effect(() => {
    checkScreenSize();
    checkPlugins();
    window.addEventListener('resize', checkScreenSize);
    return () => window.removeEventListener('resize', checkScreenSize);
  });

  // run checkPlugins every 25 seconds
  setInterval(checkPlugins, 25000);
</script>

{#if isScreenSupported || forceShowApp}
  <BackendInitializer onStatusChange={handleStatusChange}>
    <AuthGuard serverStatus={serverStatus} serverError={serverError}>
      <div class="app-container">
        <TopNav {views} {activeView} {setActiveView} />
        
        <div class="main-container">
          <Sidebar {views} {activeView} {setActiveView} />
          <MainContent {activeView} />
        </div>
      </div>
    </AuthGuard>
  </BackendInitializer>
{:else}
  <ScreenNotSupported onContinueAnyway={handleContinueAnyway} />
{/if}

<style>
  .app-container {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100vw;
    overflow: hidden;
  }
  
  .main-container {
    display: flex;
    flex: 1;
    overflow: hidden;
  }
</style>