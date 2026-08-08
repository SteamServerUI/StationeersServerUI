<script>
  import { onMount } from 'svelte';
  import { apiSSE, apiFetch, hasPermission } from '../../services/api-v7';
  import { notify } from '../../services/activity';

  let tab=$state('console'), consoleLines=$state([]), events=$state([]), query=$state(''), paused=$state(false), autoscroll=$state(true), connected=$state(false), command=$state(''), sending=$state(false), viewport;
  let consoleStream,eventStream;
  let lines=$derived((tab==='console'?consoleLines:events).filter(line=>!query||line.text.toLowerCase().includes(query.toLowerCase())));

  onMount(()=>{
    consoleStream=connect('/api/v3/streams/console',line=>consoleLines=[...consoleLines,line].slice(-1500));
    eventStream=connect('/api/v3/streams/events',line=>events=[...events,line].slice(-1000));
    return()=>{consoleStream?.close();eventStream?.close();};
  });
  function connect(path,append){
    return apiSSE(path,data=>{connected=true;if(!paused){append({id:crypto.randomUUID(),at:new Date(),text:String(data)});queueMicrotask(scroll);}},()=>connected=false);
  }
  function scroll(){if(viewport&&autoscroll)viewport.scrollTop=viewport.scrollHeight;}
  async function send(){
    const value=command.trim();if(!value)return;sending=true;
    try{const response=await apiFetch('/api/v3/SSCM/run',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({command:value})});const body=await response.json().catch(()=>({}));if(!response.ok)throw new Error(body.message||'Command rejected');command='';notify('Command sent','success',value);}
    catch(error){notify('Command failed','error',error.message);}finally{sending=false;}
  }
  function copy(){navigator.clipboard.writeText(lines.map(line=>`[${line.at.toLocaleTimeString()}] ${line.text}`).join('\n'));notify('Visible output copied','success');}
</script>

<section class="terminal-shell surface">
  <header class="terminal-toolbar">
    <div class="tabs"><button class:active={tab==='console'} onclick={()=>tab='console'}>Console <span>{consoleLines.length}</span></button><button class:active={tab==='events'} onclick={()=>tab='events'}>Detections <span>{events.length}</span></button></div>
    <div class="stream-state"><i class:online={connected}></i>{connected?'Live':'Reconnecting'}</div>
    <div class="tools"><input bind:value={query} placeholder="Filter output"><button class:active={paused} onclick={()=>paused=!paused}>{paused?'Resume':'Pause'}</button><button onclick={copy}>Copy</button></div>
  </header>
  <div class="terminal" bind:this={viewport} onscroll={()=>{if(viewport)autoscroll=viewport.scrollTop+viewport.clientHeight>=viewport.scrollHeight-24}}>
    {#if !lines.length}<div class="terminal-empty">{paused?'Output is paused.':'Waiting for output…'}</div>{/if}
    {#each lines as line (line.id)}<div class="line"><time>{line.at.toLocaleTimeString()}</time><span class:error={/error|fatal/i.test(line.text)} class:warning={/warn/i.test(line.text)}>{line.text}</span></div>{/each}
  </div>
  <footer>
    {#if hasPermission('sscm.run')}<form onsubmit={event=>{event.preventDefault();send();}}><span>›</span><input bind:value={command} placeholder="Send a Unity console command…" disabled={sending}><button disabled={sending||!command.trim()}>Send</button></form>{/if}
    <label><input type="checkbox" bind:checked={autoscroll}> Follow output</label>
  </footer>
</section>

<style>
  .terminal-shell{height:min(760px,calc(100vh - 190px));min-height:520px;display:grid;grid-template-rows:auto 1fr auto;overflow:hidden}.terminal-toolbar{display:grid;grid-template-columns:auto 1fr auto;align-items:center;gap:18px;padding:12px;border-bottom:1px solid var(--border-color)}.tabs{display:flex;gap:4px}.tabs button{background:transparent;border:0}.tabs button span{color:var(--text-muted);font-size:.7rem;margin-left:5px}.tabs button.active{background:var(--bg-active);color:var(--text-accent)}.stream-state{justify-self:center;display:flex;align-items:center;gap:7px;color:var(--text-secondary);font-size:.76rem}.stream-state i{width:7px;height:7px;border-radius:50%;background:var(--text-danger)}.stream-state i.online{background:var(--text-success);box-shadow:0 0 10px var(--text-success)}.tools{display:flex;gap:6px}.tools input{width:190px}.tools button.active{color:var(--text-warning)}.terminal{overflow:auto;background:#070a10;padding:14px 16px;font:12.5px/1.65 ui-monospace,SFMono-Regular,Consolas,monospace}.line{display:grid;grid-template-columns:82px 1fr;gap:12px}.line time{color:#536076;user-select:none}.line span{white-space:pre-wrap;word-break:break-word;color:#c9d5e8}.line span.error{color:#ff8798}.line span.warning{color:#ffd080}.terminal-empty{height:100%;display:grid;place-content:center;color:#59677d}footer{display:flex;align-items:center;justify-content:space-between;gap:20px;padding:11px 14px;border-top:1px solid var(--border-color)}footer form{display:flex;align-items:center;gap:8px;flex:1}footer form>span{color:var(--accent-primary);font-size:1.4rem}footer form input{flex:1;background:transparent;border:0}footer label{display:flex;align-items:center;gap:7px;color:var(--text-secondary);font-size:.75rem}@media(max-width:980px){.terminal-toolbar{grid-template-columns:1fr}.stream-state{display:none}.tools{flex-wrap:wrap}.tools input{flex:1}.terminal-shell{height:calc(100vh - 210px)}}
</style>
