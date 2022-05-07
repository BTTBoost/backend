DROP TABLE IF EXISTS holders_activity;
CREATE TABLE IF NOT EXISTS holders_activity
(
    holder TEXT PRIMARY KEY,
    day    BOOL,
    week   BOOL,
    month  BOOL
);
