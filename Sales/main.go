package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

	// model "Sales/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/go-sessions"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error
var templateLogin *template.Template
var tpl *template.Template
var templateSales *template.Template

type userLogin struct {
	ID       int64
	HakAkses int
	Username string
	Password []byte
}

type product struct {
	ID             int64
	Name           string
	Price          string
	Stock          string
	ProductionDate string
	ExpiredDate    string
}

type sales struct {
	ID       int64
	Fullname string
	Username string
	Password []byte
}

type transaction struct {
	ID            int64
	IDOrder       string
	IDSales       string
	StatusPayment string
	OrderDate     string
	PaymentDate   string
	TotalPayment  string
}

type transactionDetail struct {
	ID          int64
	ProductName string
	Price       string
	Qty         string
	Total       string
	IDOrder     string
	DateOrder   string
}

type salesTransaction struct {
	ID             int64
	Name           string
	Price          string
	Stock          string
	ProductionDate string
	ExpiredDate    string
	IDTransaction  int64
}

func init() {
	db, err = sql.Open("mysql", "root:dendi@/sales")
	checkErr(err)
	err = db.Ping()
	checkErr(err)
	// templateLogin = template.Must(template.ParseGlob("login/*"))
	tpl = template.Must(template.ParseGlob("admin/*"))
	templateSales = template.Must(template.ParseGlob("sales/*"))
}

func main() {
	defer db.Close()
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", login)
	http.HandleFunc("/logout", logout)

	http.HandleFunc("/home", indexAdmin)
	http.HandleFunc("/product", productList)
	http.HandleFunc("/productAdd", productAddForm)
	http.HandleFunc("/productAdd_", productAddExec)
	http.HandleFunc("/productEdit", productEditForm)
	http.HandleFunc("/productEdit_", productEditExec)
	http.HandleFunc("/productDelete", productDelete)

	http.HandleFunc("/sales", salesList)
	http.HandleFunc("/salesAdd", salesAddForm)
	http.HandleFunc("/salesAdd_", salesAddExec)
	http.HandleFunc("/salesEdit", salesEditForm)
	http.HandleFunc("/salesEdit_", salesEditExec)
	http.HandleFunc("/salesDelete", salesDelete)

	// -------- UI Sales -------- //
	http.HandleFunc("/fSales", indexSales)
	http.HandleFunc("/fSalesProduct", productListSales)
	http.HandleFunc("/fSalesBuy", orderSales)
	http.HandleFunc("/fSalesBuyAdd", orderSalesAdd)
	http.HandleFunc("/fSalesBuyProductAdd", orderSalesProductAddForm)
	http.HandleFunc("/fSalesBuyProductAdd_", orderSalesProductAddExec)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	log.Println("Server is up on 8080 port")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

// func checkErrLogin(w http.ResponseWriter, r *http.Request, err error) bool {
// 	if err != nil {

// 		fmt.Println(r.Host + r.URL.Path)

// 		http.Redirect(w, r, r.Host+r.URL.Path, 301)
// 		return false
// 	}

// 	return true
// }

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func generateIDOrder() string {
	var acakID string

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		acakID += fmt.Sprint(rand.Intn(7))
	}
	return acakID
}

func generateTime() string {
	var dt = time.Now()

	return dt.Format("01-02-2006")
}

func QueryUser(username string) userLogin {
	var users = userLogin{}
	err = db.QueryRow(`
		SELECT id, 
		hak_akses,
		username,  
		password 
		FROM login WHERE username=?`, username).
		Scan(
			&users.ID,
			&users.HakAkses,
			&users.Username,
			&users.Password,
		)
	return users
}

