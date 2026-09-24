package checkheavyfiles

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

func FindHeavyFiles() {
	fmt.Println("[BTL] Внимание! Возможно займет некоторое время... (Если надо закрыть программу - Комбинация: Ctrl + C)")
	fmt.Println("\n[BTL] Поиск самых тяжелых файлов и папок...")

	psCommand := `
	$drives = Get-CimInstance Win32_LogicalDisk -Filter "DriveType=3"
	foreach ($d in $drives) {
		$root = $d.DeviceID + "\"
		Write-Host "Сканирование диска $root..." -ForegroundColor Cyan
		Get-ChildItem -Path $root -Recurse -File -ErrorAction SilentlyContinue | 
			Sort-Object Length -Descending | 
			Select-Object -First 5 Name, Length, FullName |
			ForEach-Object {
				$sizeGB = [Math]::Round($_.Length / 1GB, 2)
				if ($sizeGB -ge 0.1) {
					"$($sizeGB) GB | $($_.FullName)"
				}
			}
	}
	`

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCommand)
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("\n[BTL] Ошибка при сканировании файлов: %v\n", err)
		return
	}

	decoder := charmap.CodePage866.NewDecoder()
	outputStr, err := decoder.String(string(outputBytes))
	if err != nil {
		outputStr = string(outputBytes)
	}

	fmt.Println("\n==================================================")
	fmt.Println("\033[32m[BTL] САМЫЕ ТЯЖЕЛЫЕ ФАЙЛЫ НА ДИСКАХ\033[0m")
	fmt.Println("==================================================")

	trimmedOutput := strings.TrimSpace(outputStr)
	if trimmedOutput == "" {
		fmt.Println(" Тяжелых файлов не обнаружено.")
	} else {
		lines := strings.Split(outputStr, "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				fmt.Println(" " + trimmed)
			}
		}
	}

	fmt.Println("==================================================")
	fmt.Println("[BTL] Действие выполнено успешно!")
	fmt.Print("\nНажмите Enter для возврата в меню...")

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
