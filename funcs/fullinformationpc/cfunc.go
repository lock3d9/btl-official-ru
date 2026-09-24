package sysinfo

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ShowSystemInfo() {
	fmt.Println("\n[BTL] Сбор подробной информации о системе...")
	psCommand := `
	[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
	$OutputEncoding = [System.Text.Encoding]::UTF8

	$os = Get-CimInstance Win32_OperatingSystem
	$cpu = Get-CimInstance Win32_Processor
	$cs = Get-CimInstance Win32_ComputerSystem
	$baseboard = Get-CimInstance Win32_BaseBoard
	$bios = Get-CimInstance Win32_BIOS
	$gpus = Get-CimInstance Win32_VideoController
	$ramSticks = Get-CimInstance Win32_PhysicalMemory
	$disks = Get-CimInstance Win32_DiskDrive
	$logicalDisks = Get-CimInstance Win32_LogicalDisk -Filter "DriveType=3"
	$soundDevices = Get-CimInstance Win32_SoundDevice
	$netAdapters = Get-CimInstance Win32_NetworkAdapter | Where-Object { $_.NetEnabled -eq $true -or $_.PhysicalAdapter -eq $true }

	$installDate = "Неизвестно"
	if ($os.InstallDate) {
		try { $installDate = [Management.ManagementDateTimeConverter]::ToDateTime($os.InstallDate).ToString('yyyy-MM-dd HH:mm:ss') } catch {}
	}

	$biosDate = "Неизвестно"
	if ($bios.ReleaseDate) {
		try { $biosDate = [Management.ManagementDateTimeConverter]::ToDateTime($bios.ReleaseDate).ToString('yyyy-MM-dd') } catch {}
	}

	Write-Output "SECTION|ОПЕРАЦИОННАЯ СИСТЕМА"
	Write-Output "ОС: $($os.Caption) ($($os.OSArchitecture))"
	Write-Output "Версия: Сборка $($os.BuildNumber) (Версия $($os.Version))"
	Write-Output "Дата установки: $installDate"
	Write-Output "Имя ПК / Домен: $($cs.Name) / $($cs.Domain)"

	Write-Output "SECTION|ЦЕНТРАЛЬНЫЙ ПРОЦЕССОР (CPU)"
	Write-Output "Процессор: $($cpu.Name)"
	Write-Output "Ядра / Потоки: $($cpu.NumberOfCores) ядер / $($cpu.NumberOfLogicalProcessors) потоков"
	Write-Output "Макс. частота: $([Math]::Round($cpu.MaxClockSpeed / 1000, 2)) GHz"
	Write-Output "Кэш L2 / L3: $([Math]::Round($cpu.L2CacheSize / 1024, 1)) MB / $([Math]::Round($cpu.L3CacheSize / 1024, 1)) MB"

	Write-Output "SECTION|МАТЕРИНСКАЯ ПЛАТА И BIOS"
	Write-Output "Плата: $($baseboard.Manufacturer) $($baseboard.Product)"
	Write-Output "Версия платы: $($baseboard.Version)"
	Write-Output "BIOS: $($bios.Manufacturer) $($bios.SMBIOSBIOSVersion) ($biosDate)"

	Write-Output "SECTION|ОПЕРАТИВНАЯ ПАМЯТЬ (RAM)"
	Write-Output "Общий объем: $([Math]::Round($cs.TotalPhysicalMemory / 1GB, 2)) GB"
	foreach ($r in $ramSticks) {
		$cap = [Math]::Round($r.Capacity / 1GB, 2)
		Write-Output "Планка: ${cap} GB | $($r.Speed) MHz | $($r.Manufacturer.Trim()) | Слот: $($r.DeviceLocator)"
	}

	Write-Output "SECTION|ВИДЕОСИСТЕМА (GPU)"
	foreach ($g in $gpus) {
		$vram = [Math]::Round($g.AdapterRAM / 1GB, 2)
		if ($vram -le 0) { $vram = [Math]::Round($g.AdapterCompatibility / 1GB, 2) }
		Write-Output "Видеокарта: $($g.Name) [VRAM: ${vram} GB]"
	}

	Write-Output "SECTION|НАКОПИТЕЛИ"
	foreach ($d in $disks) {
		$size = [Math]::Round($d.Size / 1GB, 2)
		Write-Output "Диск (Физический): $($d.Model.Trim()) — ${size} GB"
	}
	foreach ($ld in $logicalDisks) {
		$free = [Math]::Round($ld.FreeSpace / 1GB, 2)
		$total = [Math]::Round($ld.Size / 1GB, 2)
		Write-Output "Раздел [$($ld.DeviceID)] Всего: ${total} GB | Свободно: ${free} GB"
	}

	Write-Output "SECTION|ЗВУК И СЕТЬ"
	foreach ($s in $soundDevices) {
		Write-Output "Звук: $($s.Name)"
	}
	foreach ($n in $netAdapters) {
		if ($n.NetEnabled) {
			Write-Output "Сеть: $($n.Name)"
		}
	}
	`

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCommand)
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("\n[BTL] Ошибка при сборе информации: %v\n", err)
		return
	}

	outputStr := string(outputBytes)

	fmt.Println("\n==================================================")
	fmt.Println("\033[32m[BTL] ИНФОРМАЦИЯ О СИСТЕМЕ\033[0m")
	fmt.Println("==================================================")

	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "SECTION|") {
			sectionName := strings.TrimPrefix(trimmed, "SECTION|")
			fmt.Printf("\n\033[36m=== %s ===\033[0m\n", sectionName)
		} else {
			fmt.Printf("  %s\n", trimmed)
		}
	}

	fmt.Println("\n==================================================")
	fmt.Println("[BTL] Сбор данных завершен успешно!")
	fmt.Print("\nНажмите Enter для возврата в меню...")

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
