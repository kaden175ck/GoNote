package main

import (
	"bytes"
	"fmt"
	"html/template"
	"image"
	_ "image/gif" // Blank import to register GIF decoder
	"image/jpeg"  //  register JPEG decoder
	"image/png"   //  register PNG decoder
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var contents = []string{"蜀道之难，难于上青天！", "蚕丛及鱼凫，开国何茫然！", "尔来四万八千岁，不与秦塞通人烟。", "西当太白有鸟道，可以横绝峨眉巅。", "地崩山摧壮士死，然后天梯石栈相钩连。", "上有六龙回日之高标，下有冲波逆折之回川。", "黄鹤之飞尚不得过，猿猱欲度愁攀援。"}

func ChunkedTransfer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Transfer-Encoding", "chunked") // 数据分块传送。一般情况下， client接收数据长度到达Content-Length后才开始渲染，但由于此处我们是分批渲染，所以不需要（也不能）设置Content-Length
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "<html><body><ol>")
	for _, chunk := range contents {
		fmt.Fprintf(w, "<li>%s</li>", chunk) // 发送一块数据
		flusher.Flush()                      // 强制数据立刻发给对方
		time.Sleep(time.Second)              // 故意卡顿一下，前端分段渲染效果更明显
	}
	fmt.Fprintf(w, "</ol></body></html>")
}

// Server-Sent Events。基于HTTP的单向数据流技术，服务端可实时向客户端推送数据。客户端在连接断开后会自动尝试重连
func SSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/event-stream; charset=utf-8") //标识响应为事件流。charset=utf-8是为了解决中文乱码
	//w.Header().Add("Cache-Control", "no-cache")                        //防止浏览器缓存响应，确保实时性
	//w.Header().Add("Connection", "keep-alive")                         //保持连接开放，支持持续流式传输

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for _, chunk := range contents {
		fmt.Fprintf(w, "data: %s\n\n", chunk) // 发送一条数据
		flusher.Flush()                       // 强制数据立刻发给对方
		time.Sleep(time.Second)               // 故意卡顿一下，前端分段渲染效果更明显
	}

	fmt.Fprint(w, "data: [DONE]\n\n") // 结束标志
	flusher.Flush()
}

func ImageStream(w http.ResponseWriter, r *http.Request) {
	fileName := r.PathValue("file_name")
	file, err := os.Open("data/" + fileName)
	if err != nil {
		http.Error(w, "文件找不到", http.StatusNotFound)
		return
	}
	defer file.Close()

	img, formatName, err := image.Decode(file)
	if err != nil {
		http.Error(w, "图像解码失败", http.StatusInternalServerError)
		return
	}
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	rgba := image.NewRGBA(img.Bounds())
	for i := 0; i < img.Bounds().Max.X; i++ {
		for j := 0; j < img.Bounds().Max.Y; j++ {
			rgba.Set(i, j, img.At(i, j))
		}
	}

	const boundary = "--myboundary"
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary="+boundary) // 用于发送图片或视频
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	const P = 10
	step := height/P + 1 // 加1是为了防止step等于0
	for i := 0; ; i++ {
		top := i * step
		buttom := (i + 1) * step
		if buttom > img.Bounds().Max.Y {
			buttom = img.Bounds().Max.Y
		}

		segment, err := GetImgSegment(rgba, formatName, top, buttom, width)
		if err != nil {
			log.Printf("GetImgSegment error: %s", err)
			break
		}

		fmt.Fprintf(w, "\r\n%s\r\n", boundary)
		fmt.Fprintf(w, "Content-Type: image/%s\r\n", formatName)
		fmt.Fprintf(w, "Content-Length: %d\r\n\r\n", len(segment)) // 每一份数据有单独的header
		w.Write(segment)                                           // 每一份数据可以单独解码
		flusher.Flush()

		if buttom >= img.Bounds().Max.Y {
			break
		}

		time.Sleep(time.Second)
	}
}

func ReadImage(fileName string) ([]byte, error) {
	file, err := os.Open("data/" + fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	_, err = io.Copy(buffer, file)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func GetImgSegment(img image.Image, formatName string, top, buttom, width int) ([]byte, error) {
	//subImg := image.NewRGBA(img.Bounds())   // 各部分保持在原图像中的位置
	subImg := image.NewRGBA(image.Rect(0, 0, width, buttom-top))
	for i := 0; i < width; i++ {
		for j := top; j < buttom; j++ {
			//subImg.Set(i, j, img.At(i, j))   // 各部分保持在原图像中的位置
			subImg.Set(i, j-top, img.At(i, j))
		}
	}

	bs := make([]byte, 0, 2048) // 即使subImg超过了2K，bytes.Buffer也会自动扩容的
	buffer := bytes.NewBuffer(bs)
	if formatName == "png" {
		err := png.Encode(buffer, subImg) // 把subImg进行png编码，然后写入buffer
		if err != nil {
			return nil, err
		}
	} else if formatName == "jpeg" {
		err := jpeg.Encode(buffer, subImg, &jpeg.Options{Quality: 100})
		if err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

func Animation(w http.ResponseWriter, r *http.Request) {
	const boundary = "--myboundary"
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary="+boundary) // 用于发送图片或视频
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	images := []string{"prepare.png", "go.png", "run.png"}
	for _, image := range images {
		segment, err := ReadImage(image)
		if err != nil {
			http.Error(w, "read imgae error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "\r\n%s\r\n", boundary)        // 分界线
		fmt.Fprintf(w, "Content-Type: image/png\r\n") // 每一份数据有单独的header
		fmt.Fprintf(w, "Content-Length: %d\r\n\r\n", len(segment))
		w.Write(segment)
		flusher.Flush()
		time.Sleep(time.Second)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /chunk", ChunkedTransfer) // http://127.0.0.1:5678/chunk
	mux.HandleFunc("GET /server_source_event", SSE)
	mux.HandleFunc("GET /sse", func(writer http.ResponseWriter, request *http.Request) {
		tmpl, err := template.ParseFiles("./http/server/sse.html") //相对于执行go run的路径
		if err != nil {
			fmt.Println("create template failed:", err)
			return
		}
		tmpl.Execute(writer, map[string]string{"url": "http://127.0.0.1:5678/server_source_event"})
	}) // http://127.0.0.1:5678/sse
	mux.HandleFunc("GET /img/{file_name}", ImageStream) // http://127.0.0.1:5678/img/logo.png
	mux.HandleFunc("GET /animation", Animation)         // http://127.0.0.1:5678/animation
	if err := http.ListenAndServe("127.0.0.1:5678", mux); err != nil {
		panic(err)
	}

	//SplitImage("logo.png")
}

func SplitImage(fileName string) {
	file, err := os.Open("data/" + fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, formatName, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	const P = 10
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	step := height/P + 1 // 加1是为了防止step等于0
	for i := 0; i <= P; i++ {
		top := i * step
		buttom := (i + 1) * step
		if buttom > img.Bounds().Max.Y {
			buttom = img.Bounds().Max.Y
		}

		segment, err := GetImgSegment(img, formatName, top, buttom, width)
		if err != nil {
			log.Printf("GetImgSegment error: %s", err)
			break
		}

		outFile, err := os.Create("data/" + strconv.Itoa(i) + "." + formatName)
		defer outFile.Close()
		if err != nil {
			panic(err)
		}
		outFile.Write(segment)

		if buttom >= img.Bounds().Max.Y {
			break
		}
	}
}
