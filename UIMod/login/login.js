document.addEventListener('DOMContentLoaded', () => {
    // Planet creation functions
    const planetContainer = document.getElementById('planet-container');
    
    function createPlanet(container, size, orbitRadius, speed, color) {
        const orbit = document.createElement('div');
        orbit.classList.add('orbit');
        orbit.style.width = `${orbitRadius * 2}px`;
        orbit.style.height = `${orbitRadius * 2}px`;
        orbit.style.position = 'absolute';
        orbit.style.left = '50%';
        orbit.style.top = '50%';
        orbit.style.transform = 'translate(-50%, -50%)';
        
        // Add random delay to start animation at different points
        const randomDelay = -(Math.random() * speed); // Negative delay to offset start
        orbit.style.animation = `orbit ${speed}s linear infinite ${randomDelay}s`;
    
        const planet = document.createElement('div');
        planet.classList.add('planet');
        planet.style.width = `${size}px`;
        planet.style.height = `${size}px`;
        planet.style.position = 'absolute';
        planet.style.left = '0%';
        planet.style.top = '50%';
        planet.style.backgroundColor = color;
        planet.style.borderRadius = '50%';
        planet.style.boxShadow = `0 0 20px ${color}`;
        
        orbit.appendChild(planet);
        container.appendChild(orbit);
    }

    // Create planets with size, orbit radius, speed, and color
    createPlanet(planetContainer, 80, 650, 30, 'rgba(200, 100, 50, 0.7)');
    createPlanet(planetContainer, 50, 1000, 50, 'rgba(100, 200, 150, 0.5)');
    createPlanet(planetContainer, 30, 1250, 60, 'rgba(50, 150, 250, 0.6)');
    createPlanet(planetContainer, 70, 350, 30, 'rgba(200, 150, 200, 0.7)'); 

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
        const mode = document.getElementById('mode').value;

        let url;
        if (mode === 'setup') {
            url = '/api/v2/auth/setup/register';
        } else if (mode === 'changeuser') {
            url = '/api/v2/auth/adduser';
        } else {
            url = '/api/v2/auth/login';
        }

        try {
            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'HTTP2-Settings': 'AAEAAQABAAAAAQAAAAEAAAAAAAEAxABAAQAA'
                },
                body: JSON.stringify({ username, password })
            });

            const data = await response.json();
            if (response.ok) {
                if (mode === 'login') {
                    showNotification('Login Successful! Preparing launch...', 'success');
                    showPreloader();
                    const preloadSuccess = await preloadNextPage();
                    setTimeout(() => {
                        if (preloadSuccess) {
                            hidePreloader();
                            window.location.href = '/';
                        } else {
                            showNotification('Preload failed. Redirecting anyway.', 'error');
                            window.location.href = '/';
                        }
                    }, 600);
                } else if (mode === 'setup') {
                    showNotification('User registered successfully!', 'success');
                    document.getElementById('username').value = '';
                    document.getElementById('password').value = '';
                } else { // changeuser
                    const message = data.message || 'User updated successfully!';
                    showNotification(message, 'success');
                    document.getElementById('username').value = '';
                    document.getElementById('password').value = '';
                }
            } else {
                const errorMessage = data.error || 'Action failed! Please check your input.';
                showNotification(errorMessage);
            }
        } catch (error) {
            console.error('Form submission error:', error);
            showNotification('Error occurred. Please try again.');
        }
    });
    const finalizeBtn = document.getElementById('finalize-btn');
    if (finalizeBtn) {
        finalizeBtn.addEventListener('click', async () => {
            try {
                const response = await fetch('/api/v2/auth/setup/finalize', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'HTTP2-Settings': 'AAEAAQABAAAAAQAAAAEAAAAAAAEAxABAAQAA'
                    }
                });

                const data = await response.json();
                if (response.ok) {
                    showNotification(`${data.message}\n${data.restart_hint}`, 'success');
                    setTimeout(() => {
                        window.location.href = '/login';
                    }, 2000); // Give user time to read restart hint
                } else {
                    showNotification(data.error || 'Finalize failed!');
                }
            } catch (error) {
                console.error('Finalize error:', error);
                showNotification('Error finalizing setup. Please try again.');
            }
        });
    }
});