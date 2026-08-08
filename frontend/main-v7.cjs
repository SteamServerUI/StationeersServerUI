const { app, BrowserWindow, Menu, dialog, ipcMain, net, safeStorage } = require('electron');
const { autoUpdater } = require('electron-updater');
const express = require('express');
const http = require('http');
const fs = require('fs');
const path = require('path');

const ports = [28080, 28888, 29090, 27070, 26060, 35050, 34040, 30303];
const streams = new Map();
let portIndex = 0;
let server;
let mainWindow;
let localOrigin;
let credentials = {};
let trustedCertificates = {};

function credentialPath() {
  return path.join(app.getPath('userData'), 'backend-credentials.bin');
}

function trustPath() {
  return path.join(app.getPath('userData'), 'trusted-certificates.json');
}

function loadState() {
  try {
    const encrypted = fs.readFileSync(credentialPath(), 'utf8');
    if (safeStorage.isEncryptionAvailable()) {
      credentials = JSON.parse(safeStorage.decryptString(Buffer.from(encrypted, 'base64')));
    }
  } catch {
    credentials = {};
  }
  try {
    trustedCertificates = JSON.parse(fs.readFileSync(trustPath(), 'utf8'));
  } catch {
    trustedCertificates = {};
  }
}

function saveCredentials() {
  if (!safeStorage.isEncryptionAvailable()) throw new Error('OS credential encryption is unavailable');
  const encrypted = safeStorage.encryptString(JSON.stringify(credentials)).toString('base64');
  fs.writeFileSync(credentialPath(), encrypted, { mode: 0o600 });
}

function saveTrustedCertificates() {
  fs.writeFileSync(trustPath(), JSON.stringify(trustedCertificates, null, 2), { mode: 0o600 });
}

function assertRenderer(event) {
  if (!event.senderFrame.url.startsWith(localOrigin)) throw new Error('Request did not come from the SSUI renderer');
}

function backendURL(value) {
  const parsed = new URL(value);
  if (parsed.protocol !== 'https:') throw new Error('Desktop backends must use HTTPS');
  return parsed;
}

function tokenFor(url) {
  return credentials[backendURL(url).origin]?.token || '';
}

async function backendRequest(url, options = {}, token = tokenFor(url)) {
  const target = backendURL(url);
  const headers = new Headers(options.headers || {});
  if (token) headers.set('Authorization', `Bearer ${token}`);
  const response = await net.fetch(target.toString(), {
    method: options.method || 'GET',
    headers,
    body: options.body || undefined,
    redirect: 'error'
  });
  return {
    status: response.status,
    statusText: response.statusText,
    headers: Object.fromEntries(response.headers.entries()),
    body: await response.text()
  };
}

