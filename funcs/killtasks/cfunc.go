package killtasks

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

func KillFrozenApps() {
	fmt.Println("\n[BTL] Поиск зависших приложений...")

	psCommand := `
	Get-Process | Where-Object {$_.MainWindowTitle -ne "" -and $_.Responding -eq $false} | Select-Object Name, Id, MainWindowTitle
	`

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCommand)
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("\n[BTL] Ошибка при сканировании процессов: %v\n", err)
		return
	}

	decoder := charmap.CodePage866.NewDecoder()
	outputStr, err := decoder.String(string(outputBytes))
	if err != nil {
		outputStr = string(outputBytes)
	}

	trimmedOutput := strings.TrimSpace(outputStr)

	fmt.Println("\n==================================================")
	fmt.Println("\033[32m[BTL] ЗАВИСШИЕ ПРИЛОЖЕНИЯ\033[0m")
	fmt.Println("==================================================")

	if trimmedOutput == "" {
		fmt.Println(" Зависших приложений не обнаружено.")
	} else {
		fmt.Println(outputStr)
		fmt.Println("--------------------------------------------------")
		fmt.Println("[BTL] Завершаю принудительно все зависшие процессы...")
		killCmd := exec.Command("taskkill", "/F", "/FI", "STATUS eq not responding")
		killOut, err := killCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("[BTL] Не удалось завершить некоторые процессы: %v\n", err)
		} else {
			killStr, _ := decoder.String(string(killOut))
			fmt.Println(killStr)
		}
	}

	fmt.Println("==================================================")
	fmt.Println("[BTL] Действие выполнено успешно!")
	fmt.Print("\nНажмите Enter для возврата в меню...")

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
