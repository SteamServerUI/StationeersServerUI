function showNotification(message, type = 'info') {
    const notification = document.getElementById('notification');
    notification.textContent = message;
    notification.className = 'notification ' + type;
    notification.style.display = 'block';
    
    // Auto-hide after 5 seconds for success, keep longer for errors
    const duration = type === 'success' ? 5000 : 7000;
    setTimeout(() => {
        notification.style.display = 'none';
    }, duration);
}

function setButtonLoading(buttonId, isLoading) {
    const button = document.getElementById(buttonId);
    if (!button) return;
    
    if (isLoading) {
        button.disabled = true;
        button.dataset.originalText = button.textContent;
        button.textContent = '⏳ Please wait...';
        button.classList.add('loading');
    } else {
        button.disabled = false;
        button.textContent = button.dataset.originalText || button.textContent;
        button.classList.remove('loading');
    }
}

function installSLP() {
    setButtonLoading('installSLPBtn', true);
    showPopup('info', 'Installing Stationeers Launch Pad...');
    
    fetch('/api/v2/slp/install')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showPopup('success', 'Stationeers Launch Pad installed successfully! The page will refresh automatically.');
                setButtonLoading('installSLPBtn', false);
                // Reload after 3 seconds to show the success message
                setTimeout(() => window.location.reload(), 3000);
            } else {
                showPopup('error', 'Failed to install SLP:\n\n' + (data.error || 'Unknown error'));
                setButtonLoading('installSLPBtn', false);
            }
        })
        .catch(error => {
            showPopup('error', 'Failed to install SLP:\n\n' + (error.message || 'Network error'));
            setButtonLoading('installSLPBtn', false);
        });
}

function uninstallSLP() {
    if (!confirm('Are you sure you want to uninstall SLP? This will DELETE all mods too (you can always reinstall later).')) {
        return;
    }
    setButtonLoading('uninstallSLPBtn', true);
    showPopup('info', 'Uninstalling Stationeers Launch Pad...');

    fetch('/api/v2/slp/uninstall')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showPopup('success', 'Stationeers Launch Pad uninstalled successfully! The page will refresh automatically.');
                setButtonLoading('uninstallSLPBtn', false);
                setTimeout(() => window.location.reload(), 3000);
            } else {
                showPopup('error', 'Failed to uninstall SLP:\n\n' + (data.error || 'Unknown error'));
                setButtonLoading('uninstallSLPBtn', false);
            }
        })
        .catch(error => {
            showPopup('error', 'Failed to uninstall SLP:\n\n' + (error.message || 'Network error'));
            setButtonLoading('uninstallSLPBtn', false);
        });
}

let selectedModFile = null;

function handleModPackageSelection(files) {
    if (!files || files.length === 0) {
        return;
    }

    const file = files[0];
    
    // Validate file is a zip
    if (!file.name.endsWith('.zip')) {
        showPopup('error', 'Invalid file format\n\nPlease select a .zip file');
        document.getElementById('modPackageUpload').value = '';
        selectedModFile = null;
        updateFileDisplay();
        return;
    }

    // Validate filename starts with modpkg_
    if (!file.name.startsWith('modpkg_')) {
        showPopup('error', 'Invalid mod package name: Mod package filename must start with "modpkg_" Example: modpkg_2026-01-09_12-33-01-670.zip');
        document.getElementById('modPackageUpload').value = '';
        selectedModFile = null;
        updateFileDisplay();
        return;
    }

    // Check file size (limit to 500MB)
    const maxSize = 500 * 1024 * 1024;
    if (file.size > maxSize) {
        showPopup('error', 'File too large: Maximum file size is 500MB');
        document.getElementById('modPackageUpload').value = '';
        selectedModFile = null;
        updateFileDisplay();
        return;
    }

    // Store the file for later upload
    selectedModFile = file;
    updateFileDisplay();
    showNotification('✓ File selected: ' + file.name + ' (' + (file.size / 1024 / 1024).toFixed(2) + 'MB)', 'success');
}

