const themes = [
  { id: 'sjalvlysande', name: 'Självlysande' },
  { id: 'skogsgron', name: 'Skogsgrön' },
  { id: 'norrskensdrom', name: 'Norrskensdröm' },
  { id: 'solkatt', name: 'Solkatt' }
];

let current = 'sjalvlysande';

function applyTheme(id) {
  current = themes.some(theme => theme.id === id) ? id : themes[0].id;
  document.documentElement.dataset.theme = current;
  localStorage.setItem('ssui-theme-v3', current);
  return current;
}
function initTheme() { return applyTheme(localStorage.getItem('ssui-theme-v3') || current); }
function nextTheme() {
  const index = themes.findIndex(theme => theme.id === current);
  return applyTheme(themes[(index + 1) % themes.length].id);
}
function getCurrentTheme() { return current; }
function getThemes() { return themes; }
export default { applyTheme, initTheme, nextTheme, getCurrentTheme, getThemes };
