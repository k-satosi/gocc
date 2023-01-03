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
		{0, "return 0;"},
		{42, "return 42;"},
		{21, "return 5+20-4;"},
		{41, "return  12 + 34 - 5 ;"},
		{47, "return 5+6*7;"},
		{0, "return 0==1;"},
		{1, "return 0<1;"},
		{1, "return 1>0;"},
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

		token := Tokenize(v.input)
		parser := NewParser(token)
		node := parser.Program()

		Codegen(node)

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
		if exitCode != v.expected {
			t.Errorf("Failed to run program")
		}

		os.Remove(asmFile)
		os.Remove(exeFile)
	}
}
