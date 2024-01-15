//go:build unit

package parser

import (
	"io/fs"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name       string
		filename   string
		fields     Parser
		wantParsed string
		wantErr    error
	}{
		{
			name:       "Success",
			filename:   "./test/Parse_1.log",
			fields:     Parser{},
			wantParsed: "{\"game_01\":{\"total_kills\":0,\"players\":[\"Isgalamido\"],\"kills\":{},\"kills_by_means\":{}},\"game_02\":{\"total_kills\":11,\"players\":[\"Isgalamido\",\"Dono da Bola\",\"Mocinha\"],\"kills\":{\"Isgalamido\":-7},\"kills_by_means\":{\"MOD_FALLING\":1,\"MOD_ROCKET_SPLASH\":3,\"MOD_TRIGGER_HURT\":7}},\"game_03\":{\"total_kills\":4,\"players\":[\"Dono da Bola\",\"Mocinha\",\"Isgalamido\",\"Zeh\"],\"kills\":{\"Dono da Bola\":-1,\"Isgalamido\":1,\"Zeh\":-2},\"kills_by_means\":{\"MOD_FALLING\":1,\"MOD_ROCKET\":1,\"MOD_TRIGGER_HURT\":2}},\"game_04\":{\"total_kills\":105,\"players\":[\"Dono da Bola\",\"Isgalamido\",\"Zeh\",\"Assasinu Credi\"],\"kills\":{\"Assasinu Credi\":12,\"Dono da Bola\":9,\"Isgalamido\":19,\"Zeh\":20},\"kills_by_means\":{\"MOD_FALLING\":11,\"MOD_MACHINEGUN\":4,\"MOD_RAILGUN\":8,\"MOD_ROCKET\":20,\"MOD_ROCKET_SPLASH\":51,\"MOD_SHOTGUN\":2,\"MOD_TRIGGER_HURT\":9}},\"game_05\":{\"total_kills\":14,\"players\":[\"Dono da Bola\",\"Isgalamido\",\"Zeh\",\"Assasinu Credi\"],\"kills\":{\"Assasinu Credi\":-1,\"Isgalamido\":2,\"Zeh\":1},\"kills_by_means\":{\"MOD_RAILGUN\":1,\"MOD_ROCKET\":4,\"MOD_ROCKET_SPLASH\":4,\"MOD_TRIGGER_HURT\":5}},\"game_06\":{\"total_kills\":29,\"players\":[\"Fasano Again\",\"Oootsimo\",\"Isgalamido\",\"Zeh\",\"Dono da Bola\",\"UnnamedPlayer\",\"Maluquinho\",\"Assasinu Credi\",\"Mal\"],\"kills\":{\"Assasinu Credi\":1,\"Dono da Bola\":2,\"Isgalamido\":3,\"Oootsimo\":8,\"Zeh\":7},\"kills_by_means\":{\"MOD_FALLING\":1,\"MOD_MACHINEGUN\":1,\"MOD_RAILGUN\":2,\"MOD_ROCKET\":5,\"MOD_ROCKET_SPLASH\":13,\"MOD_SHOTGUN\":4,\"MOD_TRIGGER_HURT\":3}},\"game_07\":{\"total_kills\":130,\"players\":[\"Oootsimo\",\"Isgalamido\",\"Zeh\",\"Dono da Bola\",\"Mal\",\"Assasinu Credi\",\"Chessus!\",\"Chessus\"],\"kills\":{\"Assasinu Credi\":19,\"Dono da Bola\":10,\"Isgalamido\":14,\"Mal\":-3,\"Oootsimo\":20,\"Zeh\":8},\"kills_by_means\":{\"MOD_FALLING\":7,\"MOD_MACHINEGUN\":9,\"MOD_RAILGUN\":9,\"MOD_ROCKET\":29,\"MOD_ROCKET_SPLASH\":49,\"MOD_SHOTGUN\":7,\"MOD_TRIGGER_HURT\":20}},\"game_08\":{\"total_kills\":89,\"players\":[\"Oootsimo\",\"Isgalamido\",\"Zeh\",\"Dono da Bola\",\"Mal\",\"Assasinu Credi\"],\"kills\":{\"Assasinu Credi\":9,\"Dono da Bola\":1,\"Isgalamido\":20,\"Mal\":-3,\"Oootsimo\":15,\"Zeh\":12},\"kills_by_means\":{\"MOD_FALLING\":6,\"MOD_MACHINEGUN\":4,\"MOD_RAILGUN\":12,\"MOD_ROCKET\":18,\"MOD_ROCKET_SPLASH\":39,\"MOD_SHOTGUN\":1,\"MOD_TRIGGER_HURT\":9}},\"game_09\":{\"total_kills\":67,\"players\":[\"Oootsimo\",\"Isgalamido\",\"Zeh\",\"Dono da Bola\",\"Mal\",\"Assasinu Credi\",\"Chessus!\",\"Chessus\"],\"kills\":{\"Assasinu Credi\":7,\"Chessus\":8,\"Dono da Bola\":1,\"Isgalamido\":1,\"Mal\":2,\"Oootsimo\":8,\"Zeh\":12},\"kills_by_means\":{\"MOD_FALLING\":3,\"MOD_MACHINEGUN\":3,\"MOD_RAILGUN\":10,\"MOD_ROCKET\":17,\"MOD_ROCKET_SPLASH\":25,\"MOD_SHOTGUN\":1,\"MOD_TRIGGER_HURT\":8}},\"game_10\":{\"total_kills\":60,\"players\":[\"Oootsimo\",\"Dono da Bola\",\"Zeh\",\"Chessus\",\"Mal\",\"Assasinu Credi\",\"Isgalamido\"],\"kills\":{\"Assasinu Credi\":3,\"Chessus\":5,\"Dono da Bola\":3,\"Isgalamido\":5,\"Mal\":1,\"Oootsimo\":-1,\"Zeh\":7},\"kills_by_means\":{\"MOD_BFG\":2,\"MOD_BFG_SPLASH\":2,\"MOD_CRUSH\":1,\"MOD_MACHINEGUN\":1,\"MOD_RAILGUN\":7,\"MOD_ROCKET\":4,\"MOD_ROCKET_SPLASH\":1,\"MOD_TELEFRAG\":25,\"MOD_TRIGGER_HURT\":17}},\"game_11\":{\"total_kills\":20,\"players\":[\"Dono da Bola\",\"Isgalamido\",\"Zeh\",\"Oootsimo\",\"Chessus\",\"Assasinu Credi\",\"UnnamedPlayer\",\"Mal\"],\"kills\":{\"Assasinu Credi\":-3,\"Dono da Bola\":-2,\"Isgalamido\":4,\"Oootsimo\":4},\"kills_by_means\":{\"MOD_BFG_SPLASH\":3,\"MOD_CRUSH\":1,\"MOD_MACHINEGUN\":1,\"MOD_RAILGUN\":4,\"MOD_ROCKET_SPLASH\":4,\"MOD_TRIGGER_HURT\":7}},\"game_12\":{\"total_kills\":160,\"players\":[\"Isgalamido\",\"Dono da Bola\",\"Zeh\",\"Oootsimo\",\"Chessus\",\"Assasinu Credi\",\"Mal\"],\"kills\":{\"Assasinu Credi\":18,\"Chessus\":12,\"Dono da Bola\":3,\"Isgalamido\":24,\"Mal\":-7,\"Oootsimo\":12,\"Zeh\":11},\"kills_by_means\":{\"MOD_BFG\":8,\"MOD_BFG_SPLASH\":8,\"MOD_FALLING\":2,\"MOD_MACHINEGUN\":7,\"MOD_RAILGUN\":38,\"MOD_ROCKET\":25,\"MOD_ROCKET_SPLASH\":35,\"MOD_TRIGGER_HURT\":37}},\"game_13\":{\"total_kills\":6,\"players\":[\"Isgalamido\",\"Dono da Bola\",\"Zeh\",\"Oootsimo\",\"Chessus\",\"Assasinu Credi\",\"Mal\"],\"kills\":{\"Dono da Bola\":-1,\"Isgalamido\":-1,\"Oootsimo\":1,\"Zeh\":2},\"kills_by_means\":{\"MOD_BFG\":1,\"MOD_BFG_SPLASH\":1,\"MOD_ROCKET\":1,\"MOD_ROCKET_SPLASH\":1,\"MOD_TRIGGER_HURT\":2}},\"game_14\":{\"total_kills\":122,\"players\":[\"Isgalamido\",\"Dono da Bola\",\"Zeh\",\"Oootsimo\",\"Chessus\",\"Assasinu Credi\",\"Mal\"],\"kills\":{\"Assasinu Credi\":3,\"Chessus\":7,\"Dono da Bola\":1,\"Isgalamido\":22,\"Mal\":-5,\"Oootsimo\":9,\"Zeh\":4},\"kills_by_means\":{\"MOD_BFG\":5,\"MOD_BFG_SPLASH\":10,\"MOD_FALLING\":5,\"MOD_MACHINEGUN\":4,\"MOD_RAILGUN\":20,\"MOD_ROCKET\":23,\"MOD_ROCKET_SPLASH\":24,\"MOD_TRIGGER_HURT\":31}},\"game_15\":{\"total_kills\":3,\"players\":[\"Zeh\",\"Assasinu Credi\",\"Dono da Bola\",\"Fasano Again\",\"Isgalamido\",\"Oootsimo\"],\"kills\":{\"Zeh\":-3},\"kills_by_means\":{\"MOD_TRIGGER_HURT\":3}},\"game_16\":{\"total_kills\":0,\"players\":[\"Dono da Bola\",\"Oootsimo\",\"Isgalamido\",\"Assasinu Credi\",\"Zeh\"],\"kills\":{},\"kills_by_means\":{}},\"game_17\":{\"total_kills\":13,\"players\":[\"Dono da Bola\",\"Oootsimo\",\"Isgalamido\",\"Assasinu Credi\",\"Zeh\",\"UnnamedPlayer\",\"Mal\"],\"kills\":{\"Assasinu Credi\":-3,\"Dono da Bola\":-2,\"Mal\":-1},\"kills_by_means\":{\"MOD_FALLING\":3,\"MOD_RAILGUN\":2,\"MOD_ROCKET_SPLASH\":2,\"MOD_TRIGGER_HURT\":6}},\"game_18\":{\"total_kills\":7,\"players\":[\"Dono da Bola\",\"Oootsimo\",\"Isgalamido\",\"Assasinu Credi\",\"Zeh\",\"Mal\"],\"kills\":{\"Assasinu Credi\":2,\"Dono da Bola\":-1,\"Isgalamido\":1,\"Mal\":-1,\"Zeh\":2},\"kills_by_means\":{\"MOD_FALLING\":1,\"MOD_ROCKET\":1,\"MOD_ROCKET_SPLASH\":4,\"MOD_TRIGGER_HURT\":1}},\"game_19\":{\"total_kills\":95,\"players\":[\"Isgalamido\",\"Oootsimo\",\"Dono da Bola\",\"Assasinu Credi\",\"Zeh\",\"Mal\"],\"kills\":{\"Assasinu Credi\":8,\"Dono da Bola\":12,\"Isgalamido\":13,\"Mal\":2,\"Oootsimo\":10,\"Zeh\":20},\"kills_by_means\":{\"MOD_FALLING\":1,\"MOD_MACHINEGUN\":7,\"MOD_RAILGUN\":10,\"MOD_ROCKET\":27,\"MOD_ROCKET_SPLASH\":32,\"MOD_SHOTGUN\":6,\"MOD_TRIGGER_HURT\":12}},\"game_20\":{\"total_kills\":3,\"players\":[\"Isgalamido\",\"Oootsimo\",\"Dono da Bola\",\"Assasinu Credi\",\"Zeh\",\"Mal\"],\"kills\":{\"Dono da Bola\":1,\"Oootsimo\":1},\"kills_by_means\":{\"MOD_ROCKET\":1,\"MOD_ROCKET_SPLASH\":2}},\"game_21\":{\"total_kills\":131,\"players\":[\"Isgalamido\",\"Oootsimo\",\"Dono da Bola\",\"Assasinu Credi\",\"Zeh\",\"Mal\"],\"kills\":{\"Assasinu Credi\":16,\"Dono da Bola\":12,\"Isgalamido\":17,\"Mal\":6,\"Oootsimo\":21,\"Zeh\":19},\"kills_by_means\":{\"MOD_FALLING\":3,\"MOD_MACHINEGUN\":4,\"MOD_RAILGUN\":9,\"MOD_ROCKET\":37,\"MOD_ROCKET_SPLASH\":60,\"MOD_SHOTGUN\":4,\"MOD_TRIGGER_HURT\":14}}}",
			wantErr:    nil,
		},
		{
			name:       "File not found",
			filename:   "./test/Parse_2.log",
			fields:     Parser{},
			wantParsed: "",
			wantErr:    &fs.PathError{Op: "open", Path: "./test/Parse_2.log", Err: syscall.ENOENT},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := tt.fields.Parse(tt.filename)
			assert.Equal(t, tt.wantParsed, parsed)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestParser_gameKey(t *testing.T) {
	tests := []struct {
		name   string
		fields Parser
		want   string
	}{
		{
			name:   "Success",
			fields: Parser{gameCounter: 5},
			want:   "game_05",
		},
		{
			name:   "Success",
			fields: Parser{},
			want:   "game_00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := tt.fields.gameKey()
			assert.Equal(t, tt.want, key)
		})
	}
}

func TestParser_checkErrorState(t *testing.T) {
	tests := []struct {
		name   string
		fields Parser
		want   Parser
	}{
		{
			name: "Error State = true",
			fields: Parser{
				errorState:  true,
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {},
				},
			},
			want: Parser{
				errorState:  true,
				gameCounter: 1,
				log:         map[string]Game{},
			},
		},
		{
			name: "Error State = false",
			fields: Parser{
				errorState:  false,
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {},
				},
			},
			want: Parser{
				errorState:  false,
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.checkErrorState()
			assert.Equal(t, tt.want.line, tt.fields.line)
			assert.Equal(t, tt.want.errorState, tt.fields.errorState)
			assert.Equal(t, tt.want.gameCounter, tt.fields.gameCounter)
			assert.Equal(t, tt.want.log, tt.fields.log)
		})
	}
}

