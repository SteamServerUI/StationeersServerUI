<script>
  import TopNav from './components/nav/TopNav.svelte';
  import Sidebar from './components/nav/Sidebar.svelte';
  import MainContent from './components/MainContent.svelte';
  import BackendInitializerV7 from './BackendInitializerV7.svelte';
  import AuthGuardV7 from './AuthGuardV7.svelte';
  import ScreenNotSupported from './components/resuables/ScreenNotSupported.svelte';
  import { authState } from './services/api-v7';
  import './themes/theme.css';

  let activeView = $state('dashboard');
  let serverStatus = $state('checking');
  let serverError = $state(null);
  let forceShowApp = $state(false);
  let screenSupported = $state(true);

  const viewDefinitions = [
    { id: 'dashboard', name: 'Dashboard', icon: 'grid', permissions: ['server.view'] },
    { id: 'settings', name: 'Settings', icon: 'settings', permissions: [] },
    { id: 'console', name: 'Console', icon: 'terminal', permissions: ['logs.view'] },
    { id: 'logs', name: 'Logs', icon: 'file-text', permissions: ['logs.view'] },
    { id: 'gallery', name: 'Gallery', icon: 'globe', permissions: ['runfiles.view', 'plugins.view'] },
    { id: 'backups', name: 'Backups', icon: 'archive', permissions: ['backups.view'] },
    { id: 'plugins', name: 'Plugins', icon: 'plugin', permissions: ['plugins.view'] }
  ];

  let views = $derived(viewDefinitions.filter(view => !view.permissions.length || view.permissions.some(permission => $authState.permissions.includes(permission))));

  $effect(() => {
    const resize = () => screenSupported = window.innerWidth >= 1024 && window.innerHeight >= 600;
    resize();
    window.addEventListener('resize', resize);
    return () => window.removeEventListener('resize', resize);
  });

  function setActiveView(id) {
    if (views.some(view => view.id === id)) activeView = id;
  }

  function handleStatusChange(value) {
    serverStatus = value.status;
    serverError = value.error;
  }
</script>

{#if screenSupported || forceShowApp}
  <BackendInitializerV7 onStatusChange={handleStatusChange}>
    <AuthGuardV7 {serverStatus} {serverError}>
      <div class="app-container">
        <TopNav {views} {activeView} {setActiveView} />
        <div class="main-container">
          <Sidebar {views} {activeView} {setActiveView} />
          <MainContent {activeView} />
        </div>
      </div>
    </AuthGuardV7>
  </BackendInitializerV7>
{:else}
  <ScreenNotSupported onContinueAnyway={() => forceShowApp = true} />
{/if}

<style>
  .app-container { display: flex; flex-direction: column; width: 100vw; height: 100vh; overflow: hidden; }
  .main-container { display: flex; flex: 1; overflow: hidden; }
</style>
