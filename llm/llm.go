package llm

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

// handle starting the ollama process
// reads messages from rabbitmq queue
// writes to ollama process and then send response back a response queue
func Prompt(input string) (string, error) {
	cmd := exec.Command("ollama", "run", "tiny-llama")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("error creating stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("error creating stdout pipe %v", err)
	}
	//Start the process
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("error starting Ollama process: %v", err)
	}

	defer cmd.Process.Kill()

	//Send input to Ollama
	_, err = stdin.Write([]byte(input + "\n"))

	if err != nil {
		return "", fmt.Errorf("error writing to Ollama stdin: %v", err)
	}
	stdin.Close() //Close stdin to signal end of input

	//Read the response from Ollama
	scanner := bufio.NewScanner(stdout)
	var outputBuilder strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		outputBuilder.WriteString(line + "\n")
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return "", fmt.Errorf("error reading from Ollama stdout: %v", scanErr)
	}

	//Ensure the process exits cleanly
	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("error waiting for Ollama process to exit: %v", err)
	}

	//Trim and return the response
	return strings.TrimSpace(outputBuilder.String()), nil
}
