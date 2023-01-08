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
	}

	asmFile := "tmp.s"
	exeFile := "tmp"

	for _, v := range data {
		old := os.Stdout

		f, err := os.Create(asmFile)
		if err != nil {
			t.Fatal("Failed to create file")
		}
		os.Stdout = f

		compile(v.input)

		f.Close()

		os.Stdout = old

		args := []string{"-static", "-o", exeFile, asmFile}
		cmd := exec.Command("gcc", args...)
		if err := cmd.Run(); err != nil {
			t.Errorf("Failed to build program")
		}

		cmd = exec.Command("./tmp")
		cmd.Run()
		exitCode := cmd.ProcessState.ExitCode()
		t.Logf("%v => %v (expected: %v)\n", v.input, exitCode, v.expected)
		if exitCode != v.expected {
			t.Errorf("Failed to run program")
		}

		os.Remove(asmFile)
		os.Remove(exeFile)
	}
}
