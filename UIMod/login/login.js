document.addEventListener('DOMContentLoaded', () => {
    // Create planet background
    const planetContainer = document.getElementById('planet-container');
    
    // Create multiple planets
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

    // Create some planets with different characteristics
    createPlanet(80, 10, 20, 30, 'rgba(200, 100, 50, 0.7)');
    createPlanet(50, 70, 60, 45, 'rgba(100, 200, 150, 0.5)');
    createPlanet(30, 50, 80, 20, 'rgba(50, 150, 250, 0.6)');

    // Add custom CSS for planet orbits
    const styleSheet = document.createElement('style');
    styleSheet.textContent = `
        @keyframes orbit {
            0% { transform: rotate(0deg) translateX(150px) rotate(0deg); }
            100% { transform: rotate(360deg) translateX(150px) rotate(-360deg); }
        }
    `;
    document.head.appendChild(styleSheet);

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
                // Add a space-launch style transition
                document.body.style.transition = 'transform 1s ease-in';
                document.body.style.transform = 'perspective(1000px) rotateX(90deg) scale(0.5)';
                
                // Delay redirect to allow animation
                setTimeout(() => {
                    window.location.href = '/';
                }, 800);
            } else {
                const errorData = await response.json();
                const errorMessage = errorData.message || 'Login failed!';
                
                // Create and show error notification
                const errorNotification = document.createElement('div');
                errorNotification.classList.add('notification', 'error');
                errorNotification.textContent = errorMessage;
                errorNotification.style.cssText = `
                    position: fixed;
                    top: 20px;
                    left: 50%;
                    transform: translateX(-50%);
                    background-color: rgba(255, 0, 0, 0.7);
                    color: white;
                    padding: 15px;
                    border-radius: 5px;
                    z-index: 1000;
                `;
                document.body.appendChild(errorNotification);

                // Remove notification after 3 seconds
                setTimeout(() => {
                    errorNotification.remove();
                }, 3000);
            }
        } catch (error) {
            console.error('Login error:', error);
            alert('Network error. Please try again.');
        }
    });
});