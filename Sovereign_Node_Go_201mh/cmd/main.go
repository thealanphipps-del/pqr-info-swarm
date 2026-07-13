package main

import (
        "html/template"
        "net/http"
        "os"
        "path/filepath"
        "strings"
)

type FileInfo struct {
        Name string
        URL  string
}

func main() {
        storagePath := "/data/data/com.termux/files/home/Forensic_Hub_Sync/S24_Extraction_20260328/"
        http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(storagePath))))

        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                query := strings.ToLower(r.URL.Query().Get("q"))
                var files []FileInfo
                filepath.Walk(storagePath, func(path string, info os.FileInfo, err error) error {
                        if err == nil && !info.IsDir() {
                                ext := strings.ToLower(filepath.Ext(info.Name()))
                                if ext == ".jpg" || ext == ".png" || ext == ".jpeg" {
                                        if query == "" || strings.Contains(strings.ToLower(info.Name()), query) {
                                                files = append(files, FileInfo{Name: info.Name(), URL: "/images/" + info.Name()})
                                        }
                                }
                        }
                        return nil
                })

                tmpl := `
                <!DOCTYPE html>
                <html>
                <head>
                        <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=5.0, user-scalable=yes">
                        <title>S25 | Forensic Hub</title>
                        <style>
                                body{background:#000;color:#0f0;font-family:monospace;margin:0;padding-bottom:100px;}
                                #header{position:sticky;top:0;background:rgba(17,17,17,0.95);padding:15px;border-bottom:2px solid #0f0;z-index:100;backdrop-filter:blur(5px);}
                                input{background:#000;color:#0f0;border:1px solid #0f0;width:92%;padding:15px;font-size:18px;border-radius:5px;}
                                .grid{display:flex;flex-wrap:wrap;gap:8px;padding:8px;justify-content:center;}
                                .card{width:calc(50% - 10px);background:#111;border:1px solid #333;border-radius:8px;overflow:hidden;box-shadow:0 4px 10px rgba(0,0,0,0.5);}
                                img{width:100%;display:block;height:250px;object-fit:cover;transition:0.3s;}
                                .name{font-size:10px;padding:5px;display:block;color:#666;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;}
                        </style>
                </head>
                <body>
                        <div id="header">
                                <form><input type="text" name="q" placeholder="Filter Screenshots..." autofocus value="{{.Query}}"></form>
                        </div>
                        <div class="grid">
                                {{range .Files}}
                                <div class="card">
                                        <a href="{{.URL}}" target="_blank"><img src="{{.URL}}" loading="lazy"></a>
                                        <span class="name">{{.Name}}</span>
                                </div>
                                {{end}}
                        </div>
                </body>
                </html>`
                
                type PageData struct {
                    Files []FileInfo
                    Query string
                }
                
                t, _ := template.New("web").Parse(tmpl)
                t.Execute(w, PageData{Files: files, Query: query})
        })
        http.ListenAndServe(":8888", nil)
}