func login(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) != 0 {
		http.Redirect(w, r, "/home", 302)
	}
	if r.Method != "POST" {
		// http.Redirect(w, r, "/", 302)
		http.ServeFile(w, r, "index.html")
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	users := QueryUser(username)

	var password_tes = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(password))

	if password_tes == nil {
		session = sessions.Start(w, r)
		session.Set("username", users.Username)
		session.Set("hak_akses", users.HakAkses)
		if users.HakAkses == 1 {
			http.Redirect(w, r, "/home", 302)
		} else if users.HakAkses == 2 {
			http.Redirect(w, r, "/fSales", 302)
		}
	} else {
		http.Redirect(w, r, "index.html", 302)
	}
}

func logout(w http.ResponseWriter, req *http.Request) {
	session := sessions.Start(w, req)
	session.Clear()
	sessions.Destroy(w, req)
	http.Redirect(w, req, "/", 302)
}

func indexAdmin(w http.ResponseWriter, req *http.Request) {
	session := sessions.Start(w, req)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, req, "/", 301)
	}
	var data = map[string]string{
		"ses_username":  session.GetString("username"),
		"ses_hak_akses": session.GetString("hak_akses"),
	}
	err = tpl.ExecuteTemplate(w, "home.html", data)
	checkErr(err)
	return
}

func productList(w http.ResponseWriter, req *http.Request) {
	rows, e := db.Query(
		`SELECT id,
			name,
			price,
			stock,
			production_date,
			expired_date
		FROM product;
		`)
	checkErr(e)

	products := make([]product, 0)
	for rows.Next() {
		usr := product{}
		rows.Scan(&usr.ID, &usr.Name, &usr.Price, &usr.Stock, &usr.ProductionDate, &usr.ExpiredDate)
		products = append(products, usr)
	}
	log.Println(products)
	tpl.ExecuteTemplate(w, "product.html", products)
}

func productAddForm(w http.ResponseWriter, req *http.Request) {
	err = tpl.ExecuteTemplate(w, "productAdd.html", nil)
	checkErr(err)
}

func productAddExec(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		Product := product{}
		Product.Name = req.FormValue("name")
		Product.Price = req.FormValue("price")
		Product.Stock = req.FormValue("stock")
		Product.ProductionDate = req.FormValue("production_date")
		Product.ExpiredDate = req.FormValue("expired_date")
		_, err = db.Exec(
			"INSERT INTO product (name, price, stock, production_date, expired_date) VALUES (?, ?, ?, ?, ?)",
			Product.Name,
			Product.Price,
			Product.Stock,
			Product.ProductionDate,
			Product.ExpiredDate,
		)
		checkErr(err)
		http.Redirect(w, req, "/product", http.StatusSeeOther)
		return
	}
	http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
}

func productEditForm(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	rows, err := db.Query(
		`SELECT id,
		 	name,
			price,
			stock,
			production_date,
			expired_date
		FROM product
		WHERE id = ` + id + `;`)
	checkErr(err)
	Product := product{}
	for rows.Next() {
		rows.Scan(&Product.ID, &Product.Name, &Product.Price, &Product.Stock, &Product.ProductionDate, &Product.ExpiredDate)
	}
	tpl.ExecuteTemplate(w, "productEdit.html", Product)
}

func productEditExec(w http.ResponseWriter, req *http.Request) {
	_, er := db.Exec(
		"UPDATE product SET name = ?, price = ?, stock = ?, production_date = ?, expired_date = ? WHERE id = ? ",
		req.FormValue("name"),
		req.FormValue("price"),
		req.FormValue("stock"),
		req.FormValue("production_date"),
		req.FormValue("expired_date"),
		req.FormValue("id"),
	)
	checkErr(er)
	http.Redirect(w, req, "/product", http.StatusSeeOther)
}

func productDelete(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		http.Error(res, "ID not valid", http.StatusBadRequest)
	}
	_, er := db.Exec("DELETE FROM product WHERE id = ?", id)
	checkErr(er)
	http.Redirect(res, req, "/product", http.StatusSeeOther)
}

