function setupPasswordProtection() {
    const passwordPrompt = document.getElementById('password-prompt');
    const passwordForm = document.getElementById('password-form');
    const passwordInput = document.getElementById('password');
    const passwordError = document.getElementById('password-error');
    const appContainer = document.getElementById('app-container');

    passwordForm.addEventListener('submit', (e) => {
        e.preventDefault();
        const correctPassword = "admin123"; 

        if (passwordInput.value === correctPassword) {
            passwordPrompt.style.display = 'none';
            appContainer.style.display = 'block';            
            initializeApp();
        } else {
            passwordError.textContent = 'Wrong password. Please try again.';
        }
    });
}

function initializeApp() {
    loadDataForDropdowns();
    setupEventForm();
}

async function loadDataForDropdowns() {
    try {
        const [sportsRes, teamsRes, venuesRes] = await Promise.all([
            fetch('/api/v1/sports'),
            fetch('/api/v1/teams'),
            fetch('/api/v1/venues')
        ]);
        const sportsList = await sportsRes.json();
        const teamsList = await teamsRes.json();
        const venuesList = await venuesRes.json();
        populateSelect('sport', sportsList, 'Select Sport');
        populateSelect('home_team', teamsList, 'Select Home Team');
        populateSelect('away_team', teamsList, 'Select Away Team');
        populateSelect('venue', venuesList, 'Select Venue (Optional)');
    } catch (error) {
        console.error('Failed to load dropdown data:', error);
        alert('Failed to load initial data. Please check console.');
    }
}

function populateSelect(selectId, items, defaultOptionText) {
    const select = document.getElementById(selectId);
    select.innerHTML = '';
    const defaultOption = document.createElement('option');
    defaultOption.value = "";
    if (select.required) {
        defaultOption.disabled = true;
        defaultOption.selected = true;
    }
    defaultOption.textContent = defaultOptionText;
    select.appendChild(defaultOption);
    items.forEach(item => {
        const option = document.createElement('option');
        option.value = item.id;
        option.textContent = item.name;
        select.appendChild(option);
    });
}

function setupEventForm() {
    const eventForm = document.getElementById('create-event-form');
    eventForm.addEventListener('submit', handleFormSubmit);
}

async function handleFormSubmit(e) {
    e.preventDefault();
    const form = e.target;
    const messageDiv = document.getElementById('form-message');
    const eventData = {
        event_datetime: form.datetime.value + ":00Z", 
        sport_id: parseInt(form.sport.value),
        home_team_id: parseInt(form.home_team.value),
        away_team_id: parseInt(form.away_team.value),
    };
    if (form.description.value) {
        eventData.description = form.description.value;
    }
    if (form.venue.value) {
        eventData.venue_id = parseInt(form.venue.value);
    }
    try {
        const response = await fetch('/api/v1/events', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(eventData)
        });
        if (!response.ok) {
            const err = await response.json();
            throw new Error(err.error);
        }
        const data = await response.json();
        messageDiv.innerHTML = `<mark>Successfully created event with ID: ${data.id}</mark>`;
        form.reset();
        form.sport.selectedIndex = 0;
        form.home_team.selectedIndex = 0;
        form.away_team.selectedIndex = 0;
        form.venue.selectedIndex = 0;
    } catch (error) {
        console.error('Error creating event:', error);
        messageDiv.innerHTML = `<mark class="pico-color-red-550">Error: ${error.message}</mark>`;
    }
}

document.addEventListener('DOMContentLoaded', () => {
    setupPasswordProtection();
});