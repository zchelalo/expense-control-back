package domain

type MovementDetails struct {
	movement     Movement
	movementType MovementType
	category     Category
}

func NewMovementDetails(
	movement Movement,
	movementType MovementType,
	category Category,
) MovementDetails {
	return MovementDetails{
		movement:     movement,
		movementType: movementType,
		category:     category,
	}
}

func (d MovementDetails) Movement() Movement         { return d.movement }
func (d MovementDetails) MovementType() MovementType { return d.movementType }
func (d MovementDetails) Category() Category         { return d.category }
