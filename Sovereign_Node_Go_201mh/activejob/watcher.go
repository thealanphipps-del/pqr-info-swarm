package activejob
import (
    "fmt"
    "os"
    "path/filepath"
)

type VaultWatcher struct {
    VaultPath string
}

func (v *VaultWatcher) ScanForPDFs() []string {
    var files []string
    err := filepath.Walk(v.VaultPath, func(path string, info os.FileInfo, err error) error {
        if err == nil && !info.IsDir() && filepath.Ext(path) == ".pdf" {
            files = append(files, path)
        }
        return nil
    })
    if err != nil {
        fmt.Println("[!] VaultWatcher Error:", err)
    }
    return files
}
