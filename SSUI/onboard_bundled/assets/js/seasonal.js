(() => {
    const seasons = [
        {
            id: 'halloween',
            from: 1020,
            through: 1105,
            celestialBody: '/static/seasonal/halloween/celestial_halloween.webp',
            decorations: [
                '/static/seasonal/halloween/pumpkin.webp',
                '/static/seasonal/halloween/pumpkin.webp',
                '/static/seasonal/halloween/pumpkin.webp',
                '/static/seasonal/halloween/pumpkin.webp'
            ]
        },
        {
            id: 'christmas',
            from: 1201,
            through: 110,
            banner: '/static/seasonal/christmas/stationeers-winter.webp',
            celestialBody: '/static/seasonal/christmas/celestial_christmas.webp',
            decorations: [
                '/static/seasonal/christmas/c1.webp',
                '/static/seasonal/christmas/c2.webp',
                '/static/seasonal/christmas/c3.webp',
                '/static/seasonal/christmas/c4.webp'
            ]
        }
    ];

    const optOutCookie = 'ssuiSeasonOff';
    const optOutSeconds = 30 * 24 * 60 * 60;
    let applied = null;
    let originalBanner = null;

    function today(date) {
        return (date.getMonth() + 1) * 100 + date.getDate();
    }

    function isActive(season, date) {
        const value = today(date);
        if (season.from <= season.through) {
            return value >= season.from && value <= season.through;
        }
        return value >= season.from || value <= season.through;
    }

    function addDecorations(season) {
        const layer = document.createElement('div');
        layer.id = 'seasonal-layer';
        layer.className = `seasonal-layer seasonal-layer--${season.id}`;
        layer.setAttribute('aria-hidden', 'true');

        season.decorations.forEach((decoration, index) => {
            const item = decoration.startsWith('/') ? document.createElement('img') : document.createElement('span');
            item.className = `seasonal-decoration seasonal-decoration--${index + 1}`;
            if (item instanceof HTMLImageElement) {
                item.src = decoration;
                item.alt = '';
                item.draggable = false;
            } else {
                item.textContent = decoration;
            }
            layer.appendChild(item);
        });

        document.body.appendChild(layer);
    }

    function applySeason(season) {
        document.body.classList.add('seasonal-active', `seasonal-${season.id}`);
        document.documentElement.dataset.season = season.id;
        addDecorations(season);

        if (season.banner) {
            const banner = document.getElementById('banner');
            if (banner) {
                originalBanner = banner.getAttribute('src');
                banner.src = season.banner;
            }
        }
        applied = season;
    }

    function removeSeason() {
        if (!applied) return;
        document.body.classList.remove('seasonal-active', `seasonal-${applied.id}`, 'seasonal-paused');
        delete document.documentElement.dataset.season;
        document.getElementById('seasonal-layer')?.remove();
        const banner = document.getElementById('banner');
        if (banner && originalBanner) banner.src = originalBanner;
        applied = null;
    }

    function optedOut(season) {
        return document.cookie.split('; ').includes(`${optOutCookie}=${season.id}`);
    }

    function setEnabled(enabled) {
        if (!active) return;
        if (enabled) {
            document.cookie = `${optOutCookie}=; Max-Age=0; Path=/; SameSite=Lax`;
            if (!applied) applySeason(active);
        } else {
            document.cookie = `${optOutCookie}=${active.id}; Max-Age=${optOutSeconds}; Path=/; SameSite=Lax`;
            removeSeason();
        }
        publish();
    }

    function publish() {
        window.SSUISeasonal = {
            active: active?.id || null,
            enabled: Boolean(applied),
            celestialBody: applied?.celestialBody || null,
            setEnabled
        };
        document.dispatchEvent(new CustomEvent('ssui-season'));
    }

    window.addEventListener('focus', () => document.body.classList.remove('seasonal-paused'));
    window.addEventListener('blur', () => applied && document.body.classList.add('seasonal-paused'));
    document.addEventListener('visibilitychange', () => {
        document.body.classList.toggle('seasonal-paused', Boolean(applied) && document.hidden);
    });

    const preview = new URLSearchParams(window.location.search).get('season');
    const previewed = seasons.find(season => season.id === preview);
    const active = previewed || seasons.find(season => isActive(season, new Date()));
    if (active && (previewed || !optedOut(active))) applySeason(active);
    publish();
})();
