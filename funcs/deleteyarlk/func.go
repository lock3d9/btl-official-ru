package deleteyarlk

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
	"golang.org/x/sys/windows/registry"
)

func RemoveShortcutSuffix() {
	fmt.Println("\n[BTL] Отключение надписи \"Ярлык\" для новых ярлыков...")

	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer`, registry.SET_VALUE)
	if err != nil {
		printError(err)
		waitEnter()
		return
	}
	defer k.Close()

	if err := k.SetBinaryValue("link", []byte{0x00, 0x00, 0x00, 0x00}); err != nil {
		printError(err)
		waitEnter()
		return
	}

	fmt.Println("\n==================================================")
	fmt.Println("\033[32m[BTL] ПРЕФИКС \"ЯРЛЫК\" УСПЕШНО ОТКЛЮЧЕН!\033[0m")
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
