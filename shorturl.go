package main

import (
	"fmt"
	"strings"
	"strconv"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"github.com/hoisie/redis"
	"os"
	"log"
	"syscall"
	"runtime"
)

func main() {

	daemon(0, 0)
	runtime.GOMAXPROCS(1)

	port := "8887"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	http.HandleFunc("/", hello)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		fmt.Println(err)
	}
	os.Exit(0)
}


func hello (w http.ResponseWriter, req *http.Request) {

	var client redis.Client
	client.Addr			= "127.0.0.1:6379"
	client.Db 			= 0
	client.MaxPoolSize		= 1
	client.Password 		= "~xxxx"

	req.ParseForm()

	if s := req.Form.Get("s") ; len(s) > 0 {

		url := s
		_, urlOutput, _ := genShortUrl(url);

		ch := make([]chan int, len(urlOutput) * 2)
		for i, v := range urlOutput {
			ch[i] = setToCache(i, v, url, &client)
		}
		for i := range urlOutput {
			<-ch[i]
		}
		w.Write([]byte(fmt.Sprintf("<a href='http://%s/%s'>http://%s/%s</a>", "w.adango.cn", urlOutput[1], "w.adango.cn", urlOutput[1])))

	} else {
	 	path := req.URL.Path[1:]
		if err, url := getFromCache(path, &client); (err == nil) && (len(url) > 0) {
			http.Redirect(w, req, url, http.StatusFound)
		} else {
//			fmt.Println(len(url))
//			fmt.Println(err)
		}
	}
}

/**
 * 写入redis 这个redis类是惰性连接
 */
func setToCache(area int, key string, val string, client *redis.Client)  chan int {

	ch := make(chan int, 2)

	go func() {
		if _, err := client.Hset(string(strconv.AppendInt([]byte("go::url::"), int64(area), 10)), key, []byte(val)); err != nil {
			fmt.Println(err)
		}
		ch <- area
	}()
	return ch
}

/**
 * 读取redis 这个redis类是惰性连接
 */
func getFromCache(key string, client *redis.Client) (err error, val string) {

	var bu []byte
	bu, err = client.Hget("go::url::1", key);
	return err, string(bu)
}

// 生成短链接
func genShortUrl(url string) (err error, output []string, str string) {

	seed := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	h := md5.New()
	h.Write([]byte(url))
	hex := strings.ToLower(hex.EncodeToString(h.Sum(nil)))

	hexLen := len(hex)
	subHexLen := hexLen / 8

	for  i := 0; i < subHexLen; i++ {
		subHex := hex[(i * 8) : ((i + 1) * 8)]
		intb, _ := strconv.ParseInt(subHex, 16, 0)

		intb = 0x3FFFFFFF & intb
		out := []byte{}

		for i := 0; i < 6; i++ {
			val :=  0x0000003D & intb
			out = append(out, seed[val])
			intb = intb >> 5
		}
		output = append(output, string(out))
	}
	return nil, output, fmt.Sprintf("http://1010g.net/%s", output[0])
}

/**

 */
func daemon(nochdir, noclose int) int {

    var ret, ret2 uintptr
    var err syscall.Errno
 
    darwin := runtime.GOOS == "darwin"
 
    // already a daemon
    if syscall.Getppid() == 1 {
        return 0
    }
 
    // fork off the parent process
    ret, ret2, err = syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
    if err != 0 {
        return -1
    }
 
    // failure
    if ret2 < 0 {
        os.Exit(-1)
    }
 
    // handle exception for darwin
    if darwin && ret2 == 1 {
        ret = 0
    }
 
    // if we got a good PID, then we call exit the parent process.
    if ret > 0 {
        os.Exit(0)
    }
 
    /* Change the file mode mask */
    _ = syscall.Umask(0)
 
    // create a new SID for the child process
    s_ret, s_errno := syscall.Setsid()
    if s_errno != nil {
        log.Printf("Error: syscall.Setsid errno: %d", s_errno)
    }
    if s_ret < 0 {
        return -1
    }
 
    if nochdir == 0 {
        os.Chdir("/")
    }
 
    if noclose == 0 {
        f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)
        if e == nil {
            fd := f.Fd()
            syscall.Dup2(int(fd), int(os.Stdin.Fd()))
            syscall.Dup2(int(fd), int(os.Stdout.Fd()))
            syscall.Dup2(int(fd), int(os.Stderr.Fd()))
        }
    }
 
    return 0
}
 
