package storage

type packageSize struct {
	ID        string `db:"id"`
	ProductID string `db:"product_id"`
	Size      int    `db:"size"`
}

type product struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}
