package recoverysystem

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
)

func RunSFCFix() {
	fmt.Println("\n[BTL] Запуск проверки и восстановления системных файлов (SFC)...")
	fmt.Println("Пожалуйста, подождите, это может занять несколько минут.\n")

	cmd := exec.Command("sfc", "/scannow")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	fmt.Println("\n==================================================")
	if err != nil {
		fmt.Printf("\033[31m[BTL] Ошибка при выполнении проверки SFC!\033[0m\n")
		fmt.Printf("Подробности: %v\n", err)
		fmt.Println("Запустите программу от имени Администратора.")
	} else {
		fmt.Println("\033[32m[BTL] ПРОВЕРКА SFC ЗАВЕРШЕНА!\033[0m")
		fmt.Println("==================================================")
		fmt.Println("• Целостность системных файлов Windows проверена.")
		fmt.Println("• Поврежденные файлы были автоматически восстановлены.")
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
