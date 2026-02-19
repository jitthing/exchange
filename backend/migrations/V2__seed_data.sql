-- Academic events
INSERT INTO academic_events (id, type, title, start_date, end_date, priority) VALUES
    ('ev-1', 'exam',     'Economics Midterm',       '2026-03-18', '2026-03-18', 5),
    ('ev-2', 'deadline', 'Group Project Deadline',  '2026-03-24', '2026-03-24', 4),
    ('ev-3', 'holiday',  'Public Holiday',          '2026-04-03', '2026-04-05', 1)
ON CONFLICT (id) DO NOTHING;

-- Travel windows
INSERT INTO travel_windows (id, start_date, end_date, score, conflicts) VALUES
    ('w-1', '2026-03-06', '2026-03-08', 88, '[]'),
    ('w-2', '2026-03-20', '2026-03-22', 52, '["Near major deadline"]'),
    ('w-3', '2026-04-03', '2026-04-06', 95, '[]')
ON CONFLICT (id) DO NOTHING;

-- Trips
INSERT INTO trips (id, owner_id, destination, window_id, members, itinerary, estimated_cost) VALUES
    ('trip-1', 'demo-user', 'Prague', 'w-1', '["demo-user"]', '["Old Town walk","Charles Bridge sunrise"]', 220)
ON CONFLICT (id) DO NOTHING;

-- Budget entries
INSERT INTO budget_entries (id, user_id, category, amount, currency, date, trip_id, note) VALUES
    ('b-1', 'demo-user', 'living', 420, 'EUR', '2026-02-05', '', 'Rent split'),
    ('b-2', 'demo-user', 'travel', 60,  'EUR', '2026-02-08', '', 'Train to Vienna')
ON CONFLICT (id) DO NOTHING;

-- Monthly budgets
INSERT INTO monthly_budgets (user_id, budget) VALUES
    ('demo-user', 900)
ON CONFLICT (user_id) DO NOTHING;

-- Destinations
INSERT INTO destinations (city, base_travel_hrs, transport_base, hostel_night_eur, tags) VALUES
    ('Prague',    3.8, 55, 28, '["culture","city"]'),
    ('Budapest',  4.7, 47, 24, '["nightlife","city"]'),
    ('Ljubljana', 5.2, 41, 30, '["nature","city"]'),
    ('Krakow',    2.9, 50, 22, '["culture","city"]')
ON CONFLICT (city) DO NOTHING;
