package offupdatewindows

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
	"golang.org/x/sys/windows/registry"
)

func ToggleWindowsUpdate() {
	fmt.Println("\n[BTL] Проверка состояния Windows Update...")

	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows\WindowsUpdate\AU`, registry.ALL_ACCESS)
	if err != nil {
		printError(err)
		waitEnter()
		return
	}
	defer k.Close()

	val, _, err := k.GetIntegerValue("NoAutoUpdate")
	isDisabled := (err == nil && val == 1)

	if isDisabled {
		fmt.Println("[BTL] Включение Windows Update...")

		if err := k.SetDWordValue("NoAutoUpdate", 0); err != nil {
			printError(err)
			waitEnter()
			return
		}

		psScript := `
			Set-Service -Name wuauserv -StartupType Automatic -ErrorAction SilentlyContinue
			Start-Service -Name wuauserv -ErrorAction SilentlyContinue
		`
		_ = exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psScript).Run()

		fmt.Println("\n==================================================")
		fmt.Println("\033[32m[BTL] Windows Update включен!\033[0m")
		fmt.Println("==================================================")

	} else {
		fmt.Println("[BTL] Отключение Windows Update...")

		if err := k.SetDWordValue("NoAutoUpdate", 1); err != nil {
			printError(err)
			waitEnter()
			return
		}

		psScript := `
			Stop-Service -Name wuauserv -Force -ErrorAction SilentlyContinue
			Set-Service -Name wuauserv -StartupType Disabled -ErrorAction SilentlyContinue
		`
		_ = exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psScript).Run()

		fmt.Println("\n==================================================")
		fmt.Println("\033[32m[BTL] Windows Update отключен!\033[0m")
		fmt.Println("==================================================")
	}

	waitEnter()
}

func GetWindowsUpdateMenuLabel() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows\WindowsUpdate\AU`, registry.QUERY_VALUE)
	if err != nil {
		return "Отключить Windows Update"
	}
	defer k.Close()

	val, _, err := k.GetIntegerValue("NoAutoUpdate")
	if err == nil && val == 1 {
		return "Включить Windows Update"
	}

	return "Отключить Windows Update"
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
