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
    Vulcan2: {
        conditions: ['VulcanDefault', 'VulcanBrutal'],
        locations: ['VulcanSpawnVestaValley', 'VulcanSpawnEtnasFury', 'VulcanSpawnIxionsDemise', 'VulcanSpawnTitusReach']
    },
    Vulcan: {
        conditions: ['VulcanDefault', 'VulcanBrutal'],
        locations: ['VulcanSpawnVestaValley', 'VulcanSpawnEtnasFury', 'VulcanSpawnIxionsDemise', 'VulcanSpawnTitusReach']
    },
    Venus: {
        conditions: ['VenusDefault', 'VulcanBrutal (yes, VULCAN Brutal!)'],
        locations: ['VenusSpawnGaiaValley', 'VenusSpawnDaisyValley', 'VenusSpawnFaithValley', 'VenusSpawnDuskValley']
    }
};

const validDifficulties = ['Creative', 'Easy', 'Normal', 'Stationeer'];

// Get form elements
const worldIdInput = document.getElementById('WorldID');
const difficultyInput = document.getElementById('Difficulty');
const startConditionInput = document.getElementById('StartCondition');
const startLocationInput = document.getElementById('StartLocation');
const fillHintWrapper = document.getElementById('fill-hint-wraper');

// Function to check if New Terrain tab is active
function isNewTerrainTabActive() {
    const terrainButton = document.querySelector('button.section-nav-button[data-section="terrain-settings"].active');
    return !!terrainButton;
}

// Function to check if one but not all fields are filled
function shouldShowFillHint() {
    const hasDifficulty = difficultyInput.value && difficultyInput.value.trim() !== '';
    const hasStartCondition = startConditionInput.value && startConditionInput.value.trim() !== '';
    const hasStartLocation = startLocationInput.value && startLocationInput.value.trim() !== '';

    // Don't show hint for these specific combinations
    if (hasDifficulty && !hasStartCondition && !hasStartLocation) {
        return false; // Only difficulty filled
    }
    if (hasDifficulty && hasStartCondition && !hasStartLocation) {
        return false; // Difficulty and Start Condition filled
    }
    if (hasDifficulty && hasStartCondition && hasStartLocation) {
        return false; // All fields filled
    }

    // Show hint for all other combinations where at least one field is filled
    const filledInputs = [hasDifficulty, hasStartCondition, hasStartLocation].filter(Boolean).length;
    return filledInputs > 0;
}

// Function to toggle fill hint visibility
function toggleFillHint() {
    if (shouldShowFillHint()) {
        fillHintWrapper.style.display = 'flex';
    } else {
        fillHintWrapper.style.display = 'none';
    }
}

// Function to validate inputs
function validateInputs() {
    const selectedWorld = worldIdInput.value;
    const worldConfig = worldConfigs[selectedWorld];
    const inputs = [worldIdInput, difficultyInput, startConditionInput, startLocationInput];

    // Reset all validation states
    inputs.forEach(input => {
        input.classList.remove('invalid');
        const infoDiv = input.nextElementSibling;
        if (infoDiv && infoDiv.classList.contains('input-info')) {
            // Restore original text if it exists, otherwise keep current text
            infoDiv.textContent = infoDiv.getAttribute('data-original-text') || infoDiv.textContent;
        }
    });

    // Validate difficulty
    if (!difficultyInput.value || difficultyInput.value.trim() === '') {
        updateInfoText(difficultyInput, validDifficulties.join(', '));
    } else if (!validDifficulties.includes(difficultyInput.value)) {
        difficultyInput.classList.add('invalid');
        updateInfoText(difficultyInput, `❌${validDifficulties.join(', ')}`);
    }

    if (worldConfig) {
        // Validate start condition
        if (!startConditionInput.value || startConditionInput.value.trim() === '') {
            updateInfoText(startConditionInput, `${selectedWorld}: ${worldConfig.conditions.join(', ')}`);
        } else if (!worldConfig.conditions.includes(startConditionInput.value)) {
            startConditionInput.classList.add('invalid');
            updateInfoText(startConditionInput, `❌${selectedWorld}: ${worldConfig.conditions.join(', ')}`);
        }

        // Validate start location
        if (!startLocationInput.value || startLocationInput.value.trim() === '') {
            updateInfoText(startLocationInput, `${selectedWorld}: ${worldConfig.locations.join(', ')}`);
        } else if (!worldConfig.locations.includes(startLocationInput.value)) {
            startLocationInput.classList.add('invalid');
            updateInfoText(startLocationInput, `❌${selectedWorld}: ${worldConfig.locations.join(', ')}`);
        }
    }

    // Toggle fill hint visibility
    toggleFillHint();
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
worldIdInput.addEventListener('input', validateInputs);
difficultyInput.addEventListener('input', validateInputs);
startConditionInput.addEventListener('input', validateInputs);
startLocationInput.addEventListener('input', validateInputs);

// Initialize validation and fill hint on page load
document.addEventListener('DOMContentLoaded', () => {
    validateInputs();
    toggleFillHint();
});