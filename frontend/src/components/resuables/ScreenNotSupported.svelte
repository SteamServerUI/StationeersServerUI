<script>
  import { onMount } from 'svelte';
  import { fly, fade, scale } from 'svelte/transition';
  import { elasticOut, cubicOut } from 'svelte/easing';
  import { spring } from 'svelte/motion';

  // Props to allow parent to handle the override
  let { onContinueAnyway = null } = $props();

  // Animation properties
  let visible = $state(false);
  let deviceVisible = $state(false);
  
  // Spring animation for the device illustration
  const coords = spring({ x: 0, y: 0 }, {
    stiffness: 0.1,
    damping: 0.25
  });
  
  // Track mouse position to create subtle hover effects
  function handleMouseMove(event) {
    const { clientX, clientY } = event;
    const rect = event.currentTarget.getBoundingClientRect();
    const x = (clientX - rect.left - rect.width / 2) / 15;
    const y = (clientY - rect.top - rect.height / 2) / 15;
    
    coords.set({ x, y });
  }
  
  function handleContinue() {
    if (onContinueAnyway) {
      onContinueAnyway();
    }
  }
  
  onMount(() => {
    // Stagger the animations
    visible = true;
    setTimeout(() => {
      deviceVisible = true;
    }, 500);
  });
</script>

<div 
    class="container"
    role="presentation"
    onmousemove={handleMouseMove} 
    onmouseleave={() => coords.set({ x: 0, y: 0 })}
>
    {#if visible}
    <div class="message-container" transition:fade={{ duration: 800 }}>
      <div class="message-box" in:fly={{ y: 30, duration: 800, delay: 300 }}>
        <div class="message-content">
          <div class="icon" in:scale={{ duration: 600, delay: 500, easing: elasticOut }}>
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 8V12M12 16H12.01M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12Z" 
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <h2 in:fly={{ y: 20, duration: 600, delay: 600 }}>Screen too Small</h2>
          <p in:fly={{ y: 20, duration: 600, delay: 800 }}>
            Sorry, this display size is not recommended currently. If you continue, the app may not get displayed properly. Please use a device with a minimum width of 1024px and height of 600px, such as an iPad, or decrease the window scalnig. If you're using a phone, consider rotating it to landscape mode.
          </p>
          <button 
            class="continue-button" 
            onclick={handleContinue}
            in:fly={{ y: 20, duration: 600, delay: 1000 }}
          >
            I don't care, show me the page anyway
          </button>
        </div>
      </div>
    </div>
    {/if}
</div>

<style>
  .container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100vh;
    width: 100vw;
    background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
    color: #fff;
    font-family: 'Inter', system-ui, -apple-system, sans-serif;
    overflow: hidden;
    position: relative;
  }
  
  .message-container {
    position: relative;
    z-index: 10;
  }
  
  .message-box {
    background: rgba(30, 41, 59, 0.8);
    backdrop-filter: blur(12px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    box-shadow: 0 20px 80px -10px rgba(0, 0, 0, 0.5);
    border-radius: 24px;
    padding: 40px;
    width: 90vw;
    max-width: 500px;
    text-align: center;
  }
  
  .message-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
  }
  
  .icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 64px;
    height: 64px;
    background: rgba(255, 79, 79, 0.15);
    color: #ff4f4f;
    border-radius: 50%;
    margin-bottom: 8px;
  }
  
  .icon svg {
    width: 32px;
    height: 32px;
  }
  
  h2 {
    font-size: 28px;
    font-weight: 700;
    margin: 0;
    background: linear-gradient(90deg, #fff, #a5b4fc);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }
  
  p {
    font-size: 16px;
    line-height: 1.6;
    color: rgba(255, 255, 255, 0.7);
    margin: 0 0 20px;
  }
  
  .continue-button {
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    color: rgba(255, 255, 255, 0.8);
    padding: 12px 24px;
    border-radius: 12px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    backdrop-filter: blur(8px);
  }
  
  .continue-button:hover {
    background: rgba(255, 255, 255, 0.15);
    border-color: rgba(255, 255, 255, 0.3);
    color: #fff;
    transform: translateY(-1px);
  }
  
  .continue-button:active {
    transform: translateY(0);
  }
  
  @media (max-width: 600px) {
    .message-box {
      padding: 30px 20px;
    }
    
    h2 {
      font-size: 24px;
    }
    
    p {
      font-size: 14px;
    }
    
    .continue-button {
      font-size: 13px;
      padding: 10px 20px;
    }
  }
</style>