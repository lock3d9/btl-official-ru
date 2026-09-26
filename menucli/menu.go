package interfface

import (
	"fmt"
	"os"
	"syscall"
	"unicode"
	"unsafe"

	"github.com/eiannone/keyboard"
	"golang.org/x/sys/windows"
)

const menuLogo = `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣆⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢠⡀⠀⠀⠀⠀⠀⠀⠀⢀⣾⣿⡄⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠈⣿⣦⣄⠀⠀⠀⠀⢠⣿⣿⣿⣿⡀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⢹⣿⣿⣷⣤⣀⣠⣿⣿⠃⠸⣿⣧⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠈⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀⣿⣿⠀⠀⠀⠀
⠀⠀⠀⠀⠀⣤⣿⣿⣿⣿⣿⣿⣿⣿⣷⡀⢸⣿⠀⠀⠀⠀
⠀⠀⠀⠀⠀⣿⣿⣿⣿⣛⣩⣽⣿⣿⣿⣷⣸⣿⠀⠀⠀⠀
⠀⠀⠀⢀⣴⣿⣿⣿⡾⠿⠛⠛⠛⠛⠿⣿⣿⡿
⢀⣠⣾⣿⠟⠋⠉⢀⣀⣤⣤⣶⣶⣶⣦⣾⣿⠇
⠈⠙⠻⢷⣶⣴⡾⠿⠛⠉⠉⠀⠀⠈⣩⣿⠏⠀
⠀⠀⠀⠀⠀⠀`

type MenuItem struct {
	Title, Description, Warning string
}

type Category struct {
	Title string
	Items []MenuItem
}

type RECT struct {
	Left, Top, Right, Bottom int32
}

var (
	warningShown   bool
	user32         = syscall.NewLazyDLL("user32.dll")
	kernel32       = syscall.NewLazyDLL("kernel32.dll")
	procGetConsole = kernel32.NewProc("GetConsoleWindow")
	procGetMetrics = user32.NewProc("GetSystemMetrics")
	procGetWinRect = user32.NewProc("GetWindowRect")
	procMoveWin    = user32.NewProc("MoveWindow")

	procFillConsoleOutputChar = kernel32.NewProc("FillConsoleOutputCharacterW")
	procFillConsoleOutputAttr = kernel32.NewProc("FillConsoleOutputAttribute")
	procSetConsoleCursorPos   = kernel32.NewProc("SetConsoleCursorPosition")

	restorePointItem = MenuItem{
		Title:       "[РЕКОМЕНДУЕМ] Создать точку восстановления",
		Description: "Создает снимки системных файлов Windows на случай сбоев.",
		Warning:     "",
	}
)

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)

