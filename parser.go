package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type (
	Parser struct {
		line        string
		gameCounter int
		Log         map[string]Game `json:"log"`
	}

	Game struct {
		TotalKills   int            `json:"total_kills"`
		Players      []string       `json:"players"`
		Kills        map[string]int `json:"kills"`
		KillsByMeans map[string]int `json:"kills_by_means"`
	}
)

func (p *Parser) Parse(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	p.Log = make(map[string]Game)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		p.parseLine(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (p *Parser) parseLine(line string) {
	p.line = line

	// @todo: avoid to continue to other funcs when one was already parsed

	p.initGame()
	p.addPlayer()
	p.addKill()
}

func (p *Parser) gameKey() string {
	return fmt.Sprintf("game_%02d", p.gameCounter)
}

func (p *Parser) initGame() {
	if !strings.Contains(p.line, "InitGame:") {
		return
	}

	p.gameCounter++
	if _, ok := p.Log[p.gameKey()]; !ok {
		p.Log[p.gameKey()] = Game{
			Players:      make([]string, 0),
			Kills:        make(map[string]int),
			KillsByMeans: make(map[string]int),
		}
	}
}

func (p *Parser) addPlayer() {
	if !strings.Contains(p.line, "ClientUserinfoChanged:") {
		return
	}

	game, ok := p.Log[p.gameKey()]
	if !ok {
		return
	}

	matches := regexp.MustCompile(`n\\([^\\]+)\\`).FindStringSubmatch(p.line)
	if len(matches) < 2 {
		return
	}

	playerLine := matches[1]
	for _, player := range game.Players {
		if player == playerLine {
			return
		}
	}

	game.Players = append(game.Players, playerLine)
	p.Log[p.gameKey()] = game
}

func (p *Parser) addKill() {
	if !strings.Contains(p.line, "Kill:") {
		return
	}

	game := p.Log[p.gameKey()]
	game.TotalKills++
	p.Log[p.gameKey()] = game

	matches := regexp.MustCompile(`\d+ \d+ \d+: (.+?) killed (.+?) by ([^ ]+)`).FindStringSubmatch(p.line)
	if len(matches) < 4 {
		return
	}

	killer, victim, weapon := matches[1], matches[2], matches[3]
	p.addWeaponKill(weapon)

	if killer != "<world>" {
		p.addPlayerKill(killer)
		return
	}

	p.addWorldKill(victim)
}

func (p *Parser) addWeaponKill(weapon string) {
	p.Log[p.gameKey()].KillsByMeans[weapon] += 1
}

func (p *Parser) addPlayerKill(killer string) {
	p.Log[p.gameKey()].Kills[killer] += 1
}

func (p *Parser) addWorldKill(victim string) {
	p.Log[p.gameKey()].Kills[victim] -= 1
}
