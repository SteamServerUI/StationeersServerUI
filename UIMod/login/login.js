document.addEventListener('DOMContentLoaded', () => {
    // Planet creation (unchanged)
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
        const randomDelay = -(Math.random() * speed);
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
    createPlanet(planetContainer, 80, 650, 30, 'rgba(200, 100, 50, 0.7)');
    createPlanet(planetContainer, 50, 1000, 50, 'rgba(100, 200, 150, 0.5)');
    createPlanet(planetContainer, 30, 1250, 60, 'rgba(50, 150, 250, 0.6)');
    createPlanet(planetContainer, 70, 350, 30, 'rgba(200, 150, 200, 0.7)');

    // Notification function
    function showNotification(message, type = 'error') {
        const existingNotification = document.querySelector('.notification');
        if (existingNotification) existingNotification.remove();
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

    // Preloader functions
    function showPreloader() {
        document.getElementById('preloader').classList.add('show');
    }
    function hidePreloader() {
        document.getElementById('preloader').classList.remove('show');
    }
    async function preloadNextPage() {
        try {
            const response = await fetch('/static/favicon.ico', { method: 'HEAD', cache: 'force-cache' });
            return response.ok;
        } catch (error) {
            console.error('Preload failed:', error);
            return false;
        }
    }

    // Form submission
    const form = document.getElementById('loginForm');
    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const step = document.getElementById('step').value;
        const mode = document.getElementById('mode').value;

        if (step === "welcome") {
            window.location.href = '/setup?step=1';
            return;
        }

        if (step === "finalize") {
            // Return to first step of setup
            window.location.href = '/setup?step=1';
            return;
        }

        let url, body, nextStep;
        
        // Handle setup steps
        if (step >= "1" && step <= "4") {
            url = '/api/v2/saveconfig';
            const key = ["ServerName", "SaveInfo", "ServerMaxPlayers", "ServerPassword"][parseInt(step) - 1];
            body = JSON.stringify({
                [key]: document.getElementById('input-field').value
            });
            nextStep = parseInt(step) + 1;
        } else if (step === "5") { // User setup
            url = '/api/v2/auth/setup/register';
            body = JSON.stringify({
                username: document.getElementById('input-field').value,
                password: document.getElementById('password').value
            });
            nextStep = "finalize";
        } else { // Login or changeuser
            url = mode === 'changeuser' ? '/api/v2/auth/adduser' : '/auth/login';
            body = JSON.stringify({
                username: document.getElementById('input-field').value,
                password: document.getElementById('password').value
            });
        }

        try {
            showPreloader();
            const response = await fetch(url, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: body
            });
            const data = await response.json();
            
            if (response.ok) {
                if (step >= "1" && step <= "5") {
                    hidePreloader();
                    showNotification(step === "5" ? 'Admin account saved!' : 'Config saved!', 'success');
                    setTimeout(() => {
                        window.location.href = `/setup?step=${nextStep}`;
                    }, 800);
                } else if (mode === 'login') {
                    showNotification('Login Successful!', 'success');
                    await preloadNextPage();
                    setTimeout(() => {
                        hidePreloader();
                        window.location.href = '/';
                    }, 600);
                } else { // changeuser
                    hidePreloader();
                    showNotification(data.message || 'User updated!', 'success');
                    form.reset();
                }
            } else {
                hidePreloader();
                showNotification(data.error || 'Action failed!', 'error');
            }
        } catch (error) {
            hidePreloader();
            console.error('Error:', error);
            showNotification('Something went wrong!', 'error');
        }
    });

    // Skip button
    const skipBtn = document.getElementById('skip-btn');
    if (skipBtn) {
        skipBtn.addEventListener('click', () => {
            const step = document.getElementById('step').value;
            
            if (step === "welcome") {
                window.location.href = '/';
                return;
            }
            
            let nextStep;
            if (step === "finalize") {
                // Go to login page when skipping from finalize
                showNotification('Setup completed, Auth disabled!', 'success');
                setTimeout(() => window.location.href = '/', 1000);
                return;
            } else {
                nextStep = parseInt(step) + 1;
                if (nextStep > 5) {
                    nextStep = "finalize";
                }
            }
            
            window.location.href = `/setup?step=${nextStep}`;
        });
    }

    // Finalize button - now appears on the finalize page
    document.addEventListener('click', async (e) => {
        if (e.target && e.target.id === 'finalize-btn') {
            try {
                showPreloader();
                const finalizeResponse = await fetch('/api/v2/auth/setup/finalize', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' }
                });
                const data = await finalizeResponse.json();
                
                if (finalizeResponse.ok) {
                    hidePreloader();
                    showNotification(`${data.message}\n${data.restart_hint}`, 'success');
                    setTimeout(() => window.location.href = '/login', 2000);
                } else {
                    hidePreloader();
                    showNotification(data.error || 'Finalize failed!', 'error');
                }
            } catch (error) {
                hidePreloader();
                console.error('Finalize error:', error);
                showNotification('Error finalizing setup!', 'error');
            }
        }
    });
});