func TestParser_initGame(t *testing.T) {
	tests := []struct {
		name   string
		fields Parser
		want   Parser
	}{
		{
			name: "Not InitGame line",
			fields: Parser{
				line:        " 15:00 Exit: Timelimit hit.",
				errorState:  false,
				gameCounter: 0,
				log:         make(map[string]Game),
			},
			want: Parser{
				line:        " 15:00 Exit: Timelimit hit.",
				errorState:  false,
				gameCounter: 0,
				log:         make(map[string]Game),
			},
		},
		{
			name: "Success",
			fields: Parser{
				line:        "  0:00 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0\\sv_minPing\\0\\sv_maxRate\\10000\\sv_minRate\\0\\sv_hostname\\Code Miner Server\\g_gametype\\0\\sv_privateClients\\2\\sv_maxclients\\16\\sv_allowDownload\\0\\dmflags\\0\\fraglimit\\20\\timelimit\\15\\g_maxGameClients\\0\\capturelimit\\8\\version\\ioq3 1.36 linux-x86_64 Apr 12 2009\\protocol\\68\\mapname\\q3dm17\\gamename\\baseq3\\g_needpass\\0",
				errorState:  false,
				gameCounter: 0,
				log:         make(map[string]Game),
			},
			want: Parser{
				line:        "  0:00 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0\\sv_minPing\\0\\sv_maxRate\\10000\\sv_minRate\\0\\sv_hostname\\Code Miner Server\\g_gametype\\0\\sv_privateClients\\2\\sv_maxclients\\16\\sv_allowDownload\\0\\dmflags\\0\\fraglimit\\20\\timelimit\\15\\g_maxGameClients\\0\\capturelimit\\8\\version\\ioq3 1.36 linux-x86_64 Apr 12 2009\\protocol\\68\\mapname\\q3dm17\\gamename\\baseq3\\g_needpass\\0",
				errorState:  false,
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						TotalKills:   0,
						Players:      make([]string, 0),
						Kills:        make(map[string]int),
						KillsByMeans: make(map[string]int),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.initGame()
			assert.Equal(t, tt.want.line, tt.fields.line)
			assert.Equal(t, tt.want.errorState, tt.fields.errorState)
			assert.Equal(t, tt.want.gameCounter, tt.fields.gameCounter)
			assert.Equal(t, tt.want.log, tt.fields.log)
		})
	}
}

