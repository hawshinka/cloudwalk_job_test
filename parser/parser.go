package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type (
	Parser struct {
		line        string
		errorState  bool
		gameCounter int
		log         map[string]Game `json:"log"`
	}

	Game struct {
		TotalKills   int            `json:"total_kills"`
		Players      []string       `json:"players"`
		Kills        map[string]int `json:"kills"`
		KillsByMeans map[string]int `json:"kills_by_means"`
	}
)

func (p *Parser) Parse(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	p.log = make(map[string]Game)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		p.parseLine(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	out, _ := json.Marshal(p.log)
	return string(out), nil
}

func (p *Parser) parseLine(line string) {
	p.line = line

	p.checkErrorState()
	if p.initGame() || p.addPlayer() || p.addKill() {
		return
	}
}

func (p *Parser) gameKey() string {
	return fmt.Sprintf("game_%02d", p.gameCounter)
}

func (p *Parser) checkErrorState() {
	if p.errorState {
		if _, ok := p.log[p.gameKey()]; ok {
			delete(p.log, p.gameKey())
		}
	}
}

func (p *Parser) initGame() bool {
	if !strings.Contains(p.line, "InitGame:") {
		return false
	}

	p.errorState = false
	p.gameCounter++
	if _, ok := p.log[p.gameKey()]; !ok {
		p.log[p.gameKey()] = Game{
			Players:      make([]string, 0),
			Kills:        make(map[string]int),
			KillsByMeans: make(map[string]int),
		}
	}

	return true
}

func (p *Parser) addPlayer() bool {
	if p.errorState || !strings.Contains(p.line, "ClientUserinfoChanged:") {
		return false
	}

	matches := regexp.MustCompile(`ClientUserinfoChanged: \d+ n\\([^\\]+)\\`).FindStringSubmatch(p.line)
	if len(matches) < 2 {
		p.errorState = true
		return true
	}

	newPlayerName := matches[1]
	game := p.log[p.gameKey()]
	for _, existingPlayer := range game.Players {
		if existingPlayer == newPlayerName {
			return true
		}
	}

	game.Players = append(game.Players, newPlayerName)
	p.log[p.gameKey()] = game
	return true
}

func (p *Parser) addKill() bool {
	if p.errorState || !strings.Contains(p.line, "Kill:") {
		return false
	}

	matches := regexp.MustCompile(`Kill: \d+ \d+ \d+: (.+?) killed (.+?) by ([^ ]+)`).FindStringSubmatch(p.line)
	if len(matches) < 4 {
		p.errorState = true
		return true
	}

	game := p.log[p.gameKey()]
	game.TotalKills++
	p.log[p.gameKey()] = game

	killer, victim, weapon := matches[1], matches[2], matches[3]
	p.addWeaponKill(weapon)

	if killer == victim {
		return true
	}

	if killer != "<world>" {
		p.addPlayerKill(killer)
		return true
	}

	p.addWorldKill(victim)
	return true
}

func (p *Parser) addWeaponKill(weapon string) {
	p.log[p.gameKey()].KillsByMeans[weapon]++
}

func (p *Parser) addPlayerKill(killer string) {
	p.log[p.gameKey()].Kills[killer]++
	p.handleZeroKills(killer)
}

func (p *Parser) addWorldKill(victim string) {
	p.log[p.gameKey()].Kills[victim]--
	p.handleZeroKills(victim)
}

func (p *Parser) handleZeroKills(player string) {
	if _, ok := p.log[p.gameKey()].Kills[player]; ok {
		if p.log[p.gameKey()].Kills[player] == 0 {
			delete(p.log[p.gameKey()].Kills, player)
		}
	}
}
