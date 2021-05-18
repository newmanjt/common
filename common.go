package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/daddye/vips"
	"github.com/raff/godet"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func GetMin(arr []float64) (m float64) {
	for i, e := range arr {
		if i == 0 || e < m {
			m = e
		}
	}
	return
}

func GetMax(arr []float64) (m float64) {
	for i, e := range arr {
		if i == 0 || e > m {
			m = e
		}
	}
	return
}

func GetAve(arr []float64) (a float64) {
	i := 0
	var e float64
	for i, e = range arr {
		a = a + e
	}
	return a / float64(i)
}

//inspired by https://stackoverflow.com/questions/60147644/go-a-function-that-would-consume-maps-with-different-types-of-values
func GetKeys(m interface{}, typ reflect.Type) interface{} {
	if m == nil {
		return nil
	}
	if reflect.TypeOf(m).Kind() != reflect.Map {
		return nil
	}
	mapIter := reflect.ValueOf(m).MapRange()
	// mapVal := reflect.ValueOf(m).Interface()

	outputSlice := reflect.MakeSlice(reflect.SliceOf(typ), 0, 0)
	for mapIter.Next() {
		outputSlice = reflect.Append(outputSlice, mapIter.Key())
	}
	return outputSlice.Interface()
}

func ToFloatArr(arr []string) (new_arr []float64) {
	for _, el := range arr {
		new_arr = append(new_arr, ToFloat(el))
	}
	return
}

func Decimal(d decimal.Decimal) string {
	d1, _ := d.Float64()
	return ToString(d1)
}

func ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func ParseTime(t string, format string) (out_t time.Time) {
	out_t, _ = time.Parse(t, format)
	return
}

func Float(dec decimal.Decimal) (out float64) {
	out, _ = dec.Float64()
	return
}

func ToFloat(str string) (out float64) {
	if strings.Contains(str, "NaN") {
		out = -1
		return
	}
	out, err := strconv.ParseFloat(strings.TrimLeft(str, " "), 64)
	if err != nil {
		out = -1
		return
	}
	return
}

func LogMessage(source string, message string) {
	fmt.Printf("***** %s - %s: %s\n", source, strconv.FormatInt(time.Now().UnixNano(), 10), message)
}

func GetJSON(m map[string]float64) (out_string string) {
	out_map := make(map[string]string)
	for key, value := range m {
		out_map[key] = strconv.FormatFloat(value, 'f', -1, 64)
	}
	mjson, err := json.Marshal(out_map)
	if err != nil {
		return err.Error()
	}
	return string(mjson)
}

func In(l []string, el string) bool {
	for _, e := range l {
		if e == el {
			return true
		}
	}
	return false
}

func GoTo(s string, w http.ResponseWriter) {
	fmt.Fprintf(w, "<html><body><script>window.location.href='/"+s+"';</script></body></html>")
}

func LoadFile(filename string) ([]byte, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return body, err
}

func ServeImage(name string, resize bool) (outBuf []byte) {
	modelPicFile, err := os.Open(name)
	if err != nil {
		return
	}
	defer modelPicFile.Close()

	options := vips.Options{
		Width:        750,
		Height:       500,
		Crop:         false,
		Extend:       vips.EXTEND_WHITE,
		Interpolator: vips.BILINEAR,
		Gravity:      vips.CENTRE,
		Enlarge:      false,
		Quality:      90,
	}

	if resize {
		inBuf, _ := ioutil.ReadAll(modelPicFile)
		outBuf, err = vips.Resize(inBuf, options)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	} else {
		outBuf, _ = ioutil.ReadAll(modelPicFile)
	}
	return
}

func GetContentType(filename string) (contentType string) {
	if strings.HasSuffix(filename, ".eot") {
		contentType = "application/vnd.ms-fontobject"
	} else if strings.HasSuffix(filename, ".otf") {
		contentType = "application/x-font-opentype"
	} else if strings.HasSuffix(filename, ".svg") {
		contentType = "image/svg+xml"
	} else if strings.HasSuffix(filename, ".ttf") {
		contentType = "application/x-font-ttf"
	} else if strings.HasSuffix(filename, ".woff") {
		contentType = "font-woff"
	} else if strings.HasSuffix(filename, ".woff2") {
		contentType = "application/font-woff2"
	} else if strings.HasSuffix(filename, ".css") {
		contentType = "text/css"
		// contentType = ""
	} else if strings.HasSuffix(filename, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(filename, ".jpg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(filename, ".js") {
		contentType = "text/javascript"
	}
	return
}

func RunServer(handler func(http.ResponseWriter, *http.Request), port string) {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(port, nil))
}

func RunServerTLS(handler func(http.ResponseWriter, *http.Request), certFile string, keyFile string) {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServeTLS(":443", certFile, keyFile, nil))
}

func Serve(file string) (out string) {
	byteOut, err := LoadFile(file)
	if err != nil {
		out = "Error."
		return
	}
	out = string(byteOut)
	return
}

//handle SIGINT gracefully
func ManageExit() (c chan os.Signal) {
	c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	return
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func GetParam(r *http.Request, param string) (string, error) {
	params, ok := r.URL.Query()[param]
	if !ok || len(params[0]) < 1 {
		return "", errors.New(param + " missing")
	}
	return params[0], nil
}

type File struct {
	Name string
}

func GetDirectory(dir string) (files []File) {
	fs, err := os.ReadDir(dir)
	CheckError(err)
	for _, f := range fs {
		files = append(files, File{Name: f.Name()})
	}
	return
}

func GetString(files []File) (out string) {
	out += "<table>"
	out += "<tr><th>Index</th><th>Name</th></tr>"
	for i, f := range files {
		out += fmt.Sprintf("<tr><td>%d</td><td>%s</td></tr>", i, f.Name)
	}
	out += "</table>"
	return
}

//structure to request a screenshot
//we pass in the uuid of the picture and the tab to take a screenshot of
type ScreenshotRequest struct {
	ID  string
	Tab *godet.Tab
}
