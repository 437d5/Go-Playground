document.addEventListener('DOMContentLoaded', function () {
    document.getElementById('run-btn').addEventListener('click', runCode);
    document.getElementById('frmt-btn').addEventListener('click', frmtCode);
    document.getElementById('rst-btn').addEventListener('click', rstCode);

    async function runCode() {
        const code = document.getElementById('code-area').value;
        const outputArea = document.getElementById('output-area');

        outputArea.textContent = "Running...";

        try {
            const response = await fetch('http://127.0.0.1:5050/run', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ code: code })
            });

            const result = await response.json();

            if (result.error) {
                outputArea.textContent = result.error;
            } else {
                outputArea.textContent = result.output;
            }
        } catch (err) {
            outputArea.textContent = "Error: " + err.message;
        }
    }

    async function frmtCode() {
        const code = document.getElementById('code-area').value;
        const codeArea = document.getElementById('code-area');
        const outputArea = document.getElementById('output-area');

        try {
            const response = await fetch('http://127.0.0.1:5050/frmt', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ code: code })
            });

            const result = await response.json();

            if (result.error) {
                outputArea.textContent = result.error;
            } else {
                codeArea.value = result.output;
            }
        } catch (err) {
            outputArea.textContent = "Error: " + err.message;
        }
    }

    async function rstCode() {
        const codeArea = document.getElementById('code-area');

        codeArea.value = `package main

import "fmt"

func main() {
	fmt.Println("hello, world")
}`
    }
})