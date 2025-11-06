/**
 * Safely get element by ID
 * @param {string} id - Element ID
 * @returns {HTMLElement}
 * @throws {Error} - If element not found
 */
export function getElement(id) {
    const element = document.getElementById(id);
    if (!element) {
        throw new Error(`Element with id "${id}" not found`);
    }
    return element;
}

/**
 * Create a table cell element
 * @param {string} text - Cell text content
 * @param {Object} attributes - Attributes to set
 * @returns {HTMLTableCellElement}
 */
export function createTableCell(text, attributes = {}) {
    const cell = document.createElement('td');
    cell.textContent = text;
    
    Object.entries(attributes).forEach(([key, value]) => {
        if (key === 'style' && typeof value === 'object') {
            Object.assign(cell.style, value);
        } else {
            cell.setAttribute(key, value);
        }
    });
    
    return cell;
}

/**
 * Create a table row element
 * @param {Array<HTMLElement>} cells - Array of cell elements
 * @returns {HTMLTableRowElement}
 */
export function createTableRow(cells) {
    const row = document.createElement('tr');
    cells.forEach(cell => row.appendChild(cell));
    return row;
}

/**
 * Clear element content
 * @param {string|HTMLElement} element - Element ID or element itself
 */
export function clearElement(element) {
    const el = typeof element === 'string' ? getElement(element) : element;
    el.innerHTML = '';
}

/**
 * Set element text content
 * @param {string|HTMLElement} element - Element ID or element itself
 * @param {string} text - Text content
 */
export function setText(element, text) {
    const el = typeof element === 'string' ? getElement(element) : element;
    el.textContent = text;
}

