package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	var (
		res *models.Transaction
	)

	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// inisialisasi subtotal -> jumlah total transaksi keseluruhan
	totalAmount := 0
	// inisialisasi modeling transactionDetails -> nanti kita insert ke db
	details := make([]models.TransactionDetail, 0)
	// loop setiap item
	for _, item := range items {
		var productName string
		var productID, price, stock int
		// get product dapet pricing
		err := tx.QueryRow("SELECT id, name, price, stock FROM products WHERE id=$1", item.ProductID).Scan(&productID, &productName, &price, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}

		if err != nil {
			fmt.Println("error1 : " + err.Error())
			return nil, err
		}

		// hitung current total = quantity * pricing
		// ditambahin ke dalam subtotal
		subtotal := item.Quantity * price
		totalAmount += subtotal
		// kurangi jumlah stok
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, productID)
		if err != nil {
			return nil, err
		}

		// item nya dimasukkin ke transactionDetails
		details = append(details, models.TransactionDetail{
			ProductID:   productID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// insert transaction
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING ID", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// insert transaction details
	for i, detail := range details {
		details[i].TransactionID = transactionID
		var transactionDetailID int
		err := tx.QueryRow("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2,$3,$4) RETURNING ID", transactionID, detail.ProductID, detail.Quantity, detail.Subtotal).Scan(&transactionDetailID)
		if err != nil {
			return nil, err
		}
		details[i].ID = transactionDetailID
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	res = &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}

	return res, nil
}

func (repo *TransactionRepository) GetTodayReport() (*models.DailyReport, error) {
	var report models.DailyReport

	// Query total revenue dan total transaksi
	querySummary := `
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*) 
		FROM transactions 
		WHERE created_at::date = CURRENT_DATE
	`
	err := repo.db.QueryRow(querySummary).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// Query produk terlaris
	queryBestSeller := `
		SELECT p.name, COALESCE(SUM(td.quantity), 0) as sold 
		FROM transaction_details td 
		JOIN products p ON td.product_id = p.id 
		JOIN transactions t ON td.transaction_id = t.id
		WHERE t.created_at::date = CURRENT_DATE
		GROUP BY p.name 
		ORDER BY sold DESC LIMIT 1
	`
	err = repo.db.QueryRow(queryBestSeller).Scan(&report.ProdukTerlaris.Nama, &report.ProdukTerlaris.QtyTerjual)
	if err == sql.ErrNoRows {
		report.ProdukTerlaris.Nama = "-"
		report.ProdukTerlaris.QtyTerjual = 0
	} else if err != nil {
		return nil, err
	}

	return &report, nil
}
