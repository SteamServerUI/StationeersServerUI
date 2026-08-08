<script>
  import Icon from '../resuables/Icon.svelte';
  import { navigate } from '../../services/router';
  import { productState } from '../../services/product';

  let { views = [], activeView = 'home' } = $props();
  const saved=localStorage.getItem('ssui-sidebar-pinned');
  let hovered=$state(false), pinned=$state(saved===null?true:saved==='true');
  let expanded=$derived(hovered||pinned);
  let groups=$derived([...new Set(views.map(view=>view.group))]);
  function togglePin(){pinned=!pinned;localStorage.setItem('ssui-sidebar-pinned',String(pinned));}
</script>

<aside class="sidebar" class:expanded onmouseenter={()=>hovered=true} onmouseleave={()=>hovered=false}>
  <header>
    <span>Navigation</span>
    <button class="pin" class:active={pinned} onclick={togglePin} title={pinned?'Collapse sidebar':'Keep sidebar open'}><Icon name="chevron" size={16}/></button>
  </header>

  <nav>
    {#each groups as group}
      <section>
        <span class="group-name">{group}</span>
        {#each views.filter(view=>view.group===group) as view}
          <button class="side-link" class:active={activeView===view.id} onclick={()=>navigate(view.id)} title={!expanded?view.name:undefined}>
            <Icon name={view.icon} size={19}/><span>{view.name}</span>
          </button>
        {/each}
      </section>
    {/each}
  </nav>

  <footer>
    <span class="footer-dot"></span>
    <span class="footer-copy"><strong>S7 connected</strong><small>{#if $productState.backend?.version}v{$productState.backend.version} · {$productState.apiVersion}{:else}SteamServerUI v7{/if}</small></span>
  </footer>
</aside>

<style>
  .sidebar{width:var(--sidebar-collapsed-width);flex:0 0 auto;display:flex;flex-direction:column;overflow:hidden;background:color-mix(in srgb,var(--material-solid) 70%,transparent);border-right:1px solid var(--border-color);transition:width var(--transition-speed);backdrop-filter:blur(18px) saturate(120%);z-index:30}.sidebar.expanded{width:var(--sidebar-width)}
  header{height:48px;flex:0 0 auto;display:flex;align-items:center;justify-content:space-between;padding:0 11px 0 18px;border-bottom:1px solid var(--border-color);color:var(--text-muted);font-size:.62rem;font-weight:750;letter-spacing:.12em;text-transform:uppercase;white-space:nowrap}header>span{opacity:0;transition:opacity var(--transition-speed)}.expanded header>span{opacity:1}.pin{width:32px;height:32px;display:grid;place-items:center;padding:0;border:0;background:transparent;color:var(--text-muted)}.pin :global(svg){transform:rotate(0);transition:transform var(--transition-speed)}.pin.active :global(svg){transform:rotate(180deg)}
  nav{flex:1;overflow-y:auto;overflow-x:hidden;padding:10px 8px}section{padding:5px 0 8px}section+section{border-top:1px solid color-mix(in srgb,var(--border-color) 65%,transparent);padding-top:13px}.group-name{display:block;height:20px;padding:0 10px;color:var(--text-muted);font-size:.58rem;font-weight:750;letter-spacing:.13em;text-transform:uppercase;white-space:nowrap;opacity:0}.expanded .group-name{opacity:1}.side-link{width:100%;height:40px;display:flex;align-items:center;gap:12px;border:0;background:transparent;padding:0 13px;border-radius:8px;white-space:nowrap;color:var(--text-secondary);position:relative}.side-link>span{opacity:0;transition:opacity var(--transition-speed)}.expanded .side-link>span{opacity:1}.side-link:hover{color:var(--text-primary);background:var(--bg-hover)}.side-link.active{color:var(--text-primary);background:var(--bg-active)}.side-link.active::before{content:'';position:absolute;left:-8px;width:3px;height:22px;border-radius:0 4px 4px 0;background:var(--accent-primary);box-shadow:0 0 10px var(--accent-primary)}
  footer{height:54px;display:flex;align-items:center;gap:11px;padding:0 20px;border-top:1px solid var(--border-color);white-space:nowrap}.footer-dot{width:7px;height:7px;flex:0 0 auto;border-radius:50%;background:var(--text-success);box-shadow:0 0 9px color-mix(in srgb,var(--text-success) 60%,transparent)}.footer-copy{opacity:0;transition:opacity var(--transition-speed)}.expanded .footer-copy{opacity:1}.footer-copy strong,.footer-copy small{display:block}.footer-copy strong{font-size:.7rem}.footer-copy small{font-size:.6rem;color:var(--text-muted);margin-top:2px}
</style>
