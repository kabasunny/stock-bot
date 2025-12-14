ALTER TABLE stock_masters
ADD COLUMN issue_name_short VARCHAR(255),
ADD COLUMN issue_name_kana VARCHAR(255),
ADD COLUMN issue_name_english VARCHAR(255),
ADD COLUMN industry_code VARCHAR(255),
ADD COLUMN industry_name VARCHAR(255),
ADD COLUMN listed_shares_outstanding BIGINT;
