<script>
  import TopNav from './components/nav/TopNav.svelte';
  import Sidebar from './components/nav/Sidebar.svelte';
  import MainContent from './components/MainContent.svelte';
  import BackendInitializer from './BackendInitializer.svelte';
  import AuthGuard from './AuthGuard.svelte';
  import './themes/theme.css';
  import ScreenNotSupported from './components/resuables/ScreenNotSupported.svelte';

  // Track active view
  let activeView = $state('dashboard');
  
  // Views available in the app
  const views = [
    { id: 'dashboard', name: 'Dashboard', icon: 'grid' },
    { id: 'settings', name: 'Settings', icon: 'settings' },
    { id: 'console', name: 'Console', icon: 'terminal' },
    { id: 'logs', name: 'Logs', icon: 'file-text' },
  ];

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

  // Run on mount and on window resize
  $effect(() => {
    checkScreenSize();
    window.addEventListener('resize', checkScreenSize);
    return () => window.removeEventListener('resize', checkScreenSize);
  });
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