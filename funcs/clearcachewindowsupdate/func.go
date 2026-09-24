package clearcachewindowsupdate

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
)

func ClearWindowsUpdateCache() {
	fmt.Println("\n[BTL] Очистка кэша Windows Update...")

	psScript := `
		Stop-Service -Name wuauserv -Force -ErrorAction SilentlyContinue
		Stop-Service -Name bits -Force -ErrorAction SilentlyContinue
		Remove-Item -Path "$env:SystemRoot\SoftwareDistribution\Download\*" -Recurse -Force -ErrorAction SilentlyContinue
		Start-Service -Name wuauserv -ErrorAction SilentlyContinue
		Start-Service -Name bits -ErrorAction SilentlyContinue
	`

	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psScript)
	outputBytes, err := cmd.CombinedOutput()

	if err != nil {
		printError(fmt.Errorf("%v: %s", err, string(outputBytes)))
		waitEnter()
		return
	}

	fmt.Println("\n==================================================")
	fmt.Println("\033[32m[BTL] Кэш Windows Update очищен!\033[0m")
	fmt.Println("==================================================")

	waitEnter()
}

func printError(err error) {
	fmt.Println("\n==================================================")
	fmt.Printf("\033[31m[BTL] Ошибка выполнения операции!\033[0m\n")
	fmt.Printf("Подробности: %v\n", err)
	fmt.Println("==================================================")
}

func waitEnter() {
	fmt.Print("\nНажмите Enter для возврата в меню...")
	time.Sleep(200 * time.Millisecond)
	if err := keyboard.Open(); err == nil {
		defer keyboard.Close()
		for {
			_, key, err := keyboard.GetKey()
			if err != nil || key == keyboard.KeyEnter {
				break
			}
		}
	}
}
