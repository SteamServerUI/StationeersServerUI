// /static/console-manager.js

// Detection events streaming
function fetchDetectionEvents() {
    const maxMessages = 500;
    const detectionConsole = document.getElementById('detection-console');
    
    const connect = () => {
        detectionEventSource = new EventSource('/events');
        
        detectionEventSource.onmessage = event => {
            const message = document.createElement('div');
            message.className = `detection-event ${getEventClassName(event.data)}`;
            
            const timestamp = document.createElement('span');
            timestamp.className = 'event-timestamp';
            timestamp.textContent = `${new Date().toLocaleTimeString()}: `;
            
            const content = document.createElement('span');
            content.textContent = event.data;
            
            message.append(timestamp, content);
            detectionConsole.appendChild(message);
            
            while (detectionConsole.childElementCount > maxMessages) {
                detectionConsole.firstChild.remove();
            }
            detectionConsole.scrollTop = detectionConsole.scrollHeight;
            
            const detectionTab = document.getElementById('detection-tab');
            if (!detectionTab.classList.contains('active')) {
                const tabButton = document.querySelector('.tab-button[onclick*="detection-tab"]');
                tabButton.classList.add('notification');
                setTimeout(() => tabButton.classList.remove('notification'), 3000);
            }
        };
        
        detectionEventSource.onopen = () => console.log("Detection events stream connected");
        
        detectionEventSource.onerror = () => {
            console.error("Detection events stream disconnected");
            detectionEventSource.close();
            detectionEventSource = null;
            if (window.location.pathname === '/') {
                setTimeout(connect, 2000);
            }
        };
    };
    connect();
}

// Console initialization with SSE stream setup
function handleConsole() {
    const consoleElement = document.getElementById('console');
    consoleElement.innerHTML = '';
    const bootTitle = "Interface initializing...";
    const bootCompleteMessage = "Interface ready.🎮 Happy gaming! 🎮";
    const bugChance = Math.random();
    const bugMessage = "ERROR: Nuclear parts in airflow detected! Initiating repair sequence...";
    
    const funMessages = [
        "Calibrating quantum flux capacitors...",
        "Initializing player happiness modules...",
        "Checking for monsters under the server...",
        "Brewing coffee for the CPU...",
        "Charging laser sharks...",
        "Teaching AI to say 'please' and 'thank you'...",
        "Polishing pixels to a mirror shine...",
        "Convincing electrons to flow in the right direction...",
        "Rebooting atmospheric systems for the 17th time...",
        "Attempting to locate your body after that last airlock malfunction...",
        "Converting oxygen to errors at alarming efficiency...",
        "Persuading physics engine to acknowledge gravity exists...",
        "Calculating ways your base will catastrophically depressurize...",
        "Optimizing unity garbage collection (good luck with that)...",
        "Aligning planetary rotation with server tick rate...",
        "Patching holes in space-time continuum and your habitat...",
        "Convincing solar panels that 'sun' is not just a theoretical concept...",
        "Negotiating peace treaty between logic circuits and the laws of thermodynamics...",
        "Compressing atmosphere until your CPU begs for mercy...",
        "Measuring distance between you and nearest fatal bug...",
        "Attempting to explain 'pipe networks' to confused server hamsters...",
        "Calculating probability of survival (spoiler: it's low)...",
        "Wrangling rogue Unity instances back into containment...",
        "Sacrificing RAM to the gods of stable framerates...",
        "Convincing electrons to flow in the right direction... nope, the power grid's borked.",
        "Patching hull breaches with duct tape and prayers...",
        "Recalculating O2 levels... wait, why is it all CO2 now?",
        "Spinning up the fabricator... hope it doesn't eat the server this time.",
        "Debugging Unity physics... object launched into orbit, send help.",
        "Warming up the furnace... or just setting the base on fire, 50/50 shot.",
        "Rerouting pipes... because who needs logical fluid dynamics anyway?",
        "Loading terrain... oh look, it's floating 3 meters above the ground again.",
        "Processing ore... into a fine paste of lag and despair.",
        "Stabilizing frame rate... lol, just kidding, welcome to 12 FPS city.",
        "Checking for updates... new bug introduced, feature still broken!",
        "Assembling solar tracker... now it's tracking the admin instead.",
        "Balancing gas mixtures... kaboom imminent, run you fool!",
        "Spoiler: object reference not set to an instance of an object, lol",
        "Fun fact: SSUI was originally a simple powershell script",
        "Convincing server that 'out of memory' is just a state of mind.",
        "Moo, Moo! I'm a cow!",
        "Welcome home, Sir!"
    ];

    const addMessage = (text, color, style = 'normal') => {
        const div = document.createElement('div');
        div.textContent = text;
        div.style.color = color;
        div.style.fontStyle = style;
        consoleElement.appendChild(div);
        consoleElement.scrollTop = consoleElement.scrollHeight;
    };

    // Dynamically create SSCM command input
    const createCommandInput = async () => {
        try {
            // Make API call to check if SSCM is enabled
            const response = await fetch('/api/v2/SSCM/enabled', {
                method: 'GET',
                headers: {
                    'Accept': 'application/json'
                }
            });
    
            // If status is not 200, exit the function
            if (response.status !== 200) {
                console.log('SSCM is not enabled, status:', response.status);
                return;
            }
    
            // Proceed to create command input UI if status is 200
            console.log("Creating command input...");
            const commandContainer = document.createElement('div');
            commandContainer.className = 'sscm-command-container';
            
            const prompt = document.createElement('span');
            prompt.className = 'prompt';
            prompt.textContent = '>';
    
            const input = document.createElement('input');
            input.id = 'sscm-command-input';
            input.type = 'text';
            input.placeholder = 'Enter command..';
            input.setAttribute('autocomplete', 'off');
            
            const suggestions = document.createElement('div');
            suggestions.id = 'sscm-autocomplete-suggestions';
            suggestions.className = 'sscm-suggestions';
            
            commandContainer.append(prompt, input, suggestions);
            consoleElement.appendChild(commandContainer);
        } catch (error) {
            console.error('Error checking SSCM enabled status:', error);
            return; // Exit on error
        }
    };

    // Start with initializing message
    typeTextWithCallback(consoleElement, bootTitle, 30, () => {
        // Show two funny messages while connecting
        const messageIndex1 = Math.floor(Math.random() * funMessages.length);
        addMessage(funMessages[messageIndex1], '#0af', 'italic');

        let messageIndex2;
        do {
            messageIndex2 = Math.floor(Math.random() * funMessages.length);
        } while (messageIndex2 === messageIndex1);
        addMessage(funMessages[messageIndex2], '#0af', 'italic');

        // Set up the persistent console stream
        outputEventSource = new EventSource('/console');
        
        // Persistent message handler
        outputEventSource.onmessage = event => {
            const message = document.createElement('div');
            message.textContent = event.data;
            consoleElement.insertBefore(message, consoleElement.querySelector('.sscm-command-container')); // Insert before input
            // Auto-scroll only if at bottom
            if (consoleElement.scrollTop + consoleElement.clientHeight >= consoleElement.scrollHeight - 10) {
                consoleElement.scrollTop = consoleElement.scrollHeight;
            }
        };

        outputEventSource.onopen = () => {
            console.log("Console stream connected");
            finishInitialization();
        };

        outputEventSource.onerror = () => {
            console.error("Console stream disconnected");
            outputEventSource.close();
            outputEventSource = null;
            addMessage("Warning: Console stream unavailable. Retrying...", '#ff0');
            if (window.location.pathname === '/') {
                setTimeout(() => {
                    if (!outputEventSource) {
                        // Re-run setup to reconnect
                        consoleElement.innerHTML = ''; // Clear console for fresh start
                        handleConsole();
                    }
                }, 5000);
            }
        };
    });

    function finishInitialization() {
        if (bugChance < 0.05) {
            addMessage(bugMessage, 'red');
            setTimeout(() => {
                addMessage("Repair complete. Continuing initialization...", 'green');
                completeBoot();
            }, 1000);
        } else {
            completeBoot();
        }
    }

    function completeBoot() {
        setTimeout(() => {
            createCommandInput(); // Add input after boot
            addMessage(bootCompleteMessage, '#0f0');
            //addMessage("StationeersServerUI is becoming SteamServerUI!", '#ff4500');
            //addMessage("Please mind the New Terrain System warning below", '#ff4500');
            consoleElement.scrollTop = consoleElement.scrollHeight;
        }, 500);
    }
}

