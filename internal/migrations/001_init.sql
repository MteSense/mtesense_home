CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'admin',
  disabled INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS nav_groups (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  sort_order INTEGER NOT NULL DEFAULT 0,
  visible INTEGER NOT NULL DEFAULT 1,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS nav_links (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  group_id INTEGER NOT NULL REFERENCES nav_groups(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  url TEXT NOT NULL,
  icon TEXT NOT NULL DEFAULT '',
  icon_type TEXT NOT NULL DEFAULT 'emoji',
  description TEXT NOT NULL DEFAULT '',
  sort_order INTEGER NOT NULL DEFAULT 0,
  visible INTEGER NOT NULL DEFAULT 1,
  open_in_new_tab INTEGER NOT NULL DEFAULT 1,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS settings (
  key TEXT PRIMARY KEY,
  value_json TEXT NOT NULL,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO settings (key, value_json)
VALUES
  ('appearance', '{"siteTitle":"MteSense","browserTitle":"MteSense","subtitle":"Personal navigation","backgroundImage":"","defaultTheme":"dark","cardOpacity":0.34,"blurStrength":18}'),
  ('search', '{"defaultSearchEngine":"google","enabledSearchEngines":["google","bing","baidu"]}')
ON CONFLICT(key) DO NOTHING;

INSERT INTO nav_groups (title, sort_order, visible)
SELECT 'APP', 10, 1
WHERE NOT EXISTS (SELECT 1 FROM nav_groups);

INSERT INTO nav_links (group_id, title, url, icon, icon_type, description, sort_order, visible, open_in_new_tab)
SELECT g.id, 'Google', 'https://www.google.com', 'G', 'text', 'Search the web', 10, 1, 1
FROM nav_groups g
WHERE g.title = 'APP' AND NOT EXISTS (SELECT 1 FROM nav_links);

INSERT INTO nav_links (group_id, title, url, icon, icon_type, description, sort_order, visible, open_in_new_tab)
SELECT g.id, 'Baidu', 'https://www.baidu.com', '百', 'text', '中文搜索', 20, 1, 1
FROM nav_groups g
WHERE g.title = 'APP' AND NOT EXISTS (SELECT 1 FROM nav_links WHERE title = 'Baidu');