func TestParser_addPlayer(t *testing.T) {
	tests := []struct {
		name        string
		fields      Parser
		want        Parser
		expectedRes bool
	}{
		{
			name:        "Not ClientUserinfoChanged line",
			fields:      Parser{},
			want:        Parser{},
			expectedRes: false,
		},
		{
			name: "Game in error state",
			fields: Parser{
				line:        " 20:38 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\uriel/zael\\hmodel\\uriel/zael\\g_redteam\\\\g_blueteam\\\\c1\\5\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				errorState:  true,
				gameCounter: 0,
				log:         make(map[string]Game),
			},
			want: Parser{
				line:        " 20:38 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\uriel/zael\\hmodel\\uriel/zael\\g_redteam\\\\g_blueteam\\\\c1\\5\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				errorState:  true,
				gameCounter: 0,
				log:         make(map[string]Game),
			},
			expectedRes: false,
		},
		{
			name: "Error in regex",
			fields: Parser{
				line:        " 20:38 ClientUserinfoChanged: 2 ",
				errorState:  false,
				gameCounter: 0,
				log:         make(map[string]Game),
			},
			want: Parser{
				line:        " 20:38 ClientUserinfoChanged: 2 ",
				errorState:  true,
				gameCounter: 0,
				log:         make(map[string]Game),
			},
			expectedRes: true,
		},
		{
			name: "Player already exists in the list",
			fields: Parser{
				line:        " 20:38 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\uriel/zael\\hmodel\\uriel/zael\\g_redteam\\\\g_blueteam\\\\c1\\5\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				errorState:  false,
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Players: []string{"Isgalamido"},
					},
				},
			},
			want: Parser{
				line:        " 20:38 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\uriel/zael\\hmodel\\uriel/zael\\g_redteam\\\\g_blueteam\\\\c1\\5\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				errorState:  false,
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Players: []string{"Isgalamido"},
					},
				},
			},
			expectedRes: true,
		},
		{
			name: "New player to the list",
			fields: Parser{
				line:        " 20:38 ClientUserinfoChanged: 2 n\\Dono da bola\\t\\0\\model\\uriel/zael\\hmodel\\uriel/zael\\g_redteam\\\\g_blueteam\\\\c1\\5\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				errorState:  false,
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Players: []string{"Isgalamido"},
					},
				},
			},
			want: Parser{
				line:        " 20:38 ClientUserinfoChanged: 2 n\\Dono da bola\\t\\0\\model\\uriel/zael\\hmodel\\uriel/zael\\g_redteam\\\\g_blueteam\\\\c1\\5\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				errorState:  false,
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Players: []string{"Isgalamido", "Dono da bola"},
					},
				},
			},
			expectedRes: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.fields.addPlayer()
			assert.Equal(t, tt.expectedRes, res)
			assert.Equal(t, tt.want.line, tt.fields.line)
			assert.Equal(t, tt.want.errorState, tt.fields.errorState)
			assert.Equal(t, tt.want.gameCounter, tt.fields.gameCounter)
			assert.Equal(t, tt.want.log, tt.fields.log)
		})
	}
}

