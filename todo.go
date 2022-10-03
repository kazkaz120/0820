package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/labstack/echo/middleware"
	"github.com/x-color/simple-webapp/handler"

	"github.com/labstack/echo"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Date       string
	Time_start string
	Time_end   string
	To_do      string
	Which_do   string
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func WriteTasks(c echo.Context) {

	date := c.FormValue("date")
	time_start := c.FormValue("time_start")
	time_end := c.FormValue("time_end")
	to_do := c.FormValue("to_do")
	which_do := c.FormValue("which_do")

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{})

	db.Create(&Product{Date: date, Time_start: time_start, Time_end: time_end, To_do: to_do, Which_do: which_do})

	//	var pro2 []string

	/*	product := []Product{}
		db.Find(&product)
		for _, pro := range product {
			pro2 = append(pro2, pro.Date)
			pro2 = append(pro2, pro.Time_start)
			pro2 = append(pro2, pro.Time_end)
			pro2 = append(pro2, pro.To_do)
			pro2 = append(pro2, pro.Which_do)
			//	fmt.Println(pro2)
	}*/

}

type Toushi_perc struct {
	Item string
	Perc int
}
type Shouhi_perc struct {
	Item string
	Perc int
}
type Rouhi_perc struct {
	Item string
	Perc int
}

