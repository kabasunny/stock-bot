UPDATE positions SET position_account_type = 'MARGIN_NEW' WHERE position_account_type = 'MARGIN';
UPDATE orders SET position_account_type = 'MARGIN_NEW' WHERE position_account_type = 'MARGIN';