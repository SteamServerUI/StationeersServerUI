<script>
  import { onMount } from 'svelte';
  import { backendConfig, setActiveBackend, apiFetch, logout } from '../../services/api-v7';
  import { authState } from '../../services/api-v7';
  import { navigate } from '../../services/router';
  import Icon from '../resuables/Icon.svelte';
  import themeService from '../../themes/theme';
  import { notify } from '../../services/activity';

  let { views = [], activeView = 'home' } = $props();
  let showBackends=$state(false), showUser=$state(false), online=$state(false), checking=$state(false), now=$state(new Date());
  let config=$state({active:'default',backends:{default:{url:'/'}}});
  let groups=$derived([...new Set(views.map(view=>view.group))]);
  let activeGroup=$derived(views.find(view=>view.id===activeView)?.group||'Home');
  let initials=$derived(($authState.user?.username||'SS').split(/[\s_-]+/).map(part=>part[0]).join('').slice(0,2).toUpperCase());
  let clock=$derived(now.toLocaleTimeString([], {hour:'2-digit',minute:'2-digit'}));
  let date=$derived(now.toLocaleDateString([], {month:'short',day:'numeric'}));
  const unsubscribe=backendConfig.subscribe(value=>config=value);

  onMount(()=>{
    check();
    const statusTimer=setInterval(check,15000);
    const clockTimer=setInterval(()=>now=new Date(),30000);
    const closeMenus=event=>{
      if(!event.target.closest('.workspace-switcher'))showBackends=false;
      if(!event.target.closest('.user'))showUser=false;
    };
    document.addEventListener('click',closeMenus);
    return()=>{clearInterval(statusTimer);clearInterval(clockTimer);unsubscribe();document.removeEventListener('click',closeMenus);};
  });

  async function check(){checking=true;try{const response=await apiFetch('/api/v3/capabilities');online=response.ok;}catch{online=false;}finally{checking=false;}}
  function openGroup(group){const target=views.find(view=>view.group===group);if(target)navigate(target.id);}
  async function choose(id){showBackends=false;if(id!==config.active){await setActiveBackend(id);notify('Workspace changed','info',id);}}
  async function signOut(){await logout();window.location.href='/';}
  function cycleTheme(){const id=themeService.nextTheme();notify('Appearance changed','success',themeService.getThemes().find(theme=>theme.id===id)?.name||id);}
</script>

