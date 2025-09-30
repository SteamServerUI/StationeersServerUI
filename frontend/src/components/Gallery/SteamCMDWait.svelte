<script lang="js">
  // Funny messages to keep users entertained during SteamCMD operations
  const messages = $state([
    "SteamCMD's brewing coffee...",
    "Untangling virtual cables...",
    "Convincing the server to behave...",
    "SteamCMD's wrestling gremlins...",
    "Negotiating with bandwidth fairies...",
    "Defragmenting digital memories...",
    "Teaching SteamCMD new tricks...",
    "Soothing angry network packets...",
    "Doing server yoga...",
    "Persuading bits to cooperate..."
  ]);

  let currentMessageIndex = $state(0);
  let currentMessage = $state(messages[0]);

  // Rotate messages every 3 seconds
  $effect(() => {
    const interval = setInterval(() => {
      currentMessageIndex = (currentMessageIndex + 1) % messages.length;
      currentMessage = messages[currentMessageIndex];
    }, 3000);

    return () => clearInterval(interval);
  });
</script>

<div class="steamcmd-wait">
  <div class="spinner"></div>
  <p class="message">{currentMessage}</p>
  <p class="warning">Please don't poke the backend!</p>
</div>

<style>
  .steamcmd-wait {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 1rem;
    text-align: center;
    width: 100%;
  }

  .spinner {
    width: 32px;
    height: 32px;
    border: 4px solid var(--bg-secondary, #eeeeee);
    border-top: 4px solid var(--accent-primary, #2C6E49);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 1rem;
  }

  .message {
    color: var(--text-secondary, #666666);
    font-size: 0.9rem;
    margin-bottom: 0.5rem;
    font-style: italic;
  }

  .warning {
    color: var(--text-warning, #FF6B6B);
    font-size: 0.8rem;
    font-weight: bold;
    max-width: 200px;
    margin: 0;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
</style>