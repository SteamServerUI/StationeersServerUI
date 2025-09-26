// Validation configuration object
const worldConfigs = {
    Lunar: {
        conditions: ['DefaultStart', 'Brutal'],
        locations: ['LunarSpawnCraterVesper', 'LunarSpawnMontesUmbrarum', 'LunarSpawnCraterNox', 'LunarSpawnMonsArcanus']
    },
    Mars2: {
        conditions: ['DefaultStart', 'Brutal'],
        locations: ['MarsSpawnCanyonOverlook', 'MarsSpawnButchersFlat', 'MarsSpawnFindersCanyon', 'MarsSpawnHellasCrags', 'MarsSpawnDonutFlats']
    },
    Europa3: {
        conditions: ['EuropaDefault', 'EuropaBrutal'],
        locations: ['EuropaSpawnIcyBasin', 'EuropaSpawnGlacialChannel', 'EuropaSpawnBalgatanPass', 'EuropaSpawnFrigidHighlands', 'EuropaSpawnTyreValley']
    },
    MimasHerschel: {
        conditions: ['MimasDefault', 'MimasBrutal'],
        locations: ['MimasSpawnCentralMesa', 'MimasSpawnHarrietCrater', 'MimasSpawnCraterField', 'MimasSpawnDustBowl']
    },
    Vulcan: {
        conditions: ['VulcanDefault', 'VulcanBrutal'],
        locations: ['VulcanSpawnVestaValley', 'VulcanSpawnEtnasFury', 'VulcanSpawnIxionsDemise', 'VulcanSpawnTitusReach']
    },
    Venus: {
        conditions: ['VenusDefault', 'VulcanBrutal'],
        locations: ['VenusSpawnGaiaValley', 'VenusSpawnDaisyValley', 'VenusSpawnFaithValley', 'VenusSpawnDuskValley']
    }
};

const validDifficulties = ['Creative', 'Easy', 'Normal', 'Stationeer'];

// Get form elements
const worldIdInput = document.getElementById('WorldID');
const difficultyInput = document.getElementById('Difficulty');
const startConditionInput = document.getElementById('StartCondition');
const startLocationInput = document.getElementById('StartLocation');

// Function to validate inputs
function validateInputs() {
    const selectedWorld = worldIdInput.value;
    const worldConfig = worldConfigs[selectedWorld];

    // Reset all validation states
    [worldIdInput, difficultyInput, startConditionInput, startLocationInput].forEach(input => {
        input.classList.remove('invalid');
        const infoDiv = input.nextElementSibling;
        if (infoDiv && infoDiv.classList.contains('input-info')) {
            infoDiv.textContent = infoDiv.getAttribute('data-original-text') || infoDiv.textContent;
        }
    });

    // Validate difficulty
    if (!validDifficulties.includes(difficultyInput.value)) {
        difficultyInput.classList.add('invalid');
        updateInfoText(difficultyInput, `Valid options: ${validDifficulties.join(', ')}`);
    }

    if (worldConfig) {
        // Validate start condition
        if (!worldConfig.conditions.includes(startConditionInput.value)) {
            startConditionInput.classList.add('invalid');
            updateInfoText(startConditionInput, `Valid options for ${selectedWorld}: ${worldConfig.conditions.join(', ')}`);
        }

        // Validate start location
        if (!worldConfig.locations.includes(startLocationInput.value)) {
            startLocationInput.classList.add('invalid');
            updateInfoText(startLocationInput, `Valid options for ${selectedWorld}: ${worldConfig.locations.join(', ')}`);
        }
    } else if (selectedWorld) {
        worldIdInput.classList.add('invalid');
        updateInfoText(worldIdInput, `Valid worlds: ${Object.keys(worldConfigs).join(', ')}`);
    }
}

function updateInfoText(input, text) {
    const infoDiv = input.nextElementSibling;
    if (infoDiv && infoDiv.classList.contains('input-info')) {
        if (!infoDiv.getAttribute('data-original-text')) {
            infoDiv.setAttribute('data-original-text', infoDiv.textContent);
        }
        infoDiv.textContent = text;
    }
}

// Add event listeners
worldIdInput.addEventListener('change', validateInputs);
difficultyInput.addEventListener('change', validateInputs);
startConditionInput.addEventListener('change', validateInputs);
startLocationInput.addEventListener('change', validateInputs);

// Initialize validation on page load
document.addEventListener('DOMContentLoaded', validateInputs);