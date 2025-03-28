document.addEventListener('DOMContentLoaded', () => {
    // Planet creation functions (previous implementation remains the same)
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

    // Notification function (previous implementation remains the same)
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

    // Breakout Animation
    function createBreakoutAnimation() {
        const breakoutContainer = document.querySelector('.breakout-container');
        const barCount = 10;
        const colors = ['#00FFAB', '#1b1b2f', '#0a0a14'];

        for (let i = 0; i < barCount; i++) {
            const bar = document.createElement('div');
            bar.classList.add('breakout-bar');
            
            // Randomize bar properties
            bar.style.width = `${Math.random() * 100 + 50}%`;
            bar.style.height = `${Math.random() * 10 + 2}px`;
            bar.style.top = `${Math.random() * 100}%`;
            bar.style.left = `${-50 + Math.random() * 100}%`;
            bar.style.backgroundColor = colors[Math.floor(Math.random() * colors.length)];
            bar.style.transform = `rotate(${Math.random() * 360}deg)`;

            breakoutContainer.appendChild(bar);
        }

        // Trigger animation
        setTimeout(() => {
            breakoutContainer.classList.add('active');
            const bars = breakoutContainer.querySelectorAll('.breakout-bar');
            bars.forEach((bar, index) => {
                setTimeout(() => {
                    bar.classList.add('expand');
                    bar.style.transform = `rotate(${Math.random() * 720}deg) scale(${Math.random() * 2 + 1})`;
                }, index * 100);
            });
        }, 100);
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
            const response = await fetch('/', { 
                method: 'HEAD',
                cache: 'no-store'
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

                // Trigger breakout animation
                createBreakoutAnimation();
                
                // Redirect after animations
                setTimeout(() => {
                    if (preloadSuccess) {
                        hidePreloader();
                        window.location.href = '/';
                    } else {
                        showNotification('Preload failed. Redirecting anyway.', 'error');
                        window.location.href = '/';
                    }
                }, 2000);
            } else {
                const errorData = await response.json();
                const errorMessage = errorData.message || 'Login failed! Please check your credentials.';
                
                // Show error notification
                showNotification(errorMessage);
            }
        } catch (error) {
            console.error('Login error:', error);
            showNotification('Network error. Please try again.');
        }
    });
});