func TestParser_addKill(t *testing.T) {
	tests := []struct {
		name        string
		fields      Parser
		want        Parser
		expextedRes bool
	}{
		{
			name:        "Not Kill line",
			fields:      Parser{},
			want:        Parser{},
			expextedRes: false,
		},
		{
			name: "Game in error state",
			fields: Parser{
				line:        "  3:13 Kill: 3 2 6: Isgalamido killed Dono da Bola by MOD_ROCKET",
				errorState:  true,
				gameCounter: 0,
				log:         make(map[string]Game),
			},
			want: Parser{
				line:        "  3:13 Kill: 3 2 6: Isgalamido killed Dono da Bola by MOD_ROCKET",
				errorState:  true,
				gameCounter: 0,
				log:         make(map[string]Game),
			},
			expextedRes: false,
		},
		{
			name: "Error in regex",
			fields: Parser{
				line:        "  3:13 Kill: 3 2 6: Isgalamido killed Dono da Bola ",
				errorState:  false,
				gameCounter: 0,
				log:         make(map[string]Game),
			},
			want: Parser{
				line:        "  3:13 Kill: 3 2 6: Isgalamido killed Dono da Bola ",
				errorState:  true,
				gameCounter: 0,
				log:         make(map[string]Game),
			},
			expextedRes: true,
		},
		{
			name: "Success with player killing himself",
			fields: Parser{
				line:        "  2:40 Kill: 2 2 7: Isgalamido killed Isgalamido by MOD_ROCKET_SPLASH",
				errorState:  false,
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						TotalKills:   0,
						Kills:        make(map[string]int),
						KillsByMeans: make(map[string]int),
					},
				},
			},
			want: Parser{
				line:        "  2:40 Kill: 2 2 7: Isgalamido killed Isgalamido by MOD_ROCKET_SPLASH",
				errorState:  false,
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						TotalKills: 1,
						Kills:      make(map[string]int),
						KillsByMeans: map[string]int{
							"MOD_ROCKET_SPLASH": 1,
						},
					},
				},
			},
			expextedRes: true,
		},
		{
			name: "Success with player kill",
			fields: Parser{
				line:        "  3:13 Kill: 3 2 6: Isgalamido killed Dono da Bola by MOD_ROCKET",
				errorState:  false,
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						TotalKills:   0,
						Kills:        make(map[string]int),
						KillsByMeans: make(map[string]int),
					},
				},
			},
			want: Parser{
				line:        "  3:13 Kill: 3 2 6: Isgalamido killed Dono da Bola by MOD_ROCKET",
				errorState:  false,
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						TotalKills: 1,
						Kills: map[string]int{
							"Isgalamido": 1,
						},
						KillsByMeans: map[string]int{
							"MOD_ROCKET": 1,
						},
					},
				},
			},
			expextedRes: true,
		},
		{
			name: "Success with world kill",
			fields: Parser{
				line:        "  3:27 Kill: 1022 3 22: <world> killed Isgalamido by MOD_TRIGGER_HURT",
				errorState:  false,
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						TotalKills:   0,
						Kills:        make(map[string]int),
						KillsByMeans: make(map[string]int),
					},
				},
			},
			want: Parser{
				line:        "  3:27 Kill: 1022 3 22: <world> killed Isgalamido by MOD_TRIGGER_HURT",
				errorState:  false,
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						TotalKills: 1,
						Kills: map[string]int{
							"Isgalamido": -1,
						},
						KillsByMeans: map[string]int{
							"MOD_TRIGGER_HURT": 1,
						},
					},
				},
			},
			expextedRes: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.fields.addKill()
			assert.Equal(t, tt.expextedRes, res)
			assert.Equal(t, tt.want.line, tt.fields.line)
			assert.Equal(t, tt.want.errorState, tt.fields.errorState)
			assert.Equal(t, tt.want.gameCounter, tt.fields.gameCounter)
			assert.Equal(t, tt.want.log, tt.fields.log)
		})
	}
}

