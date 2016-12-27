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
       }
     )
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
