/**
 * Log error with context
 * @param {string} context - Error context (e.g., "fetchEvents")
 * @param {Error} error - Error object
 */
export function logError(context, error) {
    console.error(`[${context}]`, error);
}

/**
 * Show user-friendly error message
 * @param {HTMLElement} container - Container element to show error
 * @param {string} message - Error message
 * @param {string} className - Optional CSS class
 */
export function showError(container, message, className = 'error') {
    container.innerHTML = `<mark class="pico-color-red-550">Error: ${message}</mark>`;
}

/**
 * Show success message
 * @param {HTMLElement} container - Container element to show message
 * @param {string} message - Success message
 */
export function showSuccess(container, message) {
    container.innerHTML = `<mark>${message}</mark>`;
}

