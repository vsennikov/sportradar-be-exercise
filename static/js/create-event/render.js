import { getElement, clearElement } from '../shared/dom-utils.js';

export function renderSelectOptions(selectId, items, defaultOptionText) {
    const select = getElement(selectId);
    const isRequired = select.required;
    
    clearElement(select);
    
    const defaultOption = document.createElement('option');
    defaultOption.value = "";
    defaultOption.textContent = defaultOptionText;
    
    if (isRequired) {
        defaultOption.disabled = true;
        defaultOption.selected = true;
    }
    
    select.appendChild(defaultOption);
    
    items.forEach(item => {
        const option = document.createElement('option');
        option.value = item.id;
        option.textContent = item.name;
        select.appendChild(option);
    });
}

export function renderSuccessMessage(messageDiv, eventId) {
    messageDiv.innerHTML = `<mark>Successfully created event with ID: ${eventId}</mark>`;
}

export function renderErrorMessage(messageDiv, errorMessage) {
    messageDiv.innerHTML = `<mark class="pico-color-red-550">Error: ${errorMessage}</mark>`;
}

export function renderLoadingState(selectId, loadingText) {
    const select = getElement(selectId);
    clearElement(select);
    
    const option = document.createElement('option');
    option.value = "";
    option.disabled = true;
    option.selected = true;
    option.textContent = loadingText;
    select.appendChild(option);
}
