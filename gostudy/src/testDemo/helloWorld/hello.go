package main

import "fmt"

//TDD (测试驱动开发)是一门需要通过开发去实践的技能，通过将问题分解成更小的可测试的组件，你编写软件将会更加轻松。
const spanish = "Spanish"
const french = "French"
const helloPrefix = "Hello, "
const spanishHelloPrefix = "Hola, "
const frenchHelloPrefix = "*&, "

//const englishHelloPrefix = "Hello,"

// func Hello(name string) string {
// 	if name == "" {
// 		name = "World"
// 	}
// 	return englishHelloPrefix + name
// }

// func Hello(name string, language string) string {
// 	if name == "" {
// 		name = "World"
// 	}
// 	if language == "Spanish" {
// 		return "Hola," + name
// 	}
// 	return englishHelloPrefix + name
// }

// func Hello(name string, language string) string {
// 	if name == "" {
// 		name = "World"
// 	}

// 	if language == spanish {
// 		return spanishHelloPrefix + name
// 	}

// 	return helloPrefix + name
// }

// func Hello(name string, language string) string {
// 	if name == "" {
// 		name = "World"
// 	}

// 	prefix := helloPrefix

// 	switch language {
// 	case french:
// 		prefix = frenchHelloPrefix
// 	case spanish:
// 		prefix = spanishHelloPrefix
// 	}

// 	return prefix + name
// }

func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	return greetingPrefix(language) + name
}

func greetingPrefix(language string) (prefix string) {
	switch language {
	case french:
		prefix = frenchHelloPrefix
	case spanish:
		prefix = spanishHelloPrefix
	default:
		prefix = helloPrefix
	}
	return
}

func main() {
	fmt.Println(Hello("world", ""))
}
