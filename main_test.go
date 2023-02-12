package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestCompile(t *testing.T) {
	type testData struct {
		expected int
		input    string
	}
	data := []testData{
		{0, "int main() { return 0; }"},
		{42, "int main() {return 42; }"},
		{21, "int main() { return 5+20-4; }"},
		{41, "int main() { return  12 + 34 - 5 ; }"},
		{47, "int main() { return 5+6*7; }"},
		{0, "int main() { return 0==1; }"},
		{1, "int main() { return 0<1; }"},
		{1, "int main() { return 1>0; }"},
		{3, "int main() { int foo=3; return foo; }"},
		{8, "int main() { int foo123=3; int bar=5; return foo123+bar; }"},

		{3, "int main() { if (0) return 2; return 3; }"},
		{3, "int main() { if (1-1) return 2; return 3; }"},
		{2, "int main() { if (1) return2; return3; }"},
		{2, "int main() { if (2-1) return2; return3; }"},

		{10, "int main () { int i=0; while(i<10) i=i+1; return i; }"},

		{55, "int main() { int i=0; int j=0; for (i=0; i<=10; i=i+1) j=i+j; return j; }"},
		{3, "int main() { for (;;) return 3; return 5; }"},

		{32, "int main() { return ret32(); } int ret32() { return 32; }"},
		{7, "int main() { return add2(3,4); } int add2(int x, int y) { return x+y; }"},
		{55, "int main() { return fib(9); } int fib(int x) { if (x<=1) return 1; return fib(x-1) + fib(x-2); }"},

		{8, "int main() { int x=3; int y=5; return foo(&x, y); } int foo(int *x, int y) { return *x + y; }"},

		{3, "int main() { int x[2]; int *y=&x; *y=3; return *x; }"},
		{1, "int main() { int x[2][3]; int *y=x; *(y+1)=1; return *(*x+1); }"},

		{1, "int main() { int x[2][3]; int *y=x; y[1]=1; return x[0][1]; }"},

		{8, "int main() { int x; return sizeof(x); }"},
		{8, "int main() { int x; return sizeof x; }"},
		{8, "int main() { int *x; return sizeof(x); }"},
		{32, "int main() { int x[4]; return sizeof(x); }"},

		{3, "int x; int main() { x=3; return x; }"},
		{8, "int x; int main() { return sizeof(x); }"},

		{1, "int main() { char x=1; return x; }"},
		{1, "int main() { char x; return sizeof(x); }"},
		{10, "int main() { char x[10]; return sizeof(x); }"},

		{98, "int main() { return \"abc\"[1]; }"},
		{9, "int main() { return \"\\t\"[0]; }"},
		{2, "int main() { int x=2; { int x=3; } return x; }"},

		{1, "int main() { struct {int a; int b;} x; x.a=1; x.b=2; return x.a; }"},
	}

	exeFile := "tmp"

	for _, v := range data {
		r, w, _ := os.Pipe()

		backup := os.Stdout
		os.Stdout = w

		args := []string{"-x", "assembler", "-", "-static", "-o", exeFile}
		cmd := exec.Command("gcc", args...)

		compile(v.input)
		w.Close()

		cmd.Stdin = r

		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to build program: %v", err)
		}
		r.Close()
		os.Stdout = backup

		cmd = exec.Command("./tmp")
		cmd.Run()
		exitCode := cmd.ProcessState.ExitCode()
		t.Logf("%v => %v (expected: %v)\n", v.input, exitCode, v.expected)
		if exitCode != v.expected {
			t.Errorf("Failed to run program")
		}

		os.Remove(exeFile)
	}
}