<header class="topbar">
  <button class="brand" onclick={()=>navigate('home')} aria-label="Open Home">
    <span class="brand-mark">S7</span>
    <span class="brand-copy"><strong>SSUI</strong><small>SteamServerUI</small></span>
  </button>

  <nav class="nav-groups" aria-label="Primary navigation">
    {#each groups as group}
      <button class="group-tab" class:active={activeGroup===group} onclick={()=>openGroup(group)}>{group}</button>
    {/each}
  </nav>

  <div class="top-context">
    <div class="workspace-switcher">
      <button class="backend-button" onclick={event=>{event.stopPropagation();showBackends=!showBackends;showUser=false;}} aria-expanded={showBackends}>
        <span class="status" class:online class:checking></span>
        <span><small>Workspace</small><strong>{config.active}</strong></span>
        <span class="caret">⌄</span>
      </button>
      {#if showBackends}
        <div class="backend-popover shell-popover">
          <header><span>Saved backends</span><button onclick={()=>{showBackends=false;navigate('backends')}}>Manage</button></header>
          {#each Object.keys(config.backends) as id}
            <button class="backend-option" class:active={id===config.active} onclick={()=>choose(id)}>
              <span class="backend-state" class:online={id===config.active&&online}></span>
              <span><strong>{id}</strong><small>{config.backends[id].url||'Local backend'}</small></span>
              {#if id===config.active}<span class="check">✓</span>{/if}
            </button>
          {/each}
        </div>
      {/if}
    </div>

    <div class="datetime"><span>{date}</span><strong>{clock}</strong></div>
    <button class="icon-button" onclick={cycleTheme} title="Switch appearance"><span class="theme-orb"></span></button>
    <button class="icon-button" onclick={()=>navigate('activity')} class:active={activeView==='activity'} title="Open activity"><Icon name="bell"/></button>

    <div class="user">
      <button class="avatar" onclick={event=>{event.stopPropagation();showUser=!showUser;showBackends=false;}}>{initials}</button>
      {#if showUser}
        <div class="user-popover shell-popover">
          <header><span class="large-avatar">{initials}</span><span><strong>{$authState.user?.username}</strong><small>{$authState.user?.groupIds?.includes('system-owner')?'Owner':'Custom access'}</small></span></header>
          <button onclick={()=>{showUser=false;navigate('access')}}><Icon name="access" size={17}/>People & access</button>
          <button class="logout" onclick={signOut}><Icon name="logout" size={17}/>Sign out</button>
        </div>
      {/if}
    </div>
  </div>
</header>

<style>
  .topbar{height:var(--top-nav-height);flex:0 0 auto;display:grid;grid-template-columns:auto minmax(0,1fr) auto;align-items:center;padding:0 16px 0 12px;background:color-mix(in srgb,var(--material-solid) 82%,transparent);border-bottom:1px solid var(--border-color);box-shadow:0 8px 28px rgba(0,0,0,.12);backdrop-filter:blur(22px) saturate(125%);z-index:60}
  .brand{height:44px;display:flex;align-items:center;gap:10px;padding:4px 12px 4px 4px;margin-right:16px;background:transparent;border:0;text-align:left}.brand-mark{display:grid;place-items:center;width:36px;height:34px;border-radius:9px;background:linear-gradient(135deg,var(--accent-primary),var(--accent-secondary));color:#07101d;font-weight:900;letter-spacing:-.06em;box-shadow:0 0 18px color-mix(in srgb,var(--accent-primary) 18%,transparent)}.brand-copy strong,.brand-copy small{display:block}.brand-copy strong{font-size:.9rem;letter-spacing:.06em}.brand-copy small{font-size:.61rem;color:var(--text-muted);margin-top:1px}
  .nav-groups{height:100%;display:flex;align-items:stretch;min-width:0;overflow-x:auto;scrollbar-width:none}.nav-groups::-webkit-scrollbar{display:none}.group-tab{position:relative;flex:0 0 auto;border:0;border-radius:0;background:transparent;padding:0 14px;color:var(--text-secondary);font-size:.78rem;font-weight:650}.group-tab:hover{background:var(--bg-hover);color:var(--text-primary)}.group-tab.active{background:transparent;color:var(--text-primary)}.group-tab.active::after{content:'';position:absolute;left:14px;right:14px;bottom:0;height:2px;background:var(--accent-primary);box-shadow:0 0 10px var(--accent-primary)}
  .top-context{display:flex;align-items:center;gap:5px;margin-left:12px}.workspace-switcher,.user{position:relative}.backend-button{height:43px;display:flex;align-items:center;gap:9px;background:transparent;border-color:transparent;text-align:left;padding:5px 9px}.backend-button small,.backend-button strong{display:block}.backend-button small{font-size:.57rem;color:var(--text-muted);text-transform:uppercase;letter-spacing:.1em}.backend-button strong{font-size:.78rem;max-width:110px;overflow:hidden;text-overflow:ellipsis}.status,.backend-state{width:8px;height:8px;border-radius:50%;background:var(--text-danger);box-shadow:0 0 9px color-mix(in srgb,var(--text-danger) 55%,transparent)}.status.online,.backend-state.online{background:var(--text-success);box-shadow:0 0 9px color-mix(in srgb,var(--text-success) 55%,transparent)}.status.checking{animation:pulse 1s infinite}.caret{color:var(--text-muted)}
  .datetime{width:50px;display:grid;text-align:right;padding:0 8px;border-left:1px solid var(--border-color)}.datetime span{font-size:.58rem;color:var(--text-muted);text-transform:uppercase}.datetime strong{font-size:.75rem}.icon-button{width:38px;height:38px;display:grid;place-items:center;border:0;background:transparent;color:var(--text-secondary);padding:0}.icon-button.active{color:var(--text-accent);background:var(--bg-active)}.theme-orb{width:16px;height:16px;border-radius:50%;background:conic-gradient(from 30deg,var(--accent-primary),var(--accent-secondary),var(--accent-tertiary),var(--accent-primary));box-shadow:0 0 12px color-mix(in srgb,var(--accent-primary) 35%,transparent)}.avatar,.large-avatar{display:grid;place-items:center;border:0;background:linear-gradient(135deg,var(--accent-primary),var(--accent-secondary));color:#07101d;font-weight:850}.avatar{width:35px;height:35px;padding:0;border-radius:10px;margin-left:2px}.large-avatar{width:38px;height:38px;border-radius:11px}
  .shell-popover{position:absolute;top:50px;right:0;padding:8px;background:color-mix(in srgb,var(--material-solid) 96%,transparent);border:1px solid var(--border-color);border-radius:12px;box-shadow:var(--shadow-medium);backdrop-filter:blur(24px);z-index:100}.backend-popover{width:310px}.backend-popover>header{display:flex;align-items:center;justify-content:space-between;padding:5px 6px 9px;color:var(--text-secondary);font-size:.7rem;text-transform:uppercase;letter-spacing:.08em}.backend-popover>header button{padding:4px 7px;background:transparent;border:0;color:var(--text-accent);font-size:.7rem}.backend-option{width:100%;display:grid;grid-template-columns:9px 1fr auto;align-items:center;gap:10px;text-align:left;background:transparent;border:0;padding:10px}.backend-option strong,.backend-option small{display:block}.backend-option small{color:var(--text-muted);font-size:.66rem;margin-top:2px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}.backend-option.active{background:var(--bg-active)}.check{color:var(--text-success)}
  .user-popover{width:225px}.user-popover header{display:flex;align-items:center;gap:10px;padding:7px 7px 11px;border-bottom:1px solid var(--border-color);margin-bottom:4px}.user-popover header strong,.user-popover header small{display:block}.user-popover header small{font-size:.68rem;color:var(--text-muted);margin-top:2px}.user-popover>button{width:100%;display:flex;align-items:center;gap:9px;background:transparent;border:0;text-align:left}.user-popover .logout{color:var(--text-danger);border-top:1px solid var(--border-color);border-radius:0;margin-top:4px;padding-top:11px}
  @keyframes pulse{50%{opacity:.35}}@media(max-width:1180px){.brand-copy small,.datetime{display:none}.brand{margin-right:5px}.group-tab{padding:0 10px}.group-tab.active::after{left:10px;right:10px}}@media(max-width:900px){.brand-copy{display:none}.group-tab:nth-child(4),.group-tab:nth-child(5){display:none}.backend-button>span:nth-child(2){display:none}}
</style>
