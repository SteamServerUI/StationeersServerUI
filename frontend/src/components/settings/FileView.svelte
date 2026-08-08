<script>
  import { onMount } from 'svelte';
  import { apiFetch } from '../../services/api-v7';

  let files = $state([]);
  let selectedFile = $state(null);
  let editorMode = $state('raw'); // 'raw', 'json', 'ini', 'yaml'
  let fileContent = $state('');
  let originalContent = $state('');
  let parsedContent = $state(null);
  let loading = $state(false);
  let error = $state(null);
  let saving = $state(false);
  let notification = $state(null);
  let hasChanges = $derived(fileContent !== originalContent);

  onMount(async () => {
    await loadFileList();
  });

  const maxNestingLevels = {
    json: 2,
    yaml: 2,
    ini: 1
  };

  function showNotification(message, type = 'info') {
    notification = { message, type };
    setTimeout(() => {
      notification = null;
    }, 4000);
  }

  async function loadFileList() {
    loading = true;
    error = null;
    try {
      const response = await apiFetch('/api/v2/files');
      const data = await response.json();
      
      if (data.success) {
        files = data.data;
      } else {
        error = data.message || 'Failed to load files';
      }
    } catch (err) {
      error = 'Error loading file list: ' + err.message;
    } finally {
      loading = false;
    }
  }

  function checkDeepNesting(obj, fileType, currentDepth = 0) {
    const maxDepth = maxNestingLevels[fileType.toLowerCase()] || 2;
    if (currentDepth > maxDepth) return true;
    if (typeof obj !== 'object' || obj === null) return false;
    
    return Object.values(obj).some(value => 
      typeof value === 'object' && value !== null && 
      checkDeepNesting(value, fileType, currentDepth + 1)
    );
  }

  async function selectFile(file, mode = 'raw') {
    if (hasChanges && !confirm('You have unsaved changes. Are you sure you want to switch files?')) {
      return;
    }

    loading = true;
    error = null;
    selectedFile = file;
    editorMode = mode;
    
    try {
      const response = await apiFetch('/api/v2/files/get', {
        method: 'POST',
        body: JSON.stringify({ filename: file.filename })
      });
      
      const content = await response.text();
      fileContent = content;
      originalContent = content;
      
      if (mode !== 'raw') {
        try {
          let parsed;
          if (mode === 'json') {
            parsed = JSON.parse(content);
          } else if (mode === 'yaml') {
            parsed = parseYAML(content);
          } else {
            parsed = parseINI(content);
          }

          if (checkDeepNesting(parsed, mode)) {
            showNotification(`Cannot render ${mode.toUpperCase()} file in structured view due to deep nesting`, 'error');
            editorMode = 'raw';
            parsedContent = null;
          } else {
            parseContent(content, mode);
          }
        } catch (err) {
          showNotification(`Failed to parse ${mode.toUpperCase()}: ${err.message}`, 'error');
          editorMode = 'raw';
          parsedContent = null;
        }
      }
    } catch (err) {
      error = 'Error loading file: ' + err.message;
      selectedFile = null;
    } finally {
      loading = false;
    }
  }

  function parseContent(content, mode) {
    try {
      if (mode === 'json') {
        parsedContent = JSON.parse(content);
      } else if (mode === 'ini') {
        parsedContent = parseINI(content);
      } else if (mode === 'yaml') {
        parsedContent = parseYAML(content);
      }
    } catch (err) {
      showNotification(`Failed to parse ${mode.toUpperCase()}: ${err.message}`, 'error');
      editorMode = 'raw';
      parsedContent = null;
    }
  }

  function parseINI(content) {
    const lines = content.split('\n');
    const result = {};
    let currentSection = 'global';
    result[currentSection] = {};

    for (let line of lines) {
      line = line.trim();
      if (!line || line.startsWith(';') || line.startsWith('#')) continue;
      
      if (line.startsWith('[') && line.endsWith(']')) {
        currentSection = line.slice(1, -1);
        result[currentSection] = {};
      } else {
        const idx = line.indexOf('=');
        if (idx > 0) {
          const key = line.slice(0, idx).trim();
          const value = line.slice(idx + 1).trim();
          result[currentSection][key] = value;
        }
      }
    }
    return result;
  }

  function parseYAML(content) {
    const lines = content.split('\n');
    const result = {};
    const stack = [{ obj: result, indent: -1 }];

    for (let line of lines) {
      if (!line.trim() || line.trim().startsWith('#')) continue;
      
      const indent = line.search(/\S/);
      const trimmed = line.trim();
      
      while (stack.length > 1 && indent <= stack[stack.length - 1].indent) {
        stack.pop();
      }
      
      if (trimmed.includes(':')) {
        const idx = trimmed.indexOf(':');
        const key = trimmed.slice(0, idx).trim();
        const value = trimmed.slice(idx + 1).trim();
        
        if (value) {
          stack[stack.length - 1].obj[key] = value;
        } else {
          const newObj = {};
          stack[stack.length - 1].obj[key] = newObj;
          stack.push({ obj: newObj, indent });
        }
      }
    }
    return result;
  }

  function serializeINI(data) {
    let result = '';
    for (const [section, values] of Object.entries(data)) {
      if (section !== 'global' || Object.keys(values).length > 0) {
        if (section !== 'global') {
          result += `[${section}]\n`;
        }
        for (const [key, value] of Object.entries(values)) {
          result += `${key}=${value}\n`;
        }
        result += '\n';
      }
    }
    return result.trim();
  }

  function serializeYAML(data, indent = 0) {
    let result = '';
    for (const [key, value] of Object.entries(data)) {
      const spaces = '  '.repeat(indent);
      if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
        result += `${spaces}${key}:\n`;
        result += serializeYAML(value, indent + 1);
      } else {
        result += `${spaces}${key}: ${value}\n`;
      }
    }
    return result;
  }

  function updateParsedValue(path, value) {
    const keys = path.split('.');
    let obj = parsedContent;
    
    for (let i = 0; i < keys.length - 1; i++) {
      obj = obj[keys[i]];
    }
    obj[keys[keys.length - 1]] = value;
    
    if (editorMode === 'json') {
      fileContent = JSON.stringify(parsedContent, null, 2);
    } else if (editorMode === 'ini') {
      fileContent = serializeINI(parsedContent);
    } else if (editorMode === 'yaml') {
      fileContent = serializeYAML(parsedContent);
    }
  }

  async function saveFile() {
    if (!selectedFile) return;
    
    saving = true;
    error = null;
    
    try {
      const response = await apiFetch(`/api/v2/files/save?filename=${selectedFile.filename}`, {
        method: 'POST',
        body: fileContent
      });
      
      const data = await response.json();
      
      if (data.success) {
        originalContent = fileContent;
        showNotification(data.message || 'File saved successfully', 'success');
      } else {
        error = data.message || 'Failed to save file';
        showNotification(error, 'error');
      }
    } catch (err) {
      error = 'Error saving file: ' + err.message;
      showNotification(error, 'error');
    } finally {
      saving = false;
    }
  }

  function closeEditor() {
    if (hasChanges && !confirm('You have unsaved changes. Are you sure you want to close?')) {
      return;
    }
    selectedFile = null;
    fileContent = '';
    originalContent = '';
    parsedContent = null;
    editorMode = 'raw';
  }

  function switchMode(mode) {
    if (mode === 'raw') {
      editorMode = 'raw';
      parsedContent = null;
    } else {
      try {
        let parsed;
        if (mode === 'json') {
          parsed = JSON.parse(fileContent);
        } else if (mode === 'yaml') {
          parsed = parseYAML(fileContent);
        } else {
          parsed = parseINI(fileContent);
        }

        if (checkDeepNesting(parsed, mode)) {
          showNotification(`Cannot render ${mode.toUpperCase()} file in structured view due to deep nesting`, 'error');
          editorMode = 'raw';
          parsedContent = null;
        } else {
          parseContent(fileContent, mode);
          if (parsedContent) {
            editorMode = mode;
          }
        }
      } catch (err) {
        showNotification(`Failed to parse ${mode.toUpperCase()}: ${err.message}`, 'error');
        editorMode = 'raw';
        parsedContent = null;
      }
    }
  }

  function canUseStructuredEditor(file) {
    return ['json', 'ini', 'yaml'].includes(file.type.toLowerCase());
  }

  function getEditorLabel(type) {
    const labels = {
      json: 'JSON Editor',
      ini: 'INI Editor',
      yaml: 'YAML Editor'
    };
    return labels[type.toLowerCase()] || 'Editor';
  }
