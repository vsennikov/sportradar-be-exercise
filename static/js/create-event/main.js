import { loadDropdownData, createEvent } from './api.js';
import { renderSelectOptions, renderSuccessMessage, renderErrorMessage } from './render.js';
import { getElement } from '../shared/dom-utils.js';
import { resetForm, validateDifferentTeams, parseIntField } from '../shared/form-utils.js';
import { convertDatetimeLocalToISO } from '../shared/date-utils.js';
import { logError } from '../shared/error-handler.js';

const SELECT_IDS = ['sport', 'home_team', 'away_team', 'venue'];

function buildEventData(form) {
    const homeTeamId = parseIntField(form, 'home_team');
    const awayTeamId = parseIntField(form, 'away_team');
    
    validateDifferentTeams(homeTeamId, awayTeamId);
    
    const eventDatetime = convertDatetimeLocalToISO(form.datetime.value);
    
    const eventData = {
        event_datetime: eventDatetime,
        sport_id: parseIntField(form, 'sport'),
        home_team_id: homeTeamId,
        away_team_id: awayTeamId,
    };
    
    if (form.description.value) {
        eventData.description = form.description.value;
    }
    
    if (form.venue.value) {
        eventData.venue_id = parseIntField(form, 'venue');
    }
    
    return eventData;
}

async function handleFormSubmit(e) {
    e.preventDefault();
    
    const form = e.target;
    const messageDiv = getElement('form-message');
    
    try {
        const eventData = buildEventData(form);
        const data = await createEvent(eventData);
        
        renderSuccessMessage(messageDiv, data.id);
        resetForm(form, SELECT_IDS);
    } catch (error) {
        logError('handleFormSubmit', error);
        renderErrorMessage(messageDiv, error.message);
    }
}

function setupEventForm() {
    const eventForm = getElement('create-event-form');
    eventForm.addEventListener('submit', handleFormSubmit);
}

async function initializeApp() {
    try {
        const { sportsList, teamsList, venuesList } = await loadDropdownData();
        
        renderSelectOptions('sport', sportsList, 'Select Sport');
        renderSelectOptions('home_team', teamsList, 'Select Home Team');
        renderSelectOptions('away_team', teamsList, 'Select Away Team');
        renderSelectOptions('venue', venuesList, 'Select Venue (Optional)');
        
        setupEventForm();
    } catch (error) {
        logError('initializeApp', error);
        alert('Failed to load initial data. Please check console.');
    }
}

document.addEventListener('DOMContentLoaded', () => {
    initializeApp();
});

