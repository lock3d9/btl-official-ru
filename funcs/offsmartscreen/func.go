package offsmartscreen

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
	"golang.org/x/sys/windows/registry"
)

func DisableSmartScreen() {
	fmt.Println("\n[BTL] Отключение SmartScreen...")

	kSystem, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows\System`, registry.SET_VALUE)
	if err != nil {
		printError(err)
		waitEnter()
		return
	}
	defer kSystem.Close()

	if err := kSystem.SetDWordValue("EnableSmartScreen", 0); err != nil {
		printError(err)
		waitEnter()
		return
	}

	kExplorer, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer`, registry.SET_VALUE)
	if err != nil {
		printError(err)
		waitEnter()
		return
	}
	defer kExplorer.Close()

	if err := kExplorer.SetStringValue("SmartScreenEnabled", "Off"); err != nil {
		printError(err)
		waitEnter()
		return
	}

	fmt.Println("\n==================================================")
	fmt.Println("\033[32m[BTL] SMARTSCREEN УСПЕШНО ОТКЛЮЧЕН!\033[0m")
	fmt.Println("==================================================")
	fmt.Println("• В реестре установлен параметр EnableSmartScreen = 0.")
	fmt.Println("• В реестре установлен параметр SmartScreenEnabled = Off.")
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
