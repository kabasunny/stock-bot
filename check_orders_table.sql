-- ordersテーブルの構造を確認
\d orders;

-- 実際のカラム一覧を表示
SELECT column_name, data_type, is_nullable, column_default 
FROM information_schema.columns 
WHERE table_name = 'orders' 
ORDER BY ordinal_position;