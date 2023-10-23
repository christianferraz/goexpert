package main

const a = "Hello, World!"

var (
	b bool    = true
	c int     = 10
	d string  = "Wesley"
	e float64 = 1.2
)

func main() {
	// a = "X" // string
	println(a)
	println(b)
	println(c)
	println(x())
	println(d)
	println(e)
}

func x() string {
	return "X"
}
