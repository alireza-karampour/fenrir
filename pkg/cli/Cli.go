package cli

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"codeberg.org/bit101/go-ansi"

	"github.com/alireza-karampour/fenrir/pkg/task"
	. "github.com/alireza-karampour/fenrir/pkg/utils"
)

type DownloadableOpts int

const (
	DL_VERBOSE DownloadableOpts = 0b0001
	DL_TAR     DownloadableOpts = 0b0010
	DL_GZIP    DownloadableOpts = 0b0100
)

type Tar struct {
	Name       string
	URL        string
	Dest       string
	TargetFile string
}

func (t *Tar) ExtractFile(to string, as string, opts DownloadableOpts) (err error) {
	Println(fmt.Sprintf("extracting file %s to %s", as, to))
	defer func() {
		if err != nil {
			PrintlnErr("extraction failed")
			return
		} else {
			PrintlnOk("extraction done")
		}
	}()
	file, err := os.OpenFile(path.Join(t.Dest, t.Name), os.O_RDONLY, 0o666)
	if err != nil {
		return
	}
	defer file.Close()

	if (opts & DL_GZIP) == DL_GZIP {
		gzr, err1 := gzip.NewReader(file)
		if err1 != nil {
			return err1
		}
		defer gzr.Close()
		tr := tar.NewReader(gzr)
		for {
			hdr, err1 := tr.Next()
			if err1 != nil {
				if err1 == io.EOF {
					break
				}
				return err1
			}
			if hdr.Name == t.TargetFile {
				targetFile, err1 := os.OpenFile(path.Join(to, as), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o777)
				if err1 != nil {
					return err1
				}
				defer targetFile.Close()
				_, err = io.Copy(targetFile, tr)
				if err != nil {
					return err
				}
			}
		}
	} else {
		tr := tar.NewReader(file)
		for {
			hdr, err1 := tr.Next()
			if err1 != nil {
				if err1 == io.EOF {
					break
				}
				return err1
			}
			if hdr.Name == t.TargetFile {
				targetFile, err1 := os.OpenFile(path.Join(to, as), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o777)
				if err1 != nil {
					return err1
				}
				defer targetFile.Close()
				_, err = io.Copy(targetFile, tr)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

type Downloadable struct {
	Name     string
	Url      string
	Dest     string
	Checksum string
	TarFile  *Tar
}

func (d *Downloadable) Download(opts DownloadableOpts) error {
	dlbl := path.Join(d.Dest, d.Name)
	if _, err := os.Stat(dlbl); errors.Is(err, os.ErrNotExist) {
		ansi.Tab()
		PrintlnWarn(fmt.Sprintf("%s not found", d.Name))
		if !d.shouldDownload() {
			return nil
		}
		if d.TarFile != nil && (opts&DL_TAR) == DL_TAR {
			err = d.download(d.TarFile.URL, d.TarFile.Dest, d.TarFile.Name)
			if err != nil {
				return err
			}
			err = d.TarFile.ExtractFile(d.Dest, d.Name, opts)
			if err != nil {
				return err
			}
		} else {
			err = d.download(d.Url, d.Dest, d.Name)
			if err != nil {
				return err
			}
		}
	} else {
		ok, err := IsValidSum(path.Join(d.Dest, d.Name), d.Checksum)
		if err != nil {
			return err
		}
		if !ok {
			PrintlnWarn(fmt.Sprintf("invalid %s checksum", d.Name))
			if !d.shouldDownload() {
				return nil
			}
			err = d.download(d.Url, d.Dest, d.Name)
			if err != nil {
				return err
			}
		} else {
			if (opts & DL_VERBOSE) == DL_VERBOSE {
				PrintlnOk(fmt.Sprintf("%s found", d.Name))
			}
		}
	}
	return nil
}

func (d *Downloadable) shouldDownload() bool {
	Print(fmt.Sprintf("Download %s? (y/n)", d.Name))
	fmt.Printf(": ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	return string(input.Bytes()) != "n"
}

func (d *Downloadable) download(url string, dest string, name string) error {
	dlbl := path.Join(dest, name)
	cl := http.DefaultClient
	res, err := cl.Get(url)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dest, 0o777)
	if err != nil {
		return err
	}

	dlFile, err := os.OpenFile(dlbl, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o777)
	if err != nil {
		return err
	}
	buf := make([]byte, 1024)
	tee := io.TeeReader(res.Body, dlFile)
	accu := 0
	for {
		n, err := tee.Read(buf)
		if n > 0 {
			accu += n
			ansi.ClearLine()
			Print(fmt.Sprintf("Downloaded: %d bytes", accu))
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				ansi.NewLine()
				break
			} else {
				return err
			}
		}
	}
	err = dlFile.Sync()
	if err != nil {
		return err
	}
	PrintlnOk(fmt.Sprintf("%s download finished", name))
	dlFile.Close()
	sum, err := Sum(path.Join(dest, name))
	if err != nil {
		return err
	}
	Println(fmt.Sprintf("%s sum: %x", name, sum))
	return nil
}

func (d *Downloadable) Run(args string, stdin io.Reader) (io.Reader, error) {
	d.Download(0)
	tsk := task.NewTask(fmt.Sprintf("%s %s", path.Join(d.Dest, d.Name), args))
	tsk.SetIn(stdin)
	err := tsk.Run()
	if err != nil {
		return tsk.Err, err
	}
	return tsk.Out, nil
}