async function desktopLogin(backendUrl, username, password) {
  const origin = backendURL(backendUrl).origin;
  const response = await backendRequest(`${origin}/auth/desktop/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
  }, '');
  const data = JSON.parse(response.body || '{}');
  if (response.status < 200 || response.status >= 300) {
    throw new Error(data?.error?.message || `Login failed (${response.status})`);
  }
  credentials[origin] = { token: data.token, savedAt: new Date().toISOString() };
  saveCredentials();
  delete data.token;
  return data;
}

function registerIPC() {
  ipcMain.handle('ssui:request', async (event, request) => {
    assertRenderer(event);
    const target = backendURL(request.url);
    const token = tokenFor(target.toString());
    const publicPath = ['/api/v2/auth/setup/status', '/api/v2/auth/setup/bootstrap', '/auth/desktop/login'].includes(target.pathname);
    if (!token && !publicPath) throw new Error('This backend needs a desktop login');
    return backendRequest(target.toString(), request.options, token);
  });

  ipcMain.handle('ssui:login', async (event, request) => {
    assertRenderer(event);
    return desktopLogin(request.backendUrl, request.username, request.password);
  });

  ipcMain.handle('ssui:bootstrap', async (event, request) => {
    assertRenderer(event);
    const origin = backendURL(request.backendUrl).origin;
    const response = await backendRequest(`${origin}/api/v2/auth/setup/bootstrap`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ setupSecret: request.setupSecret, username: request.username, password: request.password })
    }, '');
    if (response.status < 200 || response.status >= 300) {
      const data = JSON.parse(response.body || '{}');
      throw new Error(data?.error?.message || `Setup failed (${response.status})`);
    }
    return desktopLogin(origin, request.username, request.password);
  });

  ipcMain.handle('ssui:logout', async (event, request) => {
    assertRenderer(event);
    const origin = backendURL(request.backendUrl).origin;
    const token = credentials[origin]?.token;
    if (token) {
      try {
        await backendRequest(`${origin}/api/v2/auth/desktop/logout`, { method: 'POST' }, token);
      } finally {
        delete credentials[origin];
        saveCredentials();
      }
    }
  });

  ipcMain.handle('ssui:trust-certificate', (event, request) => {
    assertRenderer(event);
    const origin = backendURL(request.origin).origin;
    trustedCertificates[origin] = request.fingerprint;
    saveTrustedCertificates();
  });

  ipcMain.handle('ssui:sse-start', async (event, request) => {
    assertRenderer(event);
    const target = backendURL(request.url);
    const token = tokenFor(target.toString());
    if (!token) throw new Error('This backend needs a desktop login');
    const controller = new AbortController();
    streams.set(request.id, controller);
    streamEvents(event.sender, request.id, target.toString(), token, controller.signal);
  });

  ipcMain.on('ssui:sse-stop', (_event, id) => {
    streams.get(id)?.abort();
    streams.delete(id);
  });
}

async function streamEvents(sender, id, url, token, signal) {
  try {
    const response = await net.fetch(url, { headers: { Authorization: `Bearer ${token}` }, signal });
    if (!response.ok || !response.body) throw new Error(`Event stream returned ${response.status}`);
    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';
    while (!signal.aborted) {
      const { done, value } = await reader.read();
      if (done) break;
      buffer += decoder.decode(value, { stream: true }).replace(/\r\n/g, '\n');
      let boundary;
      while ((boundary = buffer.indexOf('\n\n')) >= 0) {
        const block = buffer.slice(0, boundary);
        buffer = buffer.slice(boundary + 2);
        const data = block.split('\n').filter(line => line.startsWith('data:')).map(line => line.slice(5).trimStart()).join('\n');
        if (data) sender.send('ssui:sse-data', { id, data });
      }
    }
  } catch (error) {
    if (!signal.aborted) sender.send('ssui:sse-error', { id, error: error.message });
  } finally {
    streams.delete(id);
  }
}

function startServer() {
  const assets = path.join(process.resourcesPath, 'SSUI/onboard_bundled/v2');
  if (!fs.existsSync(assets)) throw new Error(`Could not find UI assets at ${assets}`);
  const serverApp = express();
  serverApp.use(express.static(assets));
  serverApp.use((_request, response) => response.sendFile(path.join(assets, 'index.html')));
  server = http.createServer(serverApp);
  return new Promise((resolve, reject) => listen(resolve, reject));
}

function listen(resolve, reject) {
  server.once('error', error => {
    if (error.code === 'EADDRINUSE' && portIndex < ports.length - 1) {
      portIndex++;
      server.close(() => listen(resolve, reject));
      return;
    }
    reject(error);
  });
  server.listen(ports[portIndex], '127.0.0.1', () => {
    localOrigin = `http://127.0.0.1:${ports[portIndex]}`;
    resolve();
  });
}

async function createWindow() {
  await startServer();
  mainWindow = new BrowserWindow({
    width: 1440,
    height: 900,
    minWidth: 1024,
    minHeight: 600,
    webPreferences: {
      preload: path.join(__dirname, 'preload.cjs'),
      nodeIntegration: false,
      contextIsolation: true,
      sandbox: true
    }
  });
  await mainWindow.loadURL(localOrigin);
}

function createMenu() {
  Menu.setApplicationMenu(Menu.buildFromTemplate([{
    label: 'SSUI',
    submenu: [
      { label: 'Check for updates', click: () => autoUpdater.checkForUpdatesAndNotify() },
      { label: 'Toggle developer tools', click: () => mainWindow?.webContents.toggleDevTools() },
      { type: 'separator' },
      { role: 'quit' }
    ]
  }]));
}

app.on('certificate-error', (event, _contents, url, error, certificate, callback) => {
  event.preventDefault();
  const origin = new URL(url).origin;
  const fingerprint = certificate.fingerprint256 || certificate.fingerprint;
  if (trustedCertificates[origin] === fingerprint) {
    callback(true);
    return;
  }
  callback(false);
  mainWindow?.webContents.send('ssui:certificate-error', { origin, fingerprint, error });
});

app.whenReady().then(async () => {
  loadState();
  registerIPC();
  createMenu();
  await createWindow();
  if (app.isPackaged) autoUpdater.checkForUpdatesAndNotify();
}).catch(error => {
  dialog.showErrorBox('Steam Server UI', error.message);
  app.quit();
});

app.on('window-all-closed', () => app.quit());
app.on('before-quit', () => {
  for (const controller of streams.values()) controller.abort();
  server?.close();
});
