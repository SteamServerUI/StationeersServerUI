<script>
  import { onMount, onDestroy } from 'svelte';
  import { apiFetch, hasPermission } from '../../services/api-v7';
  import { notify } from '../../services/activity';
  import { navigate } from '../../services/router';

  let status=$state(null), stats=$state(null), players=$state([]), backup=$state(null), loading=$state(true), acting=$state('');
  let background=$state(''), input, refreshTimer;

  onMount(()=>{
    refresh();
    refreshTimer=setInterval(refresh,15000);
    loadBackground();
    return()=>clearInterval(refreshTimer);
  });
  onDestroy(()=>{if(background.startsWith('blob:'))URL.revokeObjectURL(background);});

  async function json(path){
    const response=await apiFetch(path);
    if(!response.ok)throw new Error(`Request failed (${response.status})`);
    return response.json();
  }
  async function refresh(){
    try{
      const snapshot=await json('/api/v3/overview');
      const data=snapshot.data||snapshot;
      status=data.server;stats=data.system;players=data.players||[];backup=data.backup;
    }catch(error){
      notify('Overview is temporarily unavailable','error',error.message);
    }finally{loading=false;}
  }
  async function serverAction(action){
    if(action!=='start'&&!window.confirm(action==='restart'?'Restart the game server?':'Stop the game server?'))return;
    acting=action;
    try{
      const response=await apiFetch(`/api/v3/server/${action}`,{method:'POST'});
      const data=await response.json().catch(()=>({}));
      if(!response.ok)throw new Error(data?.error?.message||data?.message||`Request failed (${response.status})`);
      notify(`Server ${action} requested`,'success',data?.message||'The backend accepted the operation.');
      await refresh();
    }catch(error){notify(`Could not ${action} server`,'error',error.message);}
    finally{acting='';}
  }
  async function loadBackground(){
    try{const response=await apiFetch('/files/dashboard-background.png');if(response.ok)background=URL.createObjectURL(await response.blob());}catch{}
  }
  async function upload(event){
    const file=event.target.files?.[0];if(!file)return;
    const body=new FormData();body.append('file',file);
    try{
      const response=await apiFetch('/api/v3/settings/files/background/upload',{method:'POST',body});
      if(!response.ok)throw new Error('Upload failed');
      notify('Home background updated','success');
      if(background.startsWith('blob:'))URL.revokeObjectURL(background);
      background=URL.createObjectURL(file);
    }catch(error){notify('Background upload failed','error',error.message);}
    finally{event.target.value='';}
  }
  function playerName(player){if(player?.name)return player.name;const value=Object.values(player||{})[0];return value?.username||'Unknown player';}
  let memory=$derived(typeof stats?.memoryUsage==='number'?`${(stats.memoryUsage/1024).toFixed(1)} GB`:'—');
  let disk=$derived(typeof stats?.diskUsage==='number'?`${(100-stats.diskUsage).toFixed(0)}% free`:'—');
</script>

