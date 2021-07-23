package main

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"os/user"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type requriedinfo struct {
	Serverptc  string
	Serverip   string
	Clientport string
	Apiport    string
	Runport    string
	Dbid       string
	Dbname     string
	Dbport     string
	Dbpass     string
}

type postDetail struct {
	Id            string `json:"id"`             // 게시물 아이디
	Postdate      string `json:"postuploaddate"` // 게시물 날짜
	Title         string `json:"title"`          // 게시물 제목
	Content       string `json:"content"`        // 게시물 본문
	Thumbnail_URL string `json:"thumbURL"`       // 게시물 썸네일 URL
	Body          string `json:"body"`           // 게시물 썸네일 내용
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

func apis(api *gin.Engine, idinfo requriedinfo) { // 에러 시 패닉 사용
	api.POST("/api/upload", func(c *gin.Context) {
		var temp postDetail          // 요청 파싱할 구조체
		binderr := c.BindJSON(&temp) // 요청 파싱
		if binderr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": "Invalid or Missing Information.",
			})
			log.Panicln(binderr)
			return
		}

		if len([]rune(temp.Title)) > 200 { // 게시물 제목 200자 제한
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": "Title Length Over 200.",
			})
			log.Panicln("Title Length Over 200.")
			return
		}

		contentregex, regexerr := regexp.Compile("<(/)?([a-zA-Z]*)(\\s[a-zA-Z]*=[^>]*)?(\\s)*(/)?>") // html 태그 제거 정규식
		if regexerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(regexerr)
			return
		}

		body := contentregex.ReplaceAllString(temp.Content, "") // html 태그 제거
		if len([]rune(body)) > 400 {
			body = string([]rune(body)[0:400]) // 400자 썸네일에 쓸 본문 자름
		}

		if len([]rune(temp.Thumbnail_URL)) > 500 { // 썸네일 URL 500자 제한
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": "URL Is Too Long.",
			})
			log.Panicln("URL Is Too Long.")
			return
		}

		db, openerr := sql.Open("mysql", idinfo.Dbid+":"+idinfo.Dbpass+"@tcp("+idinfo.Serverip+":"+idinfo.Dbport+")/"+idinfo.Dbname) // db 오픈
		if openerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(openerr)
			return
		}
		defer db.Close()

		_, execerr := db.Exec("INSERT INTO postlist VALUES(0, now(), '" + html.EscapeString(temp.Title) + "', '" + html.EscapeString(temp.Content) + "', '" + html.EscapeString(temp.Thumbnail_URL) + "', '" + html.EscapeString(body) + "', " + strconv.FormatBool(temp.Secret) + ", 0 )") // SQL 실행
		if execerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(execerr)
			return
		}

		c.JSON(http.StatusCreated, gin.H{ // 생성됨 201 반환
			"code":     http.StatusCreated,
			"status":   "Post Was Created.",
			"redirect": idinfo.Serverptc + "://" + idinfo.Serverip + ":" + idinfo.Clientport,
		})
	})

	api.POST("/api/imageup", func(c *gin.Context) { // 이미지 처리 함수
		imagefile, formerr := c.GetPostForm("imagebase64") // 파싱
		if !formerr {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": "Image Value Not Found.",
			})
			log.Panicln("Image Value Not Found.")
			return
		}

		mat, _ := regexp.Compile("\\/[a-zA-Z]+;")     // 확장자 파싱 정규식
		data, _ := regexp.Compile(",[a-zA-Z0-9/+=]+") // base64 파싱 정규식
		if !mat.MatchString(imagefile) || !data.MatchString(imagefile) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln("Regex Not Matched")
			return
		}

		imgbase64 := strings.Trim(data.FindAllString(imagefile, -1)[0], ",") // base64 파싱
		exten := strings.Trim(mat.FindAllString(imagefile, -1)[0], "/;")     // 확장자 파싱

		filename := time.Now().Local().String() + "." + exten

		imgbytes, b64derr := base64.StdEncoding.DecodeString(imgbase64) // base64를 디코딩하여 이미지 바이너리 추출
		if b64derr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(b64derr)
			return
		}

		imgf, crefilerr := os.Create("uploadfiles/" + filename) // 이미지 파일 생성
		if crefilerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			imgf.Close()
			os.Remove("uploadfiles/" + filename)
			log.Panicln(crefilerr)
			return
		}
		defer imgf.Close()

		bytes, wrerr := imgf.Write(imgbytes) // 이미지 쓰기
		if wrerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(wrerr)
			return
		}

		log.Printf(filename+" %d 바이트 씀\n", bytes) // 로그

		c.JSON(http.StatusCreated, gin.H{
			"code":   http.StatusCreated,
			"status": "image upload success",
			"URL":    idinfo.Serverptc + "://" + idinfo.Serverip + ":" + idinfo.Apiport + "/images/" + filename,
		})
	})

	api.GET("/api/posts", func(c *gin.Context) { // 게시물 조회
		var name = c.Query("name")
		var page = c.Query("page")
		var id = c.Query("id")

		if name != "" && page != "" || page != "" && id != "" || id != "" && name != "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": "Too many Qurey.",
			})
			log.Printf("name: %s page: %s id: %s Too many Qurey.\n", name, page, id)
			return
		}

		var temp postDetail
		var rows *sql.Rows
		var queryerr error
		var sqldata []postDetail

		db, openerr := sql.Open("mysql", idinfo.Dbid+":"+idinfo.Dbpass+"@tcp("+idinfo.Serverip+":"+idinfo.Dbport+")/"+idinfo.Dbname)
		if openerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(openerr)
			return
		}
		defer db.Close()

		if name == "" && page == "" && id == "" {
			rows, queryerr = db.Query("select * from postlist order by post_id desc limit 0,5")
		} else if id != "" {
			idnum, iderr := strconv.Atoi(id)
			if iderr != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":   http.StatusBadRequest,
					"status": "Is Not number",
				})
				log.Panicln(iderr)
				return
			}

			if idnum < 1 {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":   http.StatusBadRequest,
					"status": "Is Not Positive number",
				})
				return
			}
			rows, queryerr = db.Query("select * from postlist where post_id = " + id)
			db.Exec("update postlist set views = views + 1 where post_id = " + id)
		} else if name != "" {
			rows, queryerr = db.Query("select * from postlist where post_title LIKE '%" + html.EscapeString(name) + "%' order by post_id desc")
		} else if page != "" {
			pagenum, _ := strconv.Atoi(page)
			if pagenum < 1 {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":   http.StatusBadRequest,
					"status": "Is Not Positive number",
				})
				return
			}
			rows, queryerr = db.Query("select * from postlist order by post_id desc limit " + strconv.Itoa((pagenum-1)*20) + "," + strconv.Itoa(pagenum*20))
		}
		defer rows.Close()

		if queryerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(queryerr)
			return
		}

		for rows.Next() {
			rowserr := rows.Scan(&temp.Id, &temp.Postdate, &temp.Title, &temp.Content, &temp.Thumbnail_URL, &temp.Body, &temp.Secret, &temp.Views)
			if rowserr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":   http.StatusInternalServerError,
					"status": "Server Error.",
				})
				log.Panicln(rowserr)
				return
			}
			temp.Title = html.UnescapeString(temp.Title)
			temp.Content = html.UnescapeString(temp.Content)
			temp.Body = html.UnescapeString(temp.Body)
			temp.Thumbnail_URL = html.UnescapeString(temp.Thumbnail_URL)

			sqldata = append(sqldata, temp)
		}

		c.JSON(http.StatusOK, sqldata)
	})

	api.POST("/api/deletepost", func(c *gin.Context) {
		postid, queryerr := c.GetQuery("postid")
		if !queryerr {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": "postid is Not Found.",
			})
			log.Panicln("postid is Not Found")
			return
		}
		_, numerr := strconv.Atoi(postid)
		if numerr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": "Is Not Number.",
			})
			log.Panicln(numerr)
			return
		}

		db, openerr := sql.Open("mysql", idinfo.Dbid+":"+idinfo.Dbpass+"@tcp("+idinfo.Serverip+":"+idinfo.Dbport+")/"+idinfo.Dbname)
		if openerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(openerr)
			return
		}
		defer db.Close()

		_, execerr := db.Exec("DELETE FROM postlist WHERE post_id = " + postid)
		if execerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(execerr)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":     http.StatusOK,
			"status":   postid + " is Deleted.",
			"redirect": idinfo.Serverptc + "://" + idinfo.Serverip + ":" + idinfo.Clientport + "/",
		})
	})

	api.GET("/api/comments", func(c *gin.Context) {
		var temp postComment
		var commentdata []postComment

		id, geterr := c.GetQuery("id")
		if !geterr {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": "id Not Found.",
			})
			log.Println("id Not Found.")
			return
		}
		_, numerr := strconv.Atoi(id)
		if numerr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": "Is Not number",
			})
			log.Panicln(numerr)
			return
		}

		db, openerr := sql.Open("mysql", idinfo.Dbid+":"+idinfo.Dbpass+"@tcp("+idinfo.Serverip+":"+idinfo.Dbport+")/"+idinfo.Dbname)
		if openerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(openerr)
			return
		}
		defer db.Close()

		rows, queryerr := db.Query("SELECT * FROM postcomment where post_id = " + id)
		if queryerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(queryerr)
			return
		}
		defer rows.Close()

		for rows.Next() {
			readerr := rows.Scan(&temp.Id, &temp.Comment_id, &temp.User_email, &temp.Upload_date, &temp.Comment, &temp.Secret)
			if readerr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":   http.StatusInternalServerError,
					"status": "Server Error.",
				})
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
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(binderr)
			return
		}

		if len([]rune(getjson.User_email)) > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": "Email Length Over 100.",
			})
			log.Panicln("Email Length Over 100.")
			return
		}
		if len([]rune(getjson.Comment)) > 2000 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": "Comment Length Over 2000.",
			})
			log.Panicln("Comment Length Over 2000.")
			return
		}

		db, openerr := sql.Open("mysql", idinfo.Dbid+":"+idinfo.Dbpass+"@tcp("+idinfo.Serverip+":"+idinfo.Dbport+")/"+idinfo.Dbname)
		if openerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(openerr)
			return
		}
		defer db.Close()

		_, exerr := db.Exec("INSERT INTO postcomment VALUES(" + getjson.Id + ", 0, '" + html.EscapeString(getjson.User_email) + "', now(), '" + html.EscapeString(getjson.Comment) + "', " + strconv.FormatBool(getjson.Secret) + ")")
		if exerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(exerr)
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"code":   http.StatusCreated,
			"status": "Upload was successful.",
		})
	})

	api.POST("/api/deletecomment", func(c *gin.Context) {
		comtid, queryerr := c.GetQuery("comtid")
		if !queryerr {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": "comtid Not Found.",
			})
			log.Panicln("comtid Not Found.")
			return
		}
		_, numerr := strconv.Atoi(comtid)
		if numerr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": "Is Not number.",
			})
			log.Panicln(numerr)
			return
		}

		db, openerr := sql.Open("mysql", idinfo.Dbid+":"+idinfo.Dbpass+"@tcp("+idinfo.Serverip+":"+idinfo.Dbport+")/"+idinfo.Dbname)
		if openerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(openerr)
			return
		}
		defer db.Close()

		_, execerr := db.Exec("DELETE FROM postcomment WHERE comment_id =" + comtid)
		if execerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
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
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(binderr)
			return
		}

		if len([]rune(temp.Title)) > 200 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": "Title Length Over 200.",
			})
			log.Panicln("Title Length Over 200.")
			return
		}

		contentregex, _ := regexp.Compile("<(/)?([a-zA-Z]*)(\\s[a-zA-Z]*=[^>]*)?(\\s)*(/)?>")
		var body = contentregex.ReplaceAllString(temp.Content, "")
		if len([]rune(body)) > 400 {
			body = string([]rune(body)[0:400])
		}
		if len([]rune(temp.Thumbnail_URL)) > 500 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": "URL Length Over 500.",
			})
			log.Panicln("thumbnail URL is too Long.")
			return
		}

		db, openerr := sql.Open("mysql", idinfo.Dbid+":"+idinfo.Dbpass+"@tcp("+idinfo.Serverip+":"+idinfo.Dbport+")/"+idinfo.Dbname)
		if openerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(openerr)
			return
		}
		defer db.Close()

		_, exeerr := db.Exec("UPDATE postlist SET post_title = '" + html.EscapeString(temp.Title) + "', post_content = '" + html.EscapeString(temp.Content) + "', thumbnail_url = '" + html.EscapeString(temp.Thumbnail_URL) + "', post_body = '" + html.EscapeString(body) + "', secret_post = " + strconv.FormatBool(temp.Secret) + " WHERE post_id = " + temp.Id)
		if exeerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Server Error.",
			})
			log.Panicln(exeerr)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":     http.StatusOK,
			"status":   "Edit was successful.",
			"redirect": idinfo.Serverptc + "://" + idinfo.Serverip + ":" + idinfo.Clientport + "/pages?id=" + temp.Id,
		})
	})
}

