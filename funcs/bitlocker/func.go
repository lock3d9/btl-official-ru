package bitlockerw

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
	"golang.org/x/sys/windows/registry"
)

func ToggleBitLockerAutoEncryption() {
	fmt.Println("\n[BTL] Проверка состояния авто-шифрования BitLocker...")

	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\BitLocker`, registry.ALL_ACCESS)
	if err != nil {
		printError(err)
		waitEnter()
		return
	}
	defer k.Close()

	val, _, err := k.GetIntegerValue("PreventDeviceEncryption")
	isDisabled := (err == nil && val == 1)

	if isDisabled {
		fmt.Println("[BTL] Включение авто-шифрования BitLocker...")

		if err := k.SetDWordValue("PreventDeviceEncryption", 0); err != nil {
			printError(err)
			waitEnter()
			return
		}

		fmt.Println("\n==================================================")
		fmt.Println("\033[32m[BTL] АВТО-ШИФРОВАНИЕ BITLOCKER ВКЛЮЧЕНО!\033[0m")
		fmt.Println("==================================================")
		fmt.Println("• В реестре установлен параметр PreventDeviceEncryption = 0.")
		fmt.Println("==================================================")

	} else {
		fmt.Println("[BTL] Отключение авто-шифрования BitLocker...")

		if err := k.SetDWordValue("PreventDeviceEncryption", 1); err != nil {
			printError(err)
			waitEnter()
			return
		}

		fmt.Println("\n==================================================")
		fmt.Println("\033[32m[BTL] АВТО-ШИФРОВАНИЕ BITLOCKER ОТКЛЮЧЕНО!\033[0m")
		fmt.Println("==================================================")
		fmt.Println("• В реестре установлен параметр PreventDeviceEncryption = 1.")
		fmt.Println("• Автоматическое шифрование дисков отключено.")
		fmt.Println("==================================================")
	}

	waitEnter()
}

func GetBitLockerMenuLabel() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\BitLocker`, registry.QUERY_VALUE)
	if err != nil {
		return "Выключить авто-шифрование BitLocker"
	}
	defer k.Close()

	val, _, err := k.GetIntegerValue("PreventDeviceEncryption")
	if err == nil && val == 1 {
		return "Включить авто-шифрование BitLocker"
	}

	return "Выключить авто-шифрование BitLocker"
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
