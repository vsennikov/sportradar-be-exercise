import { API_ENDPOINTS } from '../shared/constants.js';
import { apiGet, apiPost } from '../shared/api-client.js';

export async function fetchSports() {
    return apiGet(API_ENDPOINTS.SPORTS);
}

export async function fetchTeams() {
    return apiGet(API_ENDPOINTS.TEAMS);
}

export async function fetchVenues() {
    return apiGet(API_ENDPOINTS.VENUES);
}

export async function createEvent(eventData) {
    return apiPost(API_ENDPOINTS.EVENTS, eventData);
}

export async function loadDropdownData() {
    const [sportsList, teamsList, venuesList] = await Promise.all([
        fetchSports(),
        fetchTeams(),
        fetchVenues()
    ]);
    
    return { sportsList, teamsList, venuesList };
}

