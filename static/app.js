//
// --- 1. СТАН ДОДАТКУ ---
//
let currentPage = 1;
let currentLimit = 10;
let currentDateFilter = "";
let currentSportFilter = "";

//
// --- 2. ГОЛОВНА ФУНКЦІЯ: ОТРИМАННЯ ДАНИХ ---
//
function fetchEvents() {
    // Будуємо URL з усіма нашими параметрами стану
    let url = `/api/v1/events?page=${currentPage}&limit=${currentLimit}`;
    
    if (currentDateFilter) {
        url += `&date_from=${currentDateFilter}`;
    }
    if (currentSportFilter) {
        url += `&sport_id=${currentSportFilter}`;
    }

    // Показуємо стан "Завантаження..."
    renderTableLoading();

    fetch(url)
        .then(response => response.json())
        .then(data => {
            // Ми отримали дані. Тепер "малюємо" і таблицю, і кнопки
            renderTable(data.events);
            renderPagination(data.pagination);
        })
        .catch(error => {
            console.error('Error fetching events:', error);
            renderTableError();
        });
}

//
// --- 3. ФУНКЦІЇ РЕНДЕРИНГУ (МАЛЮВАННЯ) ---
//

// "Малює" таблицю
function renderTable(events) {
    const tableBody = document.getElementById('event-list-body');
    tableBody.innerHTML = ''; // Очищуємо

    if (events.length === 0) {
        tableBody.innerHTML = '<tr><td colspan="6" style="text-align: center;">No events found matching your criteria.</td></tr>';
        return;
    }

    events.forEach(event => {
        const row = document.createElement('tr');

        // 1. Sport
        const cellSport = document.createElement('td');
        cellSport.textContent = event.sport.name;
        row.appendChild(cellSport);

        // 2. Match (Команди)
        const cellMatch = document.createElement('td');
        cellMatch.textContent = `${event.home_team.name} vs ${event.away_team.name}`;
        row.appendChild(cellMatch);
        
        // 3. Score (Рахунок)
        const cellScore = document.createElement('td');
        // Перевіряємо, чи рахунок 'null'
        if (event.home_score != null) {
            cellScore.textContent = `${event.home_score} - ${event.away_score}`;
        } else {
            cellScore.textContent = 'N/A';
        }
        row.appendChild(cellScore);

        // 4. Venue (Стадіон)
        const cellVenue = document.createElement('td');
        cellVenue.textContent = event.venue ? event.venue.name : 'TBD';
        row.appendChild(cellVenue);

        // 5. Date & Time
        const cellTime = document.createElement('td');
        cellTime.textContent = new Date(event.event_datetime).toLocaleString();
        row.appendChild(cellTime);

        // 6. Description (Опис)
        const cellDesc = document.createElement('td');
        cellDesc.textContent = event.description ? event.description : ''; // Порожньо, якщо 'null'
        row.appendChild(cellDesc);

        tableBody.appendChild(row);
    });
}

// "Малює" кнопки пагінації
function renderPagination(pagination) {
    const container = document.getElementById('pagination-controls');
    container.innerHTML = ''; // Очищуємо

    for (let i = 1; i <= pagination.total_pages; i++) {
        const button = document.createElement('button');
        button.textContent = i;
        
        // Позначаємо поточну сторінку
        if (i === pagination.current_page) {
            button.setAttribute('aria-current', 'page');
        } else {
            // Додаємо обробник кліку для *інших* сторінок
            button.addEventListener('click', () => {
                currentPage = i; // Оновлюємо стан
                fetchEvents();   // Перезавантажуємо дані
            });
        }
        container.appendChild(button);
    }
}

// Показує стан "Завантаження..."
function renderTableLoading() {
    const tableBody = document.getElementById('event-list-body');
    tableBody.innerHTML = '<tr><td colspan="6" style="text-align: center;">Loading...</td></tr>';
}

// Показує помилку
function renderTableError() {
    const tableBody = document.getElementById('event-list-body');
    tableBody.innerHTML = '<tr><td colspan="6" style="text-align: center; color: red;">Error loading events. Please try again later.</td></tr>';
}

//
// --- 4. НАЛАШТУВАННЯ ОБРОБНИКІВ ПОДІЙ ---
//
document.addEventListener('DOMContentLoaded', () => {
    
    // 1. Отримуємо посилання на елементи форми
    const filterForm = document.getElementById('filter-form');
    const clearButton = document.getElementById('filter-clear');
    const dateInput = document.getElementById('date-filter');
    const sportInput = document.getElementById('sport-filter');

    // 2. Обробник для кнопки "Filter"
    filterForm.addEventListener('submit', (e) => {
        e.preventDefault(); // Забороняємо формі перезавантажувати сторінку
        
        // Оновлюємо глобальний стан
        currentPage = 1; // Завжди скидаємо на 1-шу сторінку при новому фільтрі
        currentDateFilter = dateInput.value;
        currentSportFilter = sportInput.value;

        // Перезавантажуємо дані
        fetchEvents();
    });

    // 3. Обробник для кнопки "Clear Filters"
    clearButton.addEventListener('click', () => {
        // Скидаємо глобальний стан
        currentPage = 1;
        currentDateFilter = "";
        currentSportFilter = "";

        // Скидаємо значення у формі
        dateInput.value = "";
        sportInput.value = "";
        
        // Перезавантажуємо дані
        fetchEvents();
    });

    // 4. Перший запуск: завантажуємо дані для сторінки 1
    fetchEvents();
});