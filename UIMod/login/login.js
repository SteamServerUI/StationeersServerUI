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

    // Notification function (unchanged)
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

    // Preloader functions (unchanged)
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

        if (step === "0") { // Welcome step
            window.location.href = '/setup?step=1';
            return;
        }

        let url, body;
        if (step === "5") { // User setup
            url = '/api/v2/auth/setup/register';
            body = JSON.stringify({
                username: document.getElementById('input-field').value,
                password: document.getElementById('password').value
            });
        } else if (step >= "1" && step <= "4") { // Config steps
            url = '/api/v2/saveconfig';
            const key = ["ServerName", "SaveInfo", "ServerMaxPlayers", "ServerPassword"][step - 1];
            body = JSON.stringify({
                [key]: document.getElementById('input-field').value
            });
        } else { // Login or changeuser
            url = mode === 'changeuser' ? '/api/v2/auth/adduser' : '/auth/login';
            body = JSON.stringify({
                username: document.getElementById('input-field').value,
                password: document.getElementById('password').value
            });
        }

        try {
            const response = await fetch(url, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: body
            });
            const data = await response.json();
            if (response.ok) {
                if (step === "5") {
                    showNotification('Admin account saved!', 'success');
                    showNextStepButton();
                } else if (step >= "1" && step <= "4") {
                    showNotification('Config saved!', 'success');
                    showNextStepButton();
                } else if (mode === 'login') {
                    showNotification('Login Successful!', 'success');
                    showPreloader();
                    const preloadSuccess = await preloadNextPage();
                    setTimeout(() => {
                        hidePreloader();
                        window.location.href = '/';
                    }, 600);
                } else { // changeuser
                    showNotification(data.message || 'User updated!', 'success');
                    form.reset();
                }
            } else {
                showNotification(data.error || 'Action failed!', 'error');
            }
        } catch (error) {
            console.error('Error:', error);
            showNotification('Something went wrong!', 'error');
        }
    });

    // Finalize button
    const finalizeBtn = document.getElementById('finalize-btn');
    if (finalizeBtn) {
        finalizeBtn.addEventListener('click', async () => {
            try {
                const finalizeResponse = await fetch('/api/v2/auth/setup/finalize', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' }
                });
                const data = await finalizeResponse.json();
                if (finalizeResponse.ok) {
                    showNotification(`${data.message}\n${data.restart_hint}`, 'success');
                    setTimeout(() => window.location.href = '/login', 2000);
                } else {
                    showNotification(data.error || 'Finalize failed!', 'error');
                }
            } catch (error) {
                console.error('Finalize error:', error);
                showNotification('Error finalizing setup!', 'error');
            }
        });
    }

    // Skip button
    const skipBtn = document.getElementById('skip-btn');
    if (skipBtn) {
        skipBtn.addEventListener('click', () => {
            const step = document.getElementById('step').value;
            const nextStep = parseInt(step) + 1;
            if (nextStep > 6) {
                showNotification('Setup skipped!', 'success');
                setTimeout(() => window.location.href = '/', 1000);
            } else {
                window.location.href = `/setup?step=${nextStep}`;
            }
        });
    }

    // Show Next Step button
    function showNextStepButton() {
        const nextBtn = document.getElementById('next-btn');
        if (nextBtn) {
            nextBtn.style.display = 'block';
            nextBtn.onclick = () => {
                const step = parseInt(document.getElementById('step').value);
                window.location.href = `/setup?step=${step + 1}`;
            };
        }
    }
});