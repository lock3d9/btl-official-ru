package autostart

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/eiannone/keyboard"
)

type AutoStartItem struct {
	Type   string `json:"Type"`
	Status string `json:"Status"`
	Name   string `json:"Name"`
	Path   string `json:"Path"`
}

func ManageAutostart() {
	fmt.Print("\033[H\033[2J")
	fmt.Println("\033[36m[BTL] Анализ автозагрузки системы...\033[0m")

	psScript := `
	[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

	$keep = 'security|defender|antivirus|avast|kaspersky|eset|malwarebytes|realtek|audio|nvidia|amd|intel|bluetooth|wi-fi|network|ctfmon|windows|microsoft|touchpad|synaptics|vgtray'
	$disable = 'update|assistant|helper|download|manager|adobe|browser|chrome|firefox|opera|yandex|discord|spotify|epic|skype|teams|slack|cortana|ccleaner|torrent|utorrent|bittorrent|steelseries|gg|epicgames|onedrive'

	$results = @()

	$regPaths = @(
		"HKCU:\Software\Microsoft\Windows\CurrentVersion\Run",
		"HKLM:\Software\Microsoft\Windows\CurrentVersion\Run"
	)

	foreach ($path in $regPaths) {
		if (Test-Path $path) {
			Get-ItemProperty -Path $path | ForEach-Object {
				$_.PSObject.Properties | Where-Object { $_.Name -notmatch '^PS' -and $_.Value } | ForEach-Object {
					$name = $_.Name
					$val = [string]$_.Value
					$status = "ПО ЖЕЛАНИЮ"

					if ($name -match $disable -or $val -match $disable) {
						$status = "СОВЕТУЕМ ОТКЛЮЧИТЬ"
					} elseif ($name -match $keep -or $val -match $keep) {
						$status = "ЛУЧШЕ ОСТАВИТЬ"
					}

					$results += [PSCustomObject]@{
						Type   = "Реестр"
						Status = $status
						Name   = $name
						Path   = $val
					}
				}
			}
		}
	}

	$startupFolders = @(
		"$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Startup",
		"$env:ProgramData\Microsoft\Windows\Start Menu\Programs\Startup"
	)

	foreach ($folder in $startupFolders) {
		if (Test-Path $folder) {
			Get-ChildItem -Path $folder -File | ForEach-Object {
				$results += [PSCustomObject]@{
					Type   = "Папка"
					Status = "СОВЕТУЕМ ОТКЛЮЧИТЬ"
					Name   = $_.Name
					Path   = $_.FullName
				}
			}
		}
	}

	ConvertTo-Json -InputObject @($results) -Compress
	`

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	outputBytes, err := cmd.Output()
	if err != nil {
		fmt.Printf("\n\033[31m[BTL] Ошибка при чтении автозагрузки: %v\033[0m\n", err)
		waitForKey()
		return
	}

	var items []AutoStartItem
	if err := json.Unmarshal(outputBytes, &items); err != nil {
		fmt.Printf("\n\033[31m[BTL] Ошибка парсинга данных: %v\033[0m\n", err)
		waitForKey()
		return
	}

	renderTable(items)
	waitForKey()
}

func renderTable(items []AutoStartItem) {
	fmt.Print("\033[H\033[2J")
	fmt.Println("\033[36m=====================================================================================\033[0m")
	fmt.Println("               \033[1;36mКОНТРОЛЬ АВТОЗАГРУЗКИ\033[0m")
	fmt.Println("\033[36m=====================================================================================\033[0m\n")

	if len(items) == 0 {
		fmt.Println("  \033[32m Активные элементы автозагрузки не обнаружены.\033[0m\n")
	} else {
		fmt.Printf(" %-15s %-22s %-10s %s\n", "РЕКОМЕНДАЦИЯ", "НАЗВАНИЕ", "ТИП", "ПУТЬ / КОМАНДА")
		fmt.Println("\033[90m-------------------------------------------------------------------------------------\033[0m")

		for _, item := range items {
			var badge string
			switch item.Status {
			case "ЛУЧШЕ ОСТАВИТЬ":
				badge = "\033[32m[ ОСТАВИТЬ ] \033[0m"
			case "СОВЕТУЕМ ОТКЛЮЧИТЬ":
				badge = "\033[31m[ ОТКЛЮЧИТЬ ]\033[0m"
			default:
				badge = "\033[33m[ ПО ЖЕЛАНИЮ]\033[0m"
			}

			name := truncate(item.Name, 20)
			path := truncate(item.Path, 45)

			fmt.Printf(" %s %-22s %-10s \033[90m%s\033[0m\n", badge, name, item.Type, path)
		}
	}

	fmt.Println("\033[90m-------------------------------------------------------------------------------------\033[0m")
	fmt.Println("\n\033[32m[BTL] Сканирование завершено успешно.\033[0m")
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen-3]) + "..."
	}
	return s
}

func waitForKey() {
	fmt.Print("\n\033[90mНажмите Enter или Esc для возврата в меню...\033[0m")
	if err := keyboard.Open(); err == nil {
		defer keyboard.Close()
		for {
			_, key, err := keyboard.GetKey()
			if err != nil || key == keyboard.KeyEnter || key == keyboard.KeyEsc {
				break
			}
		}
	}
}
