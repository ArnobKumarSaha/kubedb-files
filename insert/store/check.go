package store

import (
	"bytes"
	"fmt"
	"gomodules.xyz/go-sh"
	"strings"
)

var (
	buf *bytes.Buffer
)

func main() {
	shSession := getCommand("newer-0", "demo", 3)
	output, err := shSession.Output()
	if err != nil {
		fmt.Printf("cant get output, err %s\n", err)
		return
	}

	errOutput := buf.String()
	if errOutput != "" {
		fmt.Printf("failed to execute command, stderr: %s", errOutput)
		return
	}

	outStr := string(output)
	fmt.Println(outStr + "\n\n")
	slice := strings.Split(outStr, "\n")
	for i, _ := range slice {
		if i == 0 {
			continue
		}
		fmt.Printf("%s , ", slice[i])
	}
}

func checkOwnership() {

}

func getCommand(pod, ns string, column int) *sh.Session {
	shell := sh.NewSession()
	shell.ShowCMD = false
	buf = &bytes.Buffer{}
	shell.Stderr = buf

	podName := fmt.Sprintf("pod/%s", pod)
	kubectlCommand := []interface{}{
		"exec", "-n", ns, podName, "-c", "mongodb", "--",
	}
	mgCommand := []interface{}{
		"ls", "-l", "/data/db",
	}
	finalCommand := append(kubectlCommand, mgCommand...)

	awkCommand := fmt.Sprintf("{print $%d}", column)
	return shell.Command("kubectl", finalCommand...).Command("awk", awkCommand)
}
