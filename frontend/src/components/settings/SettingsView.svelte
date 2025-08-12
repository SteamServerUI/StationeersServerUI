<script>
  import { onMount } from 'svelte';
  import DetectionManager from './DetectionManager.svelte';
  import ConfigManager from './ConfigManager.svelte';
  // State management
  let activeSidebarTab = $state('SSUI Settings'); // Default to General tab in sidebar
  
  // Handle sidebar tab selection
  function selectSidebarTab(tab) {
    activeSidebarTab = tab;
  }
</script>

<div class="settings-container">
  <div class="settings-sidebar">
    <button 
      class="settings-nav {activeSidebarTab === 'Detection Manager' ? 'active' : ''}"
      onclick={() => selectSidebarTab('Detection Manager')}>Detection Manager</button>
    <button 
      class="settings-nav {activeSidebarTab === 'Config Manager' ? 'active' : ''}" 
      onclick={() => selectSidebarTab('Config Manager')}>Config Manager</button>
  </div>
  
  <div class="settings-content">
    {#if activeSidebarTab === 'Detection Manager'}
      <DetectionManager />
    {:else if activeSidebarTab === 'Config Manager'}
      <ConfigManager />
    {/if}
  </div>
</div>

<style>
  .settings-container {
    display: flex;
    gap: 2rem;
    height: 100%;
  }
  
  .settings-sidebar {
    width: 180px;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .settings-nav {
    text-align: left;
    padding: 0.75rem 1rem;
    background-color: transparent;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: all var(--transition-speed) ease;
  }
  
  .settings-nav:hover {
    background-color: var(--bg-hover);
  }
  
  .settings-nav.active {
    background-color: var(--bg-active);
    color: var(--accent-primary);
  }
  
  .settings-content {
    flex: 1;
    background-color: var(--bg-secondary);
    border-radius: 4px;
    padding: 1.5rem;
    overflow-y: auto;
  }
  
  @media (max-width: 768px) {
    .settings-container {
      flex-direction: column;
    }
    
    .settings-sidebar {
      width: 100%;
      flex-direction: row;
      overflow-x: auto;
      padding-bottom: 0.5rem;
    }
  }
</style>