<section class="home" style={background?`--home-image:url("${background}")`:''}>
  <div class="home-wash"></div>
  <div class="home-content">
    <header class="welcome">
      <div><span class="eyebrow">Your server, at a glance</span><h1>{status?.isRunning?'Everything is running.':'Ready when you are.'}</h1><p>{status?.uuid?`Server ${status.uuid}`:'SteamServerUI is connected and waiting for server state.'}</p></div>
      <button class="background-button" onclick={()=>input?.click()}>Change background</button>
      <input bind:this={input} onchange={upload} type="file" accept="image/png,image/jpeg,image/webp" hidden>
    </header>

    <article class="hero surface" class:running={status?.isRunning}>
      <div class="hero-state"><span class="state-orb"></span><div><span class="eyebrow">Game server</span><h2>{loading?'Checking…':status?.isRunning?'Online':'Stopped'}</h2><p>{players.length} connected {players.length===1?'player':'players'}</p></div></div>
      {#if hasPermission('server.control')}
        <div class="hero-actions">
          {#if status?.isRunning}
            <button class="secondary" disabled={!!acting} onclick={()=>serverAction('stop')}>{acting==='stop'?'Stopping…':'Stop'}</button>
            <button class="primary" disabled={!!acting} onclick={()=>serverAction('restart')}>Restart</button>
          {:else}
            <button class="primary" disabled={!!acting} onclick={()=>serverAction('start')}>{acting==='start'?'Starting…':'Start server'}</button>
          {/if}
        </div>
      {/if}
    </article>

    <div class="dashboard-board surface">
      <div class="dashboard-grid">
        <article class="health">
          <header><div><span class="eyebrow">Host health</span><h3>{stats?.osName||'Backend host'}</h3></div><span>{stats?.uptime||'Uptime unavailable'}</span></header>
          <div class="metrics"><div><strong>{typeof stats?.cpuUsage==='number'?`${stats.cpuUsage.toFixed(0)}%`:'—'}</strong><span>CPU</span><i style={`--value:${stats?.cpuUsage||0}%`}></i></div><div><strong>{memory}</strong><span>Memory</span><i style="--value:52%"></i></div><div><strong>{disk}</strong><span>Storage</span><i style={`--value:${100-(stats?.diskUsage||100)}%`}></i></div></div>
        </article>
        <article class="players">
          <header><div><span class="eyebrow">Players</span><h3>Connected now</h3></div><button onclick={()=>navigate('console')}>Open console</button></header>
          {#if players.length}<div class="player-list">{#each players.slice(0,5) as player}<div><span class="player-avatar">{playerName(player).slice(0,1).toUpperCase()}</span><strong>{playerName(player)}</strong><small>Online</small></div>{/each}</div>{:else}<div class="quiet"><span>No players connected</span><small>The server is yours for the moment.</small></div>{/if}
        </article>
        <article class="backup">
          <span class="eyebrow">Protection</span><h3>{backup?.isRunning?'Backup in progress':backup?.systemReady?'Backup system ready':'Backup needs attention'}</h3><p>{backup?.isLoopRunning?'Automatic backup scheduling is active.':'Review the backup schedule and storage before you need it.'}</p><button onclick={()=>navigate('backups')}>Manage backups</button>
        </article>
        <article class="quick">
          <span class="eyebrow">Jump back in</span>
          <div><button onclick={()=>navigate('logs')}><strong>Logs</strong><small>See what just happened</small></button><button onclick={()=>navigate('game')}><strong>Game setup</strong><small>Launch arguments and config</small></button><button onclick={()=>navigate('activity')}><strong>Activity</strong><small>Operations and failures</small></button></div>
        </article>
      </div>
    </div>
  </div>
</section>

<style>
  .home{position:relative;min-height:100%;border-radius:16px;overflow:hidden;background-image:var(--home-image,linear-gradient(135deg,color-mix(in srgb,var(--canvas-glow) 58%,var(--canvas)),var(--canvas) 58%,color-mix(in srgb,var(--accent-secondary) 10%,var(--canvas))));background-size:cover;background-position:center}.home-wash{position:absolute;inset:0;background:linear-gradient(180deg,rgba(5,8,14,.08),color-mix(in srgb,var(--canvas) 78%,transparent));backdrop-filter:saturate(110%)}.home-content{position:relative;padding:clamp(26px,3.5vw,46px);max-width:1280px;margin:auto}.welcome{display:flex;justify-content:space-between;align-items:end;gap:20px;margin-bottom:24px}.welcome h1{font-size:clamp(2rem,4vw,3.5rem);letter-spacing:-.055em;line-height:1;margin:7px 0 10px;max-width:760px}.welcome p{color:var(--text-secondary);margin:0}.background-button{background:rgba(7,12,20,.28);backdrop-filter:blur(12px)}
  .hero{display:flex;align-items:center;justify-content:space-between;gap:25px;padding:22px 26px;margin-bottom:14px}.hero-state{display:flex;align-items:center;gap:17px}.state-orb{width:42px;height:42px;border-radius:50%;background:var(--text-muted);box-shadow:inset 0 1px 2px rgba(255,255,255,.3),0 0 22px rgba(130,140,160,.22)}.hero.running .state-orb{background:var(--text-success);box-shadow:inset 0 1px 2px rgba(255,255,255,.6),0 0 28px color-mix(in srgb,var(--text-success) 45%,transparent)}.hero h2{font-size:1.7rem;margin:2px 0}.hero p{color:var(--text-secondary);margin:0}.hero-actions{display:flex;gap:8px}.hero-actions .primary{background:var(--accent-primary);color:#07101d;border:0;font-weight:800;padding:.72rem 1.05rem}.hero-actions .secondary{background:transparent}
  .dashboard-board{overflow:hidden}.dashboard-grid{display:grid;grid-template-columns:1.22fr .9fr}.dashboard-grid article{padding:23px 25px;min-width:0}.health{border-right:1px solid var(--border-color);border-bottom:1px solid var(--border-color)}.players{border-bottom:1px solid var(--border-color)}.backup{border-right:1px solid var(--border-color)}.dashboard-grid header{display:flex;align-items:start;justify-content:space-between;gap:18px}.dashboard-grid h3{font-size:1.16rem;margin:4px 0}.dashboard-grid header>span,.backup p{color:var(--text-secondary);font-size:.8rem}.metrics{display:grid;grid-template-columns:repeat(3,1fr);gap:16px;margin-top:23px}.metrics div{display:grid;gap:4px}.metrics strong{font-size:1.3rem}.metrics span{font-size:.68rem;color:var(--text-muted);text-transform:uppercase;letter-spacing:.08em}.metrics i{height:3px;background:var(--bg-hover);border-radius:4px;overflow:hidden}.metrics i::after{content:'';display:block;height:100%;width:var(--value);background:var(--accent-primary);border-radius:4px}.players header button{background:transparent}.player-list{display:grid;margin-top:12px}.player-list>div{display:grid;grid-template-columns:34px 1fr auto;align-items:center;gap:10px;padding:8px 0;border-top:1px solid var(--border-color)}.player-avatar{display:grid;place-items:center;width:30px;height:30px;border-radius:8px;background:var(--bg-active);color:var(--text-accent)}.player-list small{color:var(--text-success)}.quiet{min-height:94px;display:grid;place-content:center;text-align:center;color:var(--text-secondary)}.quiet small{color:var(--text-muted);margin-top:4px}.backup p{max-width:520px}.backup button{margin-top:10px}.quick>div{display:grid;margin-top:8px}.quick button{display:flex;justify-content:space-between;align-items:center;text-align:left;background:transparent;border:0;border-top:1px solid var(--border-color);border-radius:0;padding:11px 2px}.quick button small{color:var(--text-muted)}
  @media(max-width:1050px){.dashboard-grid{grid-template-columns:1fr}.dashboard-grid article{border-right:0;border-bottom:1px solid var(--border-color)}.dashboard-grid article:last-child{border-bottom:0}.welcome{align-items:start;flex-direction:column}.hero{align-items:flex-start;flex-direction:column}.hero-actions{width:100%}.hero-actions button{flex:1}}
</style>
