<script>
  import DashboardView from './Dashboard/DashboardView.svelte';
  import AppSettings from './settings/AppSettings.svelte';
  import RunfileSettings from './settings/RunfileSettings.svelte';
  import BackendSettings from './settings/BackendSettings.svelte';
  import FileView from './settings/FileView.svelte';
  import IdentitySettings from './settings/IdentitySettings.svelte';
  import LogsView from './views/LogsView.svelte';
  import ConsoleView from './views/ConsoleView.svelte';
  import ActivityView from './views/ActivityView.svelte';
  import PluginsView from './plugins/PluginsView.svelte';
  import GalleryView from './views/GalleryView.svelte';
  import BackupsView from './views/BackupsView.svelte';

  let { activeView='home' }=$props();
  const meta={
    console:['Console','Live game and backend output'],logs:['Logs','Investigate events without losing the signal'],activity:['Activity','Recent operations, notices and failures'],
    game:['Game configuration','Shape the server launch without fighting a runfile'],gallery:['Gallery','Discover game definitions and experimental extensions'],files:['Files','Edit the files this server actually uses'],
    backups:['Backups','Create and restore known-good server states'],access:['People & access','Users, groups, sessions, tokens and accountability'],
    settings:['Settings','Configure SteamServerUI'],backends:['Backends','Saved server workspaces and connection health'],plugins:['Plugins','Unsafe experimental native extensions']
  };
</script>

<main class="main-content">
  {#if activeView!=='home'}
    <header class="page-header"><div><span class="eyebrow">{activeView==='backups'?'Protect':activeView==='access'?'Access':activeView==='settings'||activeView==='backends'||activeView==='plugins'?'System':activeView==='game'||activeView==='gallery'||activeView==='files'?'Configure':'Operate'}</span><h1>{meta[activeView]?.[0]}</h1><p>{meta[activeView]?.[1]}</p></div>{#if activeView==='plugins'}<div class="danger-pill">Experimental · full system access</div>{/if}</header>
  {/if}
  <div class="view" class:home={activeView==='home'}>
    {#if activeView==='home'}<DashboardView/>
    {:else if activeView==='console'}<ConsoleView/>
    {:else if activeView==='logs'}<LogsView/>
    {:else if activeView==='activity'}<ActivityView/>
    {:else if activeView==='game'}<RunfileSettings/>
    {:else if activeView==='gallery'}<GalleryView/>
    {:else if activeView==='files'}<FileView/>
    {:else if activeView==='backups'}<BackupsView/>
    {:else if activeView==='access'}<IdentitySettings/>
    {:else if activeView==='settings'}<AppSettings activeSidebarTab="SSUI Settings"/>
    {:else if activeView==='backends'}<BackendSettings/>
    {:else if activeView==='plugins'}<div class="plugin-wrap"><div class="danger-banner"><strong>Danger, dragons ahead.</strong><p>Plugins are trusted native code with the same operating-system access as SSUI. There is no sandbox or safe production capability model yet.</p></div><PluginsView/></div>{/if}
  </div>
</main>

<style>
  .main-content{flex:1;min-width:0;overflow:auto;padding:28px clamp(22px,3vw,42px) 42px}.page-header{display:flex;justify-content:space-between;align-items:end;gap:20px;max-width:1320px;margin:0 auto 22px}.page-header h1{font-size:clamp(1.65rem,2.5vw,2.25rem);letter-spacing:-.035em;margin:4px 0}.page-header p{color:var(--text-secondary);margin:0}.view{max-width:1320px;margin:0 auto;min-height:calc(100% - 110px)}.view.home{max-width:1440px;height:100%}.danger-pill{color:var(--text-warning);border:1px solid color-mix(in srgb,var(--text-warning) 35%,transparent);background:color-mix(in srgb,var(--text-warning) 8%,transparent);border-radius:999px;padding:8px 13px;font-size:.78rem;font-weight:700}.danger-banner{padding:18px 20px;margin-bottom:16px;border:1px solid color-mix(in srgb,var(--text-warning) 38%,transparent);background:color-mix(in srgb,var(--text-warning) 9%,var(--material));border-radius:var(--radius-md);box-shadow:0 0 30px color-mix(in srgb,var(--text-warning) 8%,transparent)}.danger-banner strong{color:var(--text-warning);font-size:1.05rem}.danger-banner p{color:var(--text-secondary);margin:5px 0 0}@media(max-width:900px){.main-content{padding:18px 16px 28px}.page-header{align-items:start;flex-direction:column}.danger-pill{align-self:start}}
</style>
