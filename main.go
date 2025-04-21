package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Инвентарь
type Item struct {
	Name        string
	Description string
	CanTake     bool
	CanUseOn    map[string]string
}

// Локация
type Location struct {
	Name        string
	Description string
	Items       map[string]*Item
	Exits       map[string]string
}

type Player struct {
	CurrentLocation *Location
	Inventory       map[string]*Item
}

// Состояние игры
type Game struct {
	Locations   map[string]*Location
	Player      *Player
	GlobalFlags map[string]bool
}

type gameCase struct {
	step    int
	command string
	answer  string
}

var game0cases = [][]gameCase{
	{
		{1, "осмотреться", "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ. можно пройти - коридор"},
		{2, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{3, "идти комната", "ты в своей комнате. можно пройти - коридор"},
		{4, "осмотреться", "на столе: ключи, конспекты, на стуле: рюкзак. можно пройти - коридор"},
		{5, "надеть рюкзак", "вы надели: рюкзак"},
		{6, "взять ключи", "предмет добавлен в инвентарь: ключи"},
		{7, "взять конспекты", "предмет добавлен в инвентарь: конспекты"},
		{8, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{9, "применить ключи дверь", "дверь открыта"},
		{11, "идти улица", "на улице весна. можно пройти - домой"},
	},
	{
		{1, "осмотреться", "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ. можно пройти - коридор"},
		{2, "завтракать", "неизвестная команда"},
		{3, "идти комната", "нет пути в комната"},
		{4, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{5, "применить ключи дверь", "нет предмета в инвентаре - ключи"},
		{6, "идти комната", "ты в своей комнате. можно пройти - коридор"},
		{7, "осмотреться", "на столе: ключи, конспекты, на стуле: рюкзак. можно пройти - коридор"},
		{8, "взять ключи", "некуда класть"},
		{9, "надеть рюкзак", "вы надели: рюкзак"},
		{10, "осмотреться", "на столе: ключи, конспекты. можно пройти - коридор"},
		{11, "взять ключи", "предмет добавлен в инвентарь: ключи"},
		{12, "взять телефон", "нет такого"},
		{13, "взять ключи", "нет такого"},
		{14, "осмотреться", "на столе: конспекты. можно пройти - коридор"},
		{15, "взять конспекты", "предмет добавлен в инвентарь: конспекты"},
		{16, "осмотреться", "пустая комната. можно пройти - коридор"},
		{17, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{18, "идти кухня", "кухня, ничего интересного. можно пройти - коридор"},
		{19, "осмотреться", "ты находишься на кухне, на столе: чай, надо идти в универ. можно пройти - коридор"},
		{20, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{21, "идти улица", "дверь закрыта"},
		{22, "применить ключи дверь", "дверь открыта"},
		{23, "применить телефон шкаф", "нет предмета в инвентаре - телефон"},
		{24, "применить ключи шкаф", "не к чему применить"},
		{25, "идти улица", "на улице весна. можно пройти - домой"},
	},
}

func initGame() *Game {
	items := map[string]*Item{
		"рюкзак":    {Name: "рюкзак", Description: "рюкзак для вещей", CanTake: true},
		"ключи":     {Name: "ключи", Description: "ключи от двери", CanTake: true, CanUseOn: map[string]string{"дверь": "дверь открыта", "шкаф": "не к чему применить"}},
		"конспекты": {Name: "конспекты", Description: "учебные конспекты", CanTake: true},
		"чай":       {Name: "чай", Description: "чашка чая", CanTake: true},
	}

	locations := map[string]*Location{
		"кухня": {
			Name:        "кухня",
			Description: "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ",
			Items:       map[string]*Item{"чай": items["чай"]},
			Exits:       map[string]string{"коридор": "коридор"},
		},
		"коридор": {
			Name:        "коридор",
			Description: "ничего интересного",
			Items:       map[string]*Item{},
			Exits:       map[string]string{"кухня": "кухня", "комната": "комната", "улица": "улица"},
		},
		"комната": {
			Name:        "комната",
			Description: "ты в своей комнате",
			Items:       map[string]*Item{"рюкзак": items["рюкзак"], "ключи": items["ключи"], "конспекты": items["конспекты"]},
			Exits:       map[string]string{"коридор": "коридор"},
		},
		"улица": {
			Name:        "улица",
			Description: "на улице весна",
			Items:       map[string]*Item{},
			Exits:       map[string]string{"домой": "коридор"},
		},
	}

	player := &Player{
		CurrentLocation: locations["кухня"],
		Inventory:       make(map[string]*Item),
	}

	return &Game{
		Locations:   locations,
		Player:      player,
		GlobalFlags: map[string]bool{"doorOpen": false},
	}
}

func (g *Game) Look() string {
	loc := g.Player.CurrentLocation
	if loc.Name == "кухня" {
		if _, ok := g.Player.Inventory["рюкзак"]; ok {
			return "ты находишься на кухне, на столе: чай, надо идти в универ. можно пройти - коридор"
		}
		return loc.Description + ". можно пройти - коридор"
	}
	if loc.Name == "комната" {
		items := []string{}
		if _, ok := loc.Items["ключи"]; ok {
			items = append(items, "ключи")
		}
		if _, ok := loc.Items["конспекты"]; ok {
			items = append(items, "конспекты")
		}
		if _, ok := loc.Items["рюкзак"]; ok {
			items = append(items, "на стуле: рюкзак")
		}
		if len(items) == 0 {
			return "пустая комната. можно пройти - коридор"
		}
		itemsStr := strings.Join(items, ", ")
		if len(items) > 1 && items[len(items)-1] == "на стуле: рюкзак" {
			itemsStr = "на столе: " + strings.Join(items[:len(items)-1], ", ") + ", на стуле: рюкзак"
		} else if len(items) >= 1 && items[len(items)-1] != "на стуле: рюкзак" {
			itemsStr = "на столе: " + itemsStr
		}
		return itemsStr + ". можно пройти - коридор"
	}
	return loc.Description + ". можно пройти - " + strings.Join(keys(loc.Exits), ", ")
}

func (g *Game) Go(direction string) string {
	loc := g.Player.CurrentLocation
	if nextLocName, ok := loc.Exits[direction]; ok {
		if nextLocName == "улица" && !g.GlobalFlags["doorOpen"] {
			return "дверь закрыта"
		}
		g.Player.CurrentLocation = g.Locations[nextLocName]
		if nextLocName == "кухня" {
			return "кухня, ничего интересного. можно пройти - коридор"
		}
		return g.Locations[nextLocName].Description + ". можно пройти - " + strings.Join(keys(g.Locations[nextLocName].Exits), ", ")
	}
	return fmt.Sprintf("нет пути в %s", direction)
}

func (g *Game) Take(itemName string) string {
	if _, ok := g.Player.Inventory["рюкзак"]; !ok {
		return "некуда класть"
	}
	loc := g.Player.CurrentLocation
	if item, ok := loc.Items[itemName]; ok && item.CanTake {
		g.Player.Inventory[itemName] = item
		delete(loc.Items, itemName)
		return fmt.Sprintf("предмет добавлен в инвентарь: %s", itemName)
	}
	return "нет такого"
}

func (g *Game) Use(itemName, target string) string {
	if item, ok := g.Player.Inventory[itemName]; ok {
		if result, ok := item.CanUseOn[target]; ok {
			if itemName == "ключи" && target == "дверь" {
				g.GlobalFlags["doorOpen"] = true
			}
			return result
		}
		return "не к чему применить"
	}
	return fmt.Sprintf("нет предмета в инвентаре - %s", itemName)
}

func (g *Game) Wear(itemName string) string {
	loc := g.Player.CurrentLocation
	if item, ok := loc.Items[itemName]; ok && itemName == "рюкзак" {
		g.Player.Inventory[itemName] = item
		delete(loc.Items, itemName)
		return fmt.Sprintf("вы надели: %s", itemName)
	}
	return "нет такого"
}

func (g *Game) HandleCommand(command string) string {
	words := strings.Split(strings.ToLower(command), " ")
	if len(words) == 0 {
		return "неизвестная команда"
	}

	switch words[0] {
	case "осмотреться":
		return g.Look()
	case "идти":
		if len(words) > 1 {
			return g.Go(words[1])
		}
	case "взять":
		if len(words) > 1 {
			return g.Take(words[1])
		}
	case "применить":
		if len(words) > 2 {
			return g.Use(words[1], words[2])
		}
	case "надеть":
		if len(words) > 1 {
			return g.Wear(words[1])
		}
	}
	return "неизвестная команда"
}

func keys(m map[string]string) []string {
	result := []string{}
	for k := range m {
		result = append(result, k)
	}
	return result
}

func main() {
	game := initGame()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Введите команду: ")

	for scanner.Scan() {
		command := scanner.Text()
		response := game.HandleCommand(command)
		fmt.Println(response)
		fmt.Println("Введите команду: ")
	}
}
