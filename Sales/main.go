package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	// model "Sales/model"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error
var tpl *template.Template

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

func init() {
	db, err = sql.Open("mysql", "root:password@/sales")
	checkErr(err)
	err = db.Ping()
	checkErr(err)
	tpl = template.Must(template.ParseGlob("admin/*"))
}

func main() {
	defer db.Close()
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", index)

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

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	log.Println("Server is up on 8080 port")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	err = tpl.ExecuteTemplate(w, "home.html", nil)
	checkErr(err)
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
