package recoverysystemdism

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
)

func RunDISMFix() {
	fmt.Println("\n[BTL] Запуск восстановления образа системы (DISM)...")
	fmt.Println("Пожалуйста, подождите, это может занять длительное время.\n")

	cmd := exec.Command("dism", "/Online", "/Cleanup-Image", "/RestoreHealth")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	fmt.Println("\n==================================================")
	if err != nil {
		fmt.Printf("\033[31m[BTL] Ошибка при выполнении восстановления DISM!\033[0m\n")
		fmt.Printf("Подробности: %v\n", err)
		fmt.Println("Запустите программу от имени Администратора.")
	} else {
		fmt.Println("\033[32m[BTL] ВОССТАНОВЛЕНИЕ ОБРАЗА DISM ЗАВЕРШЕНО!\033[0m")
		fmt.Println("==================================================")
		fmt.Println("• Хранилище компонентов Windows успешно восстановлено.")
		fmt.Println("• Системный образ приведен в исправное состояние.")
	}
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