function updateFileDisplay() {
    const uploadZone = document.getElementById('modPackageUploadZone');
    const uploadBtn = document.getElementById('uploadModPackageBtn');
    
    if (!uploadZone || !uploadBtn) return;
    
    if (selectedModFile) {
        uploadZone.innerHTML = '<span class="upload-icon">✓</span>' +
                              '<div class="upload-text">File Selected</div>' +
                              '<div class="upload-subtext">' + selectedModFile.name + '</div>' +
                              '<div class="upload-subtext" style="margin-top: 5px; color: #00d4ff;">' + (selectedModFile.size / 1024 / 1024).toFixed(2) + ' MB</div>';
        uploadBtn.disabled = false;
        uploadBtn.style.opacity = '1';
    } else {
        uploadZone.innerHTML = '<span class="upload-icon">📦</span>' +
                              '<div class="upload-text">Drag & Drop Mod Package Here</div>' +
                              '<div class="upload-subtext">or click to select a .zip file</div>';
        uploadBtn.disabled = true;
        uploadBtn.style.opacity = '0.5';
    }
}

function uploadModPackage() {
    if (!selectedModFile) {
        showPopup('error', 'No file selected');
        return;
    }

    setButtonLoading('uploadModPackageBtn', true);
    showPopup('info', 'Uploading mod package...\n\n' + selectedModFile.name + ' (' + (selectedModFile.size / 1024 / 1024).toFixed(2) + 'MB)');
    updateUploadProgress(0);

    const reader = new FileReader();
    reader.onload = function(e) {
        const zipData = e.target.result;
        
        fetch('/api/v2/slp/upload', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/zip'
            },
            body: zipData
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showPopup('success', 'Mod package uploaded successfully!\n\n' + (data.message || 'The mods have been extracted and are ready to use.'));
                selectedModFile = null;
                document.getElementById('modPackageUpload').value = '';
                updateFileDisplay();
                updateUploadProgress(100);
                setTimeout(() => updateUploadProgress(0), 2000);
                setButtonLoading('uploadModPackageBtn', false);
            } else {
                showPopup('error', 'Failed to upload mod package:\n\n' + (data.error || 'Unknown error'));
                setButtonLoading('uploadModPackageBtn', false);
            }
        })
        .catch(error => {
            showPopup('error', 'Upload failed:\n\n' + (error.message || 'Network error'));
            setButtonLoading('uploadModPackageBtn', false);
        });
    };
    
    reader.onerror = function() {
        showPopup('error', 'Failed to read file');
        setButtonLoading('uploadModPackageBtn', false);
    };
    
    reader.readAsArrayBuffer(selectedModFile);
}

function updateUploadProgress(percent) {
    const progressBar = document.getElementById('modPackageUploadProgress');
    if (progressBar) {
        progressBar.style.width = percent + '%';
        if (percent > 0) {
            progressBar.classList.add('active');
        } else {
            progressBar.classList.remove('active');
        }
    }
}

// Drag and drop support
function initializeDragAndDrop() {
    const uploadZone = document.getElementById('modPackageUploadZone');
    if (!uploadZone) return;

    // Click to open file selector
    uploadZone.addEventListener('click', function() {
        document.getElementById('modPackageUpload').click();
    });

    // Drag and drop events
    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
        uploadZone.addEventListener(eventName, preventDefaults, false);
    });

    function preventDefaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }

    ['dragenter', 'dragover'].forEach(eventName => {
        uploadZone.addEventListener(eventName, highlight, false);
    });

    ['dragleave', 'drop'].forEach(eventName => {
        uploadZone.addEventListener(eventName, unhighlight, false);
    });

    function highlight(e) {
        uploadZone.classList.add('highlight');
    }

    function unhighlight(e) {
        uploadZone.classList.remove('highlight');
    }

    uploadZone.addEventListener('drop', handleDrop, false);

    function handleDrop(e) {
        const dt = e.dataTransfer;
        const files = dt.files;
        handleModPackageSelection(files);
    }
    
    // Initialize file display on load
    updateFileDisplay();
}

// Initialize on page load
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializeDragAndDrop);
} else {
    initializeDragAndDrop();
}