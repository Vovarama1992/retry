-- Справочник типов действий
CREATE TABLE action_types (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL
);

INSERT INTO action_types (name) VALUES
  ('visit');

-- Таблица действий
CREATE TABLE actions (
  id SERIAL PRIMARY KEY,                             
  action_type_id INT NOT NULL REFERENCES action_types(id),
  visit_id TEXT NOT NULL,                            -- с фронта
  source TEXT NOT NULL,                              -- utm/ref/direct
  ip_address TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT NOW()
);