</script>

<div class="file-browser">
  {#if notification}
    <div class="notification notification-{notification.type}">
      {notification.message}
    </div>
  {/if}

  {#if error}
    <div class="error-banner">
      {error}
    </div>
  {/if}

  <div class="container">
    {#if !selectedFile}
      <div class="file-list-container">
        <div class="header">
          <h2>File Browser</h2>
        </div>

        {#if loading && files.length === 0}
          <div class="loading">Loading files...</div>
        {:else if files.length === 0}
          <div class="empty">No files available</div>
        {:else}
          <table class="file-table">
            <thead>
              <tr>
                <th>Filename</th>
                <th>Type</th>
                <th>Description</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {#each files as file}
                <tr>
                  <td class="filename">{file.filename}</td>
                  <td class="type">{file.type}</td>
                  <td class="description">{file.description}</td>
                  <td class="actions">
                    <button 
                      onclick={() => selectFile(file, 'raw')} 
                      class="action-btn raw-btn"
                      disabled={loading}
                    >
                      Raw Editor
                    </button>
                    {#if canUseStructuredEditor(file)}
                      <button 
                        onclick={() => selectFile(file, file.type.toLowerCase())} 
                        class="action-btn structured-btn"
                        disabled={loading}
                      >
                        {getEditorLabel(file.type)}
                      </button>
                    {/if}
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        {/if}
      </div>
    {:else}
      <div class="editor-container">
        <div class="editor-header">
          <div class="file-info">
            <button onclick={closeEditor} class="back-btn">Back</button>
            <h2>{selectedFile.filename}</h2>
            <span class="mode-badge">{editorMode.toUpperCase()}</span>
            {#if hasChanges}
              <span class="changes-indicator">• Unsaved changes</span>
            {/if}
          </div>
          <div class="editor-actions">
            {#if canUseStructuredEditor(selectedFile)}
              <div class="mode-switcher">
                <button 
                  onclick={() => switchMode('raw')}
                  class="mode-btn"
                  class:active={editorMode === 'raw'}
                >
                  Raw
                </button>
                <button 
                  onclick={() => switchMode(selectedFile.type.toLowerCase())}
                  class="mode-btn"
                  class:active={editorMode !== 'raw'}
                >
                  Structured
                </button>
              </div>
            {/if}
            <button 
              onclick={saveFile} 
              disabled={saving || !hasChanges}
              class="save-btn"
            >
              {saving ? 'Saving...' : 'Save'}
            </button>
          </div>
        </div>

        <div class="editor-content">
          {#if editorMode === 'raw'}
            <textarea 
              bind:value={fileContent}
              class="editor"
              placeholder="File content..."
              spellcheck="false"
            ></textarea>
          {:else if editorMode === 'json' && parsedContent}
            <div class="structured-editor">
              <div class="structured-header">JSON Structure</div>
              <div class="json-tree">
                {#each Object.entries(parsedContent) as [key, value], i}
                  <div class="json-item">
                    <!-- svelte-ignore a11y_label_has_associated_control -->
                    <label class="json-key">{key}</label>
                    {#if typeof value === 'object' && value !== null}
                      <div class="json-nested">
                        {#each Object.entries(value) as [nestedKey, nestedValue], j}
                          <div class="json-item nested">
                            <!-- svelte-ignore a11y_label_has_associated_control -->
                            <label class="json-key">{nestedKey}</label>
                            <input 
                              type="text"
                              value={nestedValue}
                              oninput={(e) => updateParsedValue(`${key}.${nestedKey}`, e.currentTarget.value)}
                              class="json-input"
                            />
                          </div>
                        {/each}
                      </div>
                    {:else}
                      <input 
                        type="text"
                        value={value}
                        oninput={(e) => updateParsedValue(key, e.currentTarget.value)}
                        class="json-input"
                      />
                    {/if}
                  </div>
                {/each}
              </div>
            </div>
          {:else if editorMode === 'ini' && parsedContent}
            <div class="structured-editor">
              <div class="structured-header">INI Structure</div>
              <div class="ini-editor">
                {#each Object.entries(parsedContent) as [section, values], i}
                  <div class="ini-section">
                    {#if section !== 'global'}
                      <div class="section-header">[{section}]</div>
                    {:else}
                      <div class="section-header">Global Settings</div>
                    {/if}
                    <div class="section-content">
                      {#each Object.entries(values) as [key, value], j}
                        <div class="ini-item">
                          <!-- svelte-ignore a11y_label_has_associated_control -->
                          <label class="ini-key">{key}</label>
                          <input 
                            type="text"
                            value={value}
                            oninput={(e) => updateParsedValue(`${section}.${key}`, e.currentTarget.value)}
                            class="ini-input"
                          />
                        </div>
                      {/each}
                    </div>
                  </div>
                {/each}
              </div>
            </div>
          {:else if editorMode === 'yaml' && parsedContent}
            <div class="structured-editor">
              <div class="structured-header">YAML Structure</div>
              <div class="yaml-editor">
                {#each Object.entries(parsedContent) as [key, value], i}
                  <div class="yaml-item">
                    <!-- svelte-ignore a11y_label_has_associated_control -->
                    <label class="yaml-key">{key}</label>
                    {#if typeof value === 'object' && value !== null}
                      <div class="yaml-nested">
                        {#each Object.entries(value) as [nestedKey, nestedValue], j}
                          <div class="yaml-item nested">
                            <!-- svelte-ignore a11y_label_has_associated_control -->
                            <label class="yaml-key">{nestedKey}</label>
                            <input 
                              type="text"
                              value={nestedValue}
                              oninput={(e) => updateParsedValue(`${key}.${nestedKey}`, e.currentTarget.value)}
                              class="yaml-input"
                            />
                          </div>
                        {/each}
                      </div>
                    {:else}
                      <input 
                        type="text"
                        value={value}
                        oninput={(e) => updateParsedValue(key, e.currentTarget.value)}
                        class="yaml-input"
                      />
                    {/if}
                  </div>
                {/each}
              </div>
            </div>
          {/if}
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .file-browser {
    width: 100%;
    height: 100%;
    background: var(--bg-tertiary);
    color: var(--text-primary);
    display: flex;
    flex-direction: column;
    border-radius: 8px;
    box-shadow: var(--shadow-light);
    position: relative;
  }

  .notification {
    position: absolute;
    top: 80px;
    right: 20px;
    padding: 12px 20px;
    border-radius: 6px;
    font-weight: 500;
    z-index: 1000;
    box-shadow: var(--shadow-medium);
    animation: slideIn 0.3s ease-out;
  }

  @keyframes slideIn {
    from {
      transform: translateX(400px);
      opacity: 0;
    }
    to {
      transform: translateX(0);
      opacity: 1;
    }
  }

  .notification-success {
    background: var(--accent-primary);
    color: var(--bg-primary);
  }

  .notification-error {
    background: var(--text-warning);
    color: var(--bg-primary);
  }

  .notification-info {
    background: var(--accent-secondary);
    color: var(--bg-primary);
  }

  .error-banner {
    background: var(--text-warning);
    color: var(--bg-secondary);
    padding: 12px 20px;
    font-weight: 500;
  }

  .container {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .file-list-container {
    padding: 20px;
    max-width: 1200px;
    margin: 0 auto;
    width: 100%;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
  }

  h2 {
    margin: 0;
    font-size: 24px;
    color: var(--text-primary);
  }

  .loading, .empty {
    text-align: center;
    padding: 40px;
    color: var(--text-secondary);
  }

  .file-table {
    width: 100%;
    border-collapse: collapse;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    overflow: hidden;
    background-color: var(--bg-secondary);
  }

  .file-table thead {
    background: var(--bg-tertiary);
  }

  .file-table th {
    padding: 12px 16px;
    text-align: left;
    font-weight: 600;
    color: var(--text-secondary);
    font-size: 14px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .file-table td {
    padding: 12px 16px;
    border-top: 1px solid var(--border-color);
  }

  .file-table tbody tr {
    transition: background 0.2s;
  }

  .file-table tbody tr:hover {
    background: var(--bg-hover);
  }

  .filename {
    font-weight: 500;
    font-family: monospace;
    color: var(--text-accent);
  }

  .mode-badge {
    display: inline-block;
    padding: 4px 12px;
    background: var(--accent-tertiary);
    color: var(--text-secondary);
    border-radius: 12px;
    font-size: 12px;
    font-weight: 500;
  }

  .description {
    color: var(--text-secondary);
  }

  .action-btn {
    padding: 6px 12px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 13px;
    font-weight: 500;
    transition: all 0.2s;
  }

  .raw-btn {
    background: var(--bg-hover);
    color: var(--text-primary);
    border: 1px solid var(--border-color);
  }

  .raw-btn:hover:not(:disabled) {
    background: var(--bg-active);
  }

  .structured-btn {
    background: var(--accent-primary);
    color: var(--bg-primary);
  }

  .structured-btn:hover:not(:disabled) {
    background: var(--accent-secondary);
  }

  .action-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .editor-container {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  .editor-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    background: var(--bg-secondary);
    border-bottom: 1px solid var(--border-color);
  }

  .file-info {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .editor-actions {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .mode-switcher {
    display: flex;
    background: var(--bg-tertiary);
    border-radius: 4px;
    padding: 2px;
  }

  .mode-btn {
    padding: 6px 12px;
    background: transparent;
    border: none;
    color: var(--text-secondary);
    cursor: pointer;
    font-size: 13px;
    border-radius: 3px;
    transition: all 0.2s;
  }

  .mode-btn.active {
    background: var(--accent-primary);
    color: var(--bg-primary);
  }

  .mode-btn:hover:not(.active) {
    background: var(--bg-hover);
  }

  .back-btn {
    padding: 8px 12px;
    background: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    color: var(--text-primary);
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    transition: background 0.2s;
  }

  .back-btn:hover {
    background: var(--bg-hover);
  }

  .changes-indicator {
    color: var(--text-warning);
    font-size: 14px;
    font-weight: 500;
  }

  .save-btn {
    padding: 8px 24px;
    background: var(--accent-primary);
    color: var(--bg-primary);
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    transition: background 0.2s;
  }

  .save-btn:hover:not(:disabled) {
    background: var(--accent-secondary);
  }

  .save-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .editor-content {
    flex: 1;
    padding: 20px;
    overflow: auto;
  }

  .editor {
    width: 100%;
    height: 100%;
    padding: 16px;
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    color: var(--text-primary);
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 14px;
    line-height: 1.6;
    resize: none;
    outline: none;
  }

  .editor:focus {
    border-color: var(--accent-primary);
    box-shadow: 0 0 0 3px var(--shadow-light);
  }

  .structured-editor {
    height: 100%;
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    overflow: auto;
  }

  .structured-header {
    padding: 16px;
    background: var(--bg-tertiary);
    border-bottom: 1px solid var(--border-color);
    font-weight: 600;
    font-size: 14px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--text-secondary);
  }

  .json-tree, .yaml-editor {
    padding: 16px;
  }

  .json-item, .yaml-item {
    margin-bottom: 16px;
  }

  .json-item.nested, .yaml-item.nested {
    margin-left: 24px;
    margin-bottom: 12px;
  }

  .json-key, .yaml-key {
    display: block;
    margin-bottom: 6px;
    font-weight: 500;
    color: var(--text-accent);
    font-size: 13px;
  }

  .json-input, .yaml-input, .ini-input {
    width: 100%;
    padding: 8px 12px;
    background: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    border-radius: 4px;
    color: var(--text-primary);
    font-family: monospace;
    font-size: 13px;
  }

  .json-input:focus, .yaml-input:focus, .ini-input:focus {
    outline: none;
    border-color: var(--accent-primary);
    box-shadow: 0 0 0 2px var(--shadow-light);
  }

  .json-nested, .yaml-nested {
    margin-top: 12px;
    padding-left: 16px;
    border-left: 2px solid var(--border-color);
  }

  .ini-editor {
    padding: 16px;
  }

  .ini-section {
    margin-bottom: 24px;
  }

  .section-header {
    padding: 8px 12px;
    background: var(--bg-tertiary);
    border-radius: 4px;
    font-weight: 600;
    color: var(--text-accent);
    margin-bottom: 12px;
    font-family: monospace;
  }

  .section-content {
    padding-left: 12px;
  }

  .ini-item {
    margin-bottom: 12px;
  }

  .ini-key {
    display: block;
    margin-bottom: 6px;
    font-weight: 500;
    color: var(--text-secondary);
    font-size: 13px;
    font-family: monospace;
  }
</style>