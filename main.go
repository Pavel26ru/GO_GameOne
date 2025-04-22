package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Item struct {
	Name        string
	Description string
	CanTake     bool
	CanUseOn    map[string]string
	DisplayName string
	OnUseEffect map[string]func(*Game)
}

type Location struct {
	Name              string
	Description       string
	AltDescription    string
	GoDescription     string
	Items             map[string]*Item
	Exits             map[string]string
	EmptyDescription  string
	RequiresCondition string
}

type Player struct {
	CurrentLocation *Location
	Inventory       map[string]*Item
}

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
	[]gameCase{
		{1, "осмотреться", "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ. можно пройти - коридор"}, // действие осмотреться
		{2, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},                                          // действие идти
		{3, "идти комната", "ты в своей комнате. можно пройти - коридор"},
		{4, "осмотреться", "на столе: ключи, конспекты, на стуле: рюкзак. можно пройти - коридор"},
		{5, "надеть рюкзак", "вы надели: рюкзак"},                 // действие надеть
		{6, "взять ключи", "предмет добавлен в инвентарь: ключи"}, // действие взять
		{7, "взять конспекты", "предмет добавлен в инвентарь: конспекты"},
		{8, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{9, "применить ключи дверь", "дверь открыта"}, // действие применить
		{11, "идти улица", "на улице весна. можно пройти - домой"},
	},

	[]gameCase{
		{1, "осмотреться", "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ. можно пройти - коридор"},
		{2, "завтракать", "неизвестная команда"},  // придёт топать в универ голодным :(
		{3, "идти комната", "нет пути в комната"}, // через стены ходить нельзя
		{4, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{5, "применить ключи дверь", "нет предмета в инвентаре - ключи"},
		{6, "идти комната", "ты в своей комнате. можно пройти - коридор"},
		{7, "осмотреться", "на столе: ключи, конспекты, на стуле: рюкзак. можно пройти - коридор"},
		{8, "взять ключи", "некуда класть"}, // надо взять рюкзак сначала
		{9, "надеть рюкзак", "вы надели: рюкзак"},
		{10, "осмотреться", "на столе: ключи, конспекты. можно пройти - коридор"}, // состояние изменилось
		{11, "взять ключи", "предмет добавлен в инвентарь: ключи"},
		{12, "взять телефон", "нет такого"},                                // неизвестный предмет
		{13, "взять ключи", "нет такого"},                                  // предмента уже нет в комнатеы - мы его взяли
		{14, "осмотреться", "на столе: конспекты. можно пройти - коридор"}, // состояние изменилось
		{15, "взять конспекты", "предмет добавлен в инвентарь: конспекты"},
		{16, "осмотреться", "пустая комната. можно пройти - коридор"}, // состояние изменилось
		{17, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{18, "идти кухня", "кухня, ничего интересного. можно пройти - коридор"},
		{19, "осмотреться", "ты находишься на кухне, на столе: чай, надо идти в универ. можно пройти - коридор"}, // состояние изменилось
		{20, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{21, "идти улица", "дверь закрыта"},                                  //условие не удовлетворено
		{22, "применить ключи дверь", "дверь открыта"},                       //состояние изменилось
		{23, "применить телефон шкаф", "нет предмета в инвентаре - телефон"}, // нет предмета
		{24, "применить ключи шкаф", "не к чему применить"},                  // предмет есть, но применить его к этому нельзя
		{25, "идти улица", "на улице весна. можно пройти - домой"},
	},
}

func initGame() *Game {
	items := map[string]*Item{
		"рюкзак": {
			Name:        "рюкзак",
			Description: "рюкзак для вещей",
			CanTake:     true,
			DisplayName: "на стуле: рюкзак",
			CanUseOn:    map[string]string{},
			OnUseEffect: map[string]func(*Game){},
		},
		"ключи": {
			Name:        "ключи",
			Description: "ключи от двери",
			CanTake:     true,
			CanUseOn:    map[string]string{"дверь": "дверь открыта", "шкаф": "не к чему применить"},
			DisplayName: "ключи",
			OnUseEffect: map[string]func(*Game){
				"дверь": func(g *Game) { g.GlobalFlags["doorOpen"] = true },
			},
		},
		"конспекты": {
			Name:        "конспекты",
			Description: "учебные конспекты",
			CanTake:     true,
			DisplayName: "конспекты",
			CanUseOn:    map[string]string{},
			OnUseEffect: map[string]func(*Game){},
		},
		"чай": {
			Name:        "чай",
			Description: "чашка чая",
			CanTake:     true,
			DisplayName: "чай",
			CanUseOn:    map[string]string{},
			OnUseEffect: map[string]func(*Game){},
		},
	}

	locations := map[string]*Location{
		"кухня": {
			Name:              "кухня",
			Description:       "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ",
			AltDescription:    "ты находишься на кухне, на столе: чай, надо идти в универ",
			GoDescription:     "кухня, ничего интересного",
			Items:             map[string]*Item{"чай": items["чай"]},
			Exits:             map[string]string{"коридор": "коридор"},
			EmptyDescription:  "кухня, ничего интересного",
			RequiresCondition: "",
		},
		"коридор": {
			Name:              "коридор",
			Description:       "ничего интересного",
			GoDescription:     "ничего интересного",
			Items:             map[string]*Item{},
			Exits:             map[string]string{"кухня": "кухня", "комната": "комната", "улица": "улица"},
			EmptyDescription:  "ничего интересного",
			RequiresCondition: "",
		},
		"комната": {
			Name:              "комната",
			Description:       "ты в своей комнате",
			GoDescription:     "ты в своей комнате",
			Items:             map[string]*Item{"рюкзак": items["рюкзак"], "ключи": items["ключи"], "конспекты": items["конспекты"]},
			Exits:             map[string]string{"коридор": "коридор"},
			EmptyDescription:  "пустая комната",
			RequiresCondition: "",
		},
		"улица": {
			Name:              "улица",
			Description:       "на улице весна",
			GoDescription:     "на улице весна",
			Items:             map[string]*Item{},
			Exits:             map[string]string{"домой": "коридор"},
			EmptyDescription:  "на улице весна",
			RequiresCondition: "doorOpen",
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

	items := []string{}
	for _, item := range loc.Items {
		items = append(items, item.DisplayName)
	}

	var desc string
	if len(items) == 0 {
		desc = loc.EmptyDescription
	} else if loc.Name == "комната" {
		nonBackpack := []string{}
		hasBackpack := false
		for _, item := range items {
			if item == "на стуле: рюкзак" {
				hasBackpack = true
			} else {
				nonBackpack = append(nonBackpack, item)
			}
		}
		if len(nonBackpack) > 0 {
			desc = "на столе: " + strings.Join(nonBackpack, ", ")
			if hasBackpack {
				desc += ", на стуле: рюкзак"
			}
		} else if hasBackpack {
			desc = "на стуле: рюкзак"
		}
	} else {
		if _, ok := g.Player.Inventory["рюкзак"]; ok && loc.AltDescription != "" {
			desc = loc.AltDescription
		} else {
			desc = loc.Description
		}
	}

	exits := keys(loc.Exits)
	exitsStr := "можно пройти - " + strings.Join(exits, ", ")

	return fmt.Sprintf("%s. %s", desc, exitsStr)
}

func (g *Game) Go(direction string) string {
	loc := g.Player.CurrentLocation
	if nextLocName, ok := loc.Exits[direction]; ok {
		nextLoc := g.Locations[nextLocName]
		if nextLoc.RequiresCondition != "" && !g.GlobalFlags[nextLoc.RequiresCondition] {
			return "дверь закрыта"
		}
		g.Player.CurrentLocation = nextLoc
		return nextLoc.GoDescription + ". можно пройти - " + strings.Join(keys(nextLoc.Exits), ", ")
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
			if effect, ok := item.OnUseEffect[target]; ok {
				effect(g)
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
	desiredOrder := []string{"кухня", "комната", "улица", "коридор", "домой"}
	result := []string{}

	for _, key := range desiredOrder {
		if _, ok := m[key]; ok {
			result = append(result, key)
		}
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
