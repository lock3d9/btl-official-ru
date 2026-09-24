package disabledefender

import (
	"fmt"
	"os/exec"

	"github.com/eiannone/keyboard"
)

func DisableDefender() {
	fmt.Println("\n[BTL] Отключение Windows Defender...")

	psScript := fmt.Sprintf(`
		Set-MpPreference -DisableBehaviorMonitoring $true
		Set-MpPreference -DisableScriptScanning $true
		Set-MpPreference -DisableIOAVProtection $true
		Set-MpPreference -DisableIntrusionPreventionSystem $true
	`)

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	outputBytes, err := cmd.CombinedOutput()

	fmt.Println("\n==================================================")
	if err != nil {
		fmt.Printf("\033[31m[BTL] Ошибка отключение Windows Defender\033[0m\n")
		fmt.Printf("Подробности: %v\n", string(outputBytes))
		fmt.Println("\nПопробуйте запусть программу от имени Администратора.")
	} else {
		fmt.Println("\033[32m[BTL] Windows Defender был отключен\033[0m")
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
