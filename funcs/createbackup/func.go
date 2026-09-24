package createbackup

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
)

func CreateRestorePoint() {
	fmt.Println("\n[BTL] Создание точки восстановления системы...")

	now := time.Now().Format("02.01.2006 15:04:05")
	pointName := fmt.Sprintf("btlbackup %s", now)

	psScript := fmt.Sprintf(`
		Enable-ComputerRestore -Drive "C:\" -ErrorAction SilentlyContinue
		Checkpoint-Computer -Description "%s" -RestorePointType "MODIFY_SETTINGS"
	`, pointName)

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	outputBytes, err := cmd.CombinedOutput()

	fmt.Println("\n==================================================")
	if err != nil {
		fmt.Printf("\033[31m[BTL] Ошибка создания точки восстановления!\033[0m\n")
		fmt.Printf("Подробности: %v\n", string(outputBytes))
		fmt.Println("\nЗапустите программу от имени Администратора.")
	} else {
		fmt.Println("\033[32m[BTL] Точка восстановления успешно создана!\033[0m")
		fmt.Printf("Имя: %s\n", pointName)
	}
	fmt.Println("==================================================")

	fmt.Print("\nНажмите Enter для возврата в меню...")

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
