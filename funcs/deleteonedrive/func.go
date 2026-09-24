package deleteonedriv

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/eiannone/keyboard"
)

func RemoveOneDrive() {
	fmt.Println("\n[BTL] Запуск процесса удаления Microsoft OneDrive...")

	cmdKill := exec.Command("taskkill", "/f", "/im", "OneDrive.exe")
	_ = cmdKill.Run()

	sys32Path := filepath.Join(os.Getenv("SystemRoot"), "System32", "OneDriveSetup.exe")
	sysWOW64Path := filepath.Join(os.Getenv("SystemRoot"), "SysWOW64", "OneDriveSetup.exe")
	localAppDataPath := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "OneDrive", "Update", "OneDriveSetup.exe")

	var setupPath string
	if _, err := os.Stat(sysWOW64Path); err == nil {
		setupPath = sysWOW64Path
	} else if _, err := os.Stat(sys32Path); err == nil {
		setupPath = sys32Path
	} else if _, err := os.Stat(localAppDataPath); err == nil {
		setupPath = localAppDataPath
	}

	fmt.Println("\n==================================================")
	if setupPath != "" {
		fmt.Println("[BTL] Выполнение деинсталляции...")
		cmdUninstall := exec.Command(setupPath, "/uninstall")
		err := cmdUninstall.Run()

		if err != nil {
			fmt.Printf("\033[31m[BTL] Ошибка при вызове деинсталлятора: %v\033[0m\n", err)
		} else {
			fmt.Println("\033[32m[BTL] Microsoft OneDrive успешно удален!\033[0m")
		}
	} else {
		fmt.Println("\033[33m[BTL] Исполняемый файл OneDriveSetup.exe не найден. Возможно, OneDrive уже удален.\033[0m")
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
