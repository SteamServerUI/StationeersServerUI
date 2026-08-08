<script>
  import Icon from '../resuables/Icon.svelte';
  import { navigate } from '../../services/router';
  let { views = [], activeView = 'home' } = $props();
  let hovered = $state(false);
  let pinned = $state(localStorage.getItem('ssui-rail-pinned') === 'true');
  let expanded = $derived(hovered || pinned);
  let groups = $derived([...new Set(views.map(view => view.group))]);
  function togglePin(){ pinned=!pinned; localStorage.setItem('ssui-rail-pinned', String(pinned)); }
</script>

<aside class="rail surface" class:expanded onmouseenter={() => hovered=true} onmouseleave={() => hovered=false}>
  <button class="brand" onclick={() => navigate('home')} aria-label="Open Home"><span class="mark">S7</span><span class="brand-name">SteamServerUI</span></button>
  <nav>
    {#each groups as group}
      <section>
        <span class="group-name">{group}</span>
        {#each views.filter(view => view.group === group) as view}
          <button class="rail-link" class:active={activeView === view.id} onclick={() => navigate(view.id)} title={!expanded ? view.name : undefined}>
            <Icon name={view.icon}/><span>{view.name}</span>
          </button>
        {/each}
      </section>
    {/each}
  </nav>
  <button class="pin" class:active={pinned} onclick={togglePin}><span>{pinned ? 'Unpin rail' : 'Pin rail'}</span><Icon name="chevron" size={16}/></button>
</aside>

<style>
  .rail{width:var(--sidebar-collapsed-width);margin:10px 0 10px 10px;border-radius:22px;display:flex;flex-direction:column;overflow:hidden;transition:width var(--transition-speed);z-index:20;flex:0 0 auto}
  .rail.expanded{width:var(--sidebar-width)}.brand{height:58px;display:flex;align-items:center;gap:11px;border:0;background:transparent;padding:10px 14px;white-space:nowrap}
  .mark{display:grid;place-items:center;width:42px;height:36px;border-radius:12px;background:linear-gradient(135deg,var(--accent-primary),var(--accent-secondary));color:#07101d;font-weight:900;letter-spacing:-.06em;box-shadow:0 0 24px color-mix(in srgb,var(--accent-primary) 25%,transparent)}
  .brand-name{font-weight:740;opacity:0;transition:opacity var(--transition-speed)}.expanded .brand-name{opacity:1}
  nav{flex:1;overflow-y:auto;overflow-x:hidden;padding:5px 9px 10px}section+section{margin-top:10px}.group-name{display:block;height:18px;padding:0 12px;color:var(--text-muted);font-size:.62rem;font-weight:750;letter-spacing:.12em;text-transform:uppercase;white-space:nowrap;opacity:0}.expanded .group-name{opacity:1}
  .rail-link{width:100%;height:42px;display:flex;align-items:center;gap:13px;border:0;background:transparent;padding:0 13px;white-space:nowrap;color:var(--text-secondary);position:relative}.rail-link span{opacity:0;transition:opacity var(--transition-speed)}.expanded .rail-link span{opacity:1}.rail-link:hover{color:var(--text-primary)}.rail-link.active{color:var(--text-primary);background:var(--bg-active)}.rail-link.active::before{content:'';position:absolute;left:2px;width:3px;height:20px;border-radius:8px;background:var(--accent-primary);box-shadow:0 0 12px var(--accent-primary)}
  .pin{margin:8px;height:38px;display:flex;align-items:center;justify-content:space-between;gap:10px;border:0;background:transparent;color:var(--text-muted);white-space:nowrap}.pin span{opacity:0}.expanded .pin span{opacity:1}.pin :global(svg){transition:transform var(--transition-speed)}.pin.active :global(svg){transform:rotate(180deg)}
</style>
