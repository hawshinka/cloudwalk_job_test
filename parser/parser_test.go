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
		name     string
		filename string
		fields   Parser
		wantErr  error
	}{
		{
			name:     "Success",
			filename: "./test/Parse_1.log",
			fields:   Parser{},
			wantErr:  nil,
		},
		{
			name:     "File not found",
			filename: "./test/Parse_2.log",
			fields:   Parser{},
			wantErr:  &fs.PathError{Op: "open", Path: "./test/Parse_2.log", Err: syscall.ENOENT},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.Parse(tt.filename)
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
				Log: map[string]Game{
					"game_01": {},
				},
			},
			want: Parser{
				errorState:  true,
				gameCounter: 1,
				Log:         map[string]Game{},
			},
		},
		{
			name: "Error State = false",
			fields: Parser{
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {},
				},
			},
			want: Parser{
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
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
			assert.Equal(t, tt.want.Log, tt.fields.Log)
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
				Log:         make(map[string]Game),
			},
			want: Parser{
				line:        " 15:00 Exit: Timelimit hit.",
				errorState:  false,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
		},
		{
			name: "Success",
			fields: Parser{
				line:        "  0:00 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0\\sv_minPing\\0\\sv_maxRate\\10000\\sv_minRate\\0\\sv_hostname\\Code Miner Server\\g_gametype\\0\\sv_privateClients\\2\\sv_maxclients\\16\\sv_allowDownload\\0\\dmflags\\0\\fraglimit\\20\\timelimit\\15\\g_maxGameClients\\0\\capturelimit\\8\\version\\ioq3 1.36 linux-x86_64 Apr 12 2009\\protocol\\68\\mapname\\q3dm17\\gamename\\baseq3\\g_needpass\\0",
				errorState:  false,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
			want: Parser{
				line:        "  0:00 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0\\sv_minPing\\0\\sv_maxRate\\10000\\sv_minRate\\0\\sv_hostname\\Code Miner Server\\g_gametype\\0\\sv_privateClients\\2\\sv_maxclients\\16\\sv_allowDownload\\0\\dmflags\\0\\fraglimit\\20\\timelimit\\15\\g_maxGameClients\\0\\capturelimit\\8\\version\\ioq3 1.36 linux-x86_64 Apr 12 2009\\protocol\\68\\mapname\\q3dm17\\gamename\\baseq3\\g_needpass\\0",
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
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
			assert.Equal(t, tt.want.Log, tt.fields.Log)
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
				Log:         make(map[string]Game),
			},
			want: Parser{
				line:        " 20:38 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\uriel/zael\\hmodel\\uriel/zael\\g_redteam\\\\g_blueteam\\\\c1\\5\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				errorState:  true,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
			expectedRes: false,
		},
		{
			name: "Error in regex",
			fields: Parser{
				line:        " 20:38 ClientUserinfoChanged: 2 ",
				errorState:  false,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
			want: Parser{
				line:        " 20:38 ClientUserinfoChanged: 2 ",
				errorState:  true,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
			expectedRes: true,
		},
		{
			name: "Player already exists in the list",
			fields: Parser{
				line:        " 20:38 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\uriel/zael\\hmodel\\uriel/zael\\g_redteam\\\\g_blueteam\\\\c1\\5\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
					"game_01": {
						Players: []string{"Isgalamido"},
					},
				},
			},
			want: Parser{
				line:        " 20:38 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\uriel/zael\\hmodel\\uriel/zael\\g_redteam\\\\g_blueteam\\\\c1\\5\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
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
				Log: map[string]Game{
					"game_01": {
						Players: []string{"Isgalamido"},
					},
				},
			},
			want: Parser{
				line:        " 20:38 ClientUserinfoChanged: 2 n\\Dono da bola\\t\\0\\model\\uriel/zael\\hmodel\\uriel/zael\\g_redteam\\\\g_blueteam\\\\c1\\5\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
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
			assert.Equal(t, tt.want.Log, tt.fields.Log)
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
				Log:         make(map[string]Game),
			},
			want: Parser{
				line:        "  3:13 Kill: 3 2 6: Isgalamido killed Dono da Bola by MOD_ROCKET",
				errorState:  true,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
			expextedRes: false,
		},
		{
			name: "Error in regex",
			fields: Parser{
				line:        "  3:13 Kill: 3 2 6: Isgalamido killed Dono da Bola ",
				errorState:  false,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
			want: Parser{
				line:        "  3:13 Kill: 3 2 6: Isgalamido killed Dono da Bola ",
				errorState:  true,
				gameCounter: 0,
				Log:         make(map[string]Game),
			},
			expextedRes: true,
		},
		{
			name: "Success with player killing himself",
			fields: Parser{
				line:        "  2:40 Kill: 2 2 7: Isgalamido killed Isgalamido by MOD_ROCKET_SPLASH",
				errorState:  false,
				gameCounter: 1,
				Log: map[string]Game{
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
				Log: map[string]Game{
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
				Log: map[string]Game{
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
				Log: map[string]Game{
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
				Log: map[string]Game{
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
				Log: map[string]Game{
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
			assert.Equal(t, tt.want.Log, tt.fields.Log)
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
				Log: map[string]Game{
					"game_01": {
						TotalKills:   0,
						KillsByMeans: make(map[string]int),
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				Log: map[string]Game{
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
				Log: map[string]Game{
					"game_01": {
						KillsByMeans: map[string]int{
							"MOD_ROCKET": 1,
						},
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				Log: map[string]Game{
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
				Log: map[string]Game{
					"game_01": {
						KillsByMeans: map[string]int{
							"MOD_ROCKET": 1,
						},
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				Log: map[string]Game{
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
			assert.Equal(t, tt.want.Log, tt.fields.Log)
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
				Log: map[string]Game{
					"game_01": {
						Kills: make(map[string]int),
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				Log: map[string]Game{
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
				Log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
						},
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				Log: map[string]Game{
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
				Log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
						},
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				Log: map[string]Game{
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
				Log: map[string]Game{
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
				Log: map[string]Game{
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
			assert.Equal(t, tt.want.Log, tt.fields.Log)
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
				Log: map[string]Game{
					"game_01": {
						Kills: make(map[string]int),
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				Log: map[string]Game{
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
				Log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": -1,
						},
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				Log: map[string]Game{
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
				Log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
						},
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				Log: map[string]Game{
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
				Log: map[string]Game{
					"game_01": {
						Kills: map[string]int{
							"Isgalamido": 1,
						},
					},
				},
			},
			want: Parser{
				gameCounter: 1,
				Log: map[string]Game{
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
			assert.Equal(t, tt.want.Log, tt.fields.Log)
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
				Log: map[string]Game{
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
				Log: map[string]Game{
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
				Log: map[string]Game{
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
				Log: map[string]Game{
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
			assert.Equal(t, tt.want.Log, tt.fields.Log)
		})
	}
}