func CreateTasks(c echo.Context) (string, string, string, int, int, int, []string, []int, []string, []int, []string, []int) {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{})

	var pro_sum_toushi int
	var pro_sum_shouhi int
	var pro_sum_rouhi int

	product := []Product{}

	checkdate_month := c.FormValue("month")

	/*データベースから「投資」の項目を取得し、時間の合計値を取る*/
	db.Where("date LIKE ?", checkdate_month+"%").Where("which_do = ?", "投資").Find(&product)
	//	db.Where("which_do = ?", "投資").Delete(&product)
	var todo_arr []string
	var todo_arr_shohi []string
	var todo_arr_rouhi []string

	for _, pro := range product {

		time_start_head := pro.Time_start[:2]
		time_start_bottom := pro.Time_start[3:]
		time_end_head := pro.Time_end[:2]
		time_end_bottom := pro.Time_end[3:]

		time_start_head_int, _ := strconv.Atoi(time_start_head)
		time_start_bottom_int, _ := strconv.Atoi(time_start_bottom)
		time_end_head_int, _ := strconv.Atoi(time_end_head)
		time_end_bottom_int, _ := strconv.Atoi(time_end_bottom)

		time_start := time_start_head_int*60 + time_start_bottom_int
		time_end := time_end_head_int*60 + time_end_bottom_int

		if time_start > time_end {
			time_end = time_end + 24*60
		}

		time_differ := time_end - time_start

		//		fmt.Println(time_differ)

		pro_sum_toushi += time_differ

		todo_arr = append(todo_arr, pro.To_do)

	}

	m := make(map[string]bool)
	uniq := []string{}
	for _, ele := range todo_arr {
		if !m[ele] {
			m[ele] = true
			uniq = append(uniq, ele)
		}
	}

	//fmt.Println(uniq)
	//fmt.Println(todo_arr)

	//	uniq_value := len(uniq)
	uniq_count := make([]int, len(uniq))
	for n := 0; n < len(uniq); n++ {
		for p := 0; p < len(todo_arr); p++ {
			if uniq[n] == todo_arr[p] {
				uniq_count[n]++
			}

		}
	}

	var sum_uniq_count int

	for _, x := range uniq_count {
		sum_uniq_count += x
	}

	parc_uniq_count := make([]int, len(uniq))

	for y := 0; y < len(uniq_count); y++ {
		parc_uniq_count[y] = uniq_count[y] * 100 / sum_uniq_count
	}

	//	fmt.Println(uniq_count)
	//	fmt.Println(parc_uniq_count)

	//	toushi_perc := []Toushi_perc{}

	/*	for i := 0; i < len(uniq); i++ {
			toushi_perc[i].Item = uniq[i]
			toushi_perc[i].Perc = parc_uniq_count[i]
		}
		fmt.Println(toushi_perc)*/

	/*データベースから「消費」の項目を取得し、時間の合計値を取る*/
	db.Where("date LIKE ?", checkdate_month+"%").Where("which_do = ?", "消費").Find(&product)
	//	db.Where("which_do = ?", "投資").Delete(&product)
	for _, pro := range product {
		//pro2 = append(pro2, pro.Time_start)
		//pro2 = append(pro2, pro.Time_end)

		time_start_head := pro.Time_start[:2]
		time_start_bottom := pro.Time_start[3:]
		time_end_head := pro.Time_end[:2]
		time_end_bottom := pro.Time_end[3:]

		time_start_head_int, _ := strconv.Atoi(time_start_head)
		time_start_bottom_int, _ := strconv.Atoi(time_start_bottom)
		time_end_head_int, _ := strconv.Atoi(time_end_head)
		time_end_bottom_int, _ := strconv.Atoi(time_end_bottom)

		time_start := time_start_head_int*60 + time_start_bottom_int
		time_end := time_end_head_int*60 + time_end_bottom_int

		if time_start > time_end {
			time_end = time_end + 24*60
		}

		time_differ := time_end - time_start

		//		fmt.Println(time_differ)

		pro_sum_shouhi += time_differ
		//		fmt.Println(pro.To_do)

		todo_arr_shohi = append(todo_arr_shohi, pro.To_do)

	}
	m_shohi := make(map[string]bool)
	uniq_shohi := []string{}
	for _, ele := range todo_arr_shohi {
		if !m_shohi[ele] {
			m_shohi[ele] = true
			uniq_shohi = append(uniq_shohi, ele)
		}
	}

	//	fmt.Println(uniq_shohi)
	//	uniq_value := len(uniq)
	uniq_count_shohi := make([]int, len(uniq_shohi))
	for n := 0; n < len(uniq_shohi); n++ {
		for p := 0; p < len(todo_arr_shohi); p++ {
			if uniq_shohi[n] == todo_arr_shohi[p] {
				uniq_count_shohi[n]++
			}

		}
	}

	var sum_uniq_count_shohi int

	for _, x := range uniq_count_shohi {
		sum_uniq_count_shohi += x
	}

	parc_uniq_count_shohi := make([]int, len(uniq_shohi))

	for y := 0; y < len(uniq_count_shohi); y++ {
		parc_uniq_count_shohi[y] = uniq_count_shohi[y] * 100 / sum_uniq_count_shohi
	}

	//	fmt.Println(parc_uniq_count_shohi)

	/*データベースから「浪費」の項目を取得し、時間の合計値を取る*/
	db.Where("date LIKE ?", checkdate_month+"%").Where("which_do = ?", "浪費").Find(&product)
	//	db.Where("which_do = ?", "投資").Delete(&product)
	for _, pro := range product {

		time_start_head := pro.Time_start[:2]
		time_start_bottom := pro.Time_start[3:]
		time_end_head := pro.Time_end[:2]
		time_end_bottom := pro.Time_end[3:]

		time_start_head_int, _ := strconv.Atoi(time_start_head)
		time_start_bottom_int, _ := strconv.Atoi(time_start_bottom)
		time_end_head_int, _ := strconv.Atoi(time_end_head)
		time_end_bottom_int, _ := strconv.Atoi(time_end_bottom)

		time_start := time_start_head_int*60 + time_start_bottom_int
		time_end := time_end_head_int*60 + time_end_bottom_int

		if time_start > time_end {
			time_end = time_end + 24*60
		}

		time_differ := time_end - time_start

		//		fmt.Println(time_differ)

		pro_sum_rouhi += time_differ
		//		fmt.Println(pro.To_do)

		todo_arr_rouhi = append(todo_arr_rouhi, pro.To_do)

	}

	m_rouhi := make(map[string]bool)
	uniq_rouhi := []string{}
	for _, ele := range todo_arr_rouhi {
		if !m_rouhi[ele] {
			m_rouhi[ele] = true
			uniq_rouhi = append(uniq_rouhi, ele)
		}
	}

	//	fmt.Println(uniq_rouhi)

	//	login_user := new(model.User)
	//	fmt.Println(login_user.Name)

	//	uniq_value := len(uniq)
	uniq_count_rouhi := make([]int, len(uniq_rouhi))
	for n := 0; n < len(uniq_rouhi); n++ {
		for p := 0; p < len(todo_arr_rouhi); p++ {
			if uniq_rouhi[n] == todo_arr_rouhi[p] {
				uniq_count_rouhi[n]++
			}

		}
	}

	var sum_uniq_count_rouhi int

	for _, x := range uniq_count_rouhi {
		sum_uniq_count_rouhi += x
	}

	parc_uniq_count_rouhi := make([]int, len(uniq_rouhi))

	for y := 0; y < len(uniq_count_rouhi); y++ {
		parc_uniq_count_rouhi[y] = uniq_count_rouhi[y] * 100 / sum_uniq_count_rouhi
	}

	//	fmt.Println(parc_uniq_count_rouhi)

	/*順位付けを行う構造体*/
	type Rank struct {
		What string
		Sum  int
		Moji string
	}

	rank := []Rank{
		{What: "pro_sum_toushi", Sum: pro_sum_toushi, Moji: "投資"},
		{What: "pro_sum_shouhi", Sum: pro_sum_shouhi, Moji: "消費"},
		{What: "pro_sum_rouhi", Sum: pro_sum_rouhi, Moji: "浪費"},
	}

	/*時間順に並べ替えを行う*/
	sort.Slice(rank, func(i, j int) bool { return rank[i].Sum < rank[j].Sum })
	fmt.Printf("1位:%+v\n", rank[2].Moji)
	fmt.Printf("2位:%+v\n", rank[1].Moji)
	fmt.Printf("3位:%+v\n", rank[0].Moji)
	fmt.Printf("1位:%+v\n", rank[2].Sum)
	fmt.Printf("2位:%+v\n", rank[1].Sum)
	fmt.Printf("3位:%+v\n", rank[0].Sum)

	return rank[2].Moji, rank[1].Moji, rank[0].Moji, rank[2].Sum, rank[1].Sum, rank[0].Sum, uniq, parc_uniq_count, uniq_shohi, parc_uniq_count_shohi, uniq_rouhi, parc_uniq_count_rouhi

}

