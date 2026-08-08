<script>
  import { fade } from 'svelte/transition';
  import DashboardView from './Dashboard/DashboardView.svelte';
  import SettingsView from './settings/SettingsViewV7.svelte';
  import LogsView from './views/LogsView.svelte';
  import ConsoleView from './views/ConsoleView.svelte';
  import PluginsView from './plugins/PluginsView.svelte';
  import GalleryView from './views/GalleryView.svelte';
  import BackupsView from './views/BackupsView.svelte';

  /**
   * @typedef {Object} Props
   * @property {string} [activeView]
   */

  /** @type {Props} */
  let { activeView = 'dashboard' } = $props();

  // View metadata for headers
  const viewContent = {
    dashboard: {
      title: 'Dashboard',
      description: 'Server overview and statistics'
    },
    servers: {
      title: 'Servers',
      description: 'Manage your game servers'
    },
    settings: {
      title: 'Settings',
      description: 'Configure application settings'
    },
    logs: {
      title: 'Logs',
      description: 'View server logs and events'
    },
    console: {
      title: 'Console',
      description: 'View server console output'
    },
    gallery: {
      title: 'Gallery',
      description: 'Browse runfiles & plugins'
    },
    plugins: {
      title: 'Plugins',
      description: 'Manage plugins'
    },
    backups: {
      title: 'Backups',
      description: 'Manage backups'
    }
  };
</script>

<main class="main-content">
  {#if activeView !== 'dashboard'}
    <div class="view-header">
      <h1>{viewContent[activeView].title}</h1>
      <p class="description">{viewContent[activeView].description}</p>
    </div>
  {/if}

  <div class="view-container">
    {#if activeView === 'dashboard'}
      <div class="view-content" in:fade={{ duration: 350, delay: 5 }} out:fade={{ duration: 200 }}>
        <DashboardView />
      </div>
    {:else if activeView === 'settings'}
      <div class="view-content" in:fade={{ duration: 350, delay: 5 }} out:fade={{ duration: 200 }}>
        <SettingsView />
      </div>
    {:else if activeView === 'logs'}
      <div class="view-content" in:fade={{ duration: 350, delay: 5 }} out:fade={{ duration: 200 }}>
        <LogsView />
      </div>
    {:else if activeView === 'console'}
      <div class="view-content" in:fade={{ duration: 350, delay: 5}} out:fade={{ duration: 200 }}>
        <ConsoleView />
      </div>
    {:else if activeView === 'gallery'}
      <div class="view-content" in:fade={{ duration: 350, delay: 5 }} out:fade={{ duration: 200 }}>
        <GalleryView />
      </div>
    {:else if activeView === 'plugins'}
      <div class="view-content" in:fade={{ duration: 350, delay: 5 }} out:fade={{ duration: 200 }}>
        <PluginsView />
      </div>
    {:else if activeView === 'backups'}
      <div class="view-content" in:fade={{ duration: 350, delay: 5 }} out:fade={{ duration: 200 }}>
        <BackupsView />
      </div>
    {/if}
  </div>
</main>

<style>
  .main-content {
    flex: 1;
    overflow-y: auto;
    padding: 1rem;
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  .view-header {
    margin-bottom: 1.5rem;
  }

  .view-header h1 {
    margin: 0 0 0.5rem 0;
    font-size: 1.8rem;
    font-weight: 500;
  }

  .description {
    color: var(--text-secondary); 
    margin: 0;
  }

  .view-container {
    flex: 1;
    position: relative;
  }

  .view-content {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    width: 100%;
    height: 100%;
  }
</style>