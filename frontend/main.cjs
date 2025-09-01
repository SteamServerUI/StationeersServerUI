const { app, BrowserWindow, Menu, dialog } = require('electron');
const { autoUpdater } = require('electron-updater');
const path = require('path');
const fs = require('fs');
const express = require('express');
const https = require('https');
const forge = require('node-forge');
const os = require('os');

// Configure auto-updater
// Only check for updates in production builds
if (process.env.NODE_ENV !== 'development') {
  autoUpdater.setFeedURL({
    provider: 'github',
    owner: 'SteamServerUI',
    repo: 'SteamServerUI'
  });

  // Disable automatic installation on Linux for security/stability
  if (process.platform === 'linux') {
    autoUpdater.autoInstallOnAppQuit = false;
    autoUpdater.autoDownload = true;
  }
}

// Auto-updater event handlers
autoUpdater.on('checking-for-update', () => {
  console.log('Checking for update...');
});

autoUpdater.on('update-available', (info) => {
  console.log('Update available:', info.version);
  // Optional: Show notification to user
  dialog.showMessageBox({
    type: 'info',
    title: 'Update Available',
    message: `A new version (${info.version}) is available. It will be downloaded in the background.`,
    buttons: ['OK']
  });
});

autoUpdater.on('update-not-available', (info) => {
  console.log('Update not available.');
});

autoUpdater.on('error', (err) => {
  console.error('Auto-updater error:', err);
});

autoUpdater.on('download-progress', (progressObj) => {
  let log_message = "Download speed: " + progressObj.bytesPerSecond;
  log_message = log_message + ' - Downloaded ' + progressObj.percent + '%';
  log_message = log_message + ' (' + progressObj.transferred + "/" + progressObj.total + ')';
  console.log(log_message);
});

autoUpdater.on('update-downloaded', (info) => {
  console.log('Update downloaded');

  if (process.platform === 'linux') {
    // On Linux, show manual installation instructions
    const updatePath = path.join(require('os').homedir(), '.cache', 'steamserverui-updater', 'pending');
    dialog.showMessageBox({
      type: 'info',
      title: 'Update Downloaded',
      message: `Update v${info.version} has been downloaded!\n\nFor security reasons, please install manually:\n\n1. Close this application\n2. Open terminal and run:\n   sudo dpkg -i "${updatePath}/SSUI-Desktop-v${info.version}-linux.deb"\n\nOr double-click the downloaded .deb file in your file manager.`,
      buttons: ['Open Download Folder', 'Later', 'Quit App']
    }).then((result) => {
      if (result.response === 0) {
        // Open the download folder
        require('electron').shell.openPath(updatePath);
      } else if (result.response === 2) {
        // Quit the app so user can install manually
        app.quit();
      }
    });
  } else {
    // Windows
    dialog.showMessageBox({
      type: 'info',
      title: 'Update Ready',
      message: 'Update has been downloaded. The application will restart to apply the update.',
      buttons: ['Restart Now', 'Later']
    }).then((result) => {
      if (result.response === 0) {
        autoUpdater.quitAndInstall();
      }
    });
  }
});

// Certificate generation function
function generateSelfSignedCertificate() {
  console.log('Generating self-signed certificate...');
  
  // Generate a key pair
  const keys = forge.pki.rsa.generateKeyPair(2048);
  
  // Create a certificate
  const cert = forge.pki.createCertificate();
  cert.publicKey = keys.publicKey;
  cert.serialNumber = '01';
  cert.validity.notBefore = new Date();
  cert.validity.notAfter = new Date();
  cert.validity.notAfter.setFullYear(cert.validity.notBefore.getFullYear() + 1);
  
  const attrs = [
    { name: 'commonName', value: 'SteamServerUIDesktop' },
    { name: 'countryName', value: 'SE' },
    { shortName: 'ST', value: 'CA' },
    { name: 'localityName', value: 'Local' },
    { name: 'organizationName', value: 'SteamServerUI' },
    { shortName: 'OU', value: 'ElectronApp' }
  ];
  
  cert.setSubject(attrs);
  cert.setIssuer(attrs);
  cert.setExtensions([
    {
      name: 'basicConstraints',
      cA: true
    },
    {
      name: 'keyUsage',
      keyCertSign: true,
      digitalSignature: true,
      nonRepudiation: true,
      keyEncipherment: true,
      dataEncipherment: true
    },
    {
      name: 'extKeyUsage',
      serverAuth: true,
      clientAuth: true,
      codeSigning: true,
      emailProtection: true,
      timeStamping: true
    },
    {
      name: 'nsCertType',
      client: true,
      server: true,
      email: true,
      objsign: true,
      sslCA: true,
      emailCA: true,
      objCA: true
    },
    {
      name: 'subjectAltName',
      altNames: [
        { type: 2, value: 'localhost' },
        { type: 2, value: '127.0.0.1' },
        { type: 7, ip: '127.0.0.1' },
        { type: 7, ip: '::1' }
      ]
    }
  ]);
  
  // Sign the certificate
  cert.sign(keys.privateKey);
  
  // Convert to PEM format
  const certPem = forge.pki.certificateToPem(cert);
  const keyPem = forge.pki.privateKeyToPem(keys.privateKey);
  
  return {
    cert: certPem,
    key: keyPem
  };
}