function setupLogStreams({ consoleId, streamUrls, maxMessages, messageClass }) {
    const consoleElement = document.getElementById(consoleId);
    if (!consoleElement) {
        console.error(`Console element with ID '${consoleId}' not found.`);
        return;
    }

    // Clear the console initially
    consoleElement.innerHTML = '';

    const connectStream = (streamUrl) => {
        const eventSource = new EventSource(streamUrl);

        eventSource.onmessage = event => {
            const message = document.createElement('div');
            let finalClass = messageClass;

            // Check event.data for specific log levels and modify class
            if (event.data.includes('/INFO')) {
                finalClass += '-info';
            } else if (event.data.includes('/WARN')) {
                finalClass += '-warn';
            } else if (event.data.includes('/ERROR')) {
                finalClass += '-error';
            }

            message.classList.add("log-console-element", finalClass);

            const content = document.createElement('span');
            content.textContent = event.data;

            message.append(content);
            consoleElement.appendChild(message);

            // Limit the number of messages
            while (consoleElement.childElementCount > maxMessages) {
                consoleElement.firstChild.remove();
            }

            // Auto-scroll to the bottom
            consoleElement.scrollTop = consoleElement.scrollHeight;
        };

        eventSource.onopen = () => {
            console.log(`Stream ${streamUrl} connected for console ${consoleId}`);
        };

        eventSource.onerror = () => {
            console.error(`Stream ${streamUrl} disconnected for console ${consoleId}`);
            eventSource.close();
            if (window.location.pathname === '/') {
                setTimeout(() => connectStream(streamUrl), 5000); // Reconnect after 5 seconds
            }
        };
    };

    // Connect to all provided stream URLs
    streamUrls.forEach(url => connectStream(url));
}