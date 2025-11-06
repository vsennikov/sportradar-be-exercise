import { state } from './state.js';
import { API_ENDPOINTS } from '../shared/constants.js';
import { apiGet } from '../shared/api-client.js';
import { renderTableLoading, renderTableError } from './render.js';

export async function fetchSports() {
    return apiGet(API_ENDPOINTS.SPORTS);
}

export function buildEventsUrl() {
    const url = new URL(API_ENDPOINTS.EVENTS, window.location.origin);
    url.searchParams.set('page', state.currentPage.toString());
    url.searchParams.set('limit', state.currentLimit.toString());

    if (state.currentDateFilter) {
        url.searchParams.set('date_from', state.currentDateFilter);
    }

    if (state.currentSportFilter) {
        url.searchParams.set('sport_id', state.currentSportFilter);
    }

    return url.toString();
}

export async function fetchEvents() {
    const url = buildEventsUrl();
    renderTableLoading();

    try {
        return await apiGet(url);
    } catch (error) {
        renderTableError();
        throw error;
    }
}

