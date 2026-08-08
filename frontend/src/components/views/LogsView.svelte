<script>
  import { onMount } from 'svelte';
  import { apiSSE } from '../../services/api-v7';
  import { notify } from '../../services/activity';

  const sources={backend:['All backend','/api/v3/streams/logs/backend'],info:['Info','/api/v3/streams/logs/info'],warn:['Warnings','/api/v3/streams/logs/warn'],error:['Errors','/api/v3/streams/logs/error'],debug:['Debug','/api/v3/streams/logs/debug'],game:['Game','/api/v3/streams/console']};
  let source=$state('info'), query=$state(''), logs=$state([]), paused=$state(false), connected=$state(false), stream;
  let visible=$derived(logs.filter(entry=>(entry.source===source||source==='backend')&&(!query||entry.text.toLowerCase().includes(query.toLowerCase()))));

  onMount(()=>{open();return()=>stream?.close();});
  function open(){
    stream?.close();connected=false;
    stream=apiSSE(sources[source][1],data=>{connected=true;if(!paused)logs=[...logs,{id:crypto.randomUUID(),source,at:new Date(),text:String(data)}].slice(-2500);},()=>connected=false);
  }
  function choose(id){source=id;open();}
  function exportLogs(){const blob=new Blob([visible.map(item=>`${item.at.toISOString()} [${item.source.toUpperCase()}] ${item.text}`).join('\n')],{type:'text/plain'});const url=URL.createObjectURL(blob);const a=document.createElement('a');a.href=url;a.download=`ssui-${source}-logs.txt`;a.click();URL.revokeObjectURL(url);notify('Logs exported','success');}
</script>

<section class="log-workspace surface">
  <aside><span class="eyebrow">Sources</span>{#each Object.entries(sources) as [id,value]}<button class:active={source===id} onclick={()=>choose(id)}><span>{value[0]}</span><small>{logs.filter(item=>item.source===id).length}</small></button>{/each}</aside>
  <div class="log-main">
    <header><div><i class:online={connected}></i><strong>{sources[source][0]}</strong><span>{connected?'Streaming':'Reconnecting…'}</span></div><div><input bind:value={query} placeholder="Search this stream"><button class:active={paused} onclick={()=>paused=!paused}>{paused?'Resume':'Pause'}</button><button onclick={exportLogs}>Export</button><button onclick={()=>logs=[]}>Clear</button></div></header>
    <div class="log-table">
      {#if !visible.length}<div class="empty">No matching entries yet.</div>{/if}
      {#each visible as entry (entry.id)}<div class="log-row"><time>{entry.at.toLocaleTimeString()}</time><span class="level">{entry.source}</span><p class:error={/error|fatal/i.test(entry.text)} class:warning={/warn/i.test(entry.text)}>{entry.text}</p></div>{/each}
    </div>
  </div>
</section>

<style>
  .log-workspace{height:min(760px,calc(100vh - 190px));min-height:520px;display:grid;grid-template-columns:190px 1fr;overflow:hidden}aside{padding:18px 10px;border-right:1px solid var(--border-color)}aside .eyebrow{display:block;padding:0 10px 12px}aside button{width:100%;display:flex;justify-content:space-between;background:transparent;border:0;text-align:left;color:var(--text-secondary)}aside button.active{background:var(--bg-active);color:var(--text-accent)}aside small{color:var(--text-muted)}.log-main{display:grid;grid-template-rows:auto 1fr;min-width:0}.log-main header{display:flex;align-items:center;justify-content:space-between;gap:15px;padding:12px 14px;border-bottom:1px solid var(--border-color)}.log-main header>div{display:flex;align-items:center;gap:9px}.log-main header i{width:7px;height:7px;border-radius:50%;background:var(--text-danger)}.log-main header i.online{background:var(--text-success);box-shadow:0 0 9px var(--text-success)}.log-main header span{color:var(--text-muted);font-size:.73rem}.log-main header input{width:210px}.log-main header button.active{color:var(--text-warning)}.log-table{overflow:auto;background:#080b11;padding:7px 0;font:12.5px/1.55 ui-monospace,SFMono-Regular,Consolas,monospace}.log-row{display:grid;grid-template-columns:85px 68px 1fr;gap:10px;padding:5px 15px;border-bottom:1px solid rgba(255,255,255,.025)}.log-row:hover{background:rgba(255,255,255,.025)}.log-row time{color:#536076}.level{color:#7393ba;text-transform:uppercase;font-size:.67rem;padding-top:2px}.log-row p{margin:0;color:#cad5e7;white-space:pre-wrap;word-break:break-word}.log-row p.error{color:#ff8798}.log-row p.warning{color:#ffd080}.empty{height:100%;display:grid;place-content:center;color:#59677d}@media(max-width:980px){.log-workspace{grid-template-columns:1fr;height:calc(100vh - 210px)}aside{display:flex;overflow:auto;border-right:0;border-bottom:1px solid var(--border-color);padding:8px}aside .eyebrow{display:none}aside button{width:auto;white-space:nowrap}.log-main header{align-items:flex-start;flex-direction:column}.log-main header>div:last-child{width:100%;flex-wrap:wrap}.log-main header input{flex:1}}
</style>
