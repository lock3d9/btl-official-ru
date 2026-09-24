package offblurscreen

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
	"golang.org/x/sys/windows/registry"
)

func DisableLockScreenBlur() {
	fmt.Println("\n[BTL] Отключение размытия на экране блокировки...")

	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows\System`, registry.SET_VALUE)
	if err != nil {
		printError(err)
		waitEnter()
		return
	}
	defer k.Close()

	if err := k.SetDWordValue("DisableBlurring", 1); err != nil {
		printError(err)
		waitEnter()
		return
	}

	fmt.Println("\n==================================================")
	fmt.Println("\033[32m[BTL] РАЗМЫТИЕ НА ЭКРАНЕ БЛОКИРОВКИ ОТКЛЮЧЕНО!\033[0m")
	fmt.Println("==================================================")
	fmt.Println("• В реестре установлен параметр DisableBlurring = 1.")
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