// Get or create certificate
function getCertificate() {
  const certDir = path.join(os.homedir(), '.steamserverui');
  const certPath = path.join(certDir, 'cert.pem');
  const keyPath = path.join(certDir, 'key.pem');
  
  // Create directory if it doesn't exist
  if (!fs.existsSync(certDir)) {
    fs.mkdirSync(certDir, { recursive: true });
  }
  
  // Check if certificate files exist and are valid
  if (fs.existsSync(certPath) && fs.existsSync(keyPath)) {
    try {
      const cert = fs.readFileSync(certPath, 'utf8');
      const key = fs.readFileSync(keyPath, 'utf8');
      
      // Check if certificate is still valid (not expired)
      const certObj = forge.pki.certificateFromPem(cert);
      const now = new Date();
      if (now < certObj.validity.notAfter) {
        console.log('Using existing certificate');
        return { cert, key };
      } else {
        console.log('Certificate expired, generating new one');
      }
    } catch (err) {
      console.log('Error reading existing certificate, generating new one:', err.message);
    }
  }
  
  // Generate new certificate
  const certData = generateSelfSignedCertificate();
  
  // Save certificate files
  fs.writeFileSync(certPath, certData.cert);
  fs.writeFileSync(keyPath, certData.key);
  
  console.log('Certificate saved to:', certDir);
  return certData;
}

// Static file server
let server;
let currentPortIndex = 0;
const HTTPS_PORTS = [28443, 28889, 29443, 27443, 26443, 35443, 34443, 30443, 34943];

function startServer() {
  const expressApp = express();

  // Determine the assets directory path
  let assetsPath;

  // In production
  const prodPath = path.join(process.resourcesPath, 'UIMod/onboard_bundled/v2');

  // Check if path exists
  if (fs.existsSync(prodPath)) {
    assetsPath = prodPath;
  } else {
    dialog.showErrorBox('Error', 'Could not find assets directory: ' + prodPath);
    app.exit(1);
    return false;
  }

  console.log('Serving assets from:', assetsPath);

  // Serve static files from the assets directory
  expressApp.use(express.static(assetsPath));

  // Get certificate
  const certData = getCertificate();

  // Create HTTPS server
  const httpsOptions = {
    key: certData.key,
    cert: certData.cert
  };

  server = https.createServer(httpsOptions, expressApp);

  return new Promise((resolve, reject) => {
    server.listen(HTTPS_PORTS[currentPortIndex], () => {
      console.log(`HTTPS Server running at https://localhost:${HTTPS_PORTS[currentPortIndex]}`);
      resolve(true);
    });

    server.on('error', (err) => {
      console.error('Server error:', err);
      if (err.code === 'EADDRINUSE' && currentPortIndex < HTTPS_PORTS.length - 1) {
        currentPortIndex++;
        console.log(`Port ${HTTPS_PORTS[currentPortIndex - 1]} in use, trying ${HTTPS_PORTS[currentPortIndex]}`);
        server.close();
        startServer().then(resolve).catch(reject);
      } else {
        let errorMsg = err.code === 'EADDRINUSE' ? 'All default HTTPS ports are in use. Please free up a port and try again.' : err.message;
        dialog.showErrorBox('Server Error', errorMsg);
        reject(err);
      }
    });
  });
}

async function createWindow() {
  // Start the local server first
  const serverStarted = await startServer();
  if (!serverStarted) return;

  const win = new BrowserWindow({
    width: 1920,
    height: 1080,
    webPreferences: {
      nodeIntegration: false,
      webSecurity: true
    }
  });

  // Load from local HTTPS server instead of HTTP
  console.log(`Loading UI from: https://localhost:${HTTPS_PORTS[currentPortIndex]}/index.html`);
  win.loadURL(`https://localhost:${HTTPS_PORTS[currentPortIndex]}/index.html`);

  // For debugging
  // win.webContents.openDevTools();
}

