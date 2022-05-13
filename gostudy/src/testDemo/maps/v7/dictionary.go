package main

import _ "errors"

type Dictionary map[string]string

// var (
// 	ErrNotFound   = errors.New("could not find the word you were looking for")
// 	ErrWordExists = errors.New("cannot add word because it already exists")
// )
//Dictionary 型错误
const (
	ErrNotFound         = DictionaryErr("could not find the word you were looking for")
	ErrWordExists       = DictionaryErr("cannot add word because it already exists")
	ErrWordDoesNotExist = DictionaryErr("cannot update word because it does not exist")
)

//我们将错误声明为常量，这需要我们创建自己的 DictionaryErr 类型来实现 error 接口
type DictionaryErr string

func (e DictionaryErr) Error() string {
	return string(e)
}

func (d Dictionary) Search(word string) (string, error) {
	definition, ok := d[word]
	if !ok {
		return "", ErrNotFound
	}

	return definition, nil
}

func (d Dictionary) Add(word, definition string) error {
	_, err := d.Search(word)
	switch err {
	case ErrNotFound:
		d[word] = definition
	case nil:
		return ErrWordExists
	default:
		return err
	}
	return nil
}

func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)

	switch err {
	case ErrNotFound:
		return ErrWordDoesNotExist
	case nil:
		d[word] = definition
	default:
		return err
	}

	return nil
}

//Go 的 map 有一个内置函数 delete。它需要两个参数。第一个是这个 map，第二个是要删除的键。
func (d Dictionary) Delete(word string) {
	delete(d, word)
}
