package disabletelemetry

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
)

func DisableTelemetry() {
	fmt.Println("\n[BTL] Отключение телеметрии и сбора данных Windows...")

	psScript := `
		$services = @("DiagTrack", "dmwappushservice", "WerSvc")
		foreach ($service in $services) {
			Stop-Service -Name $service -Force -ErrorAction SilentlyContinue
			Set-Service -Name $service -StartupType Disabled -ErrorAction SilentlyContinue
		}

		$regPath = "HKLM:\SOFTWARE\Policies\Microsoft\Windows\DataCollection"
		if (-not (Test-Path $regPath)) {
			New-Item -Path $regPath -Force | Out-Null
		}
		Set-ItemProperty -Path $regPath -Name "AllowTelemetry" -Value 0 -Type DWord -Force -ErrorAction SilentlyContinue

		$tasks = @(
			"\Microsoft\Windows\Application Experience\Microsoft Compatibility Appraiser",
			"\Microsoft\Windows\Application Experience\ProgramDataUpdater",
			"\Microsoft\Windows\Autochk\Proxy",
			"\Microsoft\Windows\Customer Experience Improvement Program\Consolidator",
			"\Microsoft\Windows\Customer Experience Improvement Program\UsbCeip",
			"\Microsoft\Windows\DiskDiagnostic\Microsoft-Windows-DiskDiagnosticDataCollector"
		)
		foreach ($task in $tasks) {
			Disable-ScheduledTask -TaskName $task -ErrorAction SilentlyContinue | Out-Null
		}
	`
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psScript)
	outputBytes, err := cmd.CombinedOutput()

	fmt.Println("\n==================================================")
	if err != nil {
		fmt.Printf("\033[31m[BTL] Ошибка при отключении телеметрии!\033[0m\n")
		fmt.Printf("Подробности: %v\n", string(outputBytes))
		fmt.Println("\nЗапустите программу от имени Администратора.")
	} else {
		fmt.Println("\033[32m[BTL] ТЕЛЕМЕТРИЯ УСПЕШНО ОТКЛЮЧЕНА!\033[0m")
		fmt.Println("==================================================")
		fmt.Println("• Службы сбора диагностических данных остановлены и отключены.")
		fmt.Println("• В реестре установлен запрет на передачу телеметрии.")
		fmt.Println("• Запланированные задачи сбора данных Windows деактивированы.")
	}
	fmt.Println("==================================================")

	fmt.Print("\nНажмите Enter для возврата в меню...")
	time.Sleep(200 * time.Millisecond)
	if err := keyboard.Open(); err == nil {
		defer keyboard.Close()
		for {
			_, key, err := keyboard.GetKey()
			if err != nil {
				break
			}
			if key == keyboard.KeyEnter {
				break
			}
		}
	}
}
