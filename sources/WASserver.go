package main

import (
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

/*
const serverpct string = "http"
const serverid string = "127.0.0.1" // 개발 용
const clientport string = "3000"
const apiport string = "36530"
const dbid string = "root"
const dbport string = "3306"
const dbname string = "devblog"
const runport string = ":36530"
*/

const serverpct string = "https"
const serverid string = "anend.site" //서버 용
const clientport string = "443"
const apiport string = "53373"
const dbid string = "blogSQLmaster"
const dbport string = "53374"
const dbname string = "myblogDB"
const runport string = ":50000"

type postDetail struct {
	Id            string `json:"id"`
	Postdate      string `json:"postuploaddate"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Thumbnail_URL string `json:"thumbURL"`
	Body          string `json:"body"`
	Secret        bool   `json:"secret"`
	Views         int    `json:"views"`
}

type postComment struct {
	Id          string `json:"postid"`
	Comment_id  int64  `json:"comtid"`
	User_email  string `json:"email"`
	Upload_date string `json:"date"`
	Comment     string `json:"comment"`
	Secret      bool   `json:"secret"`
}

func servererr(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":   http.StatusInternalServerError,
		"status": "Server error",
	})
}

func clienterr(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":   http.StatusBadRequest,
		"status": "Client error",
	})
}

func apis(api *gin.Engine, dbpass string) {
	api.POST("/api/upload", func(c *gin.Context) {
		var temp postDetail
		binderr := c.BindJSON(&temp)
		if binderr != nil {
			clienterr(c)
			log.Panicln(binderr)
			return
		}

		if len([]rune(temp.Title)) > 200 {
			temp.Title = string([]rune(temp.Title)[0:200])
		}

		contentregex, _ := regexp.Compile("<(/)?([a-zA-Z]*)(\\s[a-zA-Z]*=[^>]*)?(\\s)*(/)?>")
		var body = []rune(contentregex.ReplaceAllString(temp.Content, ""))
		if len(body) > 400 {
			body = body[0:400]
		}
		if len([]rune(temp.Thumbnail_URL)) > 500 {
			clienterr(c)
			log.Println("thumbnail URL is too Long.")
			return
		}

		db, openerr := sql.Open("mysql", dbid+":"+dbpass+"@tcp("+serverid+":"+dbport+")/"+dbname)

		defer db.Close()
		if openerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "upload failed",
			})
			log.Panicln(openerr)
			return
		}

		_, exeerr := db.Exec("INSERT INTO postlist VALUES(0, now(), '" + html.EscapeString(temp.Title) + "', '" + html.EscapeString(temp.Content) + "', '" + html.EscapeString(temp.Thumbnail_URL) + "', '" + string(body) + "', " + strconv.FormatBool(temp.Secret) + ", 0 )")
		if exeerr != nil {
			servererr(c)
			log.Panicln(exeerr)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":     http.StatusOK,
			"status":   "Upload was successful.",
			"redirect": serverpct + "://" + serverid + ":" + clientport,
		})
	})

	api.POST("/api/imageup", func(c *gin.Context) {
		imagefile, formerr := c.GetPostForm("imagebase64")
		if !formerr {
			clienterr(c)
			log.Println("Image Value Not found.")
			return
		}

		mat, _ := regexp.Compile("\\/[a-zA-Z]+;")
		data, _ := regexp.Compile(",[a-zA-Z0-9/+=]+")
		if !mat.MatchString(imagefile) || !data.MatchString(imagefile) {
			servererr(c)
			log.Panicln("Regex Not Matched")
			return
		}

		imgbase64 := strings.Trim(data.FindAllString(imagefile, -1)[0], ",")
		exten := strings.Trim(mat.FindAllString(imagefile, -1)[0], "/;")

		t := time.Now()
		log.Println(t)

		name := strconv.Itoa(t.Year()) + "-" +
			strconv.Itoa(int(t.Month())) + "-" +
			strconv.Itoa(t.Day()) + "_" +
			strconv.Itoa(t.Hour()) + ":" +
			strconv.Itoa(t.Minute()) + ":" +
			strconv.Itoa(t.Second()) + "_" +
			exten //YYYY-MM-DD_HH:MM:SS_png

		filename := "uploadfiles/" + name + "." + exten

		imgbytes, b64derr := base64.StdEncoding.DecodeString(imgbase64)
		if b64derr != nil {
			servererr(c)
			log.Panicln(b64derr)
			return
		}

		imgf, crefilerr := os.Create(filename)
		defer imgf.Close()
		if crefilerr != nil {
			servererr(c)
			imgf.Close()
			os.Remove(filename)
			log.Panicln(crefilerr)
			return
		}

		bytes, wrerr := imgf.Write(imgbytes)
		if wrerr != nil {
			servererr(c)
			log.Panicln(wrerr)
			return
		}

		log.Printf(filename+" %d 바이트 씀\n", bytes)
		c.JSON(http.StatusOK, gin.H{
			"code":   http.StatusOK,
			"status": "image upload success",
			"URL":    serverpct + "://" + serverid + ":" + apiport + "/images/" + name + "." + exten, // 개발 용
		})
	})

	api.GET("/api/posts", func(c *gin.Context) {
		var name = c.Query("name")
		var page = c.Query("page")
		var id = c.Query("id")

		if name != "" && page != "" || page != "" && id != "" || id != "" && name != "" {
			clienterr(c)
			log.Println("Too many Qureys.")
			return
		}

		var temp postDetail
		var rows *sql.Rows
		var queryerr error
		var sqldata []postDetail

		db, openerr := sql.Open("mysql", dbid+":"+dbpass+"@tcp("+serverid+":"+dbport+")/"+dbname)

		defer db.Close()
		if openerr != nil {
			servererr(c)
			log.Panicln(openerr)
			return
		}

		if name == "" && page == "" && id == "" {
			rows, queryerr = db.Query("select * from postlist order by post_id desc limit 0,5")
		} else if id != "" {
			rows, queryerr = db.Query("select * from postlist where post_id = " + id)
			db.Exec("update postlist set views = views + 1 where post_id = " + id)
		} else if name != "" {
			rows, queryerr = db.Query("select * from postlist where post_title LIKE '%" + name + "%' order by post_id desc")
		} else if page != "" {
			pagenum, _ := strconv.Atoi(page)
			rows, queryerr = db.Query("select * from postlist order by post_id desc limit " + strconv.Itoa((pagenum-1)*20) + "," + strconv.Itoa(pagenum*20))
		}
		defer rows.Close()
		if queryerr != nil {
			servererr(c)
			log.Panicln(queryerr)
			return
		}

		for rows.Next() {
			rowserr := rows.Scan(&temp.Id, &temp.Postdate, &temp.Title, &temp.Content, &temp.Thumbnail_URL, &temp.Body, &temp.Secret, &temp.Views)
			if rowserr != nil {
				servererr(c)
				log.Panicln(rowserr)
				return
			}
			temp.Title = html.UnescapeString(temp.Title)
			temp.Content = html.UnescapeString(temp.Content)
			temp.Body = html.UnescapeString(temp.Body)

			sqldata = append(sqldata, temp)
		}

		c.JSON(http.StatusOK, sqldata)
	})

	api.POST("/api/deletepost", func(c *gin.Context) {
		postid, queryerr := c.GetQuery("postid")
		if !queryerr {
			clienterr(c)
			log.Println("id is Not Found")
			return
		}

		db, openerr := sql.Open("mysql", dbid+":"+dbpass+"@tcp("+serverid+":"+dbport+")/"+dbname)

		defer db.Close()
		if openerr != nil {
			servererr(c)
			log.Panicln(openerr)
			return
		}

		_, execerr := db.Exec("DELETE FROM postlist WHERE post_id = " + postid)
		if execerr != nil {
			servererr(c)
			log.Panicln(execerr)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":     http.StatusOK,
			"status":   postid + " is Deleted.",
			"redirect": serverpct + "://" + serverid + ":" + clientport + "/",
		})
	})

	api.GET("/api/comments", func(c *gin.Context) {
		var temp postComment
		var commentdata []postComment

		id, geterr := c.GetQuery("id")
		if !geterr {
			clienterr(c)
			log.Println("id Not Found.")
			return
		}

		db, openerr := sql.Open("mysql", dbid+":"+dbpass+"@tcp("+serverid+":"+dbport+")/"+dbname)

		defer db.Close()
		if openerr != nil {
			servererr(c)
			log.Panicln(openerr)
			return
		}

		rows, queryerr := db.Query("SELECT * FROM postcomment where post_id = " + id)
		defer rows.Close()
		if queryerr != nil {
			servererr(c)
			log.Panicln(queryerr)
			return
		}

		for rows.Next() {
			readerr := rows.Scan(&temp.Id, &temp.Comment_id, &temp.User_email, &temp.Upload_date, &temp.Comment, &temp.Secret)
			if readerr != nil {
				servererr(c)
				log.Panicln(readerr)
				return
			}

			temp.User_email = html.UnescapeString(temp.User_email)
			temp.Comment = html.UnescapeString(temp.Comment)

			commentdata = append(commentdata, temp)
		}

		c.JSON(http.StatusOK, commentdata)
	})

	api.POST("/api/comments", func(c *gin.Context) {
		var getjson postComment
		binderr := c.BindJSON(&getjson)
		if binderr != nil {
			clienterr(c)
			log.Panicln(binderr)
			return
		}

		if len([]rune(getjson.User_email)) > 100 {
			getjson.User_email = string([]rune(getjson.User_email)[0:100])
		}
		if len([]rune(getjson.Comment)) > 2000 {
			getjson.Comment = string([]rune(getjson.Comment)[0:2000])
		}

		db, openerr := sql.Open("mysql", dbid+":"+dbpass+"@tcp("+serverid+":"+dbport+")/"+dbname)

		defer db.Close()
		if openerr != nil {
			servererr(c)
			log.Panicln(openerr)
			return
		}

		_, exerr := db.Exec("INSERT INTO postcomment VALUES(" + getjson.Id + ", 0, '" + html.EscapeString(getjson.User_email) + "', now(), '" + html.EscapeString(getjson.Comment) + "', " + strconv.FormatBool(getjson.Secret) + ")")
		if exerr != nil {
			servererr(c)
			log.Panicln(exerr)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   http.StatusOK,
			"status": "Upload was successful.",
		})
	})

	api.POST("/api/deletecomment", func(c *gin.Context) {
		comtid, queryerr := c.GetQuery("comtid")
		if !queryerr {
			clienterr(c)
			log.Panicln("comtid Not Found.")
			return
		}

		db, openerr := sql.Open("mysql", dbid+":"+dbpass+"@tcp("+serverid+":"+dbport+")/"+dbname)

		defer db.Close()
		if openerr != nil {
			servererr(c)
			log.Panicln(openerr)
			return
		}

		_, execerr := db.Exec("DELETE FROM postcomment WHERE comment_id =" + comtid)
		if execerr != nil {
			servererr(c)
			log.Panicln(execerr)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":   http.StatusOK,
			"status": comtid + " is Deleted.",
		})
	})

	api.POST("/api/editpost", func(c *gin.Context) {
		var temp postDetail
		binderr := c.BindJSON(&temp)
		if binderr != nil {
			clienterr(c)
			log.Panicln(binderr)
			return
		}

		if len([]rune(temp.Title)) > 200 {
			temp.Title = string([]rune(temp.Title)[0:200])
		}

		contentregex, _ := regexp.Compile("<(/)?([a-zA-Z]*)(\\s[a-zA-Z]*=[^>]*)?(\\s)*(/)?>")
		var body = []rune(contentregex.ReplaceAllString(temp.Content, ""))
		if len(body) > 400 {
			body = body[0:400]
		}
		if len([]rune(temp.Thumbnail_URL)) > 500 {
			clienterr(c)
			log.Println("thumbnail URL is too Long.")
			return
		}

		db, openerr := sql.Open("mysql", dbid+":"+dbpass+"@tcp("+serverid+":"+dbport+")/"+dbname)

		defer db.Close()
		if openerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "edit failed",
			})
			log.Panicln(openerr)
			return
		}

		_, exeerr := db.Exec("UPDATE postlist SET post_title = '" + html.EscapeString(temp.Title) + "', post_content = '" + html.EscapeString(temp.Content) + "', thumbnail_url = '" + html.EscapeString(temp.Thumbnail_URL) + "', post_body = '" + html.EscapeString(string(body)) + "', secret_post = " + strconv.FormatBool(temp.Secret) + " WHERE post_id = " + temp.Id)
		if exeerr != nil {
			servererr(c)
			log.Panicln(exeerr)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":     http.StatusOK,
			"status":   "Edit was successful.",
			"redirect": serverpct + "://" + serverid + ":" + clientport + "/pages?id=" + temp.Id,
		})
	})
}

func handles(pagehandle *gin.Engine) {
	pagehandle.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", []byte(`
			<html lang = "ko">
				<head>
					<meta charset = "utf-8">
					<title>Anend's API</title>
				</head>
				<body>
					<h1>제 개인 서버를 운영하기 위한 API 서버입니다.</h1>
					<p>개인 서버 외에 API 서버 호출은 불가능합니다.</p>
					<p>서버로: <a href=`+serverpct+`://`+serverid+`:`+clientport+`>`+serverpct+`://`+serverid+`</a></p>
				</body>
			</html>
		`))
	})
}

func staticfiles(sf *gin.Engine) {
	sf.StaticFile("favicon.ico", "favicon.ico")
	sf.Static("/images", "uploadfiles/")
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "Content-Type, Origin")
		c.Header("Access-Control-Allow-Origin", serverpct+"://"+serverid /*+":"+clientport 개발용*/)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, HEAD, OPTIONS")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Cache-Control", "public, max-age=31536000")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		if c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
			return
		}
		c.Next()
	}
}

func initengine(dbpass string) *gin.Engine {
	engine := gin.Default()
	engine.Use(Middleware())

	apis(engine, dbpass)
	handles(engine)
	staticfiles(engine)

	return engine
}

func checkcreatefile(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}

func main() {
	var dbpass *string = flag.String("p", "", "-p [DBpassword] requried")

	flag.Parse()

	if *dbpass == "" {
		fmt.Println("-p [DBpassword] requried")
		return
	}

	db, _ := sql.Open("mysql", dbid+":"+*dbpass+"@tcp("+serverid+":"+dbport+")/"+dbname)

	if db.Ping() != nil {
		fmt.Println("Database Password is not match.")
		db.Close()
		return
	}
	db.Close()

	t := time.Now()
	log.Println(t.Local().String() + " Anend API 웹서버를 가동합니다.")

	os.Mkdir("Logs", os.FileMode(0774))
	logsownerr := os.Chown("Logs", 1000, 1000)
	checkcreatefile(logsownerr)

	os.Mkdir("WASfiles", os.FileMode(0774))
	WASownerr := os.Chown("WASfiles", 1000, 1000)
	checkcreatefile(WASownerr)

	data, logerr := os.OpenFile(
		"Logs/"+
			strconv.Itoa(t.Year())+
			strconv.Itoa(int(t.Month()))+
			strconv.Itoa(t.Day())+"-"+
			strconv.Itoa(t.Hour())+
			strconv.Itoa(t.Minute())+
			strconv.Itoa(t.Second())+
			"-log.log",
		os.O_WRONLY|os.O_CREATE|os.O_SYNC,
		os.FileMode(0774))
	logown := data.Chown(1000, 1000)
	checkcreatefile(logown)

	defer data.Close()
	checkcreatefile(logerr)

	log.SetOutput(data)
	log.Println("Server start.")

	gin.SetMode(gin.ReleaseMode)
	mainengine := initengine(*dbpass)

	if runport == ":36530" {
		mainengine.Run(runport)
	} else {
		mainengine.RunTLS(runport, "/etc/letsencrypt/live/anend.site/cert.pem", "/etc/letsencrypt/live/anend.site/privkey.pem")
	}
}
