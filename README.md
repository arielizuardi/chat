# chat
Make chat application with Golang based on Go Programming Blueprint Book By Mat Ryer with some modified oauth2 library and using Godeps dependency

## Requirements

[Go 1.7](https://golang.org/dl/)
Set your $GOPATH


Install [Godep](https://github.com/tools/godep) dependency tools,
```
go get github.com/tools/godep
cd <your/chat/path>
godep install
go build -o chat
```

To RUN the apps
```
./chat --addr=":3000"
```

Open Your Browser and login with your Google Account
```
http://localhost:3000
```

# Tutorial

## Membuat aplikasi Chat menggunakan Golang, OAuth2, dan Gravatar



**Requirements: **

* Golang 1.7

  * Installation `brew install go`
  * Set your `$GOPATH` and Workspace. Read [here](https://golang.org/doc/code.html)

* Godep https://github.com/tools/godep - This is tools to manage your dependency.

  > Godep helps build packages reproducibly by fixing their dependencies.

* Your favourite Editor - I use [Atom](https://atom.io/), with [Goplus](https://atom.io/packages/go-plus) plugins



### HTTP Server pertama!

Buat folder `<$GOPATH>/src/github.com/<username>/chat`

Dan didalam folder tersebut buat file `main.go`

```go
package main

import (
    "log"
	"net/http"
)

func main() {
     http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
       w.Write([]byte(`
         <html>
           <head>
             <title>Chat</title>
￼￼￼￼￼￼￼￼￼			</head>
		   <body>
				Let's chat!
			 </body>
         </html>
	  `))
     })
     // start the web server
     if err := http.ListenAndServe(":8080", nil); err != nil {
       log.Fatal("ListenAndServe:", err)
     }
}
```

`http.HandleFunc` memetakan path `/` dengan *function* yang menangani path tersebut

`http.Request` adalah struktur data yang merepresentasikan *request* dari client.

`http.ResponseWriter` mengumpulkan *value-value* yang dikembalikan *server* sebagai *response*

`http.ListenAndServe` menentukan di port berapa server harus me-*listen*, disini kita menggunakan port 8080

Mari kita jalankan aplikasi kita, `go build -o chat`, `./chat` dan buka browser di `http://localhost:8080`



### Templates dan pengenalan Interface

Templates membantu kita untuk menghasilkan layout yang dinamis, dalam tutorial ini kita akan menggunakan library `text/template` dari Go

Buat folder `templates` dan di dalam folder tersebut buat file `chat.html`

```html
<html>
     <head>
       <title>Chat</title>
     </head>
	<body>
		Let's chat (from template)
     </body>
 </html>
```

Sekarang kita memiliki file html yang soap digunakan, namun kita harus meng-*compile* tersebut agar bisa digunakan oleh server

Di `main.go` , tepat di atas method main buatlah sebuah `struct` dengan nama `templateHandler` dengan method `ServeHTTP`

`Struct` adalah tipe data di Go yang digunakan untuk merepresentasikan kumpulan satu atau beberapa nilai. `Struct` dapat memiliki *method*.  Secara teknikal Go bukanlah sebuah *OOP language*. Namun `type` dan `method` memungkinkan kita membuat kode kita dengan *Object Oriented Style*. Struct dapat diibaratkan sebagai `class` dalam OOP language seperti Java.

Mari kita membuat struct pertama kita.

```go
// templateHandler merepresentasikan sebuah template
// sync.Once memastikan template hanya di compile dan diekseskusi sekali saja
type templateHandler struct {
  once sync.Once
  filename string
  templ *template.Template
}

// ServeHTTP menangani request dari http client
// (t *templateHandler) adalah method receiver,
func (t *templateHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
    t.once.Do(func() {
       t.templ = template.Must(template.ParseFiles(filepath.Join("templates",t.filename)))
	})
    t.templ.Execute(w, nil)
}
```

`templateHandler` memiliki method `ServeHTTP` dimana parameter yang di masukkan sama seperti method `http.HandleFunc` yang sebelumnya telah kita buat. Method `ServeHTTP` memenuhi *interface*`http.HandlerFunc`. Maka kita dapat menggunakan method `ServeHTTP` di dalam method `http.HandleFunc`.

Interface di dalam Go tidak didefinisikan secara explisit, namun di dalam Go, bila suatu struct mempunyai method yang sama dengan interface yang di deklarasikan, maka Go mengganggap tipe tersebut merupakan bagian dari tipe interface tersebut.

Di dalam file `main.go` tambahkan baris ini

```go
func main() {
     // root
     http.Handle("/", &templateHandler{filename: "chat.html"})
     // start the web server
     if err := http.ListenAndServe(":8080", nil); err != nil {
       log.Fatal("ListenAndServe:", err)
     }
}
```



## Pemodelan Chat Room dan Clients pada Server

Kita akan membuat satu public room besar untuk menampung semua client/user pada aplikasi chat kita.

`room` type akan bertanggung jawab untuk mengelola koneksi `client` dan rute pesan yang masuk dan keluar.

`client` type merepresentasikan koneksi tunggal



### Pemodelan Client

Buat file `client.go`

`go get github.com/gorilla/websocket`

```go
package main

import (
	"github.com/gorilla/websocket"
)

// client merepresentasikan single chatting user
type client struct {
  	// socket adalah web socket untuk client ini
  	socket *websocket.Conn
  	// send, pesan akan dikirim ke channel ini
  	send chan []byte
  	// room adalah room dimana client melakukan chatting
  	room *room
}
```



`socket` menyimpan *reference* atau penunjuk ke websocket. Tujuannya adalah agar server bisa berkomunikasi dengan client secara **dua arah**. Berikut diagram yang menjelaskan skema WebSocket



![alt text](https://docs.oracle.com/cd/E55956_01/doc.11123/user_guide/content/images/general/websocket_sequence.png "WebSocket Diagram")



`send` pesan akan dikirim ke dalam *channel* ini. Dan akan diteruskan melalui websocket dan diterima oleh client.

`room` menyimpan referensi ruangan dimana kita melakukan chatting, disini pula kita akan meneruskan pesan kita ke semua client



Supaya client dapat mengirim dan menerima pesan maka kita harus menambahkan method pada client

```go
func (c *client) read() {
  	for {
      // baris ini akan membaca message yang dikirim client melalui socket dan
      // meneruskan message tersebut ke forward channel
      if _, msg, err := c.socket.ReadMessage(); if err == nil {
        c.room.forward <- msg
      } else {
        break  
      }
  	}

    c.socket.Close()
}


func (c *client) write() {
   // method write secara terus menerus akan menerima msg dari channel send dan menuliskan
   // kembali ke socket
   for msg := range c.send {
       if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
          break
       }
   }
   c.socket.Close()
}
```





### Pemodelan Room

Buat file `room.go`

```go
package main

type room struct {
  	// forward adalah channel yang menampung pesan, sebelum di kirim ke seluruh client
  	forward chan[]byte
    // join adalah channel untuk client yang masuk ke dalam room
    join chan *client
    // leave adalah channel untuk client yang akan meninggalkan room
    leave chan *client
    // clients menyimpan semua client di dalam room ini
    clients map[*client] bool
}
```



Kenapa kita membutuhkan `join` dan `leave` channel? Karena ada kemungkinan di saat yang bersamaan go routine mencoba untuk mengakses `map[*client] bool`dan mengakibatkan *corrupt memory*



###  Concurrency Programming using Idiomatic Go

>Large programs are often made up of many smaller sub-programs. For example a web server handles requests made from web browsers and serves up HTML web pages in response. Each request is handled like a small program.
>
>**It would be ideal for programs like these to be able to run their smaller components at the same time** (in the case of the web server to handle multiple requests). **Making progress on more than one task simultaneously is known as concurrency.** Go has rich support for concurrency using goroutines and channels. https://www.golang-book.com/books/intro/10



![alt text](https://talks.golang.org/2012/waza/gophercomplex1.jpg "Gopher")

### Go Routines

> A goroutine is a **function that is capable of running concurrently** with other functions.
>
> example. https://golang.org/src/net/http/server.go 2272-2293

Kita akan menggunakan fitur concurrency dalam bahasa Go yaitu `select`



```go
func (r *room) run() {
     for {
       select {
       case client := <-r.join:
          // joining
         r.clients[client] = true
       case client := <-r.leave:
         // leaving
         delete(r.clients, client)
         close(client.send)
       case msg := <-r.forward:
         // forward message to all clients
         for client := range r.clients {
           select {
           case client.send <- msg:
             // send the message
           default:
             // failed to send
             delete(r.clients, client)
             close(client.send)
￼	         }    
         }
       }      
    }
}
```

Salah satu alasan kita menggunakan channel join dan leave adalah untuk memastikan r.clients map hanya bisa dimodifikasi oleh satu hal per satuan waktu .

### Channels

> Channels provide a **way for two goroutines to communicate with one another** and synchronize their execution. https://www.golang-book.com/books/intro/10

​								

### Mengubah Room menjadi HTTPHandler

Tambahkan baris berikut di `room.go`

```go
const (
  socketBufferSize  = 1024
  messageBufferSize = 256
)


var upgrader = &websocket.Upgrader{
    ReadBufferSize:socketBufferSize,
    WriteBufferSize: socketBufferSize
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {

  socket, err := upgrader.Upgrade(w, req, nil)

  if err != nil {
    log.Fatal("ServeHTTP:", err)
    return
  }

  client := &client{
    socket: socket,
    send:   make(chan []byte, messageBufferSize),
    room:   r,
  }

  r.join <- client

  defer func() { r.leave <- client }()

// panggil method write di go routine yang lain
  go client.write()

// read akan memblok operasi dari fungsi ini (keeping the connection alive), sampai koneksi ini diputus
  client.read()
}


// newRoom membuat instansiasi room baru
func newRoom() *room {
  return &room{
    forward: make(chan []byte),
    join:    make(chan *client),
    leave:   make(chan *client),
    clients: make(map[*client]bool),
  }
}

```

​		
Dengan menambahkan `ServeHTTP` method, artinya type room dapat berperan sebagai handler.

Untuk menggunakan websocket, kita harus mengupgrade koneksi HTTP kita. Hal ini dicapai dengan menggunakan method `upgrader.Upgrade`



### Mengupdate main.go dan memanggil room

Tambahkan

```
import (
	"flag"
	"log"
)
```

```go
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
     t.once.Do(func() {
       t.templ = template.Must(template.ParseFiles(filepath.Join("templates",t.filename)))
	 })
     t.templ.Execute(w, r) // tambahkan line ini, kita memasukkan object http.Request yang nantinya akan dipanggil di template kita r.Host
}
```

```go
func main() {
  	 var addr = flag.String("addr", ":8080", "The addr of the application.")
     flag.Parse() // parse the flags

     r := newRoom()
     http.Handle("/", &templateHandler{filename: "chat.html"})
     http.Handle("/room", r)
     // get the room going
     go r.run()
     // start the web server
     log.Println("Starting web server on", *addr)
     if err := http.ListenAndServe(*addr, nil); err != nil {
       log.Fatal("ListenAndServe:", err)
     }
}

```





### Sisi FrontEnd

```html
<html>
  <head>
    <title>Chat</title>
    <style>
      input { display: block; }
      ul    { list-style: none; }
    </style>
  </head>
  <body>

    <ul id="messages"></ul>
    <form id="chatbox">
      <textarea></textarea>
      <input type="submit" value="Send" />
    </form>

    <script src="//ajax.googleapis.com/ajax/libs/jquery/1.11.1/jquery.min.js"></script>
    <script>
      $(function(){
        var socket = null;
        var msgBox = $("#chatbox textarea");
        var messages = $("#messages");
        $("#chatbox").submit(function(){
          if (!msgBox.val()) return false;
          if (!socket) {
            alert("Error: There is no socket connection.");
            return false;
          }
          socket.send(msgBox.val());
          msgBox.val("");
          return false;
        });
        if (!window["WebSocket"]) {
          alert("Error: Your browser does not support web sockets.")
        } else {
          socket = new WebSocket("ws://{{.Host}}/room");
          socket.onclose = function() {
            alert("Connection has been closed.");
          }
          socket.onmessage = function(e) {
            messages.append($("<li>").text(e.data));
          }
        }
      });
    </script>
  </body>
</html>
```



## Jalankan Program

```
go build -o chat
./chat -addr=":3000"
```

