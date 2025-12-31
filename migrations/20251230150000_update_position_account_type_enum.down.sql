UPDATE positions SET position_account_type = 'MARGIN' WHERE position_account_type = 'MARGIN_NEW';
UPDATE orders SET position_account_type = 'MARGIN' WHERE position_account_type = 'MARGIN_NEW';