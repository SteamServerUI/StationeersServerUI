document.addEventListener('DOMContentLoaded', () => {
    // Planet creation functions
    const planetContainer = document.getElementById('planet-container');
    
    function createPlanet(size, x, y, speed, color) {
        const planet = document.createElement('div');
        planet.classList.add('planet');
        planet.style.width = `${size}px`;
        planet.style.height = `${size}px`;
        planet.style.position = 'absolute';
        planet.style.left = `${x}%`;
        planet.style.top = `${y}%`;
        planet.style.backgroundColor = color;
        planet.style.borderRadius = '50%';
        planet.style.animation = `orbit ${speed}s linear infinite`;
        planet.style.boxShadow = `0 0 20px ${color}`;
        
        planetContainer.appendChild(planet);
    }

    createPlanet(80, 10, 20, 30, 'rgba(200, 100, 50, 0.7)');
    createPlanet(50, 70, 60, 45, 'rgba(100, 200, 150, 0.5)');
    createPlanet(30, 50, 80, 20, 'rgba(50, 150, 250, 0.6)');

    // Notification function
    function showNotification(message, type = 'error') {
        const existingNotification = document.querySelector('.notification');
        if (existingNotification) {
            existingNotification.remove();
        }

        const notification = document.createElement('div');
        notification.classList.add('notification', type);
        notification.textContent = message;

        document.body.appendChild(notification);
        
        notification.offsetHeight;
        
        notification.classList.add('show');

        setTimeout(() => {
            notification.classList.remove('show');
            setTimeout(() => notification.remove(), 500);
        }, 3000);
    }


    // Preloader management
    function showPreloader() {
        const preloader = document.getElementById('preloader');
        preloader.classList.add('show');
    }

    function hidePreloader() {
        const preloader = document.getElementById('preloader');
        preloader.classList.remove('show');
    }

    // Preload next page
    async function preloadNextPage() {
        try {
            const response = await fetch('/static/favicon.ico', { // as there is no actual endpoint to check login status right now and this is a dummy preloader (secretly), we just stick with this for now.
                method: 'HEAD',
                cache: 'force-cache'
            });
            return response.ok;
        } catch (error) {
            console.error('Preload failed:', error);
            return false;
        }
    }

    // Login form submission
    document.getElementById('loginForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        try {
            const response = await fetch('/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'HTTP2-Settings': 'AAEAAQABAAAAAQAAAAEAAAAAAAEAxABAAQAA' // HTTP/2 SETTINGS 
                },
                body: JSON.stringify({ username, password })
            });

            if (response.ok) {
                // Show success notification
                showNotification('Login Successful! Preparing launch...', 'success');

                // Start preloading
                showPreloader();
                const preloadSuccess = await preloadNextPage();

                // Redirect after animations
                setTimeout(() => {
                    if (preloadSuccess) {
                        hidePreloader();
                        window.location.href = '/';
                    } else {
                        showNotification('Preload failed. Redirecting anyway.', 'error');
                        window.location.href = '/';
                    }
                }, 600); // fake loading time to avoid flickering, will be replaced with a proper preloader later
            } else {
                const errorData = await response.json();
                const errorMessage = errorData.message || 'Login failed! Please check your credentials.';
                
                // Show error notification
                showNotification(errorMessage);
            }
        } catch (error) {
            console.error('Login error:', error);   
            showNotification('Login error. Please try again.');
        }
    });
});