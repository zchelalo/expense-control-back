INSERT INTO movement_types (
  id,
  name,
  key,
  description
)
VALUES
  (
    '7a6478f0-8a80-43de-9fef-e4e44d77d111',
    'Ingreso',
    'income',
    'Movimientos que incrementan el saldo disponible'
  ),
  (
    '9cb65b31-b6c5-489f-8c57-2076ce2e4222',
    'Egreso',
    'expense',
    'Movimientos que disminuyen el saldo disponible'
  )
ON CONFLICT (key) DO NOTHING
;