func TestParser_addWeaponKill(t *testing.T) {
	tests := []struct {
		name   string
		weapon string
		fields Parser
		want   Parser
	}{
		{
			name:   "Success",
			weapon: "MOD_ROCKET",
			fields: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						TotalKills:   0,
						KillsByMeans: make(map[string]int),
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						KillsByMeans: map[string]int{
							"MOD_ROCKET": 1,
						},
					},
				},
			},
		},
		{
			name:   "Success with same weapon",
			weapon: "MOD_ROCKET",
			fields: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						KillsByMeans: map[string]int{
							"MOD_ROCKET": 1,
						},
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						KillsByMeans: map[string]int{
							"MOD_ROCKET": 2,
						},
					},
				},
			},
		},
		{
			name:   "Success with different weapon",
			weapon: "MOD_ROCKET_SPLASH",
			fields: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						KillsByMeans: map[string]int{
							"MOD_ROCKET": 1,
						},
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						KillsByMeans: map[string]int{
							"MOD_ROCKET":        1,
							"MOD_ROCKET_SPLASH": 1,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.addWeaponKill(tt.weapon)
			assert.Equal(t, tt.want.line, tt.fields.line)
			assert.Equal(t, tt.want.errorState, tt.fields.errorState)
			assert.Equal(t, tt.want.gameCounter, tt.fields.gameCounter)
			assert.Equal(t, tt.want.log, tt.fields.log)
		})
	}
}

