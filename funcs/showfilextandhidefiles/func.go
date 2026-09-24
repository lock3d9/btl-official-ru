package showfilextandhidefiles

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
	"golang.org/x/sys/windows/registry"
)

func ShowHiddenFiles() {
	fmt.Println("\n[BTL] Включение отображения скрытых файлов и папок...")

	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, registry.SET_VALUE)
	if err != nil {
		printError(err)
		waitEnter()
		return
	}
	defer k.Close()

	if err := k.SetDWordValue("Hidden", 1); err != nil {
		printError(err)
		waitEnter()
		return
	}

	fmt.Println("\n==================================================")
	fmt.Println("\033[32m[BTL] СКРЫТЫЕ ФАЙЛЫ И ПАПКИ ВКЛЮЧЕНЫ!\033[0m")
	fmt.Println("==================================================")
	fmt.Println("• В реестре установлен параметр Hidden = 1.")
	fmt.Println("==================================================")

	waitEnter()
}

func ShowFileExtensions() {
	fmt.Println("\n[BTL] Включение отображения расширений файлов...")

	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, registry.SET_VALUE)
	if err != nil {
		printError(err)
		waitEnter()
		return
	}
	defer k.Close()

	if err := k.SetDWordValue("HideFileExt", 0); err != nil {
		printError(err)
		waitEnter()
		return
	}

	fmt.Println("\n==================================================")
	fmt.Println("\033[32m[BTL] РАСШИРЕНИЯ ФАЙЛОВ ВКЛЮЧЕНЫ!\033[0m")
	fmt.Println("==================================================")
	fmt.Println("• В реестре установлен параметр HideFileExt = 0.")
	fmt.Println("==================================================")

	waitEnter()
}

func printError(err error) {
	fmt.Println("\n==================================================")
	fmt.Printf("\033[31m[BTL] Ошибка выполнения операции!\033[0m\n")
	fmt.Printf("Подробности: %v\n", err)
	fmt.Println("Запустите программу от имени Администратора.")
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
