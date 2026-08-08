<script>
  import AppSettings from './AppSettings.svelte';
  import RunfileSettings from './RunfileSettings.svelte';
  import BackendSettings from './BackendSettings.svelte';
  import FileView from './FileView.svelte';
  import IdentitySettings from './IdentitySettings.svelte';
  import { hasPermission } from '../../services/api-v7';

  const tabs = [
    { id: 'app', label: 'SSUI Settings', permission: 'settings.view' },
    { id: 'game', label: 'Game Settings', permission: 'runfiles.view' },
    { id: 'identity', label: 'People & Access', permission: null },
    { id: 'backends', label: 'Backends', permission: null },
    { id: 'files', label: 'Files', permission: 'files.read' }
  ];
  let active = $state('app');
  let visibleTabs = $derived(tabs.filter(tab => !tab.permission || hasPermission(tab.permission)));
</script>

<div class="settings-container">
  <div class="settings-sidebar">
    {#each visibleTabs as tab}
      <button class:active={active === tab.id} onclick={() => active = tab.id}>{tab.label}</button>
    {/each}
  </div>
  <div class="settings-content">
    {#if active === 'app'}<AppSettings activeSidebarTab="SSUI Settings" />
    {:else if active === 'game'}<RunfileSettings />
    {:else if active === 'identity'}<IdentitySettings />
    {:else if active === 'backends'}<BackendSettings />
    {:else if active === 'files'}<FileView />{/if}
  </div>
</div>

<style>
  .settings-container { display: flex; gap: 1rem; height: 100%; }
  .settings-sidebar { width: 190px; display: flex; flex-direction: column; gap: .4rem; }
  .settings-sidebar button { padding: .7rem .8rem; border: 0; border-radius: 5px; background: transparent; color: var(--text-primary); text-align: left; cursor: pointer; }
  .settings-sidebar button:hover { background: var(--bg-hover); }
  .settings-sidebar button.active { background: var(--bg-active); color: var(--accent-primary); }
  .settings-content { flex: 1; min-width: 0; padding: 1.25rem; overflow-y: auto; border-radius: 6px; background: var(--bg-secondary); }
  @media (max-width: 800px) { .settings-container { flex-direction: column; } .settings-sidebar { width: 100%; flex-direction: row; overflow-x: auto; } }
</style>