func TestParser_addPlayerKill(t *testing.T) {
	tests := []struct {
		name   string
		killer string
		fields Parser
		want   Parser
	}{
		{
			name:   "Success",
			killer: "Isgalamido",
			fields: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: make(map[string]int),
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
						},
					},
				},
			},
		},
		{
			name:   "Success with a second hit from the same player",
			killer: "Isgalamido",
			fields: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
						},
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 2,
						},
					},
				},
			},
		},
		{
			name:   "Success with a hit from a different player",
			killer: "Zeh",
			fields: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
						},
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
							"Zeh":        1,
						},
					},
				},
			},
		},
		{
			name:   "Success with a hit from a player with -1 kills",
			killer: "Zeh",
			fields: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
							"Zeh":        -1,
						},
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.addPlayerKill(tt.killer)
			assert.Equal(t, tt.want.line, tt.fields.line)
			assert.Equal(t, tt.want.errorState, tt.fields.errorState)
			assert.Equal(t, tt.want.gameCounter, tt.fields.gameCounter)
			assert.Equal(t, tt.want.log, tt.fields.log)
		})
	}
}

func TestParser_addWorldKill(t *testing.T) {
	tests := []struct {
		name   string
		victim string
		fields Parser
		want   Parser
	}{
		{
			name:   "Success",
			victim: "Isgalamido",
			fields: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: make(map[string]int),
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": -1,
						},
					},
				},
			},
		},
		{
			name:   "Success with a second hit to the same player",
			victim: "Isgalamido",
			fields: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": -1,
						},
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": -2,
						},
					},
				},
			},
		},
		{
			name:   "Success with a hit to a different player",
			victim: "Zeh",
			fields: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
						},
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
							"Zeh":        -1,
						},
					},
				},
			},
		},
		{
			name:   "Success with a hit to a player with 1 kill",
			victim: "Isgalamido",
			fields: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
						},
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: make(map[string]int),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.addWorldKill(tt.victim)
			assert.Equal(t, tt.want.line, tt.fields.line)
			assert.Equal(t, tt.want.errorState, tt.fields.errorState)
			assert.Equal(t, tt.want.gameCounter, tt.fields.gameCounter)
			assert.Equal(t, tt.want.log, tt.fields.log)
		})
	}
}