func salesList(w http.ResponseWriter, req *http.Request) {
	rows, e := db.Query(
		`SELECT id,
			fullname,
			username,
			password
		FROM sales;
		`)
	checkErr(e)

	Sales := make([]sales, 0)
	for rows.Next() {
		sales_ := sales{}
		rows.Scan(&sales_.ID, &sales_.Fullname, &sales_.Username, &sales_.Password)
		Sales = append(Sales, sales_)
	}
	log.Println(Sales)
	tpl.ExecuteTemplate(w, "sales.html", Sales)
}

func salesAddForm(w http.ResponseWriter, req *http.Request) {
	err = tpl.ExecuteTemplate(w, "salesAdd.html", nil)
	checkErr(err)
}

func salesAddExec(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		Sales := sales{}
		Sales.Fullname = req.FormValue("name")
		Sales.Username = req.FormValue("user")
		bPass, e := bcrypt.GenerateFromPassword([]byte(req.FormValue("password")), bcrypt.MinCost)
		checkErr(e)
		Sales.Password = bPass
		_, e = db.Exec(
			"INSERT INTO sales (fullname, username, password) VALUES (?, ?, ?)",
			Sales.Fullname,
			Sales.Username,
			Sales.Password,
		)
		checkErr(e)
		http.Redirect(w, req, "/sales", http.StatusSeeOther)
		return
	}
	http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
}

func salesEditForm(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	rows, err := db.Query(
		`SELECT id,
			 fullname,
			 username
		FROM sales
		WHERE id = ` + id + `;`)
	checkErr(err)
	sls := sales{}
	for rows.Next() {
		rows.Scan(&sls.ID, &sls.Fullname, &sls.Username)
	}
	tpl.ExecuteTemplate(w, "salesEdit.html", sls)
}

func salesEditExec(w http.ResponseWriter, req *http.Request) {
	bPass, e := bcrypt.GenerateFromPassword([]byte(req.FormValue("password")), bcrypt.MinCost)
	checkErr(e)
	_, er := db.Exec(
		"UPDATE sales SET fullname = ?, username = ?, password = ? WHERE id = ? ",
		req.FormValue("name"),
		req.FormValue("user"),
		bPass,
		req.FormValue("id"),
	)
	checkErr(er)
	http.Redirect(w, req, "/sales", http.StatusSeeOther)
}

func salesDelete(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		http.Error(res, "ID not valid.", http.StatusBadRequest)
	}
	_, er := db.Exec("DELETE FROM sales WHERE id = ?", id)
	checkErr(er)
	http.Redirect(res, req, "/sales", http.StatusSeeOther)
}

// --------- For UI Sales --------- //
func indexSales(w http.ResponseWriter, req *http.Request) {
	session := sessions.Start(w, req)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, req, "/", 301)
	}
	var data = map[string]string{
		"ses_username":  session.GetString("username"),
		"ses_hak_akses": session.GetString("hak_akses"),
	}
	err = templateSales.ExecuteTemplate(w, "homeSales.html", data)
	checkErr(err)
	return
}

func productListSales(w http.ResponseWriter, req *http.Request) {
	rows, e := db.Query(
		`SELECT id,
			name,
			price,
			stock,
			production_date,
			expired_date
		FROM product;
		`)
	checkErr(e)

	products := make([]product, 0)
	for rows.Next() {
		usr := product{}
		rows.Scan(&usr.ID, &usr.Name, &usr.Price, &usr.Stock, &usr.ProductionDate, &usr.ExpiredDate)
		products = append(products, usr)
	}
	log.Println(products)
	templateSales.ExecuteTemplate(w, "productSales.html", products)
}

func orderSales(w http.ResponseWriter, req *http.Request) {
	rows, e := db.Query(
		`SELECT
    transaction.*,
    SUM(total)
FROM
    transaction
INNER JOIN transaction_detail
on transaction.id=transaction_detail.id_transaction;
		`)
	checkErr(e)

	Transaction := make([]transaction, 0)
	for rows.Next() {
		trx := transaction{}
		rows.Scan(&trx.ID, &trx.IDOrder, &trx.IDSales, &trx.StatusPayment, &trx.OrderDate, &trx.PaymentDate, &trx.TotalPayment)
		Transaction = append(Transaction, trx)
	}
	log.Println(Transaction)
	templateSales.ExecuteTemplate(w, "orderProductSales.html", Transaction)
}

