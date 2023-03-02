SELECT c.name, l.item, l.interval_value, l.limit_value
  FROM checkers c
  JOIN limits l ON l.checker_id = c.id
  WHERE c.name IN (%s)
