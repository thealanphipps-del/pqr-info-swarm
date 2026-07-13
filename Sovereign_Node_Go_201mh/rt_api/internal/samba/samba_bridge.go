package samba

import (
	"github.com/hirochachacha/go-smb2"
	"net"
	"io/ioutil"
	"fmt"
)

type SambaConfig struct {
	Addr     string
	User     string
	Pass     string
	Share    string
}

func ExecuteSambaIO(cfg SambaConfig, path string, data []byte, write bool) (string, error) {
	conn, err := net.Dial("tcp", cfg.Addr)
	if err != nil { return "", err }
	defer conn.Close()

	d := &smb2.Dialer{
		Authenticator: &smb2.NTLMv2Authenticator{
			User:     cfg.User,
			Password: cfg.Pass,
		},
	}

	s, err := d.Dial(conn)
	if err != nil { return "", err }
	defer s.Logout()

	fs, err := s.Mount(cfg.Share)
	if err != nil { return "", err }
	defer fs.Umount()

	if write {
		err = fs.WriteFile(path, data, 0644)
		if err != nil { return "", err }
		return "WRITE_SUCCESS", nil
	} else {
		f, err := fs.Open(path)
		if err != nil { return "", err }
		defer f.Close()
		content, err := ioutil.ReadAll(f)
		if err != nil { return "", err }
		return string(content), nil
	}
}