func TestParser_handleZeroKills(t *testing.T) {
	tests := []struct {
		name   string
		victim string
		fields Parser
		want   Parser
	}{
		{
			name:   "No zero kills",
			victim: "Assasinu Credi",
			fields: Parser{
				line:        "  8:30 Kill: 1022 5 22: <world> killed Assasinu Credi by MOD_TRIGGER_HURT",
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Assasinu Credi": 2,
						},
						KillsByMeans: map[string]int{
							"MOD_ROCKET_SPLASH": 2,
						},
					},
				},
			},
			want: Parser{
				line:        "  8:30 Kill: 1022 5 22: <world> killed Assasinu Credi by MOD_TRIGGER_HURT",
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Assasinu Credi": 2,
						},
						KillsByMeans: map[string]int{
							"MOD_ROCKET_SPLASH": 2,
						},
					},
				},
			},
		},
		{
			name:   "Zero kills",
			victim: "Assasinu Credi",
			fields: Parser{
				line:        "  8:30 Kill: 1022 5 22: <world> killed Assasinu Credi by MOD_TRIGGER_HURT",
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Assasinu Credi": 0,
						},
						KillsByMeans: map[string]int{
							"MOD_ROCKET_SPLASH": 2,
						},
					},
				},
			},
			want: Parser{
				line:        "  8:30 Kill: 1022 5 22: <world> killed Assasinu Credi by MOD_TRIGGER_HURT",
				gameCounter: 1,
				log: map[string]Game{
					"game_01": {
						Kills: make(map[string]int),
						KillsByMeans: map[string]int{
							"MOD_ROCKET_SPLASH": 2,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.handleZeroKills(tt.victim)
			assert.Equal(t, tt.want.line, tt.fields.line)
			assert.Equal(t, tt.want.errorState, tt.fields.errorState)
			assert.Equal(t, tt.want.gameCounter, tt.fields.gameCounter)
			assert.Equal(t, tt.want.log, tt.fields.log)
		})
	}
}