// Create application menu with update check option
function createMenu() {
  const template = [
    {
      label: 'Options',
      submenu: [
        {
          label: 'Check for Updates',
          click: () => {
            autoUpdater.checkForUpdatesAndNotify();
          }
        },
        {
          label: 'About',
          click: () => {
            dialog.showMessageBox({
              type: 'info',
              title: 'About Steam Server UI',
              message: `SSUI Desktop ${app.getVersion()}\nCopyright © 2025 JacksonTheMaster`,
              buttons: ['OK']
            });
          }
        },
        {
          label: 'Backends',
          submenu: [
            {
              label: 'Reset',
              click: async () => {
                try {
                  // Get the main window
                  const mainWindow = BrowserWindow.getFocusedWindow() || BrowserWindow.getAllWindows()[0];
                  
                  if (mainWindow) {
                    // Clear all cookies
                    await mainWindow.webContents.session.clearStorageData({
                      storages: ['cookies', 'localstorage', 'sessionstorage']
                    });
                  
                    mainWindow.reload();
                    
                    dialog.showMessageBox(mainWindow, {
                      type: 'info',
                      title: 'Reset Complete',
                      message: 'All cookies and local storage have been cleared. The application has been reloaded.',
                      buttons: ['OK']
                    });
                  }
                } catch (error) {
                  console.error('Error resetting backend data:', error);
                  dialog.showErrorBox('Reset Failed', 'Failed to reset backend data. Please try again. Alternatively, you can manually clear your cookies and local storage from the Developer Console at ctrl+shift+i.');
                }
              }
            }
          ]
        },
        {
          label: 'Edit',
          submenu: [
            { role: 'cut' },
            { role: 'copy' },
            { role: 'paste' }
          ]
        },
        {
          label: 'View',
          submenu: [
            { role: 'reload' },
            { role: 'forcereload' },
            { role: 'toggledevtools' },
            { type: 'separator' },
            { role: 'resetzoom' },
            { role: 'zoomin' },
            { role: 'zoomout' },
            { type: 'separator' },
            { role: 'togglefullscreen' }
          ]
        },
        {
          label: 'Open Github Page',
          click: () => {
            require('electron').shell.openExternal('https://github.com/SteamServerUI/SteamServerUI');
          }
          
        }, 
        {
          label: 'Join Discord Server',
          click: () => {
            require('electron').shell.openExternal('https://discord.com/invite/8n3vN92MyJ');
          }
        },
        {
          label: 'Report Issue',
          click: () => {
            require('electron').shell.openExternal('https://github.com/SteamServerUI/SteamServerUI/issues');
          }
        }
      ]
    }
  ];

  const menu = Menu.buildFromTemplate(template);
  Menu.setApplicationMenu(menu);
}

app.commandLine.appendSwitch('ignore-certificate-errors');
app.commandLine.appendSwitch('ignore-certificate-errors-spki-list');
app.commandLine.appendSwitch('ignore-ssl-errors');

app.whenReady().then(() => {
  createWindow();
  createMenu();

  // Check for updates after app is ready (but not in development)
  if (!process.env.NODE_ENV || process.env.NODE_ENV === 'production') {
    // Check for updates immediately
    autoUpdater.checkForUpdatesAndNotify();

    // Check for updates every 30 minutes
    setInterval(() => {
      autoUpdater.checkForUpdatesAndNotify();
    }, 30 * 60 * 1000);
  }
});

app.on('window-all-closed', () => {
  app.quit();
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});

app.on('quit', () => {
  if (server) {
    server.closeAllConnections();
    server.close(() => {
      console.log('Server closed');
    });
    // Force quit after 5 seconds if server doesn't close
    setTimeout(() => {
      console.log('Forcing app exit after timeout');
      app.exit(0);
    }, 5000);
  }
});

// Minimal error handling
process.on('uncaughtException', (err) => {
  console.error('Uncaught Exception:', err);
  dialog.showErrorBox('Unexpected Error', 'An unexpected error occurred: ' + err.message);
  app.exit(1);
});

process.on('unhandledRejection', (err) => {
  console.error('Unhandled Rejection:', err);
  dialog.showErrorBox('Unexpected Error', 'An unexpected error occurred: ' + (err.message || err));
  app.exit(1);
});