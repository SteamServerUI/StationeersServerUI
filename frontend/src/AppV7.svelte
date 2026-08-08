<script>
  import { onMount } from 'svelte';
  import TopNav from './components/nav/TopNav.svelte';
  import Sidebar from './components/nav/Sidebar.svelte';
  import MainContent from './components/MainContent.svelte';
  import ToastHost from './components/resuables/ToastHost.svelte';
  import BackendInitializerV7 from './BackendInitializerV7.svelte';
  import AuthGuardV7 from './AuthGuardV7.svelte';
  import ScreenNotSupported from './components/resuables/ScreenNotSupported.svelte';
  import { authState } from './services/api-v7';
  import { initializeActivity, closeActivity } from './services/activity';
  import { loadCapabilities } from './services/product';
  import { activeRoute, initializeRouter, navigate } from './services/router';
  import themeService from './themes/theme';
  import './themes/theme.css';

  let serverStatus = $state('checking');
  let serverError = $state(null);
  let forceShowApp = $state(false);
  let screenSupported = $state(true);
  let productLoadedFor = $state(null);

  const definitions = [
    { id:'home', name:'Home', icon:'home', group:'Home', permissions:['server.view'] },
    { id:'console', name:'Console', icon:'terminal', group:'Operate', permissions:['logs.view'] },
    { id:'logs', name:'Logs', icon:'logs', group:'Operate', permissions:['logs.view'] },
    { id:'activity', name:'Activity', icon:'activity', group:'Operate', permissions:[] },
    { id:'game', name:'Game', icon:'game', group:'Configure', permissions:['runfiles.view'] },
    { id:'gallery', name:'Gallery', icon:'gallery', group:'Configure', permissions:['runfiles.view','plugins.view'] },
    { id:'files', name:'Files', icon:'files', group:'Configure', permissions:['files.read'] },
    { id:'backups', name:'Backups', icon:'backup', group:'Protect', permissions:['backups.view'] },
    { id:'access', name:'People & access', icon:'access', group:'Access', permissions:[] },
    { id:'settings', name:'Settings', icon:'settings', group:'System', permissions:['settings.view'] },
    { id:'backends', name:'Backends', icon:'backends', group:'System', permissions:[] },
    { id:'plugins', name:'Plugins', icon:'plugin', group:'System', permissions:['plugins.view'], feature:'plugins' }
  ];

  let views = $derived(definitions.filter(view =>
    (!view.feature || $authState.features[view.feature]) &&
    (!view.permissions.length || view.permissions.some(permission => $authState.permissions.includes(permission)))
  ));

  onMount(() => {
    initializeRouter();
    themeService.initTheme();
    const resize = () => screenSupported = window.innerWidth >= 768 && window.innerHeight >= 500;
    resize();
    window.addEventListener('resize', resize);
    return () => {
      window.removeEventListener('resize', resize);
      closeActivity();
    };
  });

  $effect(() => {
    if (!$authState.isAuthenticated) {
      productLoadedFor = null;
      closeActivity();
      return;
    }
    if (views.length && !views.some(view => view.id === $activeRoute)) navigate(views[0].id, true);
    const identity = $authState.user?.id || $authState.user?.username || 'session';
    if (productLoadedFor !== identity) {
      productLoadedFor = identity;
      loadCapabilities().catch(error => console.warn('Could not load backend capabilities', error));
      initializeActivity();
    }
  });

  function handleStatusChange(value) {
    serverStatus = value.status;
    serverError = value.error;
  }
</script>

{#if screenSupported || forceShowApp}
  <BackendInitializerV7 onStatusChange={handleStatusChange}>
    <AuthGuardV7 {serverStatus} {serverError}>
      <div class="app-shell">
        <TopNav {views} activeView={$activeRoute} />
        <div class="workspace">
          <Sidebar {views} activeView={$activeRoute} />
          <MainContent activeView={$activeRoute} />
        </div>
        <ToastHost />
      </div>
    </AuthGuardV7>
  </BackendInitializerV7>
{:else}
  <ScreenNotSupported onContinueAnyway={() => forceShowApp = true} />
{/if}

<style>
  .app-shell{display:flex;flex-direction:column;width:100vw;height:100vh;overflow:hidden;isolation:isolate}
  .workspace{display:flex;flex:1;min-height:0;overflow:hidden}
</style>