/*フロントエンドへ値を渡す構造体*/
type Data struct {
	Rank1_moji            string
	Rank2_moji            string
	Rank3_moji            string
	Rank1_clock           int
	Rank1_minute          int
	Rank2_clock           int
	Rank2_minute          int
	Rank3_clock           int
	Rank3_minute          int
	Uniq_toshi            []string
	Parc_uniq_count_toshi []int
	Uniq_shohi            []string
	Parc_uniq_count_shohi []int
	Uniq_rouhi            []string
	Parc_uniq_count_rouhi []int
}

func CheckDate(c echo.Context) []Product {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Product{})

	checkdate_month := c.FormValue("month")

	product := []Product{}
	//	fmt.Println(checkdate_month)

	db.Where("date LIKE ?", checkdate_month+"%").Find(&product)

	//	fmt.Println(product)

	return product
	/*for _, pro := range product {
		pro_date := pro.Date
		fmt.Println(pro_date)
	}*/

}

func CheckDate_week(c echo.Context) []Product {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Product{})

	week_num := [][]string{{"01", "07"}, {"08", "14"}, {"15", "21"}, {"22", "28"}, {"29", "31"}}

	checkdate_month := c.FormValue("month")
	checkdate_week := c.FormValue("week")
	//	fmt.Println(checkdate_week)
	checkdate_week_int, _ := strconv.Atoi(checkdate_week)

	//	checkdate_month_num := strings.Replace(checkdate_month, "-", "", -1)
	product := []Product{}
	//	fmt.Println(checkdate_month)
	week_x := week_num[checkdate_week_int-1][0]
	week_y := week_num[checkdate_week_int-1][1]

	//	fmt.Println(week_x)
	//	fmt.Println(week_y)
	//	week_x_string := strconv.Itoa(week_x)
	//	week_y_string := strconv.Itoa(week_y)
	//	check_week_x := checkdate_month_num + week_x_string
	//	check_week_y := checkdate_month_num + week_y_string

	check_week_x := checkdate_month + "-" + week_x
	check_week_y := checkdate_month + "-" + week_y
	const format1 = "2006-01-02 15:04:05"

	//	fmt.Println(check_week_x)
	//	fmt.Println(check_week_y)

	t1, _ := time.Parse(format1, check_week_x+" 00:00:05")
	t2, _ := time.Parse(format1, check_week_y+" 23:59:55")

	//	fmt.Println(t1)
	//	fmt.Println(t2)

	//	check_week_x_int, _ := strconv.Atoi(check_week_x)
	//	check_week_y_int, _ := strconv.Atoi(check_week_y)

	//	for i := week_x; i < week_y+1; i++ {
	//		j := strconv.Itoa(i)
	db.Where("date BETWEEN ? AND ?", t1, t2).Find(&product)
	//	fmt.Println(check_week_x)
	//	}
	fmt.Println(product)

	return product
	/*for _, pro := range product {
		pro_date := pro.Date
		fmt.Println(pro_date)
	}*/

}