func initVTMode() {
	stdout := windows.Handle(os.Stdout.Fd())
	var mode uint32
	if err := windows.GetConsoleMode(stdout, &mode); err == nil {
		_ = windows.SetConsoleMode(stdout, mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	}
}

func centerConsoleWindow() {
	hwnd, _, _ := procGetConsole.Call()
	if hwnd == 0 {
		return
	}

	screenWidth, _, _ := procGetMetrics.Call(uintptr(SM_CXSCREEN))
	screenHeight, _, _ := procGetMetrics.Call(uintptr(SM_CYSCREEN))

	var rWin RECT
	procGetWinRect.Call(hwnd, uintptr(unsafe.Pointer(&rWin)))

	w := rWin.Right - rWin.Left
	h := rWin.Bottom - rWin.Top

	x := (int32(screenWidth) - w) / 2
	y := (int32(screenHeight) - h) / 2

	procMoveWin.Call(hwnd, uintptr(x), uintptr(y), uintptr(w), uintptr(h), 1)
}

func clearScreen() {
	stdout := windows.Handle(os.Stdout.Fd())

	var csbi windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(stdout, &csbi); err == nil {
		var count uint32
		cellCount := uint32(csbi.Size.X) * uint32(csbi.Size.Y)
		home := uintptr(0)

		procFillConsoleOutputChar.Call(uintptr(stdout), uintptr(' '), uintptr(cellCount), home, uintptr(unsafe.Pointer(&count)))
		procFillConsoleOutputAttr.Call(uintptr(stdout), uintptr(csbi.Attributes), uintptr(cellCount), home, uintptr(unsafe.Pointer(&count)))
		procSetConsoleCursorPos.Call(uintptr(stdout), home)
	}

	fmt.Print("\033[3J\033[2J\033[H")
}

func nav(key keyboard.Key, char rune, cursor, max int) int {
	r := unicode.ToLower(char)
	switch {
	case key == keyboard.KeyArrowUp || key == keyboard.KeyArrowLeft || r == 'w' || r == 'a' || r == 'ц' || r == 'ф':
		return (cursor - 1 + max) % max
	case key == keyboard.KeyArrowDown || key == keyboard.KeyArrowRight || r == 's' || r == 'd' || r == 'ы' || r == 'в':
		return (cursor + 1) % max
	}
	return cursor
}

func printOptions(opts []string, cursor int) {
	for i, opt := range opts {
		if i == cursor {
			fmt.Printf("\033[36m➔ %s\033[0m\n", opt)
		} else {
			fmt.Printf("  %s\n", opt)
		}
	}
}

func showStartupWarning() bool {
	opts := []string{"[ Создать точку восстановления сейчас ]", "[ Пропустить и перейти в меню ]"}
	cursor := 0

	for {
		clearScreen()
		fmt.Println("\033[31m========================================================\033[0m")
		fmt.Println("\033[31m                  ВАЖНОЕ ПРЕДУПРЕЖДЕНИЕ                  \033[0m")
		fmt.Println("\033[31m========================================================\033[0m\n")
		fmt.Println("Перед использованием оптимизаций настоятельно рекомендуется")
		fmt.Println("создать точку восстановления Windows.")
		fmt.Println("Это позволит вернуть систему в исходное состояние при сбоях.\n")

		printOptions(opts, cursor)

		char, key, err := keyboard.GetKey()
		if err != nil || key == keyboard.KeyEsc {
			return false
		}
		if key == keyboard.KeyEnter {
			return cursor == 0
		}

		cursor = nav(key, char, cursor, len(opts))
	}
}

func ConfirmAction(item MenuItem) bool {
	opts := []string{"[ Да, продолжить ]", "[ Отмена / Назад ]"}
	cursor := 0

	for {
		clearScreen()
		fmt.Printf("\033[36m=== %s ===\033[0m\n\n", item.Title)
		fmt.Printf("Описание:\n%s\n\n", item.Description)
		if item.Warning != "" {
			fmt.Printf("\033[31m[ПРЕДУПРЕЖДЕНИЕ]: %s\033[0m\n\n", item.Warning)
		}
		fmt.Println("Вы действительно хотите выполнить это действие?\n")
		printOptions(opts, cursor)

		char, key, err := keyboard.GetKey()
		if err != nil || key == keyboard.KeyEsc {
			return false
		}
		if key == keyboard.KeyEnter {
			return cursor == 0
		}

		cursor = nav(key, char, cursor, len(opts))
	}
}

func showSubMenu(cat Category) (string, bool) {
	opts := make([]string, len(cat.Items)+1)
	for i, item := range cat.Items {
		opts[i] = item.Title
	}
	opts[len(cat.Items)] = "[ Назад в главное меню ]"

	cursor := 0
	for {
		clearScreen()
		fmt.Printf("\033[36m=== РАЗДЕЛ: %s ===\033[0m\n", cat.Title)
		fmt.Println("Выберите опцию (W/S — навигация, Enter — выбор, Esc — назад):\n")
		printOptions(opts, cursor)

		if cursor < len(cat.Items) {
			fmt.Println("\n\033[90m----------------------------------------\033[0m")
			fmt.Printf("\033[33mОписание:\033[0m %s\n", cat.Items[cursor].Description)
			if w := cat.Items[cursor].Warning; w != "" {
				fmt.Printf("\033[31m[ПРЕДУПРЕЖДЕНИЕ]: %s\033[0m\n", w)
			}
		}

		char, key, err := keyboard.GetKey()
		if err != nil || key == keyboard.KeyEsc {
			return "", false
		}
		if key == keyboard.KeyEnter {
			if cursor == len(cat.Items) {
				return "", false
			}
			if ConfirmAction(cat.Items[cursor]) {
				return cat.Items[cursor].Title, true
			}
		}

		cursor = nav(key, char, cursor, len(opts))
	}
}

func Showmainmenu(version, buildNum, buildDate string) string {
	initVTMode()
	centerConsoleWindow()

	if err := keyboard.Open(); err != nil {
		fmt.Printf("[BTL] Ошибка инициализации клавиатуры: %v\n", err)
		os.Exit(1)
	}
	defer keyboard.Close()

	if !warningShown {
		warningShown = true
		if showStartupWarning() {
			if ConfirmAction(restorePointItem) {
				return restorePointItem.Title
			}
		}
	}

	categories := getCategories()
	opts := make([]string, len(categories)+2)
	opts[0] = restorePointItem.Title

	for i, cat := range categories {
		opts[i+1] = "[" + cat.Title + "]"
	}
	opts[len(categories)+1] = "Выход"

	cursor := 0

	for {
		clearScreen()

		fmt.Println("\033[36m" + menuLogo + "\033[0m")
		fmt.Printf("\033[90mv%s (build %s | дата: %s)\033[0m\n\n", version, buildNum, buildDate)

		fmt.Println("Выберите раздел (W/S — навигация, Enter — выбор):")
		fmt.Println()

		printOptions(opts, cursor)

		char, key, err := keyboard.GetKey()
		if err != nil {
			break
		}

		if key == keyboard.KeyEnter {
			if cursor == 0 {
				if ConfirmAction(restorePointItem) {
					return restorePointItem.Title
				}
			} else if cursor == len(categories)+1 {
				return "Выход"
			} else {
				if res, ok := showSubMenu(categories[cursor-1]); ok {
					return res
				}
			}
		}

		cursor = nav(key, char, cursor, len(opts))
	}

	return "Выход"
}

func getCategories() []Category {
	return []Category{
		{
			Title: "Система и Восстановление",
			Items: []MenuItem{
				{"Очистка временных файлов", "Удаляет кэш системы, временные файлы обновлений и корзину.", ""},
				{"Отключение телеметрии", "Отключает службы фонового сбора данных и отправку отчетов в Microsoft.", ""},
				{"Восстановить систему через SFC", "Проверяет и восстанавливает целостность системных файлов Windows.", ""},
				{"Восстановить систему через DISM", "Восстанавливает образ системы Windows через встроенное хранилище компонентов.", ""},
				{"Переключить авто-шифрование BitLocker", "Включает или отключает авто-шифрование дисков BitLocker.", ""},
			},
		},
		{
			Title: "Windows Update",
			Items: []MenuItem{
				{"Очистить кэш Windows Update", "Удаляет загруженные файлы и кэш центра обновлений.", ""},
				{"Переключить Windows Update (Вкл/Выкл)", "Включает или отключает службу автоматического обновления Windows.", ""},
			},
		},
		{
			Title: "Оптимизация и Твики",
			Items: []MenuItem{
				{"Выключить встроенный Windows Defender", "Полностью отключает Защитник Windows из системы.", "Действие необратимо! Ваш ПК останется без встроенной антивирусной защиты."},
				{"Удалить OneDrive", "Удаляет клиент OneDrive и отвязывает синхронизацию системных папок.", "Убедитесь, что важные файлы сохранены на локальном диске."},
				{"Контроль автозагрузки", "Позволяет просмотреть и отключить программы, запускаемые вместе с Windows.", ""},
			},
		},
		{
			Title: "Тестирование и Анализ дисков",
			Items: []MenuItem{
				{"Проверить диск на чтение, запись и скорость", "Проводит тестирование выбранного накопителя путем записи и чтения временных блоков.", "Может занять продолжительное время. Не закрывайте программу во время теста."},
				{"Проверить флешку на муляж", "Проверяет реальную емкость накопителя (защита от поддельных флешек).", "Все данные на накопителе могут быть перезаписаны! Сделайте резервную копию."},
				{"Анализ «тяжелых» файлов", "Сканирует выбранный диск и выводит список самых крупных файлов.", ""},
			},
		},
		{
			Title: "Диагностика и Мониторинг",
			Items: []MenuItem{
				{"Убить зависшие приложения", "Принудительно завершает процессы, которые не отвечают на запросы ОС.", "Несохраненные данные в зависших программах будут утеряны."},
				{"Полная информация о ПК", "Выводит детальные характеристики процессора, ОЗУ, видеокарты и материнской платы.", ""},
			},
		},
		{
			Title: "Дополнительно",
			Items: []MenuItem{
				{"Показать скрытые файлы и папки", "Включает отображение скрытых и системных элементов в Проводнике.", ""},
				{"Показывать расширения файлов", "Включает отображение расширений для всех типов файлов.", ""},
				{"Убрать окончание «Ярлык»", "Отключает автоматическую добавку «Ярлык» при создании новых ярлыков.", ""},
				{"Отключить залипание клавиш", "Отключает вызов залипания клавиш при многократном нажатии Shift.", ""},
				{"Отключить SmartScreen", "Отключает встроенный фильтр проверки запускаемых программ.", ""},
				{"Выключить размытие на экране блокировки", "Убирает эффект размытия (Blur) при вводе пароля на экране блокировки.", ""},
			},
		},
	}
}
