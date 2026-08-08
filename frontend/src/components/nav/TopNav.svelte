<script>
  import { onMount } from 'svelte';
  import { backendConfig, setActiveBackend, apiFetch, logout } from '../../services/api-v7';
  import { authState } from '../../services/api-v7';
  import { navigate } from '../../services/router';
  import Icon from '../resuables/Icon.svelte';
  import themeService from '../../themes/theme';
  import { notify } from '../../services/activity';

  let showBackends=$state(false), showUser=$state(false), online=$state(false), checking=$state(false);
  let config=$state({active:'default',backends:{default:{url:'/'}}});
  const unsubscribe=backendConfig.subscribe(value=>config=value);

  onMount(()=>{ check(); const timer=setInterval(check,15000); return()=>{clearInterval(timer);unsubscribe();}; });
  async function check(){ checking=true; try{const response=await apiFetch('/api/v3/server/status');online=response.ok;}catch{online=false;}finally{checking=false;} }
  async function choose(id){showBackends=false;if(id!==config.active){await setActiveBackend(id);notify('Workspace changed','info',id);}}
  async function signOut(){await logout();window.location.href='/';}
  function cycleTheme(){const id=themeService.nextTheme();notify('Appearance changed','success',themeService.getThemes().find(theme=>theme.id===id)?.name||id);}
  let initials=$derived(($authState.user?.username||'SS').split(/[\s_-]+/).map(part=>part[0]).join('').slice(0,2).toUpperCase());
</script>

<header class="topbar">
  <div class="workspace-switcher">
    <button class="backend-button" onclick={()=>showBackends=!showBackends} aria-expanded={showBackends}>
      <span class="status" class:online class:checking></span>
      <span><small>Workspace</small><strong>{config.active}</strong></span>
      <span class="caret">⌄</span>
    </button>
    {#if showBackends}
      <div class="popover surface">
        <span class="eyebrow">Saved backends</span>
        {#each Object.keys(config.backends) as id}
          <button class:active={id===config.active} onclick={()=>choose(id)}><span>{id}</span><small>{config.backends[id].url||'Local backend'}</small></button>
        {/each}
        <button class="manage" onclick={()=>{showBackends=false;navigate('backends')}}>Manage connections</button>
      </div>
    {/if}
  </div>
  <div class="top-actions">
    <button class="icon-button" onclick={cycleTheme} title="Switch appearance"><span class="theme-orb"></span></button>
    <button class="icon-button" onclick={()=>navigate('activity')} title="Open activity"><Icon name="bell"/></button>
    <div class="user">
      <button class="avatar" onclick={()=>showUser=!showUser}>{initials}</button>
      {#if showUser}
        <div class="user-popover surface">
          <strong>{$authState.user?.username}</strong><small>{$authState.user?.groupIds?.includes('system-owner')?'Owner':'Custom access'}</small>
          <button onclick={()=>{showUser=false;navigate('access')}}><Icon name="access" size={17}/>People & access</button>
          <button class="logout" onclick={signOut}><Icon name="logout" size={17}/>Sign out</button>
        </div>
      {/if}
    </div>
  </div>
</header>

<style>
  .topbar{height:var(--top-nav-height);display:flex;align-items:center;justify-content:space-between;padding:10px 18px 0 18px;z-index:50}.workspace-switcher,.user{position:relative}.backend-button{height:48px;display:flex;align-items:center;gap:11px;background:transparent;border-color:transparent;text-align:left}.backend-button small,.backend-button strong{display:block}.backend-button small{font-size:.66rem;color:var(--text-muted);text-transform:uppercase;letter-spacing:.1em}.backend-button strong{font-size:.9rem}.status{width:9px;height:9px;border-radius:50%;background:var(--text-danger);box-shadow:0 0 12px color-mix(in srgb,var(--text-danger) 60%,transparent)}.status.online{background:var(--text-success);box-shadow:0 0 12px color-mix(in srgb,var(--text-success) 60%,transparent)}.status.checking{animation:pulse 1s infinite}.caret{color:var(--text-muted);margin-left:5px}.popover,.user-popover{position:absolute;top:55px;left:0;width:300px;padding:10px;display:grid;gap:5px}.popover .eyebrow{padding:8px}.popover button{display:flex;justify-content:space-between;text-align:left;background:transparent;border:0}.popover button small{color:var(--text-muted)}.popover .manage{margin-top:5px;border-top:1px solid var(--border-color);border-radius:0;padding-top:12px;color:var(--text-accent)}
  .top-actions{display:flex;align-items:center;gap:6px}.icon-button{width:42px;height:42px;display:grid;place-items:center;border:0;background:transparent;color:var(--text-secondary)}.theme-orb{width:18px;height:18px;border-radius:50%;background:linear-gradient(135deg,var(--accent-primary),var(--accent-secondary));box-shadow:0 0 15px color-mix(in srgb,var(--accent-primary) 55%,transparent)}.avatar{width:38px;height:38px;padding:0;border-radius:13px;background:linear-gradient(135deg,var(--accent-primary),var(--accent-secondary));color:#07101d;font-weight:850}.user-popover{left:auto;right:0;width:230px;padding:14px}.user-popover>small{color:var(--text-secondary);margin-top:-3px;margin-bottom:7px}.user-popover button{display:flex;align-items:center;gap:9px;background:transparent;border:0;text-align:left}.user-popover .logout{color:var(--text-danger);border-top:1px solid var(--border-color);border-radius:0;margin-top:4px;padding-top:12px}@keyframes pulse{50%{opacity:.35}}
</style>
