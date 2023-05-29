CREATE TABLE IF NOT EXISTS COTACAO_USD_BRL (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT,
    codein TEXT,
    name TEXT,
    high NUMERIC,
    low NUMERIC,
    var_bid NUMERIC,
    pct_change NUMERIC,
    bid NUMERIC,
    ask NUMERIC,
    timestamp TEXT,
    create_date TEXT
);