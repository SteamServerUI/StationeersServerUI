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
            window.location.href = '/'; // Redirect to root on success
        } else {
            const errorData = await response.json();
            const errorMessage = errorData.message || 'Login failed!';
            
            // Create and show error notification
            const errorNotification = document.createElement('div');
            errorNotification.classList.add('notification', 'error');
            errorNotification.textContent = errorMessage;
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