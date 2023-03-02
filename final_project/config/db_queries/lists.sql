SELECT c.name, l.item, l.list_type
  FROM checkers c
  JOIN lists l ON l.checker_id = c.id
  WHERE c.name IN (%s)
