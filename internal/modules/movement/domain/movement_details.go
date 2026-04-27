package domain

type MovementDetails struct {
	movement     Movement
	movementType MovementType
	category     Category
	account      Account
}

func NewMovementDetails(
	movement Movement,
	movementType MovementType,
	category Category,
	account Account,
) MovementDetails {
	return MovementDetails{
		movement:     movement,
		movementType: movementType,
		category:     category,
		account:      account,
	}
}

func (d MovementDetails) Movement() Movement         { return d.movement }
func (d MovementDetails) MovementType() MovementType { return d.movementType }
func (d MovementDetails) Category() Category         { return d.category }
func (d MovementDetails) Account() Account           { return d.account }