func handles(pagehandle *gin.Engine, idinfo requriedinfo) {
	pagehandle.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", []byte(`
			<html lang = "ko">
				<head>
					<meta charset = "utf-8">
					<title>Anend's API</title>
				</head>
				<body>
					<h1>제 개인 서버를 운영하기 위한 API 서버입니다.</h1>
					<p>서버로: <a href=`+idinfo.Serverptc+`://`+idinfo.Serverip+`:`+idinfo.Clientport+`>`+idinfo.Serverptc+`://`+idinfo.Serverip+`</a></p>
				</body>
			</html>
		`))
	})
}

func staticfiles(sf *gin.Engine) {
	sf.StaticFile("favicon.ico", "favicon.ico")
	sf.Static("/images", "uploadfiles/")
}

func Middleware(idinfo requriedinfo) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Origin", idinfo.Serverptc+"://"+idinfo.Serverip /*+":"+Clientport 개발용*/)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, HEAD, OPTIONS")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Cache-Control", "public, max-age=31536000")
		c.Header("Content-Encoding", "gzip")
		c.Header("Expires", time.Now().Local().Add(time.Second*86400).Format(time.RFC1123))

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

func initengine(idinfo requriedinfo) *gin.Engine {
	engine := gin.Default()
	engine.Use(Middleware(idinfo))
	engine.Use(gzip.Gzip(gzip.DefaultCompression))

	apis(engine, idinfo)
	handles(engine, idinfo)
	staticfiles(engine)

	return engine
}

func main() { // 에러 시 페이탈 사용
	/* 필수 정보 파싱 */
	var idinfo requriedinfo
	currentuser, _ := user.Current()            // 현재 소유자 uid, gid 얻기
	useruid, _ := strconv.Atoi(currentuser.Uid) // string 변환
	usergid, _ := strconv.Atoi(currentuser.Gid)

	idfiles, openerr := os.Open("idfiles.cfg") // 필수 초기화 파일 읽기
	if openerr != nil {
		log.Fatalln(openerr)
	}

	idscanner := bufio.NewScanner(idfiles) // 스캐너 오픈

	for i := 0; i < 9; i++ {
		idscanner.Scan()                                                     // 다음 스캔
		reflect.ValueOf(&idinfo).Elem().Field(i).SetString(idscanner.Text()) // 파싱
	}

	closeerr := idfiles.Close() // 설정 초기화 파일 닫기
	if closeerr != nil {
		log.Fatalln(closeerr)
	}

	/* 필수 정보 검증 */
	db, connerr := sql.Open("mysql", idinfo.Dbid+":"+idinfo.Dbpass+"@tcp("+idinfo.Serverip+":"+idinfo.Dbport+")/"+idinfo.Dbname)
	if connerr != nil {
		log.Fatalln(connerr)
	}
	db.Close()

	/* 로그 생성 */
	os.Mkdir("Logs", os.FileMode(0774))              // 로그폴더 생성 폴더가 존재할 수 있으니 에러 처리 X
	logsownerr := os.Chown("Logs", useruid, usergid) // 로그폴더 현재소유자 변경
	if logsownerr != nil {
		log.Fatalln(logsownerr)
	}

	data, logerr := os.OpenFile( // 로그 생성
		"Logs/"+time.Now().Local().String()+"-log.log",
		os.O_WRONLY|os.O_CREATE|os.O_SYNC,
		os.FileMode(0774))
	if logerr != nil {
		log.Fatalln(logerr)
	}
	defer data.Close()

	ownerr := data.Chown(useruid, usergid) // 로그 현재소유자 변경
	if ownerr != nil {
		log.Fatalln(ownerr)
	}

	log.Println("Server start.")
	log.SetOutput(data) // 로그 파일로 출력 리다이렉션
	log.Println("Logging start.")

	/* 서버 시작 */
	gin.SetMode(gin.ReleaseMode) // 모드 설정 (gin.ReleaseMode 배포모드, gin.DebugMode 디버그모드)
	mainengine := initengine(idinfo)

	if idinfo.Runport == "36530" { //개발용
		runerr := mainengine.Run(":" + idinfo.Runport)
		if runerr != nil {
			fmt.Println(runerr)
			log.Fatalln(runerr)
		} else {
			fmt.Println(time.Now().Local().String() + " http Server start.")
		}
	} else {
		runerr := mainengine.RunTLS(":"+idinfo.Runport, "/etc/letsencrypt/live/anend.site/cert.pem", "/etc/letsencrypt/live/anend.site/privkey.pem") // 배포용
		if runerr != nil {
			fmt.Println(runerr)
			log.Fatalln(runerr)
		} else {
			fmt.Println(time.Now().Local().String() + " https Server start.")
		}
	}
}
