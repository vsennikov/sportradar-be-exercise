import { state } from './state.js';
import { fetchEvents, fetchSports } from './api.js';
import { renderTable, renderPagination, renderSportOptions } from './render.js';
import { getElement } from '../shared/dom-utils.js';
import { logError } from '../shared/error-handler.js';

export async function fetchEventsHandler(pageNumber = null) {
    if (pageNumber !== null) {
        state.currentPage = pageNumber;
    }
    
    try {
        const data = await fetchEvents();
        if (data) {
            renderTable(data.events);
            renderPagination(data.pagination, fetchEventsHandler);
        }
    } catch (error) {
        logError('fetchEventsHandler', error);
    }
}

function handleFilterSubmit(e) {
    e.preventDefault();

    const dateInput = getElement('date-filter');
    const sportInput = getElement('sport-filter');

    state.currentPage = 1;
    state.currentDateFilter = dateInput.value;
    state.currentSportFilter = sportInput.value;

    fetchEventsHandler();
}

function handleClearFilters() {
    const dateInput = getElement('date-filter');
    const sportInput = getElement('sport-filter');

    state.currentPage = 1;
    state.currentDateFilter = "";
    state.currentSportFilter = "";

    dateInput.value = "";
    sportInput.value = "";

    fetchEventsHandler();
}

function setupEventHandlers() {
    const filterForm = getElement('filter-form');
    const clearButton = getElement('filter-clear');

    filterForm.addEventListener('submit', handleFilterSubmit);
    clearButton.addEventListener('click', handleClearFilters);
}

async function loadSportsFilter() {
    try {
        const sports = await fetchSports();
        renderSportOptions(sports);
    } catch (error) {
        logError('loadSportsFilter', error);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    setupEventHandlers();
    loadSportsFilter();
    fetchEventsHandler();
});

