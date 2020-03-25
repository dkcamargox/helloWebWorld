package main 

import (
  "database/sql"
  "fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	
)


type answ struct {
	id int
	nome string
}

func checkErr( err error ) {
	if err != nil {
		panic(err)
	}
}
func goDotEnvVariable(key string) string {

  // load .env file
  err := godotenv.Load(".env")

  checkErr(err)

  return os.Getenv(key)
}
var (
	host = goDotEnvVariable("HOST")
	port, errI = strconv.Atoi(goDotEnvVariable("PORT"))
	user=goDotEnvVariable("USER")
	password=goDotEnvVariable("PASSWORD")
	dbname=goDotEnvVariable("DB_NAME")
)
func selectAll ( table string, db *sql.DB ) []answ {
	var row_res_strc answ
	var slice_row_res_strc []answ
	row_res, err := db.Query(fmt.Sprintf("SELECT * FROM %s",table))
	checkErr(err)
	for row_res.Next() {
		
		err = row_res.Scan(&row_res_strc.id, &row_res_strc.nome)
		checkErr(err)
		slice_row_res_strc =append(slice_row_res_strc, row_res_strc)
		
	}
	
	return slice_row_res_strc
}

func get_things_from_bank () []answ {
	db, err := sql.Open("postgres", 
	fmt.Sprintf("host=%s "+
							"port=%d "+
							"user=%s "+
							"password=%s "+
							"dbname=%s "+
							"sslmode=disable",
								host, port, user, password, dbname))
	checkErr(err)

	res_bank := selectAll( "fabricantes", db )

	defer db.Close()

	err = db.Ping()
	checkErr(err)
	return res_bank
}
func parseSliceOfStructToJSONString ( toParse []answ ) string {
	var parsed string
	for index, value := range toParse {
		if index == (len(toParse)-1){
			parsed = fmt.Sprintf("%s{\"id\":\"%d\",\n\"nome\":\"%s\"}",parsed,value.id, value.nome)
		}else {
			parsed = fmt.Sprintf("%s{\"id\":\"%d\",\n\"nome\":\"%s\"},",parsed,value.id, value.nome)
		}
		
	}
	parsed = fmt.Sprintf("[%s]", parsed)
	return parsed
}
func handler(w http.ResponseWriter, r *http.Request) {
	
	fmt.Fprintln(w, parseSliceOfStructToJSONString(get_things_from_bank()))
}
func main() {
	checkErr(errI)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}