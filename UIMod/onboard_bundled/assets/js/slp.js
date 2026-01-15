function installSLP() {
    fetch('/api/v2/slp/install')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                alert('SLP installed successfully!');
                window.location.reload();
            } else {
                alert('Failed to install SLP: ' + data.error);
            }
        })
        .catch(error => {
            alert('Failed to install SLP: ' + error);
        });
}

function uninstallSLP() {
    fetch('/api/v2/slp/uninstall')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                alert('SLP uninstalled successfully!');
                window.location.reload();
            } else {
                alert('Failed to uninstall SLP: ' + data.error);
            }
        })
        .catch(error => {
            alert('Failed to uninstall SLP: ' + error);
        });
}