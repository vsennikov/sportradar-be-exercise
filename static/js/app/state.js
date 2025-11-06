import { DEFAULT_PAGINATION } from '../shared/constants.js';

export const state = {
    currentPage: DEFAULT_PAGINATION.PAGE,
    currentLimit: DEFAULT_PAGINATION.LIMIT,
    currentDateFilter: "",
    currentSportFilter: ""
};
