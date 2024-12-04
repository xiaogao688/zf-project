package main

import (
	"errors"
	"fmt"
	"os"
)

// 模拟打开文件的错误
func openFile() error {
	return errors.New("failed to open file")
}

func main() {
	if err := C(); err != nil {
		fmt.Println(err)
		fmt.Println(errors.Unwrap(err))
		fmt.Println(errors.Unwrap(errors.Unwrap(err)))
	}
	_ = errJoin()
	errIS()
	errAs()
}

func A() error {
	return fmt.Errorf("this is A")
}

func B() error {
	if err := A(); err != nil {
		return fmt.Errorf("this is B: %w", err)
	}

	return nil
}

func C() error {
	if err := B(); err != nil {
		return fmt.Errorf("this is C: %w", err)
	}

	return nil
}

func errJoin() error {
	err1 := errors.New("err 1")
	err2 := errors.New("err 2")
	err3 := errors.New("err 3")
	err4 := errors.Join(err1, err2, err3)
	fmt.Println(err4)
	return err4
}

func errIS() {
	err := os.ErrNotExist
	//err = os.ErrClosed
	err = fmt.Errorf("add err: %w", err)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("File does not exist")
	} else {
		fmt.Println("file is close")
	}

	err1 := errors.New("err 1")
	err2 := errors.New("err 2")
	err3 := errors.New("err 3")
	err4 := errors.Join(err1, err2, err3)
	if errors.Is(err4, err1) {
		fmt.Println("exist err 1")
	} else {
		fmt.Println("not exist err 1")
	}
}

type e struct {
	A string
	B int64
	C bool
}

func (e *e) Error() string {
	return fmt.Sprintf("e: A=%s, B=%d, C=%v", e.A, e.B, e.C)
}

func errAs() {
	err := fmt.Errorf("wrapping: %w", &e{
		A: "aaa",
		B: 1234,
		C: true,
	})

	//var pathErr *e
	pathErr := &e{}
	if errors.As(err, &pathErr) {
		fmt.Printf("Failed to %s %d %v\n", pathErr.A, pathErr.B, pathErr.C)
	}
}
