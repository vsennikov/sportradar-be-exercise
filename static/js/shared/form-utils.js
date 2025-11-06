/**
 * Get form data as object
 * @param {HTMLFormElement} form - Form element
 * @returns {Object} - Form data object
 */
export function getFormData(form) {
    const formData = new FormData(form);
    const data = {};
    
    for (const [key, value] of formData.entries()) {
        data[key] = value;
    }
    
    return data;
}

/**
 * Reset form and set select elements to first option
 * @param {HTMLFormElement} form - Form element
 * @param {Array<string>} selectIds - Array of select element IDs to reset
 */
export function resetForm(form, selectIds = []) {
    form.reset();
    selectIds.forEach(id => {
        const select = form.querySelector(`#${id}`);
        if (select) {
            select.selectedIndex = 0;
        }
    });
}

/**
 * Validate required fields
 * @param {HTMLFormElement} form - Form element
 * @param {Array<string>} requiredFields - Array of required field names
 * @returns {Object} - { isValid: boolean, errors: Array<string> }
 */
export function validateForm(form, requiredFields = []) {
    const errors = [];
    
    requiredFields.forEach(fieldName => {
        const field = form.querySelector(`[name="${fieldName}"]`);
        if (!field || !field.value.trim()) {
            errors.push(`${fieldName} is required`);
        }
    });
    
    return {
        isValid: errors.length === 0,
        errors
    };
}

/**
 * Validate that two team IDs are different
 * @param {number} homeTeamId - Home team ID
 * @param {number} awayTeamId - Away team ID
 * @throws {Error} - If teams are the same
 */
export function validateDifferentTeams(homeTeamId, awayTeamId) {
    if (homeTeamId === awayTeamId) {
        throw new Error('Home team and away team must be different');
    }
}

/**
 * Parse integer from form field value
 * @param {HTMLFormElement} form - Form element
 * @param {string} fieldName - Field name
 * @returns {number} - Parsed integer
 * @throws {Error} - If value is invalid
 */
export function parseIntField(form, fieldName) {
    const field = form.querySelector(`[name="${fieldName}"]`);
    if (!field) {
        throw new Error(`Field "${fieldName}" not found`);
    }
    
    const value = parseInt(field.value);
    if (isNaN(value)) {
        throw new Error(`Invalid value for "${fieldName}"`);
    }
    
    return value;
}

