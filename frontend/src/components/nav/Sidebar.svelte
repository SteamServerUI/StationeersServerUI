<script>
  /**
   * @typedef {Object} Props
   * @property {any} [views]
   * @property {string} [activeView]
   * @property {any} setActiveView
   */

  /** @type {Props} */
  let { views = [], activeView = 'dashboard', setActiveView } = $props();
    
    let isHovered = $state(false);
    let isPinned = $state(false);
    
    function togglePin() {
      isPinned = !isPinned;
    }
    
    function handleMouseEnter() {
      isHovered = true;
    }
    
    function handleMouseLeave() {
      isHovered = false;
    }
    
    let isExpanded = $derived(isHovered || isPinned);
  </script>
  
  <aside 
    class="sidebar" 
    class:expanded={isExpanded}
    onmouseenter={handleMouseEnter}
    onmouseleave={handleMouseLeave}
  >
    <div class="pin-container">
      <button class="pin-button" class:active={isPinned} onclick={togglePin}>
        {#if isPinned}
          📌
        {:else}
          📍
        {/if}
      </button>
    </div>
    
    <nav class="sidebar-nav">
      {#each views as view}
        <button 
          class="sidebar-button" 
          class:active={activeView === view.id}
          onclick={() => setActiveView(view.id)}
        >
          <span class="sidebar-icon">
            {#if view.icon === 'grid'}
              📊
            {:else if view.icon === 'settings'}
              ⚙️
            {:else if view.icon === 'file-text'}
              📝
            {:else if view.icon === 'terminal'}
              >
            {:else if view.icon === 'globe'}
              🌐
            {:else if view.icon === 'archive'}
              📦
            {/if}
          </span>
          <span class="sidebar-text">{view.name}</span>
        </button>
      {/each}
    </nav>
  </aside>
  
  <style>
    .sidebar {
      width: var(--sidebar-collapsed-width);
      height: 100%;
      background-color: var(--bg-secondary);
      border-right: 1px solid var(--border-color);
      display: flex;
      flex-direction: column;
      overflow: hidden;
      transition: width var(--transition-speed) ease;
      z-index: 5;
    }
    
    .sidebar.expanded {
      width: var(--sidebar-width);
    }
    
    .pin-container {
      display: flex;
      justify-content: flex-end;
      padding: 0.5rem;
    }
    
    .pin-button {
      padding: 0.25rem;
      width: 2rem;
      height: 2rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: transparent;
      border: none;
      cursor: pointer;
      opacity: 0.6;
      transition: opacity var(--transition-speed) ease;
    }
    
    .pin-button:hover, .pin-button.active {
      opacity: 1;
    }
    
    .sidebar-nav {
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
      padding: 0.5rem;
      overflow-y: auto;
    }
    
    .sidebar-button {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 0.75rem;
      border: none;
      background-color: transparent;
      border-radius: 4px;
      cursor: pointer;
      transition: all var(--transition-speed) ease;
      justify-content: flex-start;
      width: 100%;
      text-align: left;
    }
    
    .sidebar-button:hover {
      background-color: var(--bg-hover);
    }
    
    .sidebar-button.active {
      background-color: var(--bg-active);
      color: var(--accent-primary);
    }
    
    .sidebar-icon {
      font-size: 1.2rem;
      min-width: 1.5rem;
      display: flex;
      justify-content: center;
    }
    
    .sidebar-text {
      white-space: nowrap;
      opacity: 0;
      transition: opacity var(--transition-speed) ease;
    }
    
    .expanded .sidebar-text {
      opacity: 1;
    }
  </style>