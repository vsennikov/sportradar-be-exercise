import { getElement, createTableCell, createTableRow, clearElement } from '../shared/dom-utils.js';

const TABLE_BODY_ID = 'event-list-body';
const PAGINATION_CONTAINER_ID = 'pagination-controls';

export function renderSportOptions(sports) {
    const select = getElement('sport-filter');
    clearElement(select);
    
    const defaultOption = document.createElement('option');
    defaultOption.value = '';
    defaultOption.textContent = 'All Sports';
    select.appendChild(defaultOption);
    
    sports.forEach(sport => {
        const option = document.createElement('option');
        option.value = sport.id;
        option.textContent = sport.name;
        select.appendChild(option);
    });
}

export function renderTable(events) {
    const tableBody = getElement(TABLE_BODY_ID);
    clearElement(tableBody);

    if (events.length === 0) {
        const emptyRow = createTableRow([
            createTableCell('No events found matching your criteria.', {
                colspan: '6',
                style: { textAlign: 'center' }
            })
        ]);
        tableBody.appendChild(emptyRow);
        return;
    }

    events.forEach(event => {
        const row = renderEventRow(event);
        tableBody.appendChild(row);
    });
}

function renderEventRow(event) {
    const cells = [
        createTableCell(event.sport.name),
        createTableCell(`${event.home_team.name} vs ${event.away_team.name}`),
        createTableCell(
            event.home_score != null 
                ? `${event.home_score} - ${event.away_score}` 
                : 'N/A'
        ),
        createTableCell(event.venue ? event.venue.name : 'TBD'),
        createTableCell(new Date(event.event_datetime).toLocaleString()),
        createTableCell(event.description || '')
    ];

    return createTableRow(cells);
}

export function renderPagination(pagination, fetchEventsHandler) {
    const container = getElement(PAGINATION_CONTAINER_ID);
    clearElement(container);

    for (let i = 1; i <= pagination.total_pages; i++) {
        const button = renderPaginationButton(i, pagination.current_page, fetchEventsHandler);
        container.appendChild(button);
    }
}

function renderPaginationButton(pageNumber, activePage, fetchEventsHandler) {
    const button = document.createElement('button');
    button.textContent = pageNumber;

    if (pageNumber === activePage) {
        button.setAttribute('aria-current', 'page');
    } else {
        button.addEventListener('click', () => {
            fetchEventsHandler(pageNumber);
        });
    }

    return button;
}

export function renderTableLoading() {
    const tableBody = getElement(TABLE_BODY_ID);
    const loadingRow = createTableRow([
        createTableCell('Loading...', {
            colspan: '6',
            style: { textAlign: 'center' }
        })
    ]);
    clearElement(tableBody);
    tableBody.appendChild(loadingRow);
}

export function renderTableError() {
    const tableBody = getElement(TABLE_BODY_ID);
    const errorRow = createTableRow([
        createTableCell('Error loading events. Please try again later.', {
            colspan: '6',
            style: { textAlign: 'center', color: 'red' }
        })
    ]);
    clearElement(tableBody);
    tableBody.appendChild(errorRow);
}