func orderSalesAdd(w http.ResponseWriter, req *http.Request) {
	Transaction := transaction{}
	Transaction.IDOrder = generateIDOrder()
	Transaction.IDSales = "2"
	Transaction.StatusPayment = "Waiting to Payment"
	Transaction.OrderDate = generateTime()
	Transaction.PaymentDate = "-"
	_, err = db.Exec(
		"INSERT INTO transaction (id_order, id_sales, status_payment, order_date, payment_date) VALUES (?, ?, ?, ?, ?)",
		Transaction.IDOrder,
		Transaction.IDSales,
		Transaction.StatusPayment,
		Transaction.OrderDate,
		Transaction.PaymentDate,
	)
	checkErr(err)
	http.Redirect(w, req, "/fSalesBuy", http.StatusSeeOther)
	return
}

func orderSalesProductAddForm(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	rowss, es := db.Query(`SELECT id FROM transaction WHERE id= ` + id + `;`)
	checkErr(es)

	rows, e := db.Query(
		`SELECT id,
			name,
			price,
			stock,
			production_date,
			expired_date
		FROM product;
		`)
	checkErr(e)

	SalesTransaction := make([]salesTransaction, 0)
	for rows.Next() {
		usr := salesTransaction{}
		rowss.Next()
		rows.Scan(&usr.ID, &usr.Name, &usr.Price, &usr.Stock, &usr.ProductionDate, &usr.ExpiredDate)
		rowss.Scan(&usr.IDTransaction)
		SalesTransaction = append(SalesTransaction, usr)
	}
	log.Println(SalesTransaction)
	templateSales.ExecuteTemplate(w, "orderProductSalesAdd.html", SalesTransaction)
}

func checkStock(name string) product {
	var Products = product{}
	err = db.QueryRow(`
		SELECT stock
		FROM product WHERE name=?`, name).
		Scan(
			&Products.Stock,
		)
	return Products
}

func checkPrice(name string) product {
	var Products = product{}
	err = db.QueryRow(`
		SELECT price
		FROM product WHERE name=?`, name).
		Scan(
			&Products.Price,
		)
	return Products
}

// func total() product {
// 	var Products = product{}
// 	err = db.QueryRow(`
// 		SELECT price, qty
// 		FROM transacntion_detail WHERE id_transaction=1`).
// 		Scan(
// 			&Products.Price,
// 		)
// 	return Products
// }

func orderSalesProductAddExec(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		// checkstock()
		TransactionDetail := transactionDetail{}
		TransactionDetail.ProductName = req.FormValue("product")
		var priceCheck = checkPrice(TransactionDetail.ProductName)
		TransactionDetail.Qty = req.FormValue("qty")
		id := req.FormValue("id")
		// var asd = 0
		var DateOrderx = generateTime()

		var stockCheck = checkStock(TransactionDetail.ProductName)
		if TransactionDetail.Qty <= stockCheck.Stock {
			_, err = db.Exec(
				"INSERT INTO transaction_detail (product_name, price, qty, id_transaction, date_order) VALUES (?, ?, ?, ?, ?)",
				TransactionDetail.ProductName,
				priceCheck.Price,
				TransactionDetail.Qty,
				id,
				DateOrderx,
			)
			checkErr(err)

			http.Redirect(w, req, "/fSalesBuy", http.StatusSeeOther)
			return

		} else {
			// TODO : create handler UI
			log.Println("qty kurang")
			http.Redirect(w, req, "/fSalesBuy", http.StatusSeeOther)
			return
		}
	}
	http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
}
