<script>
    import ConsoleCard from './cards/ConsoleCard.svelte';
    import SystemInfoCard from './cards/SystemInfoCard.svelte';
    import LogsCard from './cards/LogsCard.svelte';
    import QuickActionsCardAlternative from './cards/QuickActionsCardAlternative.svelte';
    import { apiFetch } from '../../services/api-v7';
    import { onMount, onDestroy } from 'svelte';

    let backgroundImageUrl = '';
    let isImageLoaded = false;
    let fileInput;
    let isUploading = false;
    let uploadError = '';

    onMount(async () => {
        try {
            const response = await apiFetch('/files/dashboard-background.png');
            if (response.ok) {
                const blob = await response.blob();
                backgroundImageUrl = URL.createObjectURL(blob);

                const img = new Image();
                img.src = backgroundImageUrl;
                img.onload = () => {
                    isImageLoaded = true;
                };
            } else {
                // If 404 or other error, keep backgroundImageUrl empty but still allow the component to render
                console.warn('Background image not found or failed to fetch:', response.status);
                isImageLoaded = true;
            }
        } catch (error) {
            console.error('Failed to fetch background image:', error);
            isImageLoaded = true;
        }
    });

    onDestroy(() => {
        if (backgroundImageUrl.startsWith('blob:')) {
            URL.revokeObjectURL(backgroundImageUrl);
        }
    });

    async function handleFileSelect(event) {
        const file = event.target.files?.[0];
        if (!file) return;

        if (file.type !== 'image/png') {
            uploadError = 'Only PNG files are allowed';
            setTimeout(() => uploadError = '', 3000);
            return;
        }

        isUploading = true;
        uploadError = '';

        try {
            const formData = new FormData();
            formData.append('file', file);

            const response = await apiFetch('/api/v2/settings/files/background/upload', {
                method: 'POST',
                body: formData
            });

            const result = await response.json();

            if (result.status === 'success') {
                isImageLoaded = false;
                window.location.reload();
            } else {
                uploadError = result.message || 'Upload failed';
                setTimeout(() => uploadError = '', 3000);
            }
        } catch (error) {
            uploadError = 'Upload failed: ' + error.message;
            setTimeout(() => uploadError = '', 3000);
        } finally {
            isUploading = false;
            if (fileInput) fileInput.value = '';
        }
    }

    function triggerFileSelect() {
        fileInput?.click();
    }
</script>

<div
    class="dashboard-container"
    class:image-loaded={isImageLoaded}
    style={backgroundImageUrl ? `background-image: url(${backgroundImageUrl});` : ''}
>
    <div class="dashboard-grid">
        <QuickActionsCardAlternative />
        <SystemInfoCard />
        <ConsoleCard />
        <LogsCard />
    </div>

    <input
        type="file"
        accept="image/png"
        bind:this={fileInput}
        on:change={handleFileSelect}
        style="display: none;"
    />

    <button
        class="background-upload-btn"
        on:click={triggerFileSelect}
        disabled={isUploading}
        title="Change background image"
    >
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
            <circle cx="8.5" cy="8.5" r="1.5"/>
            <polyline points="21 15 16 10 5 21"/>
        </svg>
        <span class="btn-text">Change Background</span>
        {#if isUploading}
            <span class="loading-spinner"></span>
        {/if}
    </button>

    {#if uploadError}
        <div class="error-toast">
            {uploadError}
        </div>
    {/if}
</div>

<style>
    .dashboard-container {
        position: relative;
        min-height: 100%;
        width: 100%;
        background-size: cover;
        background-repeat: no-repeat;
        background-position: center;
        border-radius: 8px;
        opacity: 0;
        transition: opacity 0.5s ease-in;
    }

    .dashboard-container.image-loaded {
        opacity: 1;
    }

    .dashboard-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        grid-auto-rows: minmax(0, 400px);
        gap: 1.5rem;
        padding: 1rem;
    }

    .background-upload-btn {
        position: fixed;
        bottom: 1.5rem;
        right: 1.5rem;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1rem;
        background: var(--bg-secondary);
        color: var(--text-primary);
        border: 1px solid var(--border-color);
        border-radius: 8px;
        cursor: pointer;
        transition: all var(--transition-speed);
        box-shadow: var(--shadow-medium);
        font-size: 0.875rem;
        font-weight: 500;
        backdrop-filter: blur(10px);
        z-index: 1000;
    }

    .background-upload-btn:hover:not(:disabled) {
        background: var(--bg-hover);
        border-color: var(--accent-primary);
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }

    .background-upload-btn:active:not(:disabled) {
        transform: translateY(0);
    }

    .background-upload-btn:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .icon {
        width: 1.25rem;
        height: 1.25rem;
        color: var(--accent-primary);
        flex-shrink: 0;
    }

    .btn-text {
        white-space: nowrap;
    }

    .loading-spinner {
        width: 1rem;
        height: 1rem;
        border: 2px solid var(--border-color);
        border-top-color: var(--accent-primary);
        border-radius: 50%;
        animation: spin 0.8s linear infinite;
    }

    .error-toast {
        position: fixed;
        bottom: 5rem;
        right: 1.5rem;
        padding: 0.75rem 1rem;
        background: var(--bg-secondary);
        color: var(--text-warning);
        border: 1px solid var(--text-warning);
        border-radius: 8px;
        font-size: 0.875rem;
        box-shadow: var(--shadow-medium);
        backdrop-filter: blur(10px);
        z-index: 1001;
        animation: slideIn 0.3s ease-out;
    }

    @keyframes spin {
        to { transform: rotate(360deg); }
    }

    @keyframes slideIn {
        from {
            opacity: 0;
            transform: translateY(1rem);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    @media (max-width: 1024px) {
        .dashboard-grid {
            grid-template-columns: 1fr;
        }
        
        .background-upload-btn {
            bottom: 1rem;
            right: 1rem;
        }
        
        .btn-text {
            display: none;
        }
        
        .background-upload-btn {
            padding: 0.75rem;
        }
    }

    @media (min-width: 1025px) and (max-width: 1400px) {
        .dashboard-grid {
            grid-template-columns: repeat(2, 1fr);
        }
    }
</style>