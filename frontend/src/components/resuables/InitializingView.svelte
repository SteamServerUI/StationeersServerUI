<script>
  export let serverStatus = 'checking';
  export let errorMessage = null;
  export let statusDescOffline = null;

  function isRunningInElectron() {
  const userAgentCheck = typeof navigator === 'object' && 
    typeof navigator.userAgent === 'string' && 
    navigator.userAgent.indexOf('Electron') >= 0;
  return userAgentCheck;
  }

  statusDescOffline = isRunningInElectron() 
  ? "Poking Backend"
  : "Backend Offline";
</script>

<div class="initializing">
  <div class="content-container">
    <div class="logo-container">
      <div class="spinner-container">
        <div class="loading-spinner"></div>
      </div>
    </div>
    
    <div class="status-container">
      <h2 class="status-title">
        {#if serverStatus === 'checking'}
          Connecting to Backend
        {:else if serverStatus === 'online'}
          Initializing Backend
        {:else if serverStatus === 'offline'}
          {statusDescOffline}
        {:else if serverStatus === 'error'}
          Backend error
        {:else if serverStatus === 'cert-error'}
          Backend Certificate or https error
        {:else if serverStatus === 'unreachable'}
          Server not found
        {/if}
      </h2>
      
      <div class="status-indicator">
        <span class="dot dot1"></span>
        <span class="dot dot2"></span>
        <span class="dot dot3"></span>
      </div>
    </div>
    
    {#if errorMessage}
      <div class="error-message">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" y1="8" x2="12" y2="12"></line>
          <line x1="12" y1="16" x2="12.01" y2="16"></line>
        </svg>
        <p>{errorMessage}</p>
      </div>
    {/if}
  </div>
</div>

<style>
  .initializing {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background-color: var(--bg-primary);
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  }
  
  .content-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 3rem;
    border-radius: 12px;
    background-color: var(--bg-secondary);
    box-shadow: var(--shadow-medium);
    max-width: 90%;
    width: 420px;
    transition: all var(--transition-speed) ease;
  }
  
  .logo-container {
    margin-bottom: 2rem;
    position: relative;
  }
  
  .spinner-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100px;
    width: 100px;
  }
  
  .loading-spinner {
    position: relative;
    width: 64px;
    height: 64px;
    border-radius: 50%;
    border: 3px solid transparent;
    border-top-color: var(--accent-primary);
    border-right-color: var(--accent-secondary);
    border-bottom-color: var(--accent-tertiary);
    animation: spin 1.2s linear infinite;
  }
  
  .loading-spinner:before {
    content: "";
    position: absolute;
    top: 5px;
    left: 5px;
    right: 5px;
    bottom: 5px;
    border-radius: 50%;
    border: 3px solid transparent;
    border-top-color: var(--accent-tertiary);
    animation: spin 1.8s linear infinite;
  }
  
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
  
  .status-container {
    text-align: center;
    margin-bottom: 1.5rem;
    width: 100%;
  }
  
  .status-title {
    font-size: 1.25rem;
    font-weight: 500;
    color: var(--text-primary);
    margin: 0 0 1rem 0;
  }
  
  .status-indicator {
    display: flex;
    justify-content: center;
    gap: 6px;
    margin-top: 0.5rem;
  }
  
  .dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background-color: var(--accent-primary);
    opacity: 0.4;
  }
  
  .dot1 {
    animation: pulse 1.5s infinite ease-in-out;
  }
  
  .dot2 {
    animation: pulse 1.5s infinite ease-in-out 0.5s;
  }
  
  .dot3 {
    animation: pulse 1.5s infinite ease-in-out 1s;
  }
  
  @keyframes pulse {
    0%, 100% { opacity: 0.4; transform: scale(1); }
    50% { opacity: 1; transform: scale(1.2); }
  }
  
  .error-message {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-top: 1rem;
    padding: 1rem;
    background-color: rgba(255, 82, 82, 0.1);
    color: var(--text-warning);
    border-left: 4px solid var(--text-warning);
    border-radius: 4px;
    width: 100%;
    font-size: 0.95rem;
    box-shadow: var(--shadow-light);
  }
  
  .error-message svg {
    min-width: 20px;
    color: var(--text-warning);
  }
  
  .error-message p {
    margin: 0;
    line-height: 1.4;
  }

  @media (max-width: 480px) {
    .content-container {
      padding: 2rem;
      width: 90%;
    }
    
    .spinner-container {
      height: 80px;
      width: 80px;
    }
    
    .loading-spinner {
      width: 50px;
      height: 50px;
    }
  }
</style>