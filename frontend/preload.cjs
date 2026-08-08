const { contextBridge, ipcRenderer } = require('electron');
const crypto = require('crypto');

const streams = new Map();

ipcRenderer.on('ssui:sse-data', (_event, message) => {
  streams.get(message.id)?.onMessage(message.data);
});

ipcRenderer.on('ssui:sse-error', (_event, message) => {
  streams.get(message.id)?.onError(new Error(message.error));
});

contextBridge.exposeInMainWorld('ssuiDesktop', {
  request: request => ipcRenderer.invoke('ssui:request', request),
  login: request => ipcRenderer.invoke('ssui:login', request),
  bootstrap: request => ipcRenderer.invoke('ssui:bootstrap', request),
  logout: request => ipcRenderer.invoke('ssui:logout', request),
  trustCertificate: request => ipcRenderer.invoke('ssui:trust-certificate', request),
  onCertificateError: callback => {
    const listener = (_event, value) => callback(value);
    ipcRenderer.on('ssui:certificate-error', listener);
    return () => ipcRenderer.removeListener('ssui:certificate-error', listener);
  },
  openEventStream: ({ url, onMessage, onError = console.error }) => {
    const id = crypto.randomUUID();
    streams.set(id, { onMessage, onError });
    ipcRenderer.invoke('ssui:sse-start', { id, url }).catch(error => onError(error));
    return {
      close: () => {
        streams.delete(id);
        ipcRenderer.send('ssui:sse-stop', id);
      }
    };
  }
});
