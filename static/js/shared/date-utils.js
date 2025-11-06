/**
 * Convert datetime-local input value to ISO string
 * @param {string} datetimeLocalValue - Value from datetime-local input
 * @returns {string} - ISO string representation
 */
export function convertDatetimeLocalToISO(datetimeLocalValue) {
    if (!datetimeLocalValue) {
        throw new Error('Datetime value is required');
    }
    
    const date = new Date(datetimeLocalValue);
    
    if (isNaN(date.getTime())) {
        throw new Error('Invalid datetime value');
    }
    
    return date.toISOString();
}

