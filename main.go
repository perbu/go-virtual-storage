package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/perbu/go-virtual-storage/gstorage"
	"github.com/perbu/go-virtual-storage/lstorage"
	"go.uber.org/zap"
	"io"
	"math/rand"
)

type VirtualStorage interface {
	List() ([]string, error)
	ListDir(dir string) ([]string, error)
	Open(name string) (io.ReadCloser, error)
	Create(name string) (io.WriteCloser, error)
	Remove(name string) error
	Rename(oldName, newName string) error
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	gs, err := gstorage.NewClient(context.TODO())
	if err != nil {
		sugar.Fatalf("failed to create gstorage client: %v", err)
	}
	err = exercise(gs, sugar, "google")
	if err != nil {
		sugar.Fatalf("exercise failed with google storage: %v", err)
	}
	ls := lstorage.New()
	err = exercise(ls, sugar, "local")
	if err != nil {
		sugar.Fatalf("exercise failed with local storage: %v", err)
	}

}

func exercise(v VirtualStorage, sugar *zap.SugaredLogger, annotation string) error {
	fname := randomString(10)
	fhw, err := v.Create(fname)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	sugar.Infow("created file", "name", fname, "annotation", annotation)
	blocks := make([][]byte, 0)
	dataWritten := make([]byte, 0)
	for i := 0; i < 10; i++ {
		block := randomBytes(4096)
		blocks = append(blocks, block)
		dataWritten = append(dataWritten, block...)
	}
	written := 0
	for i := 0; i < 10; i++ {
		n, err := fhw.Write(blocks[i])
		if err != nil {
			return fmt.Errorf("write (block: %d: %w", i, err)
		}
		written += n
		sugar.Infow("wrote block", "block", i, "written", written, "annotation", annotation)
	}
	err = fhw.Close()
	if err != nil {
		return fmt.Errorf("close: %w", err)
	}
	dataRead := make([]byte, 0)
	fhr, err := v.Open(fname)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	readTotal := 0
	for {
		b := make([]byte, 4096)
		n, err := fhr.Read(b)
		if err != nil {
			if err == io.EOF {
				sugar.Infow("read EOF", "file", fname, "readTotal", readTotal, "annotation", annotation)
				break
			}
			return fmt.Errorf("read: %w", err)
		}
		readTotal += n
		dataRead = append(dataRead, b[:n]...)
	}
	sugar.Infow("read", "file", fname, "readTotal", readTotal, "annotation", annotation)

	// compare dataRead and dataWritten
	if len(dataRead) != len(dataWritten) {
		return fmt.Errorf("dataRead and dataWritten have different lengths")
	}
	for i := range dataRead {
		if dataRead[i] != dataWritten[i] {
			return fmt.Errorf("dataRead and dataWritten differ at index %d", i)
		}
	}
	sugar.Infow("dataRead and dataWritten are equal", "annotation", annotation)
	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz-ABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789"

// Generate some random bytes and return them
func randomBytes(le int) []byte {
	b := make([]byte, le)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

func randomString(le int) string {
	return string(randomBytes(le))
}
