function switchTab(tabId) {
    // Hide all tabs
    document.querySelectorAll('.tab').forEach(tab => {
      tab.classList.remove('active');
    });
    
    // Show the selected tab
    document.getElementById(tabId).classList.add('active');
    
    // Update tab buttons
    document.querySelectorAll('.tab-button').forEach(button => {
      button.classList.remove('active');
    });
    
    // Activate the clicked button
    document.querySelector(`.tab-button[onclick*="${tabId}"]`).classList.add('active');
  }
  
  // Export for Svelte
  export { switchTab };