func CheckDate_day(c echo.Context) []Product {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Product{})

	checkdate_day := c.FormValue("date")

	product := []Product{}
	//	fmt.Println(checkdate_day)

	db.Where("date = ?", checkdate_day).Find(&product)

	//	fmt.Println(product)

	return product

}

func dbDelete(id int) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Product{})

	var product Product
	db.First(&product, id)
	db.Delete(&product)
}

func main() {

	t := &Template{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/static", "public/assets")
	e.Renderer = t

	e.File("/", "public/index.html")
	e.File("/signup", "public/signup.html")
	e.POST("/signup", handler.Signup)
	e.File("/login", "public/login.html")
	e.POST("/login", handler.Login)
	e.File("/todos", "public/todos.html")

	e.File("/writedata", "public/writedata.html")
	e.POST("/write", func(c echo.Context) error {
		var data Data
		WriteTasks(c)
		return c.Render(http.StatusOK, "writedata.html", data)
	})
	e.File("/top", "public/top.html")
	e.GET("/output_juni", func(c echo.Context) error {
		/*cはechoの変数で使っているので使えない*/
		a, b, d, e, f, g, h, i, j, k, l, m := CreateTasks(c)

		/*CreateTasksの返り値(種別と時間)が順位順になってa,b,d,e,f,gに代入*/
		var data Data
		data.Rank1_moji = a
		data.Rank2_moji = b
		data.Rank3_moji = d
		var time_clock1 int
		var time_minute1 int
		time_clock1 = e / 60
		time_minute1 = e % 60
		//		fmt.Println(time_clock1)
		//		fmt.Println(time_minute1)

		data.Rank1_clock = time_clock1
		data.Rank1_minute = time_minute1

		var time_clock2 int
		var time_minute2 int
		time_clock2 = f / 60
		time_minute2 = f % 60
		//		fmt.Println(time_clock2)
		//		fmt.Println(time_minute2)

		data.Rank2_clock = time_clock2
		data.Rank2_minute = time_minute2

		var time_clock3 int
		var time_minute3 int
		time_clock3 = g / 60
		time_minute3 = g % 60
		//		fmt.Println(time_clock3)
		//		fmt.Println(time_minute3)

		data.Rank3_clock = time_clock3
		data.Rank3_minute = time_minute3

		data.Uniq_toshi = h
		data.Parc_uniq_count_toshi = i
		data.Uniq_shohi = j
		data.Parc_uniq_count_shohi = k
		data.Uniq_rouhi = l
		data.Parc_uniq_count_rouhi = m
		//		var toushi_perc Toushi_perc
		//		toushi_perc.Item = h.Item
		//		toushi_perc.Perc = h.Perc
		return c.Render(http.StatusOK, "output_window.html", data)
	})

	e.POST("/output_juni_month", func(c echo.Context) error {
		/*cはechoの変数で使っているので使えない*/
		a, b, d, e, f, g, h, i, j, k, l, m := CreateTasks(c)

		/*CreateTasksの返り値(種別と時間)が順位順になってa,b,d,e,f,gに代入*/
		var data Data
		data.Rank1_moji = a
		data.Rank2_moji = b
		data.Rank3_moji = d
		var time_clock1 int
		var time_minute1 int
		time_clock1 = e / 60
		time_minute1 = e % 60
		data.Rank1_clock = time_clock1
		data.Rank1_minute = time_minute1

		var time_clock2 int
		var time_minute2 int
		time_clock2 = f / 60
		time_minute2 = f % 60
		data.Rank2_clock = time_clock2
		data.Rank2_minute = time_minute2

		var time_clock3 int
		var time_minute3 int
		time_clock3 = g / 60
		time_minute3 = g % 60
		data.Rank3_clock = time_clock3
		data.Rank3_minute = time_minute3

		data.Uniq_toshi = h
		data.Parc_uniq_count_toshi = i
		data.Uniq_shohi = j
		data.Parc_uniq_count_shohi = k
		data.Uniq_rouhi = l
		data.Parc_uniq_count_rouhi = m

		//		var toushi_perc Toushi_perc
		//		toushi_perc.Item = h.Item
		//		toushi_perc.Perc = h.Perc
		fmt.Println(data)

		return c.Render(http.StatusOK, "output_window.html", data)
	})

	e.File("/checkdate_month", "public/checkdate_month.html")
	e.POST("/checkdate_month_detail", func(c echo.Context) error {
		detail := CheckDate(c)
		for _, check := range detail {
			fmt.Println(check.Date)
		}
		return c.Render(http.StatusOK, "checkdate_month_detail.html", detail)
	})
	e.POST("/delete/:id", func(c echo.Context) error {

		n := c.Param("id")
		//	fmt.Println(n)
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		dbDelete(id)
		return c.Render(http.StatusOK, "checkdate_month.html", 1)
	})

	e.File("/checkdate_week", "public/checkdate_week.html")
	e.POST("/checkdate_week_detail", func(c echo.Context) error {
		detail := CheckDate_week(c)
		for _, check := range detail {
			fmt.Println(check.Date)
		}
		return c.Render(http.StatusOK, "checkdate_week_detail.html", detail)
	})
	e.POST("/delete_week/:id", func(c echo.Context) error {

		n := c.Param("id")
		//	fmt.Println(n)
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		dbDelete(id)
		return c.Render(http.StatusOK, "checkdate_week.html", 1)
	})

	e.File("/checkdate_day", "public/checkdate_day.html")
	e.POST("/checkdate_day_detail", func(c echo.Context) error {
		//		fmt.Println(c)
		//		fmt.Println(c.FormValue("date"))
		detail := CheckDate_day(c)
		for _, check := range detail {
			fmt.Println(check.Date)
		}
		return c.Render(http.StatusOK, "checkdate_day_detail.html", detail)
	})
	e.POST("/delete_day/:id", func(c echo.Context) error {

		n := c.Param("id")
		//	fmt.Println(n)
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		dbDelete(id)
		return c.Render(http.StatusOK, "checkdate_day.html", 1)
	})

	e.Start(":8080")

}
