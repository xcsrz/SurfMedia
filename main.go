package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	PORT          = "8888"
	VERBOSE       = true
	thumbnailSize = 300
	pageSize      = 100
)

var filetypes = map[string]string{}
var files []string
var threads int = 0
var knowkinds = map[string]string{
	"jpg":  "image",
	"jpeg": "image",
	"png":  "image",
	"gif":  "image",
	"gifv": "image",
	"bmp":  "image",
	"mov":  "video",
	"mkv":  "video",
	"mp4":  "video",
}
var headerTypes = map[string]string{
	"css":  "text/css",
	"js":   "text/javascript",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"mov":  "video/quicktime",
	"mp4":  "video/mp4",
}
var rootpath string

func main() {
	path, err := os.Getwd()
	checkErr(err)
	if err != nil {
		os.Exit(2)
	}
	rootpath = path
	threads = 1
	runDir(path)
	for threads > 0 {
		time.Sleep(1 * time.Second)
	}
	fmt.Println("found", len(filetypes), "files")
	if VERBOSE {
		// fmt.Println(filetypes)
	}
	for k, _ := range filetypes {
		files = append(files, k)
	}
	sort.Strings(files)
	serve()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	}
}

func fileExt(filename string) (ext string) {
	parts := strings.Split(filename, ".")
	return parts[len(parts)-1]
}

func fileKind(filename string) (ext string) {
	ext = knowkinds[fileExt(filename)]
	return
}

func runDir(path string) {
	if VERBOSE {
		fmt.Println("digging into", path)
	}
	files, err := ioutil.ReadDir(path)
	checkErr(err)
	for _, file := range files {
		// dig down into nested directories
		if file.IsDir() {
			threads += 1
			go runDir(path + string(os.PathSeparator) + file.Name())
		} else {
			kind := fileKind(file.Name())
			if len(kind) > 0 {
				// strings.Replace(path, os.Getwd(), "", 1)
				filetypes[strings.Replace(path, rootpath, "", 1)+string(os.PathSeparator)+file.Name()] = kind
			}

		}
	}
	threads -= 1
	if VERBOSE {
		fmt.Println("ending run for", path)
	}
}

func noCacheHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "max-age=0, public, must-revalidate, proxy-revalidate")
		w.Header().Add("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		h.ServeHTTP(w, r)
	})
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	if VERBOSE {
		fmt.Println("answering response for", r.URL)
	}
	output := ""
	pageNum, err := strconv.Atoi(r.URL.Query().Get("page"))
	checkErr(err)
	if err != nil {
		pageNum = 0
	}
	for i := (pageNum * pageSize); i < int(math.Min(float64(pageSize*(pageNum+1)), float64(len(files)))); i++ {
		output += `<div class="thumbnail-block">`
		if filetypes[files[i]] == "image" {
			output += fmt.Sprintf(`
				<a href="/file%s" rel="prettyPhoto[gallery]" title="%s">
					<img src="/file%s" title="%d" style="max-height: %dpx; max-width: %dpx" />
				</a>`, files[i], files[i], files[i], i, thumbnailSize, thumbnailSize)
		} else if filetypes[files[i]] == "video" {
			output += fmt.Sprintf(`
				<a href="/file%s?custom=true" rel="prettyPhoto[gallery]" title="%s">
					<img src="/static/video.png" width="120px" />
				</a>`, files[i], files[i])
		} else {
			output += fmt.Sprintf(`<img src="/static/error.jpg" width="%dpx" />`, thumbnailSize)
		}
		output += `</div>`
	}
	w.Write(wrapPage(output, pageNum))
}

func staticHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fileContents, found := staticFiles[r.URL.String()]; found {
			if contentType, found := headerTypes[fileExt(r.URL.String())]; found {
				w.Header().Set("Content-Type", contentType)
			}
			if VERBOSE {
				fmt.Println("serving static file", r.URL.String(), "(", len(fileContents), "bytes )")
			}
			w.Write(fileContents)
			return
		}
		if VERBOSE {
			fmt.Println("file NOT found:", r.URL.String())
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))

	})
}

func serve() {
	fs := noCacheHandler(http.FileServer(http.Dir(".")))

	http.Handle("/file/", http.StripPrefix("/file/", fs))
	http.HandleFunc("/page/", pageHandler)
	http.Handle("/thumb/", http.StripPrefix("/thumb/", thumbnailHandler()))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler()))
	fmt.Println("Listening on port", PORT, " ...")
	fmt.Println(http.ListenAndServe(":"+PORT, nil))
}

func wrapPage(output string, pageNum int) []byte {
	pages := pageLinks(pageNum)
	output = `<!DOCTYPE html>
<html>
<head>
<title>SurfMedia</title>
<link href="/static/prettyphoto/css/prettyPhoto.css" rel="stylesheet">
<script type="text/javascript" src="/static/jquery.js"></script>
<script type="text/javascript" src="/static/prettyphoto/js/jquery.prettyPhoto.js"></script>
<style>
.thumbnail-block {
	display: inline;
	float: left;
	width: ` + strconv.Itoa(thumbnailSize) + `px;
	height: ` + strconv.Itoa(thumbnailSize) + `px;
	text-align: center;
}
.pagelink {
	background-image: url(/static/emptydot.png);
	color: #eee;
}
.pagelink, .pagejumplink {
	background-repeat: no-repeat;
	display: inline-block;
	cursor: pointer;
	width: 22px;
	height: 22px;
	text-align: centers;
	padding: 0px;
	margin: 0px;
	font: 12px arial;
	padding-top: 4px;
}
.pagelinkrow {
	text-align: center;
	clear: both;
	padding: 0px;
	margin: 0px;
}
</style>
</head>
<body>
` + pages + `<div style="text-align:center">` + output + `</div>` + pages + `
<script type="text/javascript" charset="utf-8">
  $(document).ready(function(){
    $("a[rel^='prettyPhoto']").prettyPhoto({social_tools: false, theme: 'dark_rounded', custom_markup: '<video src="{path}" style="max-width: 1000px;" controls >'});
  });
</script>
</body>
</html>`
	return []byte(output)
}

func pageLinks(pageNum int) string {
	var links []string
	var numPages = int(math.Ceil(float64(len(files) / pageSize)))
	if pageNum > 0 {
		links = append(links, `<span class="pagejumplink" style="background-image: url(/static/prev.png);" onclick="location.href='?page=`+strconv.Itoa(pageNum-1)+`'">&nbsp;</span>`)
	}
	for i := 0; i <= numPages-1; i++ {
		links = append(links, fmt.Sprintf(`<span class="pagelink" onclick="location.href='?page=`+strconv.Itoa(i)+`'">%d</span>`, i+1))
	}
	if ((pageNum + 1) * pageSize) < len(files) {
		links = append(links, `<span class="pagejumplink" style="background-image: url(/static/next.png);" onclick="location.href='?page=`+strconv.Itoa(pageNum+1)+`'">&nbsp;</span>`)
	}
	return `<div class="pagelinkrow">` + strings.Join(links, "&nbsp;") + `</div><div style="text-align: center"> Page ` + strconv.Itoa(pageNum) + ` of ` + strconv.Itoa(numPages) + `</